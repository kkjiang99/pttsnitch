package ptt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis/v7"
)

type Settings struct {
	Host            string
	AllpostURL      string
	FchLnkExHr      int
	FetchExtPageNum int
	LinksChCap      int
	PostsChCap      int
	RedisHost       string
	RedisPort       int
}

type pttcrawler struct {
	host     string
	settings Settings
}

type post struct {
	Title   string
	Link    string
	Author  string
	Board   string
	Content string
}

var httpClient  *http.Client
var redisClient *redis.Client

// NewPttCrawler pttcrawler contructor
func NewPttCrawler(settings Settings) *pttcrawler {
	pttcrawler := pttcrawler{host: settings.Host, settings: settings}
	httpClient = &http.Client{}

	return &pttcrawler
}

func (pttcrawler *pttcrawler) Run(done <-chan bool) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", pttcrawler.settings.RedisHost, pttcrawler.settings.RedisPort),
		Password: "",
		DB:       0,
		DialTimeout: time.Second * 2,
	})

	// 取得post links
	postLinks := make(chan string, pttcrawler.settings.LinksChCap)
	go pttcrawler.getLinks(postLinks)

	// 爬post links，進postQ
	postQ := make(chan post, pttcrawler.settings.PostsChCap)
	for i := 0; i < 2; i++ {
		go pttcrawler.fetchPosts(postLinks, postQ)
	}

	// 從postQ取出比對關鍵字，符合的進redis notifies Queue
	for i := 0; i < 2; i++ {
		go pttcrawler.pickup(postQ)
	}

	// blocking
	for {
		select {
			case <-done:
				break
			default:
		}
	}
}

func (pttcrawler *pttcrawler) getLinks(postLinks chan<- string) {
	for {
		log.Println("getLinks")
		// 先拿allpost首頁
		page := pttcrawler.GetHTML(pttcrawler.host + pttcrawler.settings.AllpostURL)
		linkSlice := parseAPLinks(page)
		for _, v := range linkSlice {
			postLinks <- v
		}

		// 拿首頁上的<上頁>
		re := regexp.MustCompile(`\<a.+href=\"\/bbs\/ALLPOST\/index(\d+)\.html\".+上頁\<\/a\>`)
		latPage, err := strconv.Atoi(re.FindStringSubmatch(page)[1])
		if err != nil {
			log.Println(err)
		}

		// 從<上頁>往前拿n頁
		for i := latPage - 1; i >= latPage-pttcrawler.settings.FetchExtPageNum; i-- {
			func(index int) {
				page := pttcrawler.GetHTML(pttcrawler.host + fmt.Sprintf("/bbs/allpost/index%d.html", index))
				linkSlice := parseAPLinks(page)
				for _, v := range linkSlice {
					postLinks <- v
				}
			}(i)
		}
		log.Println("getLinks waiting...")
		select {
			case <-time.After(time.Second * 60):
				break
		}
	}
}

func parseAPLinks(page string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		log.Println(err)
	}

	links := []string{}
	// parse出列表文章的標題與連結所在
	doc.Find(".r-ent .title a").Each(func(i int, s *goquery.Selection) {
		lk, _ := s.Attr("href")
		tt := s.Text()
		re := regexp.MustCompile(`^.+\((.+)\)$`)
		board := re.FindStringSubmatch(string(tt))[1]
		// ALLPOST列表的文章連結都是錯的，要把URL的ALLPOST改成標題後面括號中的看板名稱ex:[黑特] 柯粉很妙 (HatePolitics)
		link := strings.Replace(lk, "ALLPOST", board, -1)
		links = append(links, link)
	})

	return links
}

func (pttcrawler *pttcrawler) fetchPosts(postLink <-chan string, postQ chan<- post) {
	for {
		select {
		case link := <-postLink:
			fetchedLnk, err := redisClient.Exists(link).Result()
			if err != nil {
				log.Println(err)
				continue
			}

			if fetchedLnk == 0 {
				log.Println("fetchPosts:" + pttcrawler.host + link)
				page := pttcrawler.GetHTML(pttcrawler.host+link)
				if page == "" {
					continue
				}
				post := parsePost(page)

				err = redisClient.Set(link, "", time.Hour*time.Duration(pttcrawler.settings.FchLnkExHr)).Err()
				if err != nil {
					log.Println(err)
				}

				postQ <- post
			}
		}
	}
}

func parsePost(page string) post {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		log.Println(err)
	}

	// parse 作者/標題/看板
	post := post{}
	doc.Find(".article-meta-tag").Each(func(i int, s *goquery.Selection) {
		metaType := s.Text()
		metaVal := s.SiblingsFiltered(".article-meta-value").Text()
		switch metaType {
		case "作者":
			post.Author = metaVal
		case "看板":
			post.Board = metaVal
		case "標題":
			post.Title = metaVal
		}
	})
	post.Content = doc.Find("#main-content").Text()
	post.Link, _ = doc.Find(".f2 a").Attr("href")

	return post
}

func getKeywords() []string {
	keywords, err := redisClient.SMembers("keywords").Result()
	if err != nil {
		log.Println(err)
	}

	return keywords
}

func (pttcrawler *pttcrawler) pickup(postQ <-chan post) {
	keywords := getKeywords()
	lst := time.Now()

	for post := range postQ {
		// 每10秒抓新的關鍵字
		if time.Now().Sub(lst) > time.Second*10 {
			keywords = getKeywords()
			log.Printf("renew keywords: %v\n", keywords)
		}

		for _, kw := range keywords {
			re := regexp.MustCompile(kw)
			if re.MatchString(post.Title) || re.MatchString(post.Content) {
				log.Printf("pickuped : kw<%s> title:<%s> lnk<%s>\n", kw, post.Title, post.Link)
				pJSON, err := json.Marshal(post)
				if err != nil {
					fmt.Println(err)
					return
				}
				redisClient.RPush("notifies", pJSON).Result()
			}
		}
		lst = time.Now()
	}
}

func (pttcrawler *pttcrawler) GetHTML(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	req.AddCookie(&http.Cookie{Name: "over18", Value: "1"})

	response, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	return string(body)
}
package ptt

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

type ClientMk http.Client

func Test正常爬到內容符合關鍵字(t *testing.T) {
	keyword := []string{"美國", "中國"}
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/allpost" {
			res.Write([]byte(apAll))
		} else {
			res.Write([]byte(postPage))
		}
	}))
	defer func() { testServer.Close() }()
	mr,_ := miniredis.Run()
	rHost := mr.Host()
	rPort, _ := strconv.Atoi(mr.Port())
	pttweb := NewPttCrawler(Settings{
		Host:            testServer.URL,
		AllpostURL:      "/allpost",
		FchLnkExHr:      0,
		FetchExtPageNum: 0,
		LinksChCap:      0,
		PostsChCap:      0,
		RedisHost:       rHost,
		RedisPort:       rPort,
	})

	mr.SAdd("keywords", keyword...)
	doneCh := make(chan bool)
	go pttweb.Run(doneCh)
	time.Sleep(time.Second*1)
	doneCh <- true

	keys := mr.Keys()
	foundPost := false
	for i := range keys {
		if strings.Contains(keys[i], "bbs") {
			foundPost = true
		}
	}
	assert.True(t, foundPost)
	ls, _ :=mr.List("notifies")
	assert.Greater(t, len(ls), 0)
}

func Test正常爬到內容但無符合關鍵字(t *testing.T) {
	keyword := []string{"爐渣米", "瘦肉精"}
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/allpost" {
			res.Write([]byte(apAll))
		} else {
			res.Write([]byte(postPage))
		}
	}))
	defer func() { testServer.Close() }()
	mr,_ := miniredis.Run()
	rHost := mr.Host()
	rPort, _ := strconv.Atoi(mr.Port())
	pttweb := NewPttCrawler(Settings{
		Host:            testServer.URL,
		AllpostURL:      "/allpost",
		FchLnkExHr:      0,
		FetchExtPageNum: 0,
		LinksChCap:      0,
		PostsChCap:      0,
		RedisHost:       rHost,
		RedisPort:       rPort,
	})

	mr.SAdd("keywords", keyword...)
	doneCh := make(chan bool)
	go pttweb.Run(doneCh)
	time.Sleep(time.Second*1)
	doneCh <- true

	keys := mr.Keys()
	foundPost := false
	for i := range keys {
		if strings.Contains(keys[i], "bbs") {
			foundPost = true
		}
	}
	assert.True(t, foundPost)
	ls, _ :=mr.List("notifies")
	assert.Equal(t, 0, len(ls))
}

func Test正常爬不到內容(t *testing.T) {
	keyword := []string{"美國", "中國"}
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/allpost" {
			res.Write([]byte(apAll))
		} else {
			res.WriteHeader(http.StatusNotFound)
		}
	}))
	defer func() { testServer.Close() }()
	mr,_ := miniredis.Run()
	rHost := mr.Host()
	rPort, _ := strconv.Atoi(mr.Port())
	pttweb := NewPttCrawler(Settings{
		Host:            testServer.URL,
		AllpostURL:      "/allpost",
		FchLnkExHr:      0,
		FetchExtPageNum: 0,
		LinksChCap:      0,
		PostsChCap:      0,
		RedisHost:       rHost,
		RedisPort:       rPort,
	})

	mr.SAdd("keywords", keyword...)
	doneCh := make(chan bool)
	go pttweb.Run(doneCh)
	time.Sleep(time.Second*1)
	doneCh <- true

	keys := mr.Keys()
	foundPost := false
	for i := range keys {
		if strings.Contains(keys[i], "bbs") {
			foundPost = true
		}
	}
	t.Log(keys)
	assert.False(t, foundPost)
	ls, _ :=mr.List("notifies")
	assert.Equal(t, 0, len(ls))
}






var apAll = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">


<meta name="viewport" content="width=device-width, initial-scale=1">

<title>看板 ALLPOST 文章列表 - 批踢踢實業坊</title>

<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-common.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-base.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-custom.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/pushstream.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-print.css" media="print">




	</head>
    <body>

<div id="topbar-container">
	<div id="topbar" class="bbs-content">
		<a id="logo" href="/bbs/">批踢踢實業坊</a>
		<span>&rsaquo;</span>
		<a class="board" href="/bbs/ALLPOST/index.html"><span class="board-label">看板 </span>ALLPOST</a>
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>

<div id="main-container">
	<div id="action-bar-container">
		<div class="action-bar">
			<div class="btn-group btn-group-dir">
				<a class="btn selected" href="/bbs/ALLPOST/index.html">看板</a>
				<a class="btn" href="/man/ALLPOST/index.html">精華區</a>
			</div>
			<div class="btn-group btn-group-paging">
				<a class="btn wide" href="/bbs/ALLPOST/index1.html">最舊</a>
				<a class="btn wide" href="/bbs/ALLPOST/index3636.html">&lsaquo; 上頁</a>
				<a class="btn wide disabled">下頁 &rsaquo;</a>
				<a class="btn wide" href="/bbs/ALLPOST/index.html">最新</a>
			</div>
		</div>
	</div>

	<div class="r-list-container action-bar-margin bbs-screen">
		<div class="search-bar">
			<form type="get" action="search" id="search-bar">
				<input class="query" type="text" name="q" value="" placeholder="搜尋文章&#x22ef;">
			</form>
		</div>






		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604423.A.A73.html">[文章] 777216濫用報掛skyexers (GossipPicket)</a>

			</div>
			<div class="meta">
				<div class="author">Anonymous.</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E6%96%87%E7%AB%A0%5D&#43;777216%E6%BF%AB%E7%94%A8%E5%A0%B1%E6%8E%9Bskyexers&#43;%28GossipPicket%29">搜尋同標題文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604448.A.166.html">[３Ｃ] 宜蘭 MONTECH君主BETA550W銅…(forsale)</a>

			</div>
			<div class="meta">
				<div class="author">windforce</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%EF%BC%93%EF%BC%A3%5D&#43;%E5%AE%9C%E8%98%AD&#43;MONTECH%E5%90%9B%E4%B8%BBBETA550W%E9%8A%85%E2%80%A6%28forsale%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Awindforce">搜尋看板內 windforce 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604480.A.A30.html">[贈送] 北屯/漂亮禮餅紙盒       (TaichungBun)</a>

			</div>
			<div class="meta">
				<div class="author">mhianni</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E8%B4%88%E9%80%81%5D&#43;%E5%8C%97%E5%B1%AF%2F%E6%BC%82%E4%BA%AE%E7%A6%AE%E9%A4%85%E7%B4%99%E7%9B%92&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28TaichungBun%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Amhianni">搜尋看板內 mhianni 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604483.A.22A.html">[３Ｃ] 全國 DuoFree真無線藍牙耳機  (forsale)</a>

			</div>
			<div class="meta">
				<div class="author">windforce</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%EF%BC%93%EF%BC%A3%5D&#43;%E5%85%A8%E5%9C%8B&#43;DuoFree%E7%9C%9F%E7%84%A1%E7%B7%9A%E8%97%8D%E7%89%99%E8%80%B3%E6%A9%9F&#43;&#43;%28forsale%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Awindforce">搜尋看板內 windforce 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604489.A.5FA.html">[震怒] 壓力測試                       (Hate)</a>

			</div>
			<div class="meta">
				<div class="author">panther6634</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E9%9C%87%E6%80%92%5D&#43;%E5%A3%93%E5%8A%9B%E6%B8%AC%E8%A9%A6&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28Hate%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Apanther6634">搜尋看板內 panther6634 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604507.A.386.html">Re: [新聞] 台鐵限速備援系統遲未裝／經…(Railway)</a>

			</div>
			<div class="meta">
				<div class="author">yymeow</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E6%96%B0%E8%81%9E%5D&#43;%E5%8F%B0%E9%90%B5%E9%99%90%E9%80%9F%E5%82%99%E6%8F%B4%E7%B3%BB%E7%B5%B1%E9%81%B2%E6%9C%AA%E8%A3%9D%EF%BC%8F%E7%B6%93%E2%80%A6%28Railway%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Ayymeow">搜尋看板內 yymeow 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604519.A.0CE.html">Fw: [情報] Conley離開泡泡 預計缺席… (UTAH-JAZZ)</a>

			</div>
			<div class="meta">
				<div class="author">FAYeeeeeeee</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E6%83%85%E5%A0%B1%5D&#43;Conley%E9%9B%A2%E9%96%8B%E6%B3%A1%E6%B3%A1&#43;%E9%A0%90%E8%A8%88%E7%BC%BA%E5%B8%AD%E2%80%A6&#43;%28UTAH-JAZZ%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3AFAYeeeeeeee">搜尋看板內 FAYeeeeeeee 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604538.A.DDA.html">[徵女] 想認識厚片女            (AllTogether)</a>

			</div>
			<div class="meta">
				<div class="author">betapla</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E5%BE%B5%E5%A5%B3%5D&#43;%E6%83%B3%E8%AA%8D%E8%AD%98%E5%8E%9A%E7%89%87%E5%A5%B3&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28AllTogether%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Abetapla">搜尋看板內 betapla 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604540.A.8F9.html">Re: [問卦] 30歲以後單身下班生活大概…(Gossiping)</a>

			</div>
			<div class="meta">
				<div class="author">appledavid</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E5%95%8F%E5%8D%A6%5D&#43;30%E6%AD%B2%E4%BB%A5%E5%BE%8C%E5%96%AE%E8%BA%AB%E4%B8%8B%E7%8F%AD%E7%94%9F%E6%B4%BB%E5%A4%A7%E6%A6%82%E2%80%A6%28Gossiping%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Aappledavid">搜尋看板內 appledavid 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604551.A.C8B.html">[心得] 語言所推甄錄取心得         (graduate)</a>

			</div>
			<div class="meta">
				<div class="author">imagine155</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E5%BF%83%E5%BE%97%5D&#43;%E8%AA%9E%E8%A8%80%E6%89%80%E6%8E%A8%E7%94%84%E9%8C%84%E5%8F%96%E5%BF%83%E5%BE%97&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28graduate%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Aimagine155">搜尋看板內 imagine155 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604621.A.346.html">[交易] 威秀電影票/全聯道爾頓點數 (BigShiLin)</a>

			</div>
			<div class="meta">
				<div class="author">susu78211</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E4%BA%A4%E6%98%93%5D&#43;%E5%A8%81%E7%A7%80%E9%9B%BB%E5%BD%B1%E7%A5%A8%2F%E5%85%A8%E8%81%AF%E9%81%93%E7%88%BE%E9%A0%93%E9%BB%9E%E6%95%B8&#43;%28BigShiLin%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Asusu78211">搜尋看板內 susu78211 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604622.A.69E.html">[贈送] 全國二手多本漫畫               (give)</a>

			</div>
			<div class="meta">
				<div class="author">Suki061</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E8%B4%88%E9%80%81%5D&#43;%E5%85%A8%E5%9C%8B%E4%BA%8C%E6%89%8B%E5%A4%9A%E6%9C%AC%E6%BC%AB%E7%95%AB&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28give%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3ASuki061">搜尋看板內 Suki061 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604630.A.960.html">心疼                                 (Diary)</a>

			</div>
			<div class="meta">
				<div class="author">Canyoufindme</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%E5%BF%83%E7%96%BC&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28Diary%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3ACanyoufindme">搜尋看板內 Canyoufindme 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604643.A.224.html">[寶寶] 4個月睡眠倒退？          (BabyMother)</a>

			</div>
			<div class="meta">
				<div class="author">oliviology</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E5%AF%B6%E5%AF%B6%5D&#43;4%E5%80%8B%E6%9C%88%E7%9D%A1%E7%9C%A0%E5%80%92%E9%80%80%EF%BC%9F&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28BabyMother%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Aoliviology">搜尋看板內 oliviology 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604663.A.074.html">[閒聊] 教授怎麼可能沒有女朋友呢(Marginalman)</a>

			</div>
			<div class="meta">
				<div class="author">fnm525</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E9%96%92%E8%81%8A%5D&#43;%E6%95%99%E6%8E%88%E6%80%8E%E9%BA%BC%E5%8F%AF%E8%83%BD%E6%B2%92%E6%9C%89%E5%A5%B3%E6%9C%8B%E5%8F%8B%E5%91%A2%28Marginalman%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Afnm525">搜尋看板內 fnm525 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604678.A.27D.html">[自介] Nekomiao                     (Vtuber)</a>

			</div>
			<div class="meta">
				<div class="author">Nekomiao</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E8%87%AA%E4%BB%8B%5D&#43;Nekomiao&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28Vtuber%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3ANekomiao">搜尋看板內 Nekomiao 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604679.A.E5E.html">[討論] 只要努力，幸福就會自然到…(WomenTalk)</a>

			</div>
			<div class="meta">
				<div class="author">notfound0506</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E8%A8%8E%E8%AB%96%5D&#43;%E5%8F%AA%E8%A6%81%E5%8A%AA%E5%8A%9B%EF%BC%8C%E5%B9%B8%E7%A6%8F%E5%B0%B1%E6%9C%83%E8%87%AA%E7%84%B6%E5%88%B0%E2%80%A6%28WomenTalk%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Anotfound0506">搜尋看板內 notfound0506 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604682.A.732.html">[閒聊] 你各位                  (Marginalman)</a>

			</div>
			<div class="meta">
				<div class="author">cpblgu</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E9%96%92%E8%81%8A%5D&#43;%E4%BD%A0%E5%90%84%E4%BD%8D&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;&#43;%28Marginalman%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Acpblgu">搜尋看板內 cpblgu 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604697.A.13C.html">[問題] 驗證遊戲檔案完整性 頻繁出現… (Steam)</a>

			</div>
			<div class="meta">
				<div class="author">RockCaveChen</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E5%95%8F%E9%A1%8C%5D&#43;%E9%A9%97%E8%AD%89%E9%81%8A%E6%88%B2%E6%AA%94%E6%A1%88%E5%AE%8C%E6%95%B4%E6%80%A7&#43;%E9%A0%BB%E7%B9%81%E5%87%BA%E7%8F%BE%E2%80%A6&#43;%28Steam%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3ARockCaveChen">搜尋看板內 RockCaveChen 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>





		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">

				<a href="/bbs/ALLPOST/M.1597604712.A.B6F.html">[揪團] 8/22(六)快樂練歌團-錢櫃台北南… (KTV)</a>

			</div>
			<div class="meta">
				<div class="author">frjsts</div>
				<div class="article-menu">

					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/ALLPOST/search?q=thread%3A%5B%E6%8F%AA%E5%9C%98%5D&#43;8%2F22%28%E5%85%AD%29%E5%BF%AB%E6%A8%82%E7%B7%B4%E6%AD%8C%E5%9C%98-%E9%8C%A2%E6%AB%83%E5%8F%B0%E5%8C%97%E5%8D%97%E2%80%A6&#43;%28KTV%29">搜尋同標題文章</a></div>

						<div class="item"><a href="/bbs/ALLPOST/search?q=author%3Afrjsts">搜尋看板內 frjsts 的文章</a></div>

					</div>

				</div>
				<div class="date"> 8/17</div>
				<div class="mark"></div>
			</div>
		</div>



	</div>


</div>



<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-32365737-1', {
    cookieDomain: 'ptt.cc',
    legacyCookieDomain: 'ptt.cc'
  });
  ga('send', 'pageview');
</script>



<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
<script src="//images.ptt.cc/bbs/v2.27/bbs.js"></script>

    </body>
</html>
`

var apNone = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">


<meta name="viewport" content="width=device-width, initial-scale=1">

<title>看板 ALLPOST 文章列表 - 批踢踢實業坊</title>

<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-common.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-base.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-custom.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/pushstream.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-print.css" media="print">




	</head>
    <body>

<div id="topbar-container">
	<div id="topbar" class="bbs-content">
		<a id="logo" href="/bbs/">批踢踢實業坊</a>
		<span>&rsaquo;</span>
		<a class="board" href="/bbs/ALLPOST/index.html"><span class="board-label">看板 </span>ALLPOST</a>
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>

<div id="main-container">
	<div id="action-bar-container">
		<div class="action-bar">
			<div class="btn-group btn-group-dir">
				<a class="btn selected" href="/bbs/ALLPOST/index.html">看板</a>
				<a class="btn" href="/man/ALLPOST/index.html">精華區</a>
			</div>
			<div class="btn-group btn-group-paging">
				<a class="btn wide" href="/bbs/ALLPOST/index1.html">最舊</a>
				<a class="btn wide" href="/bbs/ALLPOST/index3636.html">&lsaquo; 上頁</a>
				<a class="btn wide disabled">下頁 &rsaquo;</a>
				<a class="btn wide" href="/bbs/ALLPOST/index.html">最新</a>
			</div>
		</div>
	</div>

	<div class="r-list-container action-bar-margin bbs-screen">
		<div class="search-bar">
			<form type="get" action="search" id="search-bar">
				<input class="query" type="text" name="q" value="" placeholder="搜尋文章&#x22ef;">
			</form>
		</div>
	</div>


</div>



<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-32365737-1', {
    cookieDomain: 'ptt.cc',
    legacyCookieDomain: 'ptt.cc'
  });
  ga('send', 'pageview');
</script>



<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
<script src="//images.ptt.cc/bbs/v2.27/bbs.js"></script>

    </body>
</html>
`

var postPage = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">


<meta name="viewport" content="width=device-width, initial-scale=1">

<title>[新聞] 以卵擊石 朴廷桓九段拒絕參加KBS杯圍棋王戰？ - 看板 GO - 批踢踢實業坊</title>
<meta name="robots" content="all">
<meta name="keywords" content="Ptt BBS 批踢踢">
<meta name="description" content="以卵擊石 朴廷桓九段拒絕參加KBS杯圍棋王戰？
　　來源：吉林大福圓
　　韓國第39屆KBS杯圍棋王戰，將於8月19日開賽。
　　1979年創立的這項賽事，對於圍棋基礎薄弱，賽事稀少的韓國圍棋界來說，堪稱悠久
的傳統和一個大事件。
">
<meta property="og:site_name" content="Ptt 批踢踢實業坊">
<meta property="og:title" content="[新聞] 以卵擊石 朴廷桓九段拒絕參加KBS杯圍棋王戰？">
<meta property="og:description" content="以卵擊石 朴廷桓九段拒絕參加KBS杯圍棋王戰？
　　來源：吉林大福圓
　　韓國第39屆KBS杯圍棋王戰，將於8月19日開賽。
　　1979年創立的這項賽事，對於圍棋基礎薄弱，賽事稀少的韓國圍棋界來說，堪稱悠久
的傳統和一個大事件。
">
<link rel="canonical" href="https://www.ptt.cc/bbs/GO/M.1597603762.A.135.html">

<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-common.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-base.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-custom.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/pushstream.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-print.css" media="print">




	</head>
    <body>

<div id="topbar-container">
	<div id="topbar" class="bbs-content">
		<a id="logo" href="/bbs/">批踢踢實業坊</a>
		<span>&rsaquo;</span>
		<a class="board" href="/bbs/GO/index.html"><span class="board-label">看板 </span>GO</a>
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>
<div id="navigation-container">
	<div id="navigation" class="bbs-content">
		<a class="board" href="/bbs/GO/index.html">返回看板</a>
		<div class="bar"></div>
	</div>
</div>
<div id="main-container">
    <div id="main-content" class="bbs-screen bbs-content"><div class="article-metaline"><span class="article-meta-tag">作者</span><span class="article-meta-value">joanzkow (星浪)</span></div><div class="article-metaline-right"><span class="article-meta-tag">看板</span><span class="article-meta-value">GO</span></div><div class="article-metaline"><span class="article-meta-tag">標題</span><span class="article-meta-value">%s</span></div><div class="article-metaline"><span class="article-meta-tag">時間</span><span class="article-meta-value">Mon Aug 17 02:49:15 2020</span></div>
以卵擊石 朴廷桓九段拒絕參加KBS杯圍棋王戰？

　　來源：吉林大福圓

　　韓國第39屆KBS杯圍棋王戰，將於8月19日開賽。

　　1979年創立的這項賽事，對於圍棋基礎薄弱，賽事稀少的韓國圍棋界來說，堪稱悠久
的傳統和一個大事件。

　　更何況，本項比賽的冠軍獎金高達2000萬韓元。2000萬韓元折合人民幣大約為11萬元
，韓國棋手可以蜂擁而上了。

　　據說，比賽將通過預賽決出5位本賽棋手。而本賽，將由8位棋手進行。也就是說，比
賽組委會設了三名外卡棋手。

　　本次比賽的三名外卡棋手包括，申真諝九段、申旻埈九段和女子圍棋世界第一人崔精
九段。

　　現在，朴廷桓九段竟然拒絕參加如此重要的一場比賽？一時間，韓國圍棋界對此議論
紛紛。

　　關於朴廷桓九段拒絕參加比賽的原因，大約有如下幾條。

　　一是，自尊心受不了。的確，作為韓國數一數二的高手，如果被擠出外卡參賽名單，
非要從預選賽打起，這種場面也確實挺難看。自尊心，確實受不了。

　　二是，韓國棋院的報復。

　　2019年1月，在第37屆韓國KBS棋王戰決賽中，申旻埈九段以2：0戰勝了朴廷桓九段。
即便如此，朴廷桓九段也以亞軍的身份拿到了亞洲杯電視圍棋快棋賽的參賽權。

　　問題是同年6月，朴廷桓九段為了集中精力備戰與朴永訓九段的春蘭杯決賽，而自動
放棄了亞洲杯電視快棋賽的參賽權。

　　作為韓國棋院，一定會對此感到很不愉快。既然不愉快，必然會有報復。而朴廷桓九
段會不會以卵擊石，對此予以反擊？

　　當然了，最大的可能性還是因為農心杯三國圍棋擂臺賽。

　　8月18日，第21屆農心辛拉麵杯世界圍棋團體錦標賽將進行第10場比賽，由朴廷桓九
段向日本隊主將井山裕太九段發起挑戰。

　　但是，這一解釋並不充分。

　　一是，由於疫情原因，本次比賽將在網上進行。二是，就算朴廷桓九段過了井山裕太
九段一關，也很難再次突破中國棋手的關口。

　　中國的任何一位棋手，都不是朴廷桓九段能突破得了的。如果他還想像他的前輩徐奉
洙九段那樣連戰連勝，奪得比賽的冠軍，真可謂“想多了”。

　　到底什麼原因，讓朴廷桓放棄了KBS圍棋王戰？

<a href="https://sports.sina.com.cn/go/2020-08-16/doc-iivhvpwy1334335.shtml" target="_blank" rel="nofollow">https://sports.sina.com.cn/go/2020-08-16/doc-iivhvpwy1334335.shtml</a>

--
<span class="f2">※ 發信站: 批踢踢實業坊(ptt.cc), 來自: 123.192.113.200 (臺灣)
</span><span class="f2">※ 文章網址: <a href="https://www.ptt.cc/bbs/GO/M.1597603762.A.135.html" target="_blank" rel="nofollow">https://www.ptt.cc/bbs/GO/M.1597603762.A.135.html</a>
</span><span class="f2">※ 編輯: joanzkow (123.192.113.200 臺灣), 08/17/2020 02:49:57
</span><div class="push"><span class="hl push-tag">推 </span><span class="f3 hl push-userid">pan0531</span><span class="f3 push-content">: 最後一段是作者沒睡醒寫的嗎???</span><span class="push-ipdatetime"> 08/17 02:53
</span></div></div>

    <div id="article-polling" data-pollurl="/poll/GO/M.1597603762.A.135.html?cacheKey=2084-335629357&amp;offset=2188&amp;offset-sig=27e3011452ee3542232acf3c3435a6dba99be66e" data-longpollurl="/v1/longpoll?id=a8f677aef9a06c5ed7353f2b09e1e41ef74452ce" data-offset="2188"></div>



</div>



<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-32365737-1', {
    cookieDomain: 'ptt.cc',
    legacyCookieDomain: 'ptt.cc'
  });
  ga('send', 'pageview');
</script>



<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
<script src="//images.ptt.cc/bbs/v2.27/bbs.js"></script>

    </body>
</html>
`
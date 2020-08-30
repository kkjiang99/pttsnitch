package main

import (
	"log"

	"github.com/kkjiang99/pttsnitch/config"
	"github.com/kkjiang99/pttsnitch/ptt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	config := config.GetConfig()
	log.Printf("%+v\n", config)
	done := make(chan bool, 1)
	blocking := make(chan bool)

	pttweb := ptt.NewPttCrawler(ptt.Settings{
		Host:            config.PttCrawler.Host,
		AllpostURL:      config.PttCrawler.AllPostURL,
		FchLnkExHr:      config.PttCrawler.FchLnkExHr,
		FetchExtPageNum: config.PttCrawler.FetchExtPageNum,
		LinksChCap:      config.PttCrawler.LinksChCap,
		PostsChCap:      config.PttCrawler.PostChCap,
		RedisHost:       config.Redis.Host,
		RedisPort:       config.Redis.Port,
	})

	go pttweb.Run(done)
	<-blocking
}

package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Env        string     `toml:"env"`
	Redis      Redis      `toml:"redis"`
	PttCrawler PttCrawler `toml:"ptt_crawler"`
}

type Redis struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type PttCrawler struct {
	Host            string `toml:"host"`
	AllPostURL      string `toml:"allpostURL"`
	FchLnkExHr      int    `toml:"fchLnkExHr"`
	FetchExtPageNum int    `toml:"fetchExtPageNum"`
	LinksChCap      int    `toml:"linksChCap"`
	LinksSp         int    `toml:"linksSp"`
	PostChCap       int    `toml:"postChCap"`
}

func GetConfig() Config {
	config := Config{}
	if _, err := toml.DecodeFile("config/config.toml", &config); err != nil {
		log.Panicln(err)
	}
	return config
}

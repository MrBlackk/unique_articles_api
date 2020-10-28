package main

import (
	"github.com/tkanos/gonfig"
	"log"
)

type Config struct {
	ArticleSimilarity int
}

func GetConfig() Config {
	config := Config{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

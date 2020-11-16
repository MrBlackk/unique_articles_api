package main

import (
	"github.com/tkanos/gonfig"
	"log"
)

type Config struct {
	ArticleSimilarity int
}

type FilterConfig struct {
	Commons     []string
	Transitions []string
	Synonyms    [][]string
}

var Conf Config
var FilterConf FilterConfig

func init() {
	Conf = getConfig()
	FilterConf = getFilterConfig()
}

func getConfig() Config {
	config := Config{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func getFilterConfig() FilterConfig {
	config := FilterConfig{}
	err := gonfig.GetConf("filter_config.json", &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

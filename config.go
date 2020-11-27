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
	Conf = readConfig()
	FilterConf = readFilterConfig()
}

func readConfig() Config {
	config := Config{}
	err := gonfig.GetConf("configs/config.json", &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func readFilterConfig() FilterConfig {
	config := FilterConfig{}
	err := gonfig.GetConf("configs/filter_config.json", &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

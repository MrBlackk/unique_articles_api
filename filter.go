package main

import (
	stemmer "github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

func prepare(str string) []string {
	arr := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	return filter(arr)
}

func filter(words []string) []string {
	r := words[:0]
	for i := range words {
		word := strings.ToLower(words[i])
		word = stemmer.Stem(word, false)
		if isInStopLists(word) {
			continue
		}
		word = getBaseSynonymOrSelf(word)
		r = append(r, word)
	}
	return r
}

func getBaseSynonymOrSelf(word string) string {
	for _, val := range FilterConf.Synonyms {
		if contains(word, val) {
			return val[0]
		}
	}
	return word
}

func isInStopLists(word string) bool {
	return contains(word, FilterConf.Commons) || contains(word, FilterConf.Transitions)
}

func contains(needle string, hay []string) bool {
	for _, val := range hay {
		if val == needle {
			return true
		}
	}
	return false
}

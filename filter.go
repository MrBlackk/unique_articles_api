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
	arr = lowercase(arr)
	arr = filter(arr)
	arr = stem(arr)
	arr = replaceSynonyms(arr)
	return arr
}

func lowercase(words []string) []string {
	r := make([]string, len(words))
	for i, word := range words {
		r[i] = strings.ToLower(word)
	}
	return r
}

func stem(words []string) []string {
	r := make([]string, len(words))
	for i, word := range words {
		r[i] = stemmer.Stem(word, false)
	}
	return r
}

func replaceSynonyms(words []string) []string {
	r := make([]string, len(words))
	for i, word := range words {
		r[i] = getBaseSynonymOrSelf(word)
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

func filter(words []string) []string {
	r := make([]string, 0, len(words))
	for _, word := range words {
		if !isInStopLists(word) {
			r = append(r, word)
		}
	}
	return r
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

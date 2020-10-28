package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type post struct {
	Content string `json:"content"`
}

type article struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

type articles struct {
	Articles []article `json:"articles"`
}

var inMemDb = articles{[]article{
	{1, "content1"},
	{2, "content2"},
}}

func main() {
	w1 := strings.Fields("hello world")
	w2 := strings.Fields("world hello hello")

	d := Distance(w1, w2)
	fmt.Println("Distance ", d)
	fmt.Println("Proc: ", 100.0-float64(d)/float64(Max(len(w1), len(w2)))*100.0)

	config := GetConfig()
	fmt.Println(config.ArticleSimilarity)

	mux := httprouter.New()
	mux.POST("/articles", saveArticle)
	mux.GET("/articles", getArticles)
	mux.GET("/articles/:id", getArticle)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getArticle(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err1 := strconv.Atoi(idStr)
	if err1 != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	bs, err := json.Marshal(inMemDb.Articles[id-1])
	if err != nil {
		fmt.Println("error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func getArticles(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	bs, err := json.Marshal(inMemDb)
	if err != nil {
		fmt.Println("error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func saveArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check content type
	var data post
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l := len(inMemDb.Articles)
	a := article{l + 1, data.Content}
	inMemDb.Articles = append(inMemDb.Articles, a)
	bs, err := json.Marshal(a)
	if err != nil {
		fmt.Println("error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

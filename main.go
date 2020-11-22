package main

import (
	"encoding/json"
	"github.com/golang/gddo/httputil/header"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

type post struct {
	Content string `json:"content"`
}

type article struct {
	Id            int    `json:"id"`
	Content       string `json:"content"`
	DuplicatesIds []int  `json:"duplicate_article_ids"`
	prepared      []string
}

type articles struct {
	Articles []article `json:"articles"`
}

func main() {
	router := httprouter.New()
	router.POST("/articles", saveArticle)
	router.GET("/articles", getArticles)
	router.GET("/articles/:id", getArticleById)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func getArticleById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	txn := db.Txn(false)
	defer txn.Abort()
	raw, err := txn.First(ArticleTable, ArticleId, id)
	if err != nil || raw == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	writeJson(w, raw.(*article))
}

func getArticles(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	txn := db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get(ArticleTable, ArticleId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	arts := articles{[]article{}}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		a := obj.(*article)
		arts.Articles = append(arts.Articles, *a)
	}

	writeJson(w, arts)
}

func saveArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := decodeJson(w, r)
	if err != nil {
		return
	}

	a := article{NextId(), data.Content, make([]int, 0), prepare(data.Content)}
	ids, err := findDuplicatesIds(&a)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	a.DuplicatesIds = ids

	txn := db.Txn(true)
	err = txn.Insert(ArticleTable, &a)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	txn.Commit()

	writeJson(w, a)
}

func decodeJson(w http.ResponseWriter, r *http.Request) (post, error) {
	var data post
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
			return data, errors.New("")
		}
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return data, err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		http.Error(w, "Request body must only contain a single JSON object", http.StatusBadRequest)
		return data, err
	}
	return data, nil
}

func writeJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bs)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func findDuplicatesIds(a *article) ([]int, error) {
	txn := db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get(ArticleTable, ArticleId)
	if err != nil {
		return nil, err
	}

	r := make([]int, 0)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		v := obj.(*article)
		if isDuplicate(v.prepared, a.prepared) {
			r = append(r, v.Id)
			txn := db.Txn(true)
			v.DuplicatesIds = append(v.DuplicatesIds, a.Id)
			err = txn.Insert(ArticleTable, v)
			if err != nil {
				return nil, err
			}
			txn.Commit()
		}
	}

	return r, nil
}

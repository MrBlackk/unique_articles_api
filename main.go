package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:pass@/dc?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	mux := httprouter.New()
	mux.GET("/", defaultRemove)
	mux.POST("/articles", saveArticle)
	mux.GET("/articles", getArticles)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getArticles(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	articles, err := db.Query(`SELECT id, content FROM articles;`)
	if err != nil {
		log.Println(err)
	}
	defer articles.Close()

	var s, id, content string

	for articles.Next() {
		err = articles.Scan(&id, &content)
		if err != nil {
			log.Println(err)
		}
		s += id + " - " + content + "\n"
	}
	fmt.Fprint(w, s)
}

func saveArticle(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	stmt, err := db.Prepare(`INSERT INTO articles(content) VALUES ("Some_content");`)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	r, err := stmt.Exec()
	if err != nil {
		log.Println(err)
	}
	insertId, err := r.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	id := strconv.FormatInt(insertId, 10)

	fmt.Fprint(w, "Inserted id: "+id)
}

func defaultRemove(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Default page")
}

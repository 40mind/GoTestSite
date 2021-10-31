package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Article struct {
	Id     int
	Author string
	Name   string
	Text   string
}

var database *sql.DB

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := database.Query("select id, author, name, text from blog.articles where id = ?", id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	a := Article{}
	rows.Next()
	err = rows.Scan(&a.Id, &a.Author, &a.Name, &a.Text)
	if err != nil {
		fmt.Println(err)
	}
	tmpl, _ := template.ParseFiles("C:/Users/maxim/Documents/GitHub/GoTestSite/templates/article.html")
	tmpl.Execute(w, a)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("select * from blog.articles")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	articles := []Article{}

	for rows.Next() {
		a := Article{}
		err := rows.Scan(&a.Id, &a.Author, &a.Name, &a.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		articles = append(articles, a)
	}
	tmpl, _ := template.ParseFiles("C:/Users/maxim/Documents/GitHub/GoTestSite/templates/index.html")
	tmpl.Execute(w, articles)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		author := r.FormValue("author")
		name := r.FormValue("name")
		text := r.FormValue("text")

		_, err = database.Exec("insert into blog.articles (author, name, text) values (?, ?, ?)",
			author, name, text)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "C:/Users/maxim/Documents/GitHub/GoTestSite/templates/add.html")
	}
}

func main() {

	db, err := sql.Open("mysql", "root:40MinDmysqlstudy@/blog")
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/add", AddHandler)
	router.HandleFunc("/article/{id:[0-9]+}", ArticleHandler)

	http.Handle("/", router)

	fmt.Println("Сервер запущен...")
	http.ListenAndServe(":8000", nil)
}

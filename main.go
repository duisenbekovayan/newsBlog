package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Article struct {
	Id                     uint16
	Title, Anons, Fulltext string
}

var posts = []Article{}
var show_Post = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:@tcp(127.127.126.24:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Fulltext)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.127.126.24:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func showPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/show.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	db, err := sql.Open("mysql", "root:@tcp(127.127.126.24:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	show_Post = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Fulltext)
		if err != nil {
			panic(err)
		}
		show_Post = post
	}
	t.ExecuteTemplate(w, "show", show_Post)
}

func handleFunc() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")
	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}

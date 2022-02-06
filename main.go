package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

func main() {
	handleFunc()
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())

	}

	//Подключаемся к БД
	db, err := sql.Open("mysql", "gotest:!23Qweasdzxc@/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("подключился к БД")

	//Забираем данные из БД
	res, err := db.Query(fmt.Sprintf("SELECT *  FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post

		//fmt.Printf("Post: %s with id: %d\n", post.Title, post.Id)
	}

	t.ExecuteTemplate(w, "show", showPost)

}

func handleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create/", create).Methods("GET")
	router.HandleFunc("/save_article", saveArticle).Methods("POST")
	router.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":10000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())

	}

	//Подключаемся к БД
	db, err := sql.Open("mysql", "gotest:!23Qweasdzxc@/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Забираем данные из БД
	res, err := db.Query("SELECT *  FROM `articles`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)

		fmt.Printf("Post: %s with id: %d\n", post.Title, post.Id)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	//Проверяем валидность заполнения формы
	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "Заполните все поля!")
	} else {

		//Подключаемся к БД
		db, err := sql.Open("mysql", "gotest:!23Qweasdzxc@/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()
		fmt.Println("подключился к БД")

		//Отправка данных в БД
		insertDB, err := db.Query(fmt.Sprintf("INSERT INTO `articles`(`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, fullText))
		if err != nil {
			panic(err)
		}
		defer insertDB.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

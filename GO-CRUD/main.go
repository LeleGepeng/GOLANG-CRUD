package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
var tmpl = template.Must(template.ParseGlob("templates/*"))

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error!! file .env tidak bisa diakses")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPassword + "@/" + dbName
	var errDb error
	db, errDb = sql.Open("mysql", dsn)
	if errDb != nil {
		log.Fatal(errDb)
	}
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// AKSES LOCALHOST NYA DISINI
	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

type User struct {
	Id int
	Name string
	Email string
	CreatedAt string
}

func Index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	tmpl.ExecuteTemplate(w, "index.html", users)
}

func Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT * FROM users WHERE id=?", id)

	var user User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "show.html", user)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w,"create.html", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT * FROM users WHERE id=?", id)

	var user User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "update.html", user)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		_, err := db.Exec("INSERT INTO users(name,email) VALUES(?, ?)", name, email)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", 301)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		_, err := db.Exec("UPDATE users SET name=?, email=? WHERE id=?", name, email, id)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", 301)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", 301)
}


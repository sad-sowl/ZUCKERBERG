package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"

	"github.com/alexedwards/scs"
)

type User struct {
	id       int
	email    string
	password string
	username string
}

var session *scs.Session
var db *sql.DB

func isInDatabase(email string, username string) bool {
	query := fmt.Sprintf("SELECT * FROM users WHERE email='%s' OR username='%s'", email, username)

	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	if rows.Next() {
		return true
	} else {
		return false
	}

}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func logup(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("logup.html"))

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password1 := r.FormValue("password")
		password2 := r.FormValue("password2")

		//Validation of email and passwords
		if password1 != password2 || !isValidEmail(email) || isInDatabase(email, username) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			//Create new user and add it to database
			query := fmt.Sprintf("INSERT INTO users(email, password, username) VALUES('%s', '%s', '%s')", email, password1, username)

			rows, err := db.Query(query)
			if err != nil {
				panic(err)
			}

			defer rows.Close()

			http.Redirect(w, r, "/thanks", http.StatusTemporaryRedirect)
		}
	}

	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))

	if r.Method == http.MethodPost {

		var user User
		err := db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", r.FormValue("email"), r.FormValue("password")).Scan(&user.id, &user.email, &user.password, &user.username)
		if err != nil {
			panic(err.Error())
		}

		http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
	}

	tmpl.Execute(w, nil)

}

func thanks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("thanks.html"))

	tmpl.Execute(w, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("home.html"))

	//user templating here somehow

	tmpl.Execute(w, nil)
}

//stuff for session
// func put(w http.ResponseWriter, r *http.Request) {
// 	session.Put(r.Context(), "message", "Hello")
// }

// func get(w http.ResponseWriter, r *http.Request) {
// 	msg := session.GetString(r.Context(), "message")
// 	io.WriteString(w, msg)
// }
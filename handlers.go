package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/sessions"
)

type User struct {
	id       int
	email    string
	password string
	username string
}

var (
	db    *sql.DB
	key   = []byte("pL!,$C@jc)~!4>m%z&Mb;^I7OBW1X")
	store = sessions.NewCookieStore(key)
)

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
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
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

	session, _ := store.Get(r, "UserCookie")

	if r.Method == http.MethodPost {

		var user User
		err := db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", r.FormValue("email"), r.FormValue("password")).Scan(&user.id, &user.email, &user.password, &user.username)
		if err != nil {
			session.Values["authenticated"] = false
			session.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}

		session.Values["authenticated"] = true
		session.Values["username"] = user.username
		session.Values["email"] = user.email
		session.Values["id"] = user.id
		session.Values["password"] = user.password

		session.Save(r, w)

		http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
	}

	tmpl.Execute(w, nil)

}

func thanks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("thanks.html"))

	//add button for redirecting to the login page
	tmpl.Execute(w, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("home.html"))

	session, _ := store.Get(r, "UserCookie")

	//if user is authorised, read his username and put it in html
	if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
		errPage := template.Must(template.ParseFiles("error.html"))

		err := struct {
			Message string
		}{
			Message: "You need to authorise",
		}

		errPage.Execute(w, err)
		return
	}

	//user templating here somehow

	tmpl.Execute(w, nil)
}

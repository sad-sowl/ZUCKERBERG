package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id       int
	email    string
	password string
	username string
}

func isInDatabase(email string, username string) bool {

	db, err := sql.Open("mysql", "root:Qwerty0106@tcp(127.0.0.1:3306)/users")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM users WHERE email='%s' OR username='%s'", email, username)

	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	if !rows.Next() {
		return false
	} else {
		return true
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

			db, err := sql.Open("mysql", "root:Qwerty0106@tcp(127.0.0.1:3306)/users")
			if err != nil {
				panic(err)
			}

			defer db.Close()

			query := fmt.Sprintf("INSERT INTO users(email, password, username) VALUES('%s', '%s', '%s')", email, password1, username)

			rows, err := db.Query(query)
			if err != nil {
				panic(err)
			}

			defer rows.Close()
			//Create new user and add it to database
			fmt.Printf("Username: %s\nEmail: %s\nPassword: %s", username, email, password1)
			http.Redirect(w, r, "/thanks", http.StatusPermanentRedirect)
		}
	}

	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))

	tmpl.Execute(w, nil)

}

func thanks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("thanks.html"))

	tmpl.Execute(w, nil)
}

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/logup", logup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/thanks", thanks)

	fmt.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/sessions"
	"github.com/rs/xid"
)

type User struct {
	ID       string
	Email    string
	Password string
	Username string
}

type Post struct {
	ID    string
	Text  string
	Owner string
	Likes int
}

type PostsOfUser struct {
	Owner User
	Posts []Post
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

			guid := xid.New()
			//Create new user and add it to database
			query := fmt.Sprintf("INSERT INTO users(id, email, password, username) VALUES('%s', '%s', '%s', '%s')", guid.String(), email, password1, username)

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
		err := db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", r.FormValue("email"), r.FormValue("password")).Scan(&user.ID, &user.Email, &user.Password, &user.Username)
		if err != nil {
			session.Values["authenticated"] = false
			session.Save(r, w)

			http.Redirect(w, r, "", http.StatusUnauthorized)
		}

		session.Values["Authenticated"] = true
		session.Values["Username"] = user.Username
		session.Values["Email"] = user.Email
		session.Values["ID"] = user.ID
		session.Values["Password"] = user.Password

		session.Save(r, w)

		http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
	}

	tmpl.Execute(w, nil)

}

func thanks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("thanks.html"))

	if r.FormValue("login") == "Log in" {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	}
	tmpl.Execute(w, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("home.html"))

	session, _ := store.Get(r, "UserCookie")

	//if user is authorised, read his username and put it in html
	if auth, ok := session.Values["Authenticated"].(bool); !auth || !ok {
		errPage := template.Must(template.ParseFiles("error.html"))

		err := struct {
			Message string
		}{
			Message: "You need to authorise",
		}

		errPage.Execute(w, err)
		return
	}

	user := User{
		Username: session.Values["Username"].(string),
		Email:    session.Values["Email"].(string),
		Password: session.Values["Password"].(string),
		ID:       session.Values["ID"].(string),
	}

	if r.FormValue("PostButton") == "Send" {

		//add post to the database
		guid := xid.New()
		query := fmt.Sprintf("INSERT INTO posts(id, text, owner, likes) VALUES('%s', '%s', '%s', 0)", guid.String(), r.FormValue("text"), user.Username)

		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

	}

	//-----Search all post of logged in user's------------
	query := fmt.Sprintf("SELECT * FROM posts WHERE owner = '%s'", user.Username)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post

		err := rows.Scan(&post.ID, &post.Text, &post.Owner, &post.Likes)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	PostWithOwner := PostsOfUser{
		Owner: user,
		Posts: posts,
	}

	tmpl.Execute(w, PostWithOwner)
}

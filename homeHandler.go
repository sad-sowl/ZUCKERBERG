package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/rs/xid"
)

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

	// add for uploading images to post
	if r.Method == http.MethodPost {

	}
	//----When submit the post------
	if r.FormValue("PostButton") == "Send" {

		timeOfPost := time.Now().String()

		//add post to the database
		guid := xid.New()

		text := r.FormValue("text")
		text = strings.ReplaceAll(text, "'", "''")

		query := fmt.Sprintf("INSERT INTO posts(id, text, owner, likes, post_at) VALUES('%s', '%s', '%s', 0, '%s')", guid.String(), text, user.Username, timeOfPost[:19])

		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

	}
	//-----Enter here when follow someone---------
	if r.FormValue("FollowButton") == "Follow" {

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

		err := rows.Scan(&post.ID, &post.Text, &post.Owner, &post.Likes, &post.TimeOfPost)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	//reverse the posts
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].TimeOfPost > posts[j].TimeOfPost
	})

	PostWithOwner := PostsOfUser{
		Owner: user,
		Posts: posts,
	}

	tmpl.Execute(w, PostWithOwner)
}

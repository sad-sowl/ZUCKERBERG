package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func welcome(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("welcome.html"))

	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))

	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", welcome)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/login", login)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func logup(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("logup.html"))

	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))

	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/logup", logup)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/login", login)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

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
		//TODO: add checking existense of the user in the database

		if password1 != password2 || !isValidEmail(email) {
			http.RedirectHandler("", http.StatusPermanentRedirect)
		} else {

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
	http.HandleFunc("/logup", logup)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/login", login)
	http.HandleFunc("/thanks", thanks)

	fmt.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

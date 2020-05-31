package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var err error
	db, err = sql.Open("mysql", "root:*PASSWORD*@tcp(127.0.0.1:3306)/users")
	if err != nil {
		panic(err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/logup", logup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/thanks", thanks)
	http.HandleFunc("/home", home)
	http.HandleFunc("/", login)

	fmt.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

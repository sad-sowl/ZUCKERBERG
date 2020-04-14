package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/memstore"
)

func main() {

	var err error
	db, err = sql.Open("mysql", "root:Qwerty0106@tcp(127.0.0.1:3306)/users")
	if err != nil {
		panic(err)
	}

	//stuff for session
	session = scs.NewSession()
	session.Store = memstore.New()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/logup", logup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/thanks", thanks)
	http.HandleFunc("/home", home)
	//http.HandleFunc("/put", put)
	//http.HandleFunc("/get", get)

	fmt.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

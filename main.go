package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var filename string
var hardware string
var application string

func main() {
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/pull", pullHandler)
	http.HandleFunc("/overwrite", overWriteHandler)
	fmt.Println("Server is running on port 50080...")
	dbConnect()
	http.ListenAndServe(":50080", nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "POST" {
		myUrl, _ := url.Parse(r.URL.String())

		param, _ := url.ParseQuery(myUrl.RawQuery)

		name := param.Get("name")
		hard := param.Get("hard")
		app := param.Get("app")

		if name == "" || hard == "" || app == "" {
			w.Write([]byte("Missing parameters"))
			return
		}
		filename = name
		hardware = hard
		application = app
	}

}

func dbConnect() {
	db, err := sql.Open("mysql", "fileserver:fileserver@tcp(localhost:3306)/fileserver")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")
}

func pushHandler(w http.ResponseWriter, r *http.Request) {

}

func pullHandler(w http.ResponseWriter, r *http.Request) {

}

func overWriteHandler(w http.ResponseWriter, r *http.Request) {

}

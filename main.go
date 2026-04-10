package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/pull", pullHandler)
	http.HandleFunc("/overwrite", overWriteHandler)

	fmt.Println("Connecting to the database...")
	dbConnect()

	fmt.Println("Server is running on port 50080...")
	http.ListenAndServe(":50080", nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "POST" {
		recieveUrl, _ := url.Parse(r.URL.String())

		param, _ := url.ParseQuery(recieveUrl.RawQuery)

		name := param.Get("name")
		hard := param.Get("hard")
		app := param.Get("app")

		if name == "" || hard == "" || app == "" {
			w.Write([]byte("Missing parameters"))
			return
		}
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
	if r.Method == "POST" {
		recieveUrl, _ := url.Parse(r.URL.String())

		param, _ := url.ParseQuery(recieveUrl.RawQuery)

		name := param.Get("name")
		hard := param.Get("hard")
		app := param.Get("app")

		if name == "" || hard == "" || app == "" {
			w.Write([]byte("Missing parameters"))
			return
		}

		pushOperation(db, name, hard, app, w)
	}

}

func pushOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) {
	latestID, err := db.Query("SELECT max(id) FROM filepath")
	if err != nil {
		log.Fatal(err)
	}
	defer latestID.Close()

	if latestID.Next() {
		var id int
		err := latestID.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		newID := id + 1
		_, err = db.Exec("INSERT INTO filepath (id, filename, hardlayer, applayer) VALUES (?, ?, ?, ?)", newID, name, hard, app)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "File information inserted successfully with ID: %d", newID)
	}
}

func pullHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		recieveUrl, _ := url.Parse(r.URL.String())

		param, _ := url.ParseQuery(recieveUrl.RawQuery)

		name := param.Get("name")
		hard := param.Get("hard")
		app := param.Get("app")

		if name == "" || hard == "" || app == "" {
			w.Write([]byte("Missing parameters"))
			return
		}

		pullOperation(db, name, hard, app)
	}
}

func pullOperation(db *sql.DB, name string, hard string, app string) {

}

func overWriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		recieveUrl, _ := url.Parse(r.URL.String())

		param, _ := url.ParseQuery(recieveUrl.RawQuery)

		name := param.Get("name")
		hard := param.Get("hard")
		app := param.Get("app")

		if name == "" || hard == "" || app == "" {
			w.Write([]byte("Missing parameters"))
			return
		}

		overWriteOperation(db, name, hard, app)
	}
}

func overWriteOperation(db *sql.DB, name string, hard string, app string) {

}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	http.HandleFunc("/test", testHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/overwrite", overWriteHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Connecting to the database...")
	dbConnect()
	defer db.Close()

	fmt.Println("Server is running on port 50080...")
	http.ListenAndServe(":50080", nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "POST" {

		name := r.URL.Query().Get("name")
		hard := r.URL.Query().Get("hard")
		app := r.URL.Query().Get("app")

		if name == "" || hard == "" || app == "" {
			w.Write([]byte("Missing parameters"))
			return
		}
	}

}

func dbConnect() {
	data, err := sql.Open("mysql", "fileserver:fileserver@tcp(localhost:3306)/fileserver")
	if err != nil {
		log.Fatal(err)
	}

	err = data.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")
	db = data
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	UploadOperation(db, name, hard, app, w)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	rows, err := db.Query("SELECT path FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error querying file information", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			http.Error(w, "Error scanning path", http.StatusInternalServerError)
			return
		}
		if path == "" {
			http.Error(w, "No file information found for the given parameters", http.StatusNotFound)
			return
		}
	}
	DownloadOperation(db, name, hard, app, w)
}

func overWriteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	OverwriteOperation(db, name, hard, app)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}
	DeleteOperation(db, name, hard, app, w)
}

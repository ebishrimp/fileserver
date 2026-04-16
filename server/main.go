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

	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/pull", pullHandler)
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

func pushHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	pushOperation(db, name, hard, app, w)
}

// Check if a record exists

/*latestID, err := db.Query("SELECT max(id) FROM filepath")
if err != nil {
	log.Fatal(err)
}
defer latestID.Close()

if latestID.Next() {
	var id int
	err := latestID.Scan(&id)
	if err != nil {
		http.Error(w, "Error scanning ID", http.StatusInternalServerError)
		return
	}
	newID := id + 1
	_, err = db.Exec("INSERT INTO filepath (id, filename, hardlayer, applayer, path) VALUES (?, ?, ?, ?, ?)", newID, name, hard, app, "/"+hard+"/"+app+"/"+name)
	if err != nil {
		http.Error(w, "Error inserting file information", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File information inserted successfully with ID: %d", newID)
}*/

func pullHandler(w http.ResponseWriter, r *http.Request) {
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
	pullOperation(db, name, hard, app, w)
}

func pullOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) {
	pathRow, err := db.Query("SELECT path FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error querying file information", http.StatusInternalServerError)
		return
	}
	defer pathRow.Close()

	if pathRow.Next() {
		var path string
		err := pathRow.Scan(&path)
		if err != nil {
			http.Error(w, "Error scanning path", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "File path: %s", path)
	} else {
		http.Error(w, "No file information found for the given parameters", http.StatusNotFound)
	}
}

func overWriteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	overWriteOperation(db, name, hard, app)
}

func overWriteOperation(db *sql.DB, name string, hard string, app string) {

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}
	deleteOperation(db, name, hard, app, w)
}

func deleteOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) {
	selectID, err := db.Query("SELECT id FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error selecting file information", http.StatusInternalServerError)
		return
	}
	defer selectID.Close()

	if selectID.Next() {
		var id int
		err := selectID.Scan(&id)
		if err != nil {
			http.Error(w, "Error scanning ID", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("DELETE FROM filepath WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Error deleting file information", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "File information with ID: %d deleted successfully", id)
	} else {
		http.Error(w, "No file information found for the given parameters", http.StatusNotFound)
	}
}

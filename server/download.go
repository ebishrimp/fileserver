package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func DownloadOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) {
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

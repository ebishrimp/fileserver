package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
)

func DeleteOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) {
	selectID, err := db.Query("SELECT id FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error selecting file information", http.StatusInternalServerError)
		return
	}
	selectPath, err := db.Query("SELECT path FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error selecting file path", http.StatusInternalServerError)
		return
	}
	defer selectID.Close()
	defer selectPath.Close()

	//path operation
	if selectPath.Next() {
		var path string
		err := selectPath.Scan(&path)
		if err != nil {
			http.Error(w, "Error scanning path", http.StatusInternalServerError)
			return
		}
		deleteDirOp(path, w)
	}

	//id operation
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

func deleteDirOp(path string, w http.ResponseWriter) {
	err := os.Remove(path)
}

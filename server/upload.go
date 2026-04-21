package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func UploadOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) {
	duplicateCheck, err := db.Query("SELECT id FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error checking for file information", http.StatusInternalServerError)
		return
	}
	defer duplicateCheck.Close()

	if duplicateCheck.Next() {
		var path string
		err := duplicateCheck.Scan(&path)
		if err != nil {
			http.Error(w, "Error scanning path", http.StatusInternalServerError)
			return
		}
		if path != "" {
			http.Error(w, "File information already exists for the given parameters", http.StatusConflict)
			return
		}
	}

	_, err = db.Exec("INSERT INTO filepath (filename, hardlayer, applayer, path) VALUES (?, ?, ?, ?)", name, hard, app, "/"+hard+"/"+app+"/"+name)
	if err != nil {
		http.Error(w, "Error inserting file information", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File information inserted successfully")

	idRaw, err := db.Query("SELECT LAST_INSERT_ID()")
	if err != nil {
		http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
		return
	}
	defer idRaw.Close()

	if idRaw.Next() {
		var id int
		err := idRaw.Scan(&id)
		if err != nil {
			http.Error(w, "Error scanning ID", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "ID: %d", id)
	}

	env, err := user.Current()
	if err != nil {
		http.Error(w, "Error getting host user information", http.StatusInternalServerError)
		return
	}
	userEnv := env.Username
	var dirpath string = filepath.Join("/", "home", userEnv, "data", hard, app)
	uploadDirOp(dirpath, w)
	//make directory before uploading file
}

func uploadDirOp(dirpath string, w http.ResponseWriter) {
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Directory created successfully at %s", dirpath)
}

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func UploadOperation(db *sql.DB, name string, hard string, app string, w http.ResponseWriter) error {
	duplicateCheck, err := db.Query("SELECT path FROM filepath WHERE filename = ? AND hardlayer = ? AND applayer = ?", name, hard, app)
	if err != nil {
		http.Error(w, "Error checking for file information", http.StatusInternalServerError)
		return err
	}
	defer duplicateCheck.Close()

	if duplicateCheck.Next() {
		var path string
		err := duplicateCheck.Scan(&path)
		if err != nil {
			http.Error(w, "Error scanning path", http.StatusInternalServerError)
			return err
		}
		if path != "" {
			http.Error(w, "File information already exists for the given parameters", http.StatusConflict)
			return fmt.Errorf("File information already exists for the given parameters")
		}
	}

	_, err = db.Exec("INSERT INTO filepath (filename, hardlayer, applayer, path) VALUES (?, ?, ?, ?)", name, hard, app, "/"+hard+"/"+app+"/"+name)
	if err != nil {
		http.Error(w, "Error inserting file information", http.StatusInternalServerError)
		return err
	}

	fmt.Fprintf(w, "File information inserted successfully")

	idRaw, err := db.Query("SELECT LAST_INSERT_ID()")
	if err != nil {
		http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
		return err
	}
	defer idRaw.Close()

	if idRaw.Next() {
		var id int
		err := idRaw.Scan(&id)
		if err != nil {
			http.Error(w, "Error scanning ID", http.StatusInternalServerError)
			return err
		}
		fmt.Fprintf(w, "ID: %d", id)
	}

	var dirpath string = filepath.Join("/", "srv", "fileserver", hard, app)
	if err := uploadDirOp(dirpath, w); err != nil {
		return err
	}
	//make directory before uploading file
	return nil
}

func uploadDirOp(dirpath string, w http.ResponseWriter) error {
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return err
	}
	fmt.Fprintf(w, "Directory created successfully at %s", dirpath)
	return nil
}

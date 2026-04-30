package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	confparser "github.com/ebishrimp/config-file-parser-go"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var conf *confparser.Config
var configFilePath string = "./fileserver.conf.template"

// Basic configurations
var ListenPort string
var sqlUsername string
var sqlPassword string

// Request restrictions
var allowUpload bool
var allowDownload bool
var allowOverwrite bool
var allowDelete bool

var whiteList bool

// pseudo raid 0 settings
var raid0 bool
var raidpath string

// log settings
var logfile string

// IP address and network restrictions in the whitelist
var whitelistPath string = "/etc/fileserver/whitelist.conf"
var IPs *confparser.MultiConfig
var allowedIPs []net.IP
var allowedSubnets []*net.IPNet

func main() {
	configParse()
	configLoad(conf)
	fmt.Println("configs loaded successfully")

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/overwrite", overWriteHandler)
	http.HandleFunc("/delete", deleteHandler)
	fmt.Println("Handlers registered successfully")

	fmt.Println("Connecting to the database...")
	dbConnect()
	defer db.Close()

	fmt.Println("Server is running on port " + ListenPort + "...")
	if err := http.ListenAndServe(":"+ListenPort, nil); err != nil {
		log.Println("Error starting server on port " + ListenPort + ", starting on port 50080 instead")
		if err := http.ListenAndServe(":50080", nil); err != nil {
			log.Fatal(err)
		}
	}
}

func configParse() {
	f, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	configs, err := confparser.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	conf = configs
}

func configLoad(c *confparser.Config) {
	ListenPort = c.GetValue("ListenPort")

	sqlUsername = c.GetValue("mysqlUsername")
	sqlPassword = c.GetValue("mysqlPassword")

	allowUpload = c.GetValue("allowUpload") == "yes"
	allowDownload = c.GetValue("allowDownload") == "yes"
	allowOverwrite = c.GetValue("allowOverwrite") == "yes"
	allowDelete = c.GetValue("allowDelete") == "yes"

	raid0 = c.GetValue("raid0") == "yes"
	raidpath = c.GetValue("raidpath")
	if f, err := os.Stat(raidpath); os.IsNotExist(err) || !f.IsDir() {
		fmt.Println("RAID 0 path does not exist or is not a directory, disabling RAID 0")
		raid0 = false
	}

	whiteList = c.GetValue("whiteList") == "yes"
	if f, err := os.Stat("/etc/fileserver/whitelist.conf"); os.IsNotExist(err) || f.IsDir() {
		fmt.Println("White list file does not exist or is a directory, disabling white list")
		whiteList = false
	}

	logfile = c.GetValue("logfile")
	if f, err := os.Stat(logfile); os.IsNotExist(err) || f.IsDir() {
		fmt.Println("Log file does not exist or is a directory, using default log file: fileserver.log")
		logfile = "/var/log/fileserver/fileserver.log"
		os.MkdirAll("/var/log/fileserver", os.ModePerm)
		_, err := os.Create(logfile)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func IPParse() {
	if whiteList {
		f, err := os.Open(whitelistPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		allowIPs, err := confparser.ParseMultipleValues(f)
		if err != nil {
			log.Fatal(err)
		}
		IPs = allowIPs
	}
}

func IPLoad() {
	if whiteList {
		//string -> net.IP (convert IP address string to net.IP type) then store in allowedIPs and allowedSubnets (if it's a subnet)
	}
}

func dbConnect() {
	data, err := sql.Open("mysql", sqlUsername+":"+sqlPassword+"@tcp(localhost:3306)/fileserver")
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	if !allowUpload {
		http.Error(w, "Upload not allowed", http.StatusForbidden)
		return
	}

	UploadOperation(db, name, hard, app, w)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	if !allowDownload {
		http.Error(w, "Download not allowed", http.StatusForbidden)
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
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	if !allowOverwrite {
		http.Error(w, "Overwrite not allowed", http.StatusForbidden)
		return
	}

	OverwriteOperation(db, name, hard, app)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	if !allowDelete {
		http.Error(w, "Delete not allowed", http.StatusForbidden)
		return
	}

	DeleteOperation(db, name, hard, app, w)
}

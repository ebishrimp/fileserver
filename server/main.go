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
var conf *confparser.ConfigurationMap
var configFilePath string = "/etc/fileserver/fileserver.conf"

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

// pseudo raid 1 settings
var raid1 bool
var raidpath string

// log settings
var logfile string

// IP address and network restrictions in the whitelist
var whitelistPath string = "/etc/fileserver/whitelist.conf"
var IPs *confparser.ConfigurationMap
var allowedIPs []net.IP
var allowedSubnets []*net.IPNet
var domain []string

func main() {
	configParse()
	configLoad(conf)
	fmt.Println("configs loaded successfully")

	if whiteList {
		fmt.Println("IP whitelist enabled, loading IP whitelist...")
		IPParse()
		IPLoad()
		fmt.Println("IP whitelist loaded successfully")
	} else {
		fmt.Println("IP whitelist disabled")
	}

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

	configs, pErr := confparser.ParseConfig(f)
	if pErr != nil {
		log.Fatal(pErr)
	}
	conf = configs
}

func configLoad(c *confparser.ConfigurationMap) {
	ListenPort = c.String("ListenPort")

	sqlUsername = c.String("mysqlUsername")
	sqlPassword = c.String("mysqlPassword")

	Up, err := c.Bool("allowUpload")
	if err != nil {
		fmt.Println("Error parsing allowUpload, defaulting to false")
		Up = false
	}
	allowUpload = Up

	Down, err := c.Bool("allowDownload")
	if err != nil {
		fmt.Println("Error parsing allowDownload, defaulting to false")
		Down = false
	}
	allowDownload = Down

	Over, err := c.Bool("allowOverwrite")
	if err != nil {
		fmt.Println("Error parsing allowOverwrite, defaulting to false")
		Over = false
	}
	allowOverwrite = Over

	Delete, err := c.Bool("allowDelete")
	if err != nil {
		fmt.Println("Error parsing allowDelete, defaulting to false")
		Delete = false
	}
	allowDelete = Delete

	r1, err := c.Bool("raid1")
	if err != nil {
		fmt.Println("Error parsing raid1, defaulting to false")
		r1 = false
	}
	raid1 = r1

	if raid1 {
		fmt.Println("pseudo RAID 1 enabled, checking RAID 1 path...")
		raidpath = c.String("raidpath")
		if f, err := os.Stat(raidpath); os.IsNotExist(err) || !f.IsDir() {
			fmt.Println("RAID 1 path does not exist or is not a directory, disabling pseudo RAID 1")
			raid1 = false
		}
	}

	wl, err := c.Bool("whiteList")
	if err != nil {
		fmt.Println("Error parsing whiteList, defaulting to false")
		wl = false
	}
	whiteList = wl

	if f, err := os.Stat("/etc/fileserver/whitelist.conf"); os.IsNotExist(err) || f.IsDir() {
		fmt.Println("White list file does not exist or is a directory, disabling white list")
		whiteList = false
	}

	logfile = c.String("logfile")
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

		allowIPs, pErr := confparser.ParseConfig(f)
		if pErr != nil {
			log.Fatal(pErr)
		}
		IPs = allowIPs
	}
}

func IPLoad() {
	// Load allowed IPs and subnets from the whitelist configuration , and store them in the allowedIPs and allowedSubnets slices
	if whiteList {
		//address
		stringAllowedIPs := IPs.StringSlice("address")
		for i := 0; i < len(stringAllowedIPs); i++ {
			ip := net.ParseIP(stringAllowedIPs[i])
			if ip != nil {
				allowedIPs = append(allowedIPs, ip)
			} else {
				fmt.Printf("Invalid IP address in whitelist: %s\n", stringAllowedIPs[i])
			}
		}

		//subnet
		stringAllowedSubnets := IPs.StringSlice("subnet")
		for i := 0; i < len(stringAllowedSubnets); i++ {
			_, subnet, err := net.ParseCIDR(stringAllowedSubnets[i])
			if err != nil {
				fmt.Printf("Invalid subnet in whitelist: %s\n", stringAllowedSubnets[i])
			} else {
				allowedSubnets = append(allowedSubnets, subnet)
			}
		}

		//domain
		domain = IPs.StringSlice("domain")
		solved, err := NamesSolve(domain)
		if err != nil {
			fmt.Printf("Error resolving domain names: %v\n", err)
		}
		allowedIPs = append(allowedIPs, solved...)
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
	if !allowUpload {
		http.Error(w, "Upload not allowed", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if ipInfo := GetClientIP(r); !AuthorizeIP(ipInfo, w) {
		http.Error(w, "Your IP address is not allowed to access", http.StatusForbidden)
		return
	}

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
	if !allowDownload {
		http.Error(w, "Download not allowed", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if ipInfo := GetClientIP(r); !AuthorizeIP(ipInfo, w) {
		http.Error(w, "Your IP address is not allowed to access", http.StatusForbidden)
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
	if !allowOverwrite {
		http.Error(w, "Overwrite not allowed", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if ipInfo := GetClientIP(r); !AuthorizeIP(ipInfo, w) {
		http.Error(w, "Your IP address is not allowed to access", http.StatusForbidden)
		return
	}

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
	if !allowDelete {
		http.Error(w, "Delete not allowed", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if ipInfo := GetClientIP(r); !AuthorizeIP(ipInfo, w) {
		http.Error(w, "Your IP address is not allowed to access", http.StatusForbidden)
		return
	}

	name := r.URL.Query().Get("name")
	hard := r.URL.Query().Get("hard")
	app := r.URL.Query().Get("app")

	if name == "" || hard == "" || app == "" {
		w.Write([]byte("Missing parameters"))
		return
	}

	DeleteOperation(db, name, hard, app, w)
}

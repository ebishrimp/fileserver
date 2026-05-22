package main

import (
	"fmt"
	"os"
	"strconv"
)

type AccessLog struct {
	IP        string
	Operation string
	Path      string
	Error     error
}

func (log *AccessLog) WriteLog(path string) {
	var compressPath string

	if fileSizeLarge(logfile) {
		idx := 1
		exists := true
		for exists {
			if f, err := os.Stat(path + "." + strconv.Itoa(idx)); os.IsNotExist(err) || f.IsDir() {
				exists = false
				compressPath = path + "." + strconv.Itoa(idx) + ".gz"
			} else {
				idx++
			}
		}
		//write log to compressPath
		compressLog(compressPath)
	}

}

func fileSizeLarge(path string) bool {
	logfileinfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return logfileinfo.Size() > int64(maxlogfilesize)
}

func compressLog(path string) error {
	return nil
}

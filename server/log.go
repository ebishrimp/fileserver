package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type AccessLog struct {
	IP        string
	Operation string
	Path      string
}

func (logstat *AccessLog) WriteLog(path string) {
	var compressPath string
	var log *os.File

	if fileSizeLarge(path) {
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
		err := compressLog(compressPath, path)
		if err != nil {
			fmt.Printf("failed to compress logfile: %s\n, err: %s", compressPath, err)
		}
		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("failed to open the logfile. path: %s\n", path)
		}
		log = f
	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("failed to open the logfile. path: %s, err: "+err.Error()+"\n", path, err)
		}
		log = f
	}
	defer log.Close()

	writer := bufio.NewWriter(log)
	currentTime := time.Now().Format("2006-01-02T15:05:04")

	logContent := fmt.Sprintf("[%s] Client: %s, Operation: %s, Path: %s\n", currentTime, logstat.IP, logstat.Operation, logstat.Path)

	_, err := writer.WriteString(logContent)
	if err != nil {
		fmt.Println("failed to write log: %w", err)
		return
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

func compressLog(path string, logsrc string) error {
	dist, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dist.Close()

	gw, err := gzip.NewWriterLevel(dist, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gw.Close()

	src, err := os.Open(logsrc)
	if err != nil {
		return err
	}
	defer src.Close()

	if _, err := io.Copy(gw, src); err != nil {
		return err
	}

	return nil
}

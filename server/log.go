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
	Error     error
}

func (logstat *AccessLog) WriteLog(path string) {
	var compressPath string
	var log *os.File
	var errstr string

	if fileSizeLarge(path) {
		temppath := path + "." + time.Now().Format("2006-01-02")
		idx := 1
		exists := true
		for exists {
			if f, err := os.Stat(temppath + "." + strconv.Itoa(idx)); os.IsNotExist(err) || f.IsDir() {
				exists = false
				compressPath = temppath + "." + strconv.Itoa(idx) + ".gz"
			} else {
				idx++
			}
		}
		//write log to compressPath
		go func() {
			err := compressLog(compressPath, path)
			if err != nil {
				fmt.Printf("failed to compress logfile: %s\n, err: %s", compressPath, err)
			}
		}()

		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("failed to open the logfile. path: %s\n", temppath)
		}
		defer f.Close()
		log = f
	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("failed to open the logfile. path: %s, err: "+err.Error()+"\n", path, err)
		}
		defer f.Close()
		log = f
	}
	defer log.Close()

	writer := bufio.NewWriter(log)
	currentTime := time.Now().Format("2006-01-02T15:05:04")

	if logstat.Error != nil {
		errstr = logstat.Error.Error()
	} else {
		errstr = "success"
	}

	logContent := fmt.Sprintf("[%s] Client: %s, Operation: %s, Path: %s, Status: %s\n", currentTime, logstat.IP, logstat.Operation, logstat.Path, errstr)

	_, err := writer.WriteString(logContent)
	if err != nil {
		fmt.Println("failed to write log: %w", err)
		return
	}
	if err := writer.Flush(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(logContent)
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

package main

type AccessLog struct {
	IP        string
	Operation string
	Path      string
	Error     error
}

func WriteLog(log AccessLog) {

}

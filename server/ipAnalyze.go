package main

import (
	"net"
	"net/http"
)

type IPInfo struct {
	address string
	data    net.IP
}

func GetClientIP(r *http.Request) IPInfo {
	var ipinfo IPInfo
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ipinfo
	}
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return ipinfo
	}
	ipinfo.address = clientIP.String()
	ipinfo.data = clientIP
	return ipinfo
}

func AuthorizeIP(ipinfo IPInfo) bool {
	//if ip is in the whitelist, return true. when whitelist system is allowed
	return false
}

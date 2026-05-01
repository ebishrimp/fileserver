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

func AuthorizeIP(ipinfo IPInfo, w http.ResponseWriter) bool {
	if !whiteList {
		return true
	}

	if ipinfo.data == nil {
		http.Error(w, "Invalid IP address", http.StatusForbidden)
		return false
	}
	return isIPAllowed(ipinfo.data) || isSubnetAllowed(ipinfo.data)
}

func isIPAllowed(ip net.IP) bool {
	for i := 0; i < len(allowedIPs); i++ {
		if allowedIPs[i].Equal(ip) {
			return true
		}
	}
	return false
}

func isSubnetAllowed(ip net.IP) bool {
	for i := 0; i < len(allowedSubnets); i++ {
		if allowedSubnets[i].Contains(ip) {
			return true
		}
	}
	return false
}

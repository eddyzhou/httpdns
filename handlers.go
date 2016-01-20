package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func GetOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, "pong")
}

func ResolveHandler(rw http.ResponseWriter, req *http.Request) {
	dn := req.FormValue("dn")
	if dn == "" {
		log.Printf("invalid req: %v", req.URL)
		http.Error(rw, "Invalid parameter", http.StatusBadRequest)
		return
	}
	resolver := GetResolver()
	if ip, err := resolver.GetIp(dn); err == nil {
		log.Printf("dn: %v, ip: %v", dn, ip)
		rw.WriteHeader(http.StatusOK)
		io.WriteString(rw, ip)
		return
	}

	if ips, err := net.LookupIP(dn); err == nil {
		resp := ips[0].String()
		if len(ips) > 1 {
			strs := make([]string, len(ips))
			for i, v := range ips {
				strs[i] = v.String()
			}
			resp = strings.Join(strs, ";")
		}
		log.Printf("dn: %v, ip: %v", dn, resp)
		rw.WriteHeader(http.StatusOK)
		io.WriteString(rw, resp)
		return
	}

	http.NotFound(rw, req)
}

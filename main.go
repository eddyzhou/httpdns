package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

const (
	MaxOpenfile uint64 = 128000
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Setenv("GOTRACEBACK", "crash")

	lim := syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim)
	if lim.Cur < MaxOpenfile || lim.Max < MaxOpenfile {
		lim.Cur = MaxOpenfile
		lim.Max = MaxOpenfile
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
	}

	log.SetOutput(io.Writer(os.Stderr))
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	pid2file()
}

func pid2file() {
	file, err := os.OpenFile("httpdns.pid", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln("error opening file: %v", err)
		return
	}
	defer file.Close()
	file.Write([]byte(strconv.Itoa(syscall.Getpid())))
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [configFile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "listen address(0.0.0.0:80)")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
	configFile := args[0]

	InitResolver(configFile)

	http.HandleFunc("/d", GetOnly(ResolveHandler))
	http.HandleFunc("/p", GetOnly(PingHandler))
	log.Fatal(http.ListenAndServe(addr, nil))
}

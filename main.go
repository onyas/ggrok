package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/onyas/ggrok/core"
)

var client bool
var proxyServer string
var port int
var config *core.Config

func init() {
	flag.BoolVar(&client, "client", false, "start client")
	flag.StringVar(&proxyServer, "proxyServer", "", "provide server address, for example: https://proxy.yourdomain.com")
	flag.IntVar(&port, "port", -1, "provide port, for example: 8080")
	config = core.NewConfig()
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if client {
		if proxyServer != "" {
			saveToConfig(proxyServer)
			return
		}

		if port != -1 {
			startProxy(port)
		}
		log.Println("Using -proxyServer or -port args")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	s := core.NewServer()

	http.HandleFunc("/$$ggrok", s.Register)
	http.HandleFunc("/", s.Proxy)
	log.Println("Server started at port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func saveToConfig(proxyServer string) {
	config.SaveToConfig(proxyServer)
	log.Println("config success, the proxy server is ", proxyServer)
}

func startProxy(port int) {
	proxyServer := config.ReadConfig()
	if proxyServer == "" {
		log.Fatal("Config proxy server first. ggrok -client -proxyServer ")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		c := core.NewClient(proxyServer, port)
		c.Start()
	}()
	<-done

}

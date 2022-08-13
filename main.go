package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/onyas/ggrok/ggrok"
)

var addr string
var client bool

func init() {
	flag.StringVar(&addr, "serverAddr", "localhost:8080", "http service address")
	flag.BoolVar(&client, "client", false, "start client")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if client {
		done := make(chan struct{})
		go func() {
			c := ggrok.NewClient()
			c.Start()
		}()

		<-done
	}

	s := ggrok.NewServer()

	http.HandleFunc("/$$ggrok", s.Register)
	http.HandleFunc("/", s.Proxy)
	log.Fatal(http.ListenAndServe(addr, nil))
}

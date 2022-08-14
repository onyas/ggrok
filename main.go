package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/onyas/ggrok/core"
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
			defer close(done)

			c := core.NewClient()
			c.Start(3000)
		}()

		<-done
	}

	s := core.NewServer()

	http.HandleFunc("/$$ggrok", s.Register)
	http.HandleFunc("/", s.Proxy)
	log.Fatal(http.ListenAndServe(addr, nil))
}

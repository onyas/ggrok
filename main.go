package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/onyas/ggrok/core"
)

var client bool

func init() {
	flag.BoolVar(&client, "client", false, "start client")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

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
	log.Println("Server started at port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

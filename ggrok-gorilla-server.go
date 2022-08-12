package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var connections = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func register(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	connections[r.Host] = c
	log.Println("current connections: %s", connections)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
func proxy(w http.ResponseWriter, r *http.Request) {
	remoteConn := connections[r.Host]
	if remoteConn == nil {
		io.WriteString(w, "client not register")
		return
	}

	reqStr, err := captureRequestData(r)
	if err != nil {
		log.Println("captureRequestData error:", err)
	}
	log.Println("req serialized: %s", reqStr)

	// remoteConn.WriteMessage(100, []byte(reqStr))
	type WebSocketRequest struct {
		Req string
		URL string
	}
	reqRemote := WebSocketRequest{Req: reqStr, URL: r.URL.String()}

	remoteConn.WriteJSON(reqRemote)

	type WebSocketResponse struct {
		Status      string // e.g. "200 OK"
		StatusCode  int    // e.g. 200
		Proto       string // e.g. "HTTP/1.0"
		Header      map[string][]string
		Body        []byte
		ContentType string
	}
	var wsRes WebSocketResponse
	err = remoteConn.ReadJSON(&wsRes)
	if err != nil {
		log.Println("read remote client response error", err)
	}
	log.Println("remote client response: %s", wsRes)

	copyHeader(w.Header(), wsRes.Header)
	w.WriteHeader(wsRes.StatusCode)
	w.Header().Set("Content-Type", wsRes.ContentType)
	io.Copy(w, bytes.NewReader(wsRes.Body))
}

func captureRequestData(req *http.Request) (string, error) {
	var b = &bytes.Buffer{} // holds serialized representation
	var err error
	if err = req.Write(b); err != nil { // serialize request to HTTP/1.1 wire format
		return "", err
	}
	return b.String(), nil
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/$$ggrok", register)
	http.HandleFunc("/", proxy)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

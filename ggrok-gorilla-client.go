package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func captureResponseData(resp *http.Response) (string, error) {
	rump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("local response dump error ", err)
		return "", err
	}
	return string(rump), nil
}

func main() {
	type RemoteRequest struct {
		Req string
		URL string
	}

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/$$ggrok"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read message error:", err)
				continue
			}
			log.Printf("recv: %s", message)

			var websocketReq RemoteRequest
			if err := json.Unmarshal(message, &websocketReq); err != nil {
				log.Println("json.Unmarshal error", err)
				continue
			}

			var localRequest *http.Request
			r := bufio.NewReader(bytes.NewReader([]byte(websocketReq.Req)))
			if localRequest, err = http.ReadRequest(r); err != nil { // deserialize request
				log.Printf("deserialize request error", err)
				continue
			}

			//TODO: change to config
			localRequest.RequestURI = ""
			u, err := url.Parse("/ada08e16-2112-4720-8fcb-18f2f8e47c2d")
			if err != nil {
				log.Printf("parse url error", err)
			}
			localRequest.URL = u
			localRequest.URL.Scheme = "https"
			localRequest.URL.Host = "webhook.site"
			resp, err := (&http.Client{}).Do(localRequest)
			if err != nil {
				log.Println("local http request error:", err)
				continue
			}

			respStr, err := captureResponseData(resp)
			if err != nil {
				continue
			}

			type WebSocketResponse struct {
				Status      string // e.g. "200 OK"
				StatusCode  int    // e.g. 200
				Proto       string // e.g. "HTTP/1.0"
				Header      map[string][]string
				Body        []byte
				ContentType string
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("read local response error ", err)
			}
			wsRes := WebSocketResponse{Status: resp.Status, StatusCode: resp.StatusCode,
				Proto: resp.Proto, Header: resp.Header, Body: body, ContentType: resp.Header.Get("Content-Type")}

			log.Println("client send response: %s", respStr)
			c.WriteJSON(wsRes)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

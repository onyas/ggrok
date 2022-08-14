package core

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type GGrokClient struct {
	RemoteServer   string
	ProxyLocalPort int
}

func NewClient(s string, p int) *GGrokClient {
	return &GGrokClient{RemoteServer: s, ProxyLocalPort: p}
}

func (ggclient *GGrokClient) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: ggclient.RemoteServer, Path: "/$$ggrok"}
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
			websocketReq := readWebSocketReq(c)

			localRequest := socketToLocalRequest(websocketReq, ggclient.ProxyLocalPort)
			resp, err := (&http.Client{}).Do(localRequest)
			if err != nil {
				log.Println("local http request error:", err)
				continue
			}

			wsRes := localResponseToWebSocketResponse(resp)

			// log.Printf("client send response: %s \n", wsRes.Body)
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

func localResponseToWebSocketResponse(resp *http.Response) WebSocketResponse {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read local response error ", err)
	}
	resp.Body.Close()
	wsRes := WebSocketResponse{Status: resp.Status, StatusCode: resp.StatusCode,
		Proto: resp.Proto, Header: resp.Header, Body: body, ContentType: resp.Header.Get("Content-Type")}
	return wsRes
}

func readWebSocketReq(c *websocket.Conn) WebSocketRequest {
	var websocketReq WebSocketRequest
	if err := c.ReadJSON(&websocketReq); err != nil {
		log.Println("json.Unmarshal error", err)
		return websocketReq
	}
	log.Printf("recv: %s", websocketReq)

	return websocketReq
}

// deserialize request
//TODO: change to config
func socketToLocalRequest(websocketReq WebSocketRequest, port int) *http.Request {
	r := bufio.NewReader(bytes.NewReader([]byte(websocketReq.Req)))
	localRequest, err := http.ReadRequest(r)
	if err != nil {
		log.Println("deserialize request error", err)
		return localRequest
	}

	localRequest.RequestURI = ""
	u, err := url.Parse(websocketReq.URL)
	if err != nil {
		log.Println("parse url error", err)
	}
	localRequest.URL = u
	localRequest.URL.Scheme = "http"
	localRequest.URL.Host = "localhost:" + strconv.Itoa(port)
	return localRequest
}

package core

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

var connections = make(map[string]*Connection)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	gconn := &Connection{
		Socket: c,
		mu:     sync.Mutex{},
	}
	connections[r.Host] = gconn
	log.Println("current connections: ", connections)
}

func (s *Server) Proxy(w http.ResponseWriter, r *http.Request) {
	remoteConn := connections[r.Host]
	if remoteConn == nil || remoteConn.Socket == nil {
		io.WriteString(w, "client not register")
		return
	}

	wsRequest := httpRequestToWebSocketRequest(r)

	wsRes := triggerWS(remoteConn, wsRequest)

	wsResToHttpResponse(w, wsRes)
}

func copyHeader(dst http.ResponseWriter, src WebSocketResponse) {
	for k, vv := range src.Header {
		for _, v := range vv {
			dst.Header().Add(k, v)
		}
	}
	dst.WriteHeader(src.StatusCode)
	dst.Header().Set("Content-Type", src.ContentType)
}

func wsResToHttpResponse(w http.ResponseWriter, wsRes WebSocketResponse) {
	copyHeader(w, wsRes)
	io.Copy(w, bytes.NewReader(wsRes.Body))
}

func triggerWS(remoteConn *Connection, reqRemote WebSocketRequest) WebSocketResponse {
	remoteConn.mu.Lock()
	defer remoteConn.mu.Unlock()

	remoteConn.Socket.WriteJSON(reqRemote)

	var wsRes WebSocketResponse
	err := remoteConn.Socket.ReadJSON(&wsRes)
	if err != nil {
		log.Println("read remote client response error", err)
	}
	log.Println("remote client response: ", wsRes)
	return wsRes
}

func httpRequestToWebSocketRequest(r *http.Request) (ws WebSocketRequest) {
	reqStr, err := captureRequestData(r)
	if err != nil {
		log.Println("captureRequestData error:", err)
	}
	log.Println("req serialized: ", reqStr)

	reqRemote := WebSocketRequest{Req: reqStr, URL: r.URL.String()}
	return reqRemote
}

func captureRequestData(req *http.Request) (string, error) {
	var b = &bytes.Buffer{} // holds serialized representation
	var err error
	if err = req.Write(b); err != nil { // serialize request to HTTP/1.1 wire format
		return "", err
	}
	return b.String(), nil
}

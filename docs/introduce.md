# Create your own reverse proxy based on Golang and WebSocket

When we are developing, sometimes we want to expose our own developed interfaces to other developers or third-party services to facilitate our debugging and troubleshooting, so we need some mechanism to expose our local service interfaces to the Internet. This article will introduce how to realize this function through golang and WebSocket

## Why do we need to develop our own proxy services

Currently, many proxy services are available, such as ngrok and localtunnel. However, ngrok has a disadvantage: the domain name provided can only be used for a few hours, and then a new domain name needs to be generated. If you want a fixed domain name, you need to spend money. However, our own proxy can use a fixed domain name. If frontend developers use it, it is very convenient without changing the domain over time.

## ggrok introduction

![ggrok-flow](https://github.com/onyas/ggrok/blob/main/docs/flow.jpg?raw=true)

Ggrok is a proxy application implemented through golang and WebSocket. You can use the Heroku button on the GitHub [repo](https://github.com/onyas/ggrok) to deploy it conveniently, and then you can have a fixed domain name.

## How to implement

### Step1 establish a WebSocket connection between the server and the client

The server is based on the gorilla and listens for WebSocket connections

```golang
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

http.HandleFunc("/$$ggrok", s.Register)
```

The Client connect to the Server

```golang
func (ggclient *GGrokClient) Proxy() {
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

	for {
		select {
		case <-done:
			return
		}
	}
}
```

### Step2 After receiving the HTTP request, the server converts it into a websocket message and forwards it to the client


```golang
http.HandleFunc("/", s.Proxy)

func (s *Server) Proxy(w http.ResponseWriter, r *http.Request) {
	remoteConn := connections[r.Host]
	if remoteConn == nil || remoteConn.Socket == nil {
		io.WriteString(w, "client not register")
		return
	}

	wsRequest := httpRequestToWebSocketRequest(r)

	wsRes := triggerWS(remoteConn, wsRequest)
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
```

### Step3 After receiving the websocket message, the client forwards it to the localserver and returns the response of the localserver to the server

```golang
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
```

### Step 4 After receiving the response, the server returns it to the http response

```golang
func wsResToHttpResponse(w http.ResponseWriter, wsRes WebSocketResponse) {
	copyHeader(w, wsRes)
	io.Copy(w, bytes.NewReader(wsRes.Body))
}
```

So far, though ggrok, we have implemented the local service proxy and published it on the Internet. The above are some main codes. See [GitHub](https://github.com/onyas/ggrok) for details. Please create issue or PR if you have any problems, and jointly create a more robust open source system.
# 基于Golang和WebSocket打造自已的反向代理

当我们在开发的时候，有时想要把自已开发的接口暴露给其他开发者或者第三方的服务，方便我们调试和排查问题，那就需要某种机制把我们本地的服务接口暴露到互联网上，本文将要介绍如何通过Golang和WebSocket来实现这一功能

## 为什么我们需要开发自已的代理服务

目前已经有许多可用的代理服务了，比如ngrok和localtunnel,但ngrok有个缺点就是提供的域名只能用几个小时，然后需要新生成新的域名，如果想要固定域名就要花钱，但我们自已实现的代理可以用一个固定的域名，如果给前端同学来调试的话，不用改来改去，很方便。

## ggrok简介

![ggrok-flow](https://github.com/onyas/ggrok/blob/main/docs/flow.jpg?raw=true)

ggrok是通过Golang和WebSocket实现的代理应用，你可以通过Github[仓库](https://github.com/onyas/ggrok)上的Heroku按钮非常方便的部署，然后就可以拥有一个固定的域名了。

## 如何实现

### Step1 在服务器和客户端建立WebSocket连接

服务端基于gorilla，监听WebSocket连接

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

客户端连接服务端

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


### Step2 服务端收到http请求以后转成WebSocket消息转发给客户端

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

### Step3 客户端收到WebSocket消息以后转发到LocalServer,并把LocalServer的响应返回给服务端

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

### Step4 服务端收到响应以后返回给前端

```golang
func wsResToHttpResponse(w http.ResponseWriter, wsRes WebSocketResponse) {
	copyHeader(w, wsRes)
	io.Copy(w, bytes.NewReader(wsRes.Body))
}
```

至此，通过ggrok我们实现了本地服务的代理，并发布到互联网上。以上是一些主要的代码，详细的可以看[github](https://github.com/onyas/ggrok)上面的代码，有问题请提issue或者pr，共同打造更健壮的开源系统。
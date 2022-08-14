package core

type WebSocketRequest struct {
	Req string
	URL string
}

type WebSocketResponse struct {
	Status      string // e.g. "200 OK"
	StatusCode  int    // e.g. 200
	Proto       string // e.g. "HTTP/1.0"
	Header      map[string][]string
	Body        []byte
	ContentType string
}

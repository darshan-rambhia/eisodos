package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend interface {
	SetAlive(bool)
	IsAlive() bool
	GetURL() *url.URL
	GetActiveConnections() int
	Serve(http.ResponseWriter, *http.Request)
}

type backend struct {
	url          *url.URL
	alive        chan bool
	connections  chan int
	reverseProxy *httputil.ReverseProxy
}

func (b *backend) GetActiveConnections() int {
	return <-b.connections
}

func (b *backend) SetAlive(alive bool) {
	b.alive <- alive
}

func (b *backend) IsAlive() bool {
	return <-b.alive
}

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) Serve(rw http.ResponseWriter, req *http.Request) {
	b.connections <- (<-b.connections + 1)
	defer func() {
		b.connections <- (<-b.connections - 1)
	}()
	b.reverseProxy.ServeHTTP(rw, req)
}

func NewBackend(u *url.URL, rp *httputil.ReverseProxy) Backend {
	return &backend{
		url:          u,
		alive:        make(chan bool, 1),
		connections:  make(chan int, 1),
		reverseProxy: rp,
	}
}

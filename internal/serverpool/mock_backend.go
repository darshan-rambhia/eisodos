package serverpool

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/darshan-rambhia/eisodos/internal/backend"
)

type mockBackend struct {
	url               *url.URL
	proxy             *httputil.ReverseProxy
	activeConnections int
	alive             bool
}

func (b *mockBackend) GetURL() *url.URL {
	return b.url
}

func (b *mockBackend) GetActiveConnections() int {
	return b.activeConnections
}

func (b *mockBackend) IsAlive() bool {
	return b.alive
}

func (b *mockBackend) SetAlive(alive bool) {
	b.alive = alive
}

func (b *mockBackend) Serve(w http.ResponseWriter, r *http.Request) {
	b.activeConnections++
}

func newMockBackend(url *url.URL, proxy *httputil.ReverseProxy) backend.Backend {
	return &mockBackend{
		url:   url,
		proxy: proxy,
		alive: true,
	}
}

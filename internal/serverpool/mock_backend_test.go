package serverpool

import (
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestMockBackend(t *testing.T) {
	urlStr := "http://localhost:8080"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	mb := newMockBackend(u, proxy)

	t.Run("GetURL", func(t *testing.T) {
		if got := mb.GetURL(); got != u {
			t.Errorf("GetURL() = %v, want %v", got, u)
		}
	})

	t.Run("SetAlive", func(t *testing.T) {
		mb.SetAlive(true)
		if !mb.IsAlive() {
			t.Error("SetAlive(true) did not set backend to alive")
		}

		mb.SetAlive(false)
		if mb.IsAlive() {
			t.Error("SetAlive(false) did not set backend to not alive")
		}
	})
}

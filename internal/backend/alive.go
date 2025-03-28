package backend

import (
	"context"
	"log/slog"
	"net"
	"net/url"
)

func IsBackendAlive(ctx context.Context, aliveChannel chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		slog.Debug("Site unreachable", "error", err)
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}

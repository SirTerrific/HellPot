//go:build linux || darwin || freebsd

package http

import (
	"bufio"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/SirTerrific/HellPot/internal/config"
)

// The unix socket listener must serve with the configured server, not a default one:
// the Server header comes from deception.server_name and GetOnly must be honored.
func TestUnixSocketUsesConfiguredServer(t *testing.T) {
	// a short path under os.TempDir: unix socket paths are limited to ~104 bytes
	dir, err := os.MkdirTemp("", "hp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	sock := filepath.Join(dir, "s")

	config.UnixSocketPermissions = 0o600
	srv := &fasthttp.Server{
		Name:    "nginx-test",
		GetOnly: true,
		Handler: func(ctx *fasthttp.RequestCtx) { ctx.SetBodyString("ok") },
	}
	go func() { _ = listenOnUnixSocket(sock, srv) }()
	defer func() { _ = srv.Shutdown() }()

	var conn net.Conn
	for i := 0; i < 100; i++ {
		if conn, err = net.Dial("unix", sock); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("cannot connect to unix socket: %v", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))

	if _, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: x\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if got := resp.Header.Get("Server"); got != "nginx-test" {
		t.Fatalf("Server header = %q, want the configured name", got)
	}

	conn2, err := net.Dial("unix", sock)
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()
	_ = conn2.SetDeadline(time.Now().Add(5 * time.Second))
	if _, err = conn2.Write([]byte("POST / HTTP/1.1\r\nHost: x\r\nContent-Length: 0\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	resp2, err := http.ReadResponse(bufio.NewReader(conn2), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode == http.StatusOK {
		t.Fatal("POST was served: GetOnly is not honored on the unix socket")
	}
}

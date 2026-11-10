package http

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"

	"github.com/SirTerrific/HellPot/internal/config"
)

func TestBlacklistedUseragentGets404(t *testing.T) {
	nop := zerolog.Nop()
	log = &nop
	config.UseragentBlacklistMatchers = []string{"badbot"}

	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("/wp-login.php")
	ctx.Request.Header.SetUserAgent("badbot/1.0")
	hellPot(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusNotFound {
		t.Fatalf("status = %d, want 404", ctx.Response.StatusCode())
	}
}

func TestRobotsListsEveryPath(t *testing.T) {
	nop := zerolog.Nop()
	log = &nop
	config.Paths = []string{"wp-login.php", "wp-login"}

	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("/robots.txt")
	robotsTXT(&ctx)

	body := string(ctx.Response.Body())
	for _, p := range config.Paths {
		if !strings.Contains(body, "Disallow: /"+p+"\r\n") {
			t.Fatalf("robots.txt misses %q:\n%s", p, body)
		}
	}
}

package _test

import (
	"fmt"
	"mime"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_http/httpserver"
	"github.com/rskvp/qb-lib/qb_http/httpserver/httprewrite"
)

func TestServer(t *testing.T) {

	ct := mime.TypeByExtension(".js")
	fmt.Printf("ct: %s\n", ct)
	restart := false

	server := httpserver.NewHttpServer("./server", nil, nil)
	/** rewrite */
	server.Use(httprewrite.New(map[string]interface{}{
		"#IGNORE:/api/*": "",
		"#ROUTES:.":      "/index.html",
	}))

	errs := server.Configure(80, 443, "./cert/ssl.cert", "./cert/ssl.key", "./www4", false).
		All("/api/v1/*", handleAPI).
		Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	_ = qbc.Exec.Open("http://127.0.0.1")

	if restart {
		errs = server.Restart()
		if len(errs) > 0 {
			t.Error(errs)
			t.FailNow()
		}
	}

	time.Sleep(20 * time.Minute)
}

func handleAPI(ctx *fiber.Ctx) error {
	// http://localhost:9090/api/v1/sys/version
	_, _ = ctx.WriteString("{ \"response\":\"v0.1\"}")
	return nil
}

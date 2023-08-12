package _test

import (
	"fmt"
	"log"
	"mime"
	"net"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_http/httpserver/httprewrite"
)

func TestFiber(t *testing.T) {

	ct := mime.TypeByExtension(".js")
	fmt.Printf("ct: %s\n", ct)

	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	/***/
	app.Use(httprewrite.New(map[string]interface{}{
		"#IGNORE:/api/*": "",
		"#ROUTES:.":      "/index.html",
	}))

	app.Static("/", "./server/www4")

	go func() {
		log.Fatal(app.Listen(":80"))
	}()

	time.Sleep(1 * time.Second)

	_ = qbc.Exec.Open("http://127.0.0.1")

	time.Sleep(1 * time.Hour)
}

func TestFiber2(t *testing.T) {

	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	/** */
	app.Use(httprewrite.New(map[string]interface{}{
		"#IGNORE:/api/*": "",
		".":              "/index.html",
	}))

	app.Static("/", "./server/www4")

	go func() {
		ln, err := net.Listen("tcp", ":80")
		if nil == err {
			if err = app.Listener(ln); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	_ = qbc.Exec.Open("http://127.0.0.1")

	time.Sleep(1 * time.Hour)
}

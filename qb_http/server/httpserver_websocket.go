package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpWebsocket struct {

	//-- private --//
	app       *fiber.App
	cfgHost   *ConfigHost
	cfgRoutes []*httpServerConfigRouteWebsocket
	pool      map[string]*HttpWebsocketConn
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpWebsocket
//----------------------------------------------------------------------------------------------------------------------

func NewHttpWebsocket(app *fiber.App, cfgHost *ConfigHost, cfgRoutes []*httpServerConfigRouteWebsocket) *HttpWebsocket {

	instance := new(HttpWebsocket)
	instance.app = app
	instance.cfgHost = cfgHost
	instance.cfgRoutes = cfgRoutes

	instance.pool = make(map[string]*HttpWebsocketConn)

	return instance
}

func (instance *HttpWebsocket) Init() {
	app := instance.app
	routes := instance.cfgRoutes
	cfgHost := instance.cfgHost
	if nil != cfgHost.Websocket && len(routes) > 0 {
		settings := cfgHost.Websocket

		config := websocket.Config{}
		config.EnableCompression = settings.EnableCompression
		if settings.HandshakeTimeout > 0 {
			config.HandshakeTimeout = settings.HandshakeTimeout * time.Millisecond
		}
		if len(settings.Origins) > 0 {
			config.Origins = settings.Origins
		} else {
			config.Origins = []string{"*"}
		}
		if len(settings.Subprotocols) > 0 {
			config.Subprotocols = settings.Subprotocols
		}
		if settings.ReadBufferSize > 0 {
			config.ReadBufferSize = settings.ReadBufferSize
		}
		if settings.WriteBufferSize > 0 {
			config.WriteBufferSize = settings.WriteBufferSize
		}

		// handle upgrade
		app.Use(routes[0].Path, func(c *fiber.Ctx) error {
			// IsWebSocketUpgrade returns true if the client
			// requested upgrade to the WebSocket protocol.
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("allowed", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})

		// open websocket handlers
		for _, route := range routes {
			if len(route.Path) > 0 && nil != route.Handler {
				app.Get(route.Path, websocket.New(func(c *websocket.Conn) {
					if nil != route.Handler {
						ws := newConnection(c, instance.pool)
						route.Handler(ws)
						ws.Join() // lock waiting close
					}
				}, config))
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func newConnection(c *websocket.Conn, pool map[string]*HttpWebsocketConn) *HttpWebsocketConn {
	ws := NewHttpWebsocketConn(c, pool)
	return ws
}

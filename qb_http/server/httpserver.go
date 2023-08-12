package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_events"
)

// https://dev.to/koddr/go-fiber-by-examples-delving-into-built-in-functions-1p3k

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type CallbackError func(serverError *HttpServerError)
type CallbackLimitReached func(ctx *fiber.Ctx) error

type HttpServerError struct {
	Sender  interface{}
	Message string
	Error   error
	Context *fiber.Ctx
}

type IHttpHandler interface {
	Handle(ctx *fiber.Ctx) error
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpServer
//----------------------------------------------------------------------------------------------------------------------

type HttpServer struct {
	workspace string
	apps      []*fiber.App

	cfgServer         *ConfigServer
	cfgHosts          []*ConfigHost
	cfgStatic         []*ConfigStatic
	cfgCORS           *ConfigCORS
	cfgCompression    *ConfigCompression
	cfgLimiter        *ConfigLimiter
	cfgRoute          *httpServerConfigRoute
	cfgMiddleware     []*httpServerConfigRouteItem
	cfgRouteWebsocket []*httpServerConfigRouteWebsocket

	middlewares []fiber.Handler
	routers     []fiber.Router

	monitorFiles         []string
	monitor              *ServerMonitor
	callbackError        CallbackError
	callbackLimitReached CallbackLimitReached
	stopChan             chan bool
}

func NewHttpServer(workspace string, callbackError CallbackError, callbackLimit CallbackLimitReached) *HttpServer {
	instance := new(HttpServer)
	instance.workspace = qbc.Paths.Absolute(workspace)
	instance.monitorFiles = make([]string, 0)
	instance.apps = make([]*fiber.App, 0)
	instance.cfgServer = new(ConfigServer)
	instance.cfgServer.EnableRequestId = true
	instance.cfgServer.Prefork = false
	instance.cfgCORS = new(ConfigCORS)
	instance.cfgCORS.Enabled = true
	instance.cfgCompression = new(ConfigCompression)
	instance.cfgCompression.Enabled = false
	instance.cfgLimiter = new(ConfigLimiter)
	instance.cfgLimiter.Enabled = false
	instance.cfgStatic = make([]*ConfigStatic, 0)
	instance.cfgHosts = make([]*ConfigHost, 0)
	instance.cfgRoute = NewHttpServerConfigRoute()

	instance.stopChan = make(chan bool, 1)

	instance.callbackError = callbackError
	instance.callbackLimitReached = callbackLimit

	instance.middlewares = make([]fiber.Handler, 0)
	instance.routers = make([]fiber.Router, 0)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n f i g u r a t i o n
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) Configuration() *Config {
	response := new(Config)
	response.Server = instance.cfgServer
	response.Static = instance.cfgStatic
	response.Limiter = instance.cfgLimiter
	response.Hosts = instance.cfgHosts
	response.Cors = instance.cfgCORS
	response.Compression = instance.cfgCompression

	return response
}

func (instance *HttpServer) Configure(httpPort, httpsPort int, certPem, certKey, wwwRoot string, enableWebsocket bool) *HttpServer {
	// init hosts
	if httpPort > 0 || httpsPort > 0 {
		instance.cfgHosts = make([]*ConfigHost, 0) // reset
		if httpPort > 0 {
			http := new(ConfigHost)
			http.Address = fmt.Sprintf(":%v", httpPort)
			http.TLS = false
			if enableWebsocket {
				http.Websocket = new(ConfigHostWebsocket)
				http.Websocket.Enabled = true
			}
			instance.cfgHosts = append(instance.cfgHosts, http)
		}
		if httpsPort > 0 {
			https := new(ConfigHost)
			https.Address = fmt.Sprintf(":%v", httpsPort)
			https.TLS = true
			if len(certPem) == 0 {
				certPem = "./cert/ssl-cert.pem"
			}
			if len(certKey) == 0 {
				certKey = "./cert/ssl-cert.key"
			}
			https.SslCert = qbc.Paths.Concat(instance.workspace, certPem)
			https.SslKey = qbc.Paths.Concat(instance.workspace, certKey)
			if enableWebsocket {
				https.Websocket = new(ConfigHostWebsocket)
				https.Websocket.Enabled = true
			}
			instance.cfgHosts = append(instance.cfgHosts, https)
		}
		// STATIC
		if len(wwwRoot) > 0 {
			instance.cfgStatic = make([]*ConfigStatic, 0)
			static := new(ConfigStatic)
			static.Enabled = true
			static.Prefix = "/"
			static.Root = qbc.Paths.Concat(instance.workspace, wwwRoot)
			static.Compress = true
			instance.cfgStatic = append(instance.cfgStatic, static)
		}
	}
	return instance
}

func (instance *HttpServer) ConfigureFromFile(filename string) error {
	text, err := qbc.IO.ReadTextFromFile(filename)
	if nil != err {
		return err
	}
	return instance.ConfigureFromJson(text)
}

func (instance *HttpServer) ConfigureFromMap(settings map[string]interface{}) error {
	text := qbc.JSON.Stringify(settings)
	return instance.ConfigureFromJson(text)
}

func (instance *HttpServer) ConfigureFromJson(text string) error {
	var c *Config
	err := qbc.JSON.Read(text, &c)
	if nil == err && nil != c {
		if nil != c.Server {
			instance.cfgServer = c.Server
		}
		if nil != c.Hosts {
			instance.cfgHosts = c.Hosts
		}
		if nil != c.Static {
			instance.cfgStatic = c.Static
		}
		if nil != c.Compression {
			instance.cfgCompression = c.Compression
		}
		if nil != c.Limiter {
			instance.cfgLimiter = c.Limiter
		}
		if nil != c.Cors {
			instance.cfgCORS = c.Cors
		}
	}
	return err
}

func (instance *HttpServer) ConfigureServer(settings map[string]interface{}) {
	s := qbc.JSON.Stringify(settings)
	if len(s) > 0 {
		var c *ConfigServer
		err := qbc.JSON.Read(s, &c)
		if nil == err && nil != c {
			instance.cfgServer = c
		}
	}
}

func (instance *HttpServer) ConfigureCors(settings map[string]interface{}) {
	s := qbc.JSON.Stringify(settings)
	if len(s) > 0 {
		var c *ConfigCORS
		err := qbc.JSON.Read(s, &c)
		if nil == err && nil != c {
			instance.cfgCORS = c
		}
	}
}

func (instance *HttpServer) ConfigureStatic(settings ...map[string]interface{}) {
	instance.cfgStatic = make([]*ConfigStatic, 0)
	for _, setting := range settings {
		s := qbc.JSON.Stringify(setting)
		if len(s) > 0 {
			var c *ConfigStatic
			err := qbc.JSON.Read(s, &c)
			if nil == err {
				instance.cfgStatic = append(instance.cfgStatic, c)
			}
		}
	}
}

func (instance *HttpServer) ConfigureHosts(settings ...map[string]interface{}) {
	instance.cfgHosts = make([]*ConfigHost, 0)
	for _, setting := range settings {
		s := qbc.JSON.Stringify(setting)
		if len(s) > 0 {
			var c *ConfigHost
			err := qbc.JSON.Read(s, &c)
			if nil == err && nil != c {
				instance.cfgHosts = append(instance.cfgHosts, c)
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	m i d d l e w a r e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) Use(middleware fiber.Handler) {
	if nil != instance {
		if len(instance.apps) == 0 {
			instance.middlewares = append(instance.middlewares, middleware)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	r o u t i n g
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) Group(route string, callbacks ...func(ctx *fiber.Ctx) error) {
	if len(route) > 0 && len(callbacks) > 0 {
		instance.cfgRoute.Group(route, callbacks...)
	}
}

func (instance *HttpServer) HandleMiddleware(route string, handler IHttpHandler) *HttpServer {
	if nil != handler {
		return instance.Middleware(route, handler.Handle)
	}
	return instance
}

func (instance *HttpServer) Middleware(route string, callback func(ctx *fiber.Ctx) error) *HttpServer {
	if len(route) > 0 && nil != callback {
		item := new(httpServerConfigRouteItem)
		item.Path = route
		item.Handlers = append(item.Handlers, callback)

		if len(item.Handlers) > 0 {
			instance.cfgMiddleware = append(instance.cfgMiddleware, item)
		}
	}
	return instance
}

func (instance *HttpServer) HandleAll(route string, handler IHttpHandler) *HttpServer {
	if nil != handler {
		return instance.All(route, handler.Handle)
	}
	return instance
}

func (instance *HttpServer) All(route string, callbacks ...func(ctx *fiber.Ctx) error) *HttpServer {
	if len(route) > 0 && len(callbacks) > 0 {
		instance.cfgRoute.All(route, callbacks...)
	}
	return instance
}

func (instance *HttpServer) HandleGet(route string, handler IHttpHandler) *HttpServer {
	if nil != handler {
		return instance.Get(route, handler.Handle)
	}
	return instance
}
func (instance *HttpServer) Get(route string, callbacks ...func(ctx *fiber.Ctx) error) *HttpServer {
	if len(route) > 0 && len(callbacks) > 0 {
		instance.cfgRoute.Get(route, callbacks...)
	}
	return instance
}

func (instance *HttpServer) HandlePost(route string, handler IHttpHandler) *HttpServer {
	if nil != handler {
		return instance.Post(route, handler.Handle)
	}
	return instance
}

func (instance *HttpServer) Post(route string, callbacks ...func(ctx *fiber.Ctx) error) *HttpServer {
	if len(route) > 0 && len(callbacks) > 0 {
		instance.cfgRoute.Post(route, callbacks...)
	}
	return instance
}

func (instance *HttpServer) HandlePut(route string, handler IHttpHandler) *HttpServer {
	if nil != handler {
		return instance.Put(route, handler.Handle)
	}
	return instance
}

func (instance *HttpServer) Put(route string, callbacks ...func(ctx *fiber.Ctx) error) *HttpServer {
	if len(route) > 0 && len(callbacks) > 0 {
		instance.cfgRoute.Put(route, callbacks...)
	}
	return instance
}

func (instance *HttpServer) HandleDelete(route string, handler IHttpHandler) *HttpServer {
	if nil != handler {
		return instance.Delete(route, handler.Handle)
	}
	return instance
}

func (instance *HttpServer) Delete(route string, callbacks ...func(ctx *fiber.Ctx) error) *HttpServer {
	if len(route) > 0 && len(callbacks) > 0 {
		instance.cfgRoute.Delete(route, callbacks...)
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	w e b s o c k e t
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) Websocket(route string, callback func(conn *HttpWebsocketConn)) *HttpServer {
	if len(route) > 0 && nil != callback {

		item := new(httpServerConfigRouteWebsocket)
		item.Path = route
		item.Handler = callback

		instance.cfgRouteWebsocket = append(instance.cfgRouteWebsocket, item)
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	s t a r t
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) IsOpen() bool {
	if nil != instance {
		return nil != instance.stopChan && len(instance.apps) > 0
	}
	return false
}

func (instance *HttpServer) Restart() []error {
	response := make([]error, 0)
	if nil != instance && instance.IsOpen() {
		err := instance.Stop()
		if nil != err {
			response = append(response, err)
		} else {
			response = append(response, instance.Start()...)
		}
	}
	return response
}

func (instance *HttpServer) Start(settings ...map[string]interface{}) []error {
	errorList := make([]error, 0)
	instance.stopChan = make(chan bool, 1)
	if len(settings) > 0 {
		instance.ConfigureHosts(settings...)
	}
	instance.initWsHosts()
	for _, host := range instance.cfgHosts {
		err := instance.listen(host)
		if nil != err {
			errorList = append(errorList, err)
		}
	}
	instance.startSSLMonitor()

	// wait a while
	time.Sleep(1 * time.Second)

	return errorList
}

func (instance *HttpServer) Stop() (err error) {
	if nil != instance && len(instance.apps) > 0 {
		instance.stopSSLMonitor()
		for _, app := range instance.apps {
			appErr := app.Shutdown()
			if nil != appErr {
				err = appErr
			}
		}
		instance.stopChan <- true
		// reset stopChan
		instance.stopChan = nil
	}
	return
}

func (instance *HttpServer) Join() (err error) {
	if nil != instance && len(instance.apps) > 0 {
		<-instance.stopChan
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) initWsHosts() {
	for _, host := range instance.cfgHosts {
		if nil != host {
			if nil == host.Websocket {
				host.Websocket = new(ConfigHostWebsocket)
				host.Websocket.Enabled = false
			}
		}
	}
}

func (instance *HttpServer) handleServerError(c *fiber.Ctx, err error) error {
	if e := c.SendString(err.Error()); nil != e {
		return e
	}
	return c.SendStatus(500)
}

func (instance *HttpServer) handleLimitReached(c *fiber.Ctx) error {
	// request limit reached
	if nil != instance.callbackLimitReached {
		return instance.callbackLimitReached(c)
	}
	return c.SendStatus(fiber.StatusTooManyRequests)
}

func (instance *HttpServer) notifyError(message string, err error, ctx *fiber.Ctx) {
	go func() {
		if nil != instance.callbackError {
			instance.callbackError(&HttpServerError{
				Sender:  instance,
				Message: message,
				Context: ctx,
				Error:   err,
			})
		}
	}()
}

func (instance *HttpServer) listen(host *ConfigHost) error {
	app := instance.newApp(host)
	if nil != app {
		instance.apps = append(instance.apps, app)

		// get address and validate
		addr := host.Address
		if len(addr) > 0 {
			if strings.Index(addr, ":") == -1 {
				addr = ":" + addr
			}
		} else {
			return errors.New("empty_address")
		}

		// get protocol and validate
		network := "tcp" // tcp, tcp4

		var tlsConfig *tls.Config
		if host.TLS && len(host.SslKey) > 0 && len(host.SslCert) > 0 {
			pathCert := instance.absolutePath(host.SslCert)
			pathKey := instance.absolutePath(host.SslKey)
			cer, err := tls.LoadX509KeyPair(pathCert, pathKey)
			if err != nil {
				msg := qbc.Strings.Format("Error loading Certificates: '%s' '%s'", pathCert, pathKey)
				instance.notifyError(msg, err, nil)
				return qbc.Errors.Prefix(err, msg)
			}
			tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
		}

		if nil == tlsConfig {
			// STANDARD LISTENER
			ln, err := net.Listen(network, addr)
			if err != nil {
				msg := qbc.Strings.Format("Error creating listener: '%s'", addr)
				instance.notifyError(msg, err, nil)
				return qbc.Errors.Prefix(err, msg)
			} else {
				go instance.listener(app, ln, addr)
			}
		} else {
			// TLS LISTENER
			ln, err := tls.Listen(network, addr, tlsConfig)
			if err != nil {
				msg := qbc.Strings.Format("Error creating TLS listener: '%s'", addr)
				instance.notifyError(msg, err, nil)
				return qbc.Errors.Prefix(err, msg)
			} else {
				// Start server with https/ssl enabled on http://localhost:443
				go instance.listener(app, ln, addr)
			}
		}
	} else {
		// unable to create web application
		return errors.New("nil_application")
	}
	return nil
}

func (instance *HttpServer) listener(app *fiber.App, ln net.Listener, addr string) {
	if err := app.Listener(ln); err != nil {
		instance.notifyError(qbc.Strings.Format("Error Opening TLS channel: '%s'", addr), err, nil)
	}
}

func (instance *HttpServer) newApp(cfgHost *ConfigHost) *fiber.App {
	app := instance.createApp(instance.cfgServer)
	if nil != app {

		// RECOVER
		cfg := recover.Config{
			// Next defines a function to skip this middleware when returned true.
			Next: nil,
		}
		app.Use(recover.New(cfg))

		// REQUEST-ID
		if instance.cfgServer.EnableRequestId {
			app.Use(requestid.New())
		}

		// CORS
		initCORS(app, instance.cfgCORS)

		// compression
		initCompression(app, instance.cfgCompression)

		// limiter
		initLimiter(app, instance.cfgLimiter, instance.handleLimitReached)

		// prepare middlewares
		for _, middleware := range instance.middlewares {
			instance.cfgMiddleware = append(instance.cfgMiddleware, &httpServerConfigRouteItem{
				Path:     "",
				Handlers: []fiber.Handler{middleware},
			})
		}
		// Middleware
		if len(instance.cfgMiddleware) > 0 {
			instance.initMiddleware(app, instance.cfgMiddleware)
		}

		// Route
		if nil != instance.cfgRoute {
			initRoute(app, instance.cfgRoute, nil)
		}

		// websocket
		socket := NewHttpWebsocket(app, cfgHost, instance.cfgRouteWebsocket)
		socket.Init()

		// Static
		if len(instance.cfgStatic) > 0 {
			for _, static := range instance.cfgStatic {
				if static.Enabled && len(static.Root) > 0 {
					root := static.Root
					if !qbc.Paths.IsAbs(root) {
						root = qbc.Paths.Concat(instance.workspace, root)
					}
					cfgStatic := fiber.Static{
						Compress:  static.Compress,
						ByteRange: static.ByteRange,
						Browse:    static.Browse,
						Index:     static.Index,
						MaxAge:    static.MaxAge,
					}
					if static.CacheDurationSec > 0 {
						cfgStatic.CacheDuration = time.Duration(static.CacheDurationSec) * time.Second
					}
					// cfgStatic.ModifyResponse = instance.staticModifyResponse
					cfgStatic.Next = instance.staticNext
					app.Static(static.Prefix, root, cfgStatic)
				}
			}
		}
	}
	return app
}

func (instance *HttpServer) createApp(config *ConfigServer) *fiber.App {
	if nil != config {

		// creates app configuration
		appConfig := fiber.Config{
			ServerHeader:  config.ServerHeader,
			Prefork:       config.Prefork,
			CaseSensitive: config.CaseSensitive,
			StrictRouting: config.StrictRouting,
			Immutable:     config.Immutable,
		}
		if config.Concurrency > 0 {
			appConfig.Concurrency = config.Concurrency
		}
		if config.DisableKeepalive {
			appConfig.DisableKeepalive = config.DisableKeepalive
		}
		if config.DisableStartupMessage {
			appConfig.DisableStartupMessage = config.DisableStartupMessage
		}
		if config.BodyLimit > 0 {
			appConfig.BodyLimit = config.BodyLimit
		}
		if config.ReadTimeout > 0 {
			appConfig.ReadTimeout = config.ReadTimeout * time.Millisecond
		}
		if config.WriteTimeout > 0 {
			appConfig.WriteTimeout = config.WriteTimeout * time.Millisecond
		}
		if config.IdleTimeout > 0 {
			appConfig.IdleTimeout = config.IdleTimeout * time.Millisecond
		}

		// creates web application
		return fiber.New(appConfig)
	}
	return nil
}

func (instance *HttpServer) absolutePath(path string) string {
	if qbc.Paths.IsAbs(path) {
		return path
	}
	return qbc.Paths.Concat(instance.workspace, path)
}

func (instance *HttpServer) startSSLMonitor() {
	if nil != instance && nil == instance.monitor {
		instance.monitorFiles = make([]string, 0)
		for _, host := range instance.cfgHosts {
			if len(host.SslCert) > 0 && len(host.SslKey) > 0 {
				instance.monitorFiles = append(instance.monitorFiles, host.SslKey, host.SslCert)
			}
		}
		instance.monitor = NewMonitor(instance.monitorFiles)
		instance.monitor.OnFileChanged(instance.onSSLFileChanged)
		instance.monitor.Start()
	}
}

func (instance *HttpServer) stopSSLMonitor() {
	if nil != instance && nil != instance.monitor {
		instance.monitor.Stop()
		instance.monitor = nil
	}
}

func (instance *HttpServer) onSSLFileChanged(_ *qb_events.Event) {
	// fmt.Println("CHANGED FILE....")
	instance.stopSSLMonitor()
	// waite a while to ensure certificates are well stored
	time.Sleep(10 * time.Second)
	// restart server
	instance.Restart()
}

func (instance *HttpServer) staticModifyResponse(ctx *fiber.Ctx) (err error) {
	return
}

func (instance *HttpServer) staticNext(c *fiber.Ctx) (handled bool) {
	if nil != instance && nil != c {
		//empty
	}
	return
}

func (instance *HttpServer) initMiddleware(app *fiber.App, items []*httpServerConfigRouteItem) {
	for _, item := range items {
		path := item.Path
		if len(path) == 0 {
			instance.useMiddleware(app, item.Handlers[0])
		} else {
			instance.useMiddleware(app, path, item.Handlers[0])
		}
	}
}

func (instance *HttpServer) useMiddleware(app *fiber.App, args ...interface{}) {
	if nil != instance && nil != app {
		router := app.Use(args...)
		instance.routers = append(instance.routers, router)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func initCORS(app *fiber.App, corsCfg *ConfigCORS) {
	if nil != corsCfg && corsCfg.Enabled {
		config := cors.Config{}
		if corsCfg.MaxAge > 0 {
			config.MaxAge = corsCfg.MaxAge
		}
		if corsCfg.AllowCredentials {
			config.AllowCredentials = true
		}
		if len(corsCfg.AllowMethods) > 0 {
			config.AllowMethods = strings.Join(corsCfg.AllowMethods, ",")
		}
		if len(corsCfg.AllowOrigins) > 0 {
			config.AllowOrigins = strings.Join(corsCfg.AllowOrigins, ",")
		}
		if len(corsCfg.ExposeHeaders) > 0 {
			config.ExposeHeaders = strings.Join(corsCfg.ExposeHeaders, ",")
		}
		app.Use(cors.New(config))
	}
}

func initCompression(app *fiber.App, cfg *ConfigCompression) {
	if nil != cfg && cfg.Enabled {
		config := compress.Config{}
		config.Level = compress.Level(cfg.Level)

		app.Use(compress.New(config))
	}
}

func initLimiter(app *fiber.App, cfg *ConfigLimiter, handler func(ctx *fiber.Ctx) error) {
	if nil != cfg && cfg.Enabled {
		config := limiter.Config{}
		if cfg.Max > 0 {
			config.Max = cfg.Max
		}
		if cfg.Duration > 0 {
			config.Expiration = cfg.Duration
		}

		config.LimitReached = handler

		app.Use(limiter.New(config))
	}
}

func initRoute(app *fiber.App, route *httpServerConfigRoute, parent fiber.Router) {
	for k, i := range route.Data {
		initRouteItem(app, k, i, parent)
	}
}

func initGroup(app *fiber.App, group *httpServerConfigGroup, parent fiber.Router) {
	var g fiber.Router
	if nil == parent {
		g = app.Group(group.Path, group.Handlers...)
	} else {
		g = parent.Group(group.Path, group.Handlers...)
	}
	if nil != g && len(group.Children) > 0 {
		for _, c := range group.Children {
			if cc, b := c.(*httpServerConfigGroup); b {
				// children is a group
				initGroup(app, cc, g)
			} else if cc, b := c.(*httpServerConfigRoute); b {
				// children is route
				initRoute(app, cc, g)
			}
		}
	}
}

func initRouteItem(app *fiber.App, key string, item interface{}, parent fiber.Router) {
	method := qbc.Strings.SplitAndGetAt(key, "_", 0)
	switch method {
	case "GROUP":
		v := item.(*httpServerConfigGroup)
		initGroup(app, v, parent)
	case "ALL":
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.All(v.Path, v.Handlers...)
		} else {
			parent.All(v.Path, v.Handlers...)
		}
	case fiber.MethodGet:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Get(v.Path, v.Handlers...)
		} else {
			parent.Get(v.Path, v.Handlers...)
		}
	case fiber.MethodPost:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Post(v.Path, v.Handlers...)
		} else {
			parent.Post(v.Path, v.Handlers...)
		}
	case fiber.MethodOptions:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Options(v.Path, v.Handlers...)
		} else {
			parent.Options(v.Path, v.Handlers...)
		}
	case fiber.MethodPut:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Put(v.Path, v.Handlers...)
		} else {
			parent.Put(v.Path, v.Handlers...)
		}
	case fiber.MethodHead:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Head(v.Path, v.Handlers...)
		} else {
			parent.Head(v.Path, v.Handlers...)
		}
	case fiber.MethodPatch:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Patch(v.Path, v.Handlers...)
		} else {
			parent.Patch(v.Path, v.Handlers...)
		}
	case fiber.MethodDelete:
		v := item.(*httpServerConfigRouteItem)
		if nil == parent {
			app.Delete(v.Path, v.Handlers...)
		} else {
			parent.Delete(v.Path, v.Handlers...)
		}
	}
}

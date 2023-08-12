package http

import (
	"encoding/json"
	"sync"

	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_http/server"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpServerCallback struct {
	path     string
	runtime  *goja.Runtime
	this     goja.Value
	callback goja.Callable
	mux  sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewServerCallback(runtime *goja.Runtime, this goja.Value, path string, callback goja.Callable) *JsHttpServerCallback {
	instance := new(JsHttpServerCallback)
	instance.runtime = runtime
	instance.this = this
	instance.path = path
	instance.callback = callback

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpServerCallback) HandleRoute(ctx *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			_ = qbc.Strings.Format("[panic] jsNext.next: %s", r)
			// TODO: implement logger
			// fmt.Println(message)
		}
	}()
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		if nil != instance.callback && nil != ctx {
			req := WrapRequest(instance.runtime, nil, ctx)
			res := WrapResponse(instance.runtime, nil, ctx)
			next := WrapNext(instance.runtime, ctx)
			_, err = instance.handle(req.Value(), res.Value(), next.NextFunc())
		}
	}
	if nil != err {
		//panic(instance.runtime.NewTypeError(err.Error()))
	}
	return err
}

func (instance *JsHttpServerCallback) handle(req goja.Value, res goja.Value, next goja.Value) (goja.Value, error) {
	if nil != instance && nil != instance.runtime && nil != instance.callback && nil != req && nil != res {
		return instance.callback(instance.this, req, res, next)
	}
	return goja.Undefined(), nil
}

func (instance *JsHttpServerCallback) HandleWs(conn *server.HttpWebsocketConn) {
	if nil != instance && nil != instance.callback {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		conn.OnMessage(instance.onWsMessage)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpServerCallback) onWsMessage(payload *server.HttpWebsocketEventPayload) {
	if nil != instance && nil != instance.runtime && nil != instance.callback && nil != payload {
		ws := payload.Websocket
		if nil != ws && ws.IsAlive() && len(payload.Message.Data) > 0 {
			data := payload.Message.Data
			var m map[string]interface{}
			err := json.Unmarshal(data, &m)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
			req := WrapRequest(instance.runtime, data, nil)
			res := WrapResponse(instance.runtime, nil, ws)
			_, err = instance.callback(instance.this, req.Value(), res.Value())
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
		}
	}
}

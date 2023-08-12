package nio

import (
	"errors"
	"sync"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_nio"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsServer struct {
	runtime  *goja.Runtime
	nio      *qb_nio.NioServer
	object   *goja.Object
	handlers map[string]goja.Callable
	mux      sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	JsServer
//----------------------------------------------------------------------------------------------------------------------

func WrapServer(runtime *goja.Runtime, port int) goja.Value {
	instance := new(JsServer)
	instance.runtime = runtime
	instance.nio = qbc.NIO.NewServer(port)
	instance.handlers = make(map[string]goja.Callable)

	instance.object = instance.runtime.NewObject()
	instance.export()

	// register as closable
	commons.AddClosableObject(instance.runtime, instance.object)

	return instance.object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsServer) open(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		err := instance.nio.Open()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		// add message handler
		instance.nio.OnMessage(instance.onMessage)

	}
	return goja.Undefined()
}

func (instance *JsServer) close(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		instance.object = nil
		err := instance.nio.Close()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsServer) join(_ goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		// NO MORE ALLOWED
		// instance.nio.Join()
	}
	return goja.Undefined()
}

func (instance *JsServer) isOpen(call goja.FunctionCall) goja.Value {
	if nil != instance.nio && nil != instance.object {
		return instance.runtime.ToValue(instance.nio.IsOpen())
	}
	return instance.runtime.ToValue(false)
}

func (instance *JsServer) count(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		return instance.runtime.ToValue(instance.nio.ClientsCount())
	}
	return instance.runtime.ToValue(0)
}

func (instance *JsServer) clients(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		return instance.runtime.ToValue(instance.nio.ClientsId())
	}
	return instance.runtime.ToValue(0)
}

func (instance *JsServer) listen(call goja.FunctionCall) goja.Value {
	if nil != instance.nio && nil != instance.handlers {
		callback := commons.GetCallbackIfAny(call)
		if nil != callback {
			route := "*"
			if len(call.Arguments) == 2 {
				route = call.Argument(0).String()
			}
			instance.handlers[route] = callback
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsServer) onMessage(message *qb_nio.NioMessage) interface{} {
	if nil != instance && nil != instance.object && nil != instance.handlers {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		var m ModuleNioMessage
		err := qbc.JSON.Read(message.Body, &m)
		if nil != err {
			return err
		}
		if len(m.Name) == 0 {
			return errors.New("missing_route")
		}
		if v, b := instance.handlers[m.Name]; b {
			resp, err := instance.invoke(v, m.Name, m.Params)
			if nil != err {
				return err
			}
			return resp.Export()
		} else if v, b := instance.handlers["*"]; b {
			resp, err := instance.invoke(v, m.Name, m.Params)
			if nil != err {
				return err
			}
			if nil != resp {
				return resp.Export()
			}
		}
	}
	return false
}

func (instance *JsServer) invoke(callback goja.Callable, name string, params []interface{}) (goja.Value, error) {
	if nil != instance {
		defer func() {
			if r := recover(); r != nil {
				// recovered from panic

			}
		}()

		value, err := callback(instance.object, instance.runtime.ToValue(name), instance.runtime.ToValue(params))

		return value, err
	}
	return goja.Undefined(), nil
}

func (instance *JsServer) export() {
	o := instance.object
	_ = o.Set("open", instance.open)
	_ = o.Set("close", instance.close)
	_ = o.Set("isOpen", instance.isOpen)
	_ = o.Set("clients", instance.clients)
	_ = o.Set("count", instance.count)
	_ = o.Set("listen", instance.listen)
	_ = o.Set("join", instance.join)

}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

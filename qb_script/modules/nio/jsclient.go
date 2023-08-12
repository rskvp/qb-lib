package nio

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_nio"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsClient struct {
	runtime *goja.Runtime
	nio     *qb_nio.NioClient
}

//----------------------------------------------------------------------------------------------------------------------
//	JsClient
//----------------------------------------------------------------------------------------------------------------------

func WrapClient(runtime *goja.Runtime, host string, port int) goja.Value {
	instance := new(JsClient)
	instance.runtime = runtime
	instance.nio = qbc.NIO.NewClient(host, port)
	instance.nio.EnablePing = false // ping disabled (avoid continuous connect/disconnect)

	object := instance.runtime.NewObject()
	instance.export(object)

	return object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsClient) secure(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		var err error
		val := call.Argument(0).ToBoolean()
		opened := instance.nio.IsOpen()
		if opened{
			err = instance.nio.Close()
		}
		instance.nio.Secure = val
		if opened{
			err = instance.nio.Open()
		}
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsClient) enablePing(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		var err error
		val := call.Argument(0).ToBoolean()
		opened := instance.nio.IsOpen()
		if opened{
			err = instance.nio.Close()
		}
		instance.nio.EnablePing = val
		if opened{
			err = instance.nio.Open()
		}
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsClient) open(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		err := instance.nio.Open()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsClient) close(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		err := instance.nio.Close()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsClient) send(call goja.FunctionCall) goja.Value {
	if nil != instance.nio {
		if instance.nio.IsOpen() {
			if len(call.Arguments) > 0 {
				// firs parameter should be a command name, a "route"
				command := call.Argument(0).String()
				if len(command) > 0 {
					// prepare message to send
					message := &ModuleNioMessage{
						Name: command,
						Params: commons.ToArray(call.Arguments[1:]),
					}
					resp, err := instance.nio.Send(message)
					if nil != err {
						panic(instance.runtime.NewTypeError(err.Error()))
					}
					return instance.runtime.ToValue(qbc.Convert.ToString(resp.Body))
				} else {
					panic(instance.runtime.NewTypeError("missing_command_name"))
				}
			}
		} else {
			panic(instance.runtime.NewTypeError("client_disconnected"))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsClient) export(o *goja.Object) {
	_ = o.Set("secure", instance.secure)
	_ = o.Set("enablePing", instance.enablePing)
	_ = o.Set("open", instance.open)
	_ = o.Set("close", instance.close)
	_ = o.Set("send", instance.send)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

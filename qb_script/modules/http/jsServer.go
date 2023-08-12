package http

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_http/server"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpServer struct {
	runtime *goja.Runtime
	server  *server.HttpServer
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpClient
//----------------------------------------------------------------------------------------------------------------------

func WrapServer(runtime *goja.Runtime) goja.Value {
	instance := new(JsHttpServer)
	instance.runtime = runtime
	instance.server = server.NewHttpServer(commons.GetRtRoot(runtime), instance.handleServerError, nil)

	object := instance.runtime.NewObject()
	instance.export(object)

	// add closable: all closable objects must be exposed to avoid
	commons.AddClosableObject(instance.runtime, object)

	return object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d  -  c o n f i g u r a t i o n
//----------------------------------------------------------------------------------------------------------------------

// settings set qb_sms_engine configuration
func (instance *JsHttpServer) settings(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) > 0 {
			settings := commons.GetExport(call, 0)
			if nil != settings {
				if v, b := settings.(map[string]interface{}); b {
					instance.server.ConfigureFromMap(v)
				} else {
					panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
				}
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

// static set web server static configuration
func (instance *JsHttpServer) static(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) > 0 {
			settings := commons.ToArrayOfMap(commons.GetExport(call, 0))
			if len(settings) > 0 {
				instance.server.ConfigureStatic(settings...)
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

// settings set qb_sms_engine configuration
func (instance *JsHttpServer) cors(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) > 0 {
			settings := commons.GetExport(call, 0)
			if nil != settings {
				if v, b := settings.(map[string]interface{}); b {
					instance.server.ConfigureCors(v)
				} else {
					panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
				}
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

// listen start webserver
func (instance *JsHttpServer) listen(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) > 0 {
			settings := commons.ToArrayOfMap(commons.GetExport(call, 0))
			if len(settings) > 0 {
				errors := instance.server.Start(settings...)
				if len(errors) > 0 {
					panic(instance.runtime.NewTypeError(errors[0].Error()))
				}
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d  -  a c t i o n s
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpServer) join(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		// NO MORE ALLOWED
		/*
			err := instance.server.Join()
			if nil!=err{
				panic(instance.runtime.NewTypeError(err.Error()))
			}*/
	}
	return goja.Undefined()
}

func (instance *JsHttpServer) close(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		err := instance.server.Stop()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d  -  r o u t i n g
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpServer) get(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) == 2 {
			path := commons.GetString(call, 0)
			callback := commons.GetCallbackIfAny(call)
			if len(path) > 0 && nil != callback {
				handler := NewServerCallback(instance.runtime, call.This, path, callback)
				instance.server.Get(path, handler.HandleRoute)
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpServer) post(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) == 2 {
			path := commons.GetString(call, 0)
			callback := commons.GetCallbackIfAny(call)
			if len(path) > 0 && nil != callback {
				handler := NewServerCallback(instance.runtime, call.This, path, callback)
				instance.server.Post(path, handler.HandleRoute)
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpServer) all(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) == 2 {
			path := commons.GetString(call, 0)
			callback := commons.GetCallbackIfAny(call)
			if len(path) > 0 && nil != callback {
				handler := NewServerCallback(instance.runtime, call.This, path, callback)
				instance.server.All(path, handler.HandleRoute)
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d  -  w e b s o c k e t
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpServer) websocket(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) == 2 {
			path := commons.GetString(call, 0)
			callback := commons.GetCallbackIfAny(call)
			if len(path) > 0 && nil != callback {
				handler := NewServerCallback(instance.runtime, call.This, path, callback)
				instance.server.Websocket(path, handler.HandleWs)
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d  -  m i d d l e w a r e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpServer) use(call goja.FunctionCall) goja.Value {
	if nil != instance.server {
		if len(call.Arguments) == 2 {
			path := commons.GetString(call, 0)
			callback := commons.GetCallbackIfAny(call)
			if len(path) > 0 && nil != callback {
				handler := NewServerCallback(instance.runtime, call.This, path, callback)
				instance.server.Middleware(path, handler.HandleRoute)
			} else {
				panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// handleServerError receive async errors from app
func (instance *JsHttpServer) handleServerError(serverError *server.HttpServerError) {
	err := serverError.Error
	if nil != err {
		// panic(instance.runtime.NewTypeError(err.Error()))
	}
}

func (instance *JsHttpServer) export(o *goja.Object) {
	_ = o.Set("settings", instance.settings)
	_ = o.Set("static", instance.static)
	_ = o.Set("cors", instance.cors)
	_ = o.Set("join", instance.join)
	_ = o.Set("listen", instance.listen)
	_ = o.Set("close", instance.close)
	_ = o.Set("all", instance.all)
	_ = o.Set("get", instance.get)
	_ = o.Set("post", instance.post)
	_ = o.Set("websocket", instance.websocket)
	_ = o.Set("use", instance.use)

}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

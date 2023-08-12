package http

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_http/utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpClientResponse struct {
	runtime *goja.Runtime
	object  *goja.Object
	data    *utils.ResponseData
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpResponse
//----------------------------------------------------------------------------------------------------------------------

func WrapClientResponse(runtime *goja.Runtime, data *utils.ResponseData) goja.Value {
	instance := new(JsHttpClientResponse)
	instance.runtime = runtime
	instance.data = data

	instance.object = instance.runtime.NewObject()
	instance.export()

	return instance.value()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpClientResponse) length(_ goja.FunctionCall) goja.Value {
	if nil != instance.data {
		return instance.runtime.ToValue(len(instance.data.Body))
	}
	return instance.runtime.ToValue(0)
}

func (instance *JsHttpClientResponse) bytes(_ goja.FunctionCall) goja.Value {
	if nil != instance.data {
		return instance.runtime.ToValue(instance.data.Body)
	}
	return instance.runtime.ToValue([]byte{})
}

func (instance *JsHttpClientResponse) text(_ goja.FunctionCall) goja.Value {
	if nil != instance.data {
		return instance.runtime.ToValue(string(instance.data.Body))
	}
	return instance.runtime.ToValue("")
}

func (instance *JsHttpClientResponse) json(_ goja.FunctionCall) goja.Value {
	if nil != instance.data {
		var m map[string]interface{}
		err := qbc.JSON.Read(instance.data.Body, &m)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(m)
	}
	return instance.runtime.ToValue("")
}

func (instance *JsHttpClientResponse) status(_ goja.FunctionCall) goja.Value {
	if nil != instance && nil != instance.data {
		return instance.runtime.ToValue(instance.data.StatusCode)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpClientResponse) export() {
	o := instance.object
	if nil != o {
		// methods
		_ = o.Set("length", instance.length)
		_ = o.Set("bytes", instance.bytes)
		_ = o.Set("text", instance.text)
		_ = o.Set("json", instance.json)
		_ = o.Set("status", instance.status)
		// properties
		if nil != instance.data {
			_ = o.Set("size", len(instance.data.Body))
			_ = o.Set("statusCode", instance.data.StatusCode)
			_ = o.Set("header", instance.runtime.ToValue(instance.data.Header))
			_ = o.Set("body", instance.runtime.ToValue(string(instance.data.Body)))
		}
	}
}

func (instance *JsHttpClientResponse) value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

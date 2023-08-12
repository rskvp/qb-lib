package http

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_http/client"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpClient struct {
	runtime *goja.Runtime
	client  *client.HttpClient
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpClient
//----------------------------------------------------------------------------------------------------------------------

func WrapClient(runtime *goja.Runtime) goja.Value {
	instance := new(JsHttpClient)
	instance.runtime = runtime
	instance.client = client.NewHttpClient()

	object := instance.runtime.NewObject()
	instance.export(object)

	return object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpClient) addHeader(call goja.FunctionCall) goja.Value {
	if nil != instance.client {
		if len(call.Arguments) == 2 {
			key := commons.GetString(call, 0)
			value := commons.GetString(call, 1)

			instance.client.AddHeader(key, value)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpClient) removeHeader(call goja.FunctionCall) goja.Value {
	if nil != instance.client {
		if len(call.Arguments) == 2 {
			key := commons.GetString(call, 0)

			instance.client.RemoveHeader(key)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpClient) get(call goja.FunctionCall) goja.Value {
	if nil != instance.client {
		if len(call.Arguments) > 0 {
			url := commons.GetString(call, 0)
			// get data
			data, err := instance.client.Get(url)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}

			return WrapClientResponse(instance.runtime, data)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpClient) post(call goja.FunctionCall) goja.Value {
	if nil != instance.client {
		reqUrl, reqBody := getRequestParameters(call)
		if len(reqUrl)==0{
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		// post data
		data, err := instance.client.Post(reqUrl, reqBody)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		return WrapClientResponse(instance.runtime, data)
	}
	return goja.Undefined()
}

func (instance *JsHttpClient) put(call goja.FunctionCall) goja.Value {
	if nil != instance.client {
		reqUrl, reqBody := getRequestParameters(call)
		if len(reqUrl)==0{
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		// put data
		data, err := instance.client.Put(reqUrl, reqBody)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		return WrapClientResponse(instance.runtime, data)
	}
	return goja.Undefined()
}

func (instance *JsHttpClient) delete(call goja.FunctionCall) goja.Value {
	if nil != instance.client {
		reqUrl, reqBody := getRequestParameters(call)
		if len(reqUrl)==0{
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		// delete data
		data, err := instance.client.Delete(reqUrl, reqBody)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		return WrapClientResponse(instance.runtime, data)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpClient) export(o *goja.Object) {
	_ = o.Set("addHeader", instance.addHeader)
	_ = o.Set("removeHeader", instance.removeHeader)
	_ = o.Set("get", instance.get)
	_ = o.Set("post", instance.post)
	_ = o.Set("put", instance.put)
	_ = o.Set("delete", instance.delete)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func getRequestParameters(call goja.FunctionCall)(reqUrl string, reqBody interface{}){
	switch len(call.Arguments) {
	case 1:
		reqUrl = commons.GetString(call, 0)
		reqBody = nil
	case 2:
		reqUrl = commons.GetString(call, 0)
		reqBody = commons.GetExport(call, 1)
	default:

	}
	return
}
package http

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpRequest struct {
	runtime     *goja.Runtime
	object      *goja.Object
	data        bytes.Buffer
	ctx         *fiber.Ctx
	params      map[string]interface{}
	query       map[string]interface{}
	baseUrl     string
	originalUrl string
	path        string
	secure      bool
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpClient
//----------------------------------------------------------------------------------------------------------------------

func WrapRequest(runtime *goja.Runtime, data []byte, ctx *fiber.Ctx) *JsHttpRequest {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			message := qbc.Strings.Format("[panic] WrapRequest -> \"%s\"", r)
			// TODO: implement logger
			fmt.Println(message)
		}
	}()
	instance := new(JsHttpRequest)
	instance.runtime = runtime
	instance.ctx = ctx
	instance.params = make(map[string]interface{})
	instance.query = make(map[string]interface{})
	if nil != data {
		instance.data.Write(data)
	}

	if nil != instance.ctx {
		// HTTP CONTEXT
		instance.baseUrl = instance.ctx.BaseURL()
		instance.originalUrl = instance.ctx.OriginalURL()
		instance.path = instance.ctx.Path()
		instance.secure = instance.ctx.Secure()
		// body
		if body := ctx.Body(); len(body) > 0 {
			_, _ = instance.data.Write(body)
		}
		instance.params = instance.getParams(ctx) // never null

		// url query
		if path := instance.ctx.OriginalURL(); len(path) > 0 {
			uri, err := url.Parse(path)
			if nil == err {
				query := uri.Query()
				if nil != query && len(query) > 0 {
					for k, v := range query {
						if len(v) == 1 {
							instance.query[k] = v[0]
						} else {
							instance.query[k] = v
						}
					}
				}
			}
		}
	}

	if instance.data.Len() > 0 && len(instance.params) == 0 {
		if qbc.Regex.IsValidJsonObject(instance.data.String()) {
			_ = qbc.JSON.Read(instance.data.String(), &instance.params)
		}
	}

	instance.object = instance.runtime.NewObject()
	instance.export()

	return instance
}

func (instance *JsHttpRequest) Value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpRequest) length(_ goja.FunctionCall) goja.Value {
	if instance.data.Len() > 0 {
		return instance.runtime.ToValue(instance.data.Len())
	}
	return instance.runtime.ToValue(0)
}

func (instance *JsHttpRequest) bytes(_ goja.FunctionCall) goja.Value {
	if instance.data.Len() > 0 {
		return instance.runtime.ToValue(instance.data.Bytes())
	}
	return instance.runtime.ToValue([]byte{})
}

func (instance *JsHttpRequest) text(_ goja.FunctionCall) goja.Value {
	if instance.data.Len() > 0 {
		return instance.runtime.ToValue(instance.data.String())
	}
	return instance.runtime.ToValue("")
}

func (instance *JsHttpRequest) param(call goja.FunctionCall) goja.Value {
	if nil != instance.params {
		name := commons.GetString(call, 0)
		if len(name) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if v, b := instance.params[name]; b {
			return instance.runtime.ToValue(v)
		} else {
			if nil != instance.ctx && nil != instance.ctx.Route() {
				return instance.runtime.ToValue(instance.ctx.Params(name, ""))
			}
		}
	}
	return goja.Undefined()
}

// get Get data from header
func (instance *JsHttpRequest) get(call goja.FunctionCall) goja.Value {
	if nil != instance.ctx {
		headerName := commons.GetString(call, 0)
		if len(headerName) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		v := instance.ctx.Get(headerName, "")
		return instance.runtime.ToValue(v)
	}
	return goja.Undefined()
}

// getAuth retrieve Authorization header
func (instance *JsHttpRequest) getAuth(_ goja.FunctionCall) goja.Value {
	if nil != instance.ctx {
		headerName := "Authorization"
		v := instance.ctx.Get(headerName, "")
		tokens := qbc.Strings.Split(v, " ")
		response := map[string]string{
			"authorization": v,
		}
		if len(tokens) == 2 {
			mode := tokens[0]
			response["mode"] = mode
			switch mode {
			case "Bearer":
				response["token"] = tokens[1]
			case "Basic":
				response["username"] = ""
				response["password"] = ""
				data, _ := qbc.Coding.DecodeBase64(tokens[1])
				if len(data) > 0 {
					t := strings.Split(string(data), ":")
					if len(t) == 2 {
						response["username"] = t[0]
						response["password"] = t[1]
					}
				}
			}
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

// multipart return a
func (instance *JsHttpRequest) multipart(_ goja.FunctionCall) goja.Value {
	if nil != instance && nil != instance.ctx {
		form, err := instance.ctx.MultipartForm() // => *multipart.Form
		if err != nil {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		if nil != form {
			// create response object
			fileArray := make([]goja.Value, 0)
			response := map[string]interface{}{
				"files": fileArray,
				"data":  form.Value,
			}
			for key, files := range form.File {
				// Loop through files:
				for _, file := range files {
					fileArray = append(fileArray, WrapFile(instance.runtime, key, instance.ctx, file).Value())
				}
			}
			response["files"] = fileArray
			return instance.runtime.ToValue(response)
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
func (instance *JsHttpRequest) getParams(ctx *fiber.Ctx) map[string]interface{} {
	response := map[string]interface{}{}
	if nil != ctx {
		// try add form params
		if form, err := instance.ctx.MultipartForm(); nil == err && nil != form && nil != form.Value {
			for k, v := range form.Value {
				response[k] = v
			}
		}
		// try add body params
		if body := ctx.Body(); len(body) > 0 {
			if params := qbc.Convert.ToMap(body); nil != params {
				for k, v := range params {
					response[k] = v
				}
			}
		}
	}
	return response
}

func (instance *JsHttpRequest) export() {
	o := instance.object
	// methods
	_ = o.Set("length", instance.length)
	_ = o.Set("bytes", instance.bytes)
	_ = o.Set("text", instance.text)
	_ = o.Set("param", instance.param)
	_ = o.Set("get", instance.get)
	_ = o.Set("getAuth", instance.getAuth)
	_ = o.Set("multipart", instance.multipart)
	// properties
	_ = o.Set("size", instance.data.Len())
	_ = o.Set("originalUrl", instance.originalUrl)
	_ = o.Set("baseUrl", instance.baseUrl)
	_ = o.Set("path", instance.path)
	_ = o.Set("secure", instance.secure)
	if nil != instance.params {
		_ = o.Set("params", instance.params)
	}
	if nil != instance.query {
		_ = o.Set("query", instance.query)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

package http

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_http/server"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpResponse struct {
	runtime       *goja.Runtime
	object        *goja.Object
	data          bytes.Buffer
	ctxHttp       *fiber.Ctx
	ctxWs         *server.HttpWebsocketConn
	statusCode    int
	statusMessage string
	header        map[string]string
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpResponse
//----------------------------------------------------------------------------------------------------------------------

func WrapResponse(runtime *goja.Runtime, data []byte, ctx interface{}) *JsHttpResponse {
	instance := new(JsHttpResponse)
	instance.runtime = runtime
	if v, b := ctx.(*fiber.Ctx); b {
		instance.ctxHttp = v
	}
	if v, b := ctx.(*server.HttpWebsocketConn); b {
		instance.ctxWs = v
	}
	if nil != data {
		instance.data.Write(data)
	}

	instance.object = instance.runtime.NewObject()
	instance.export()

	return instance
}

func (instance *JsHttpResponse) Value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpResponse) length(_ goja.FunctionCall) goja.Value {
	if instance.data.Len() > 0 {
		return instance.runtime.ToValue(instance.data.Len())
	}
	return instance.runtime.ToValue(0)
}

func (instance *JsHttpResponse) bytes(_ goja.FunctionCall) goja.Value {
	if instance.data.Len() > 0 {
		return instance.runtime.ToValue(instance.data.Bytes())
	}
	return instance.runtime.ToValue([]byte{})
}

func (instance *JsHttpResponse) send(call goja.FunctionCall) goja.Value {
	if nil != instance.ctxHttp {
		var data interface{}
		var contentType string
		switch len(call.Arguments) {
		case 1:
			data = commons.GetExport(call, 0)
			contentType = fiber.MIMETextPlain
		case 2:
			data = commons.GetExport(call, 0)
			contentType = commons.GetString(call, 1)
		default:
			data = ""
			contentType = fiber.MIMETextPlain
		}

		if nil != data {
			// HTTP
			if nil != instance.ctxHttp {
				err := writeHttp(instance.ctxHttp, contentType, data)
				if nil != err {
					panic(instance.runtime.NewTypeError(err.Error()))
				}
			}
			// WEBSOCKET
			if nil != instance.ctxWs {
				if v, b := data.(string); b {
					instance.ctxWs.SendData([]byte(v))
				} else if v, b := data.([]byte); b {
					instance.ctxWs.SendData(v)
				} else if v, b := data.([]uint8); b {
					instance.ctxWs.SendData(v)
				}
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpResponse) text(call goja.FunctionCall) goja.Value {
	if nil != instance.ctxHttp {
		var data interface{}
		var contentType string
		switch len(call.Arguments) {
		case 1:
			data = commons.GetExport(call, 0)
			contentType = fiber.MIMETextPlain
		default:
			data = ""
			contentType = fiber.MIMETextPlain
		}

		if nil != data {
			// HTTP
			if nil != instance.ctxHttp {
				err := writeHttp(instance.ctxHttp, contentType, data)
				if nil != err {
					panic(instance.runtime.NewTypeError(err.Error()))
				}
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return instance.runtime.ToValue("")
}

func (instance *JsHttpResponse) json(call goja.FunctionCall) goja.Value {
	if nil != instance {
		value := commons.GetExport(call, 0)
		if nil != value {
			if nil != instance.ctxHttp {
				err := instance.ctxHttp.JSON(value)
				if nil != err {
					panic(instance.runtime.NewTypeError(err.Error()))
				}
			}
			if nil != instance.ctxWs {
				text := qbc.Convert.ToString(value)
				instance.ctxWs.SendData([]byte(text))
			}
		}
	}
	return goja.Undefined()
}

func (instance *JsHttpResponse) html(call goja.FunctionCall) goja.Value {
	if nil != instance.ctxHttp {
		var data interface{}
		var contentType string
		switch len(call.Arguments) {
		case 1:
			data = commons.GetExport(call, 0)
			contentType = fiber.MIMETextHTML
		default:
			data = ""
			contentType = fiber.MIMETextHTML
		}

		if nil != data {
			// HTTP
			if nil != instance.ctxHttp {
				err := writeHttp(instance.ctxHttp, contentType, data)
				if nil != err {
					panic(instance.runtime.NewTypeError(err.Error()))
				}
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return instance.runtime.ToValue("")
}

func (instance *JsHttpResponse) image(call goja.FunctionCall) goja.Value {
	if nil != instance.ctxHttp {
		var data interface{}
		var contentType string
		switch len(call.Arguments) {
		case 1:
			data = commons.GetExport(call, 0)
			contentType = "image/png"
		default:
			data = ""
			contentType = "image/png"
		}

		if nil != data {
			// try convert data
			if s, b := data.(string); b {
				bdata, err := qbc.Coding.DecodeBase64(s)
				if nil == err {
					data = bdata
				}
			}
			// HTTP
			if nil != instance.ctxHttp {
				err := writeHttp(instance.ctxHttp, contentType, data)
				if nil != err {
					panic(instance.runtime.NewTypeError(err.Error()))
				}
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return instance.runtime.ToValue("")
}

func (instance *JsHttpResponse) status(call goja.FunctionCall) goja.Value {
	if nil != instance {
		value := commons.GetInt(call, 0)
		if value > 0 {
			if nil != instance.ctxHttp {
				err := instance.ctxHttp.SendStatus(int(value))
				if nil != err {
					panic(instance.runtime.NewTypeError(err.Error()))
				}
			}
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpResponse) init() {
	o := instance.object
	if nil != o {
		_ = o.Set("size", instance.data.Len())
		_ = o.Set("header", instance.runtime.ToValue(instance.header))
		//_ = o.Set("body", instance.runtime.ToValue(instance.data.String()))
		if instance.statusCode > 0 {
			_ = o.Set("statusCode", instance.statusCode)
			_ = o.Set("statusMessage", instance.statusMessage)
			_ = o.Set("status", fmt.Sprintf("%v:%v", instance.statusCode, instance.statusMessage))
		}
	}
}

func (instance *JsHttpResponse) export() {
	o := instance.object
	_ = o.Set("length", instance.length)
	_ = o.Set("bytes", instance.bytes)

	_ = o.Set("send", instance.send)
	_ = o.Set("text", instance.text)
	_ = o.Set("json", instance.json)
	_ = o.Set("image", instance.image)
	_ = o.Set("html", instance.html)

	_ = o.Set("status", instance.status)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func writeHttp(ctx *fiber.Ctx, contentType string, data interface{}) error {
	if nil != ctx {
		var err error
		ctx.Response().Header.SetContentType(contentType)
		if v, b := data.(string); b {
			_, err = ctx.Response().BodyWriter().Write([]byte(v))
		} else if v, b := data.([]byte); b {
			_, err = ctx.Response().BodyWriter().Write(v)
		} else if v, b := data.([]uint8); b {
			_, err = ctx.Response().BodyWriter().Write(v)
		} else if v, b := data.(map[string]interface{}); b {
			raw, err := json.Marshal(v)
			if err == nil {
				_, err = ctx.Response().BodyWriter().Write(raw)
			}
		}
		return err
	}
	return nil
}
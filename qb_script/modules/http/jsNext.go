package http

import (
	"sync"

	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsHttpNext struct {
	runtime *goja.Runtime
	object  *goja.Object
	ctx     *fiber.Ctx
	mux     sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpClient
//----------------------------------------------------------------------------------------------------------------------

func WrapNext(runtime *goja.Runtime, ctx *fiber.Ctx) *JsHttpNext {
	instance := new(JsHttpNext)
	instance.runtime = runtime
	instance.ctx = ctx

	instance.object = instance.runtime.NewObject()
	instance.export()

	return instance
}

func (instance *JsHttpNext) Value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

func (instance *JsHttpNext) NextFunc() goja.Value {
	if nil != instance.object {
		return instance.object.Get("next")
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpNext) next(_ goja.FunctionCall) goja.Value {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			_ = qbc.Strings.Format("[panic] jsNext.next -> \"%s\"", r)
			// TODO: implement logger
			// fmt.Println(message)
		}
	}()
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		ctx := instance.ctx
		if nil != ctx && nil != ctx.App() {
			err := instance.ctx.Next()
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsHttpNext) export() {
	o := instance.object
	// methods
	_ = o.Set("next", instance.next)

}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

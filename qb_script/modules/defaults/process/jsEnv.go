package process

import (
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsEnv struct {
	runtime *goja.Runtime
	object  *goja.Object
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpClient
//----------------------------------------------------------------------------------------------------------------------

func NewEnv(runtime *goja.Runtime) goja.Value {
	instance := new(JsEnv)
	instance.runtime = runtime

	instance.object = instance.runtime.NewObject()
	instance.export()

	return instance.value()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsEnv) names(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		names := make([]string, 0)
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			names = append(names, pair[0])
		}
		return instance.runtime.ToValue(names)
	}

	return goja.Undefined()
}

func (instance *JsEnv) get(call goja.FunctionCall) goja.Value {
	if nil != instance {
		name := commons.GetString(call, 0)
		if len(name) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		return instance.runtime.ToValue(os.Getenv(name))
	}

	return goja.Undefined()
}

func (instance *JsEnv) set(call goja.FunctionCall) goja.Value {
	if nil != instance {
		name := commons.GetString(call, 0)
		value := commons.GetString(call, 1)
		if len(name) == 0 || len(value) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		err := os.Setenv(name, value)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}

	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsEnv) value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

func (instance *JsEnv) export() {
	o := instance.object

	// load env variables as fields
	instance.initEnv(o)

	// methods
	_ = o.Set("names", instance.names)
	_ = o.Set("get", instance.get)
	_ = o.Set("set", instance.set)
}

func (instance *JsEnv) initEnv(o *goja.Object) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		_ = o.Set(pair[0], pair[1])
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

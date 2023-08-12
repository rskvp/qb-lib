package rndutils

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "rnd-utils"

type RndUtils struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// rnd.guid()
func (instance *RndUtils) guid(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Rnd.Uuid())
}

func (instance *RndUtils) tguid(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Rnd.UuidTimestamp())
}

func (instance *RndUtils) between(call goja.FunctionCall) goja.Value {
	var num1 int64
	var num2 int64
	switch len(call.Arguments) {
	case 2:
		num1 = commons.GetInt(call, 0)
		num2 = commons.GetInt(call, 1)
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return instance.runtime.ToValue(qbc.Rnd.Between(num1, num2))
}

func (instance *RndUtils) digits(call goja.FunctionCall) goja.Value {
	var num1 int
	switch len(call.Arguments) {
	case 1:
		num1 = int(commons.GetInt(call, 0))
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return instance.runtime.ToValue(qbc.Rnd.RndDigits(num1))
}

func (instance *RndUtils) chars(call goja.FunctionCall) goja.Value {
	var num1 int
	switch len(call.Arguments) {
	case 1:
		num1 = int(commons.GetInt(call, 0))
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return instance.runtime.ToValue(qbc.Rnd.RndChars(num1))
}

func (instance *RndUtils) charsLower(call goja.FunctionCall) goja.Value {
	var num1 int
	switch len(call.Arguments) {
	case 1:
		num1 = int(commons.GetInt(call, 0))
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return instance.runtime.ToValue(qbc.Rnd.RndCharsLower(num1))
}

func (instance *RndUtils) charsUpper(call goja.FunctionCall) goja.Value {
	var num1 int
	switch len(call.Arguments) {
	case 1:
		num1 = int(commons.GetInt(call, 0))
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return instance.runtime.ToValue(qbc.Rnd.RndCharsUpper(num1))
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &RndUtils{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)

	// uuid
	_ = o.Set("guid", instance.guid)
	_ = o.Set("tguid", instance.tguid)
	// random
	_ = o.Set("digits", instance.digits)
	_ = o.Set("between", instance.between)
	_ = o.Set("chars", instance.chars)
	_ = o.Set("charsLower", instance.charsLower)
	_ = o.Set("charsUpper", instance.charsUpper)
}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

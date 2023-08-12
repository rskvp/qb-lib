package dbal

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "dbal"

type ModuleDbal struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// dbal.create(driverName, dataSourceName)
// Creates new instance of database and does not use cache.
// For cache optimization use: dbal.get(driverName, dataSourceName)
func (instance *ModuleDbal) create(call goja.FunctionCall) goja.Value {
	var driverName, dataSourceName string
	switch len(call.Arguments) {
	case 2:
		driverName = commons.GetString(call, 0)
		dataSourceName = commons.GetString(call, 1)
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	if len(driverName) > 0 && len(dataSourceName) > 0 {
		db, err := drivers.NewDatabase(driverName, dataSourceName)
		if nil != err {
			// throw back error to javascript
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		if nil != db {
			return instance.runtime.ToValue(WrapDbal(instance.runtime, db))
		}
	} else {
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return goja.Undefined()
}

func (instance *ModuleDbal) get(call goja.FunctionCall) goja.Value {
	var driverName, dataSourceName string
	switch len(call.Arguments) {
	case 2:
		driverName = commons.GetString(call, 0)
		dataSourceName = commons.GetString(call, 1)
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	if len(driverName) > 0 && len(dataSourceName) > 0 {
		db, err := drivers.Cache().Get(driverName, dataSourceName)
		if nil != err {
			// throw back error to javascript
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		if nil != db {
			return instance.runtime.ToValue(WrapDbal(instance.runtime, db))
		}
	} else {
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return goja.Undefined()
}

func (instance *ModuleDbal) reset(call goja.FunctionCall) goja.Value {
	var uid string
	switch len(call.Arguments) {
	case 1:
		uid = commons.GetString(call, 0)
	default:
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	if len(uid) > 0 {
		drivers.Cache().Remove(uid)
	} else {
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return goja.Undefined()
}

func (instance *ModuleDbal) clear(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		drivers.Cache().Clear()
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	instance := &ModuleDbal{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)

	_ = o.Set("cache", WrapDbalCache(runtime))

	_ = o.Set("create", instance.create)
	_ = o.Set("get", instance.get)
	_ = o.Set("reset", instance.reset)
	_ = o.Set("clear", instance.clear)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

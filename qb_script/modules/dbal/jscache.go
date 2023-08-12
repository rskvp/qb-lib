package dbal

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsDbalCache struct {
	cache   *drivers.CacheManager
	runtime *goja.Runtime
	object  goja.Value
}

//----------------------------------------------------------------------------------------------------------------------
//	JsDatabase
//----------------------------------------------------------------------------------------------------------------------

func WrapDbalCache(runtime *goja.Runtime) goja.Value {
	instance := new(JsDbalCache)
	instance.runtime = runtime
	instance.cache = drivers.Cache()

	return instance.getObject()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDbalCache) clear(_ goja.FunctionCall) goja.Value {
	if nil != instance.cache {
		instance.cache.Clear()
	}
	return goja.Undefined()
}

func (instance *JsDbalCache) get(call goja.FunctionCall) goja.Value {
	if nil != instance.cache {
		var driver, dsn string
		switch len(call.Arguments) {
		case 2:
			driver = commons.GetString(call, 0)
			dsn = commons.GetString(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if len(driver) > 0 && len(dsn) > 0 {
			db, err := instance.cache.Get(driver, dsn)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
			return WrapDbal(instance.runtime, db)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsDbalCache) remove(call goja.FunctionCall) goja.Value {
	if nil != instance.cache {
		var key string
		switch len(call.Arguments) {
		case 1:
			key = commons.GetString(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		instance.cache.Remove(key)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDbalCache) getObject() goja.Value {
	if nil == instance.object {
		object := instance.runtime.NewObject()

		_ = object.Set("get", instance.get)
		_ = object.Set("clear", instance.clear)
		_ = object.Set("remove", instance.remove)

		instance.object = object
	}
	return instance.object
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func exportCache(instance *JsDbalCache, o *goja.Object) {
	_ = o.Set("clear", instance.clear)

}

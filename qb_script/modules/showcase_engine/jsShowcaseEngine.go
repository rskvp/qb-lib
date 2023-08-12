package showcase_engine

import (
	"errors"
	"time"

	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_dbal"
	"github.com/rskvp/qb-lib/qb_dbal/showcase_search"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsShowcaseEngine struct {
	runtime *goja.Runtime
	object  *goja.Object
	config  interface{}
	engine  *showcase_search.ShowcaseEngine
}

//----------------------------------------------------------------------------------------------------------------------
//	JsShowcaseEngine
//----------------------------------------------------------------------------------------------------------------------

func WrapEngine(runtime *goja.Runtime, config interface{}) goja.Value {
	instance := new(JsShowcaseEngine)
	instance.runtime = runtime
	instance.config = config

	instance.object = instance.runtime.NewObject()
	instance.export()

	// add closable: all closable objects must be exposed to avoid
	commons.AddClosableObject(instance.runtime, instance.object)

	return instance.value()
}

func (instance *JsShowcaseEngine) value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsShowcaseEngine) open(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) close(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		//empty
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) setAutoResetSession(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var value bool
		switch len(call.Arguments) {
		case 1:
			value = commons.GetBool(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		instance.engine.SetAutoResetSession(value)
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) put(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var payload interface{}
		var category string
		var timestamp int64
		switch len(call.Arguments) {
		case 2:
			payload = commons.GetExport(call, 0)
			category = commons.GetString(call, 1)
			timestamp = time.Now().Unix()
		case 3:
			payload = commons.GetExport(call, 0)
			category = commons.GetString(call, 1)
			timestamp = commons.GetInt(call, 2)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		item, err := instance.engine.Put(payload, timestamp, category)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(item)
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) get(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var key string
		switch len(call.Arguments) {
		case 1:
			key = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		item, err := instance.engine.Get(key)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(item)
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) delete(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var key string
		switch len(call.Arguments) {
		case 1:
			key = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		item, err := instance.engine.Delete(key)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(item)
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) update(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var key string
		var payload interface{}
		switch len(call.Arguments) {
		case 2:
			key = commons.GetString(call, 0)
			payload = commons.GetExport(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		item, err := instance.engine.Update(key, payload)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(item)
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) query(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var sessionId string
		var limit int64
		switch len(call.Arguments) {
		case 1:
			sessionId = commons.GetString(call, 0)
			limit = 1
		case 2:
			sessionId = commons.GetString(call, 0)
			limit = commons.GetInt(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		data := instance.engine.Query(sessionId, int(limit))
		return instance.runtime.ToValue(data)
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) reset(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var sessionId string
		switch len(call.Arguments) {
		case 0:
			sessionId = ""
		case 1:
			sessionId = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if len(sessionId) > 0 {
			// SESSION
			instance.engine.ResetSession(sessionId)
		} else {
			// ALL
			instance.engine.Reset()
		}
	}
	return goja.Undefined()
}

func (instance *JsShowcaseEngine) setSessionCategoryWeight(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var sessionId, category string
		var inTime bool
		var weight int64
		switch len(call.Arguments) {
		case 4:
			sessionId = commons.GetString(call, 0)
			category = commons.GetString(call, 1)
			inTime = commons.GetBool(call, 2)
			weight = commons.GetInt(call, 3)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		data := instance.engine.SetSessionCategoryWeight(sessionId, category, inTime, int(weight))
		return instance.runtime.ToValue(data)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsShowcaseEngine) init() (err error) {
	if nil == instance.engine {
		if nil != instance.config {
			instance.engine, err = qb_dbal.NewShowcaseEngine(instance.config)
		} else {
			return errors.New("missing_configuration")
		}
	}
	return err
}

func (instance *JsShowcaseEngine) export() {
	o := instance.object

	_ = o.Set("open", instance.open)
	_ = o.Set("close", instance.close)
	_ = o.Set("setAutoResetSession", instance.setAutoResetSession)
	_ = o.Set("put", instance.put)
	_ = o.Set("get", instance.get)
	_ = o.Set("delete", instance.delete)
	_ = o.Set("update", instance.update)
	_ = o.Set("query", instance.query)
	_ = o.Set("reset", instance.reset)
	_ = o.Set("setSessionCategoryWeight", instance.setSessionCategoryWeight)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

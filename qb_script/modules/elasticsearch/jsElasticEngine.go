package elasticsearch

import (
	"errors"

	"github.com/dop251/goja"
	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
	"github.com/rskvp/qb-lib/qb_dbal/semantic_search"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsElasticEngine struct {
	runtime *goja.Runtime
	object  *goja.Object
	config  interface{}
	engine  *semantic_search.SemanticEngine
}

//----------------------------------------------------------------------------------------------------------------------
//	JsElasticEngine
//----------------------------------------------------------------------------------------------------------------------

func WrapEngine(runtime *goja.Runtime, config interface{}) goja.Value {
	instance := new(JsElasticEngine)
	instance.runtime = runtime
	instance.config = config

	instance.object = instance.runtime.NewObject()
	instance.export()

	// add closable: all closable objects must be exposed to avoid
	commons.AddClosableObject(instance.runtime, instance.object)

	return instance.value()
}

func (instance *JsElasticEngine) value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsElasticEngine) open(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsElasticEngine) close(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		//empty

	}
	return goja.Undefined()
}

func (instance *JsElasticEngine) put(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var group, key, text string
		switch len(call.Arguments) {
		case 3:
			group = commons.GetString(call, 0)
			key = commons.GetString(call, 1)
			text = commons.GetString(call, 2)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		err = instance.engine.Put(group, key, text)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsElasticEngine) get(call goja.FunctionCall) goja.Value {
	if nil != instance {
		err := instance.init()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}

		var group, text string
		var offset, count int
		switch len(call.Arguments) {
		case 1:
			group = ""
			text = commons.GetString(call, 0)
			offset = 0
			count = 0
		case 2:
			group = commons.GetString(call, 0)
			text = commons.GetString(call, 1)
			offset = 0
			count = 0
		case 4:
			group = commons.GetString(call, 0)
			text = commons.GetString(call, 1)
			offset = int(commons.GetInt(call, 2))
			count = int(commons.GetInt(call, 3))
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		data, err := instance.engine.Get(group, text, offset, count)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(data)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsElasticEngine) init() (err error) {
	if nil == instance.engine {
		if nil != instance.config {
			config := dbalcommons.NewSemanticConfig()
			err = config.Parse(instance.config)
			if nil == err {
				instance.engine, err = semantic_search.NewSemanticEngine(config)
			}
		} else {
			return errors.New("missing_configuration")
		}
	}
	return err
}

func (instance *JsElasticEngine) export() {
	o := instance.object

	_ = o.Set("open", instance.open)
	_ = o.Set("close", instance.close)
	_ = o.Set("put", instance.put)
	_ = o.Set("get", instance.get)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

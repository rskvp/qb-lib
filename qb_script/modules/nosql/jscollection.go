package nosql

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/nosql/nosqldrivers"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsCollection struct {
	collection nosqldrivers.INoSqlCollection
	runtime    *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	JsCollection
//----------------------------------------------------------------------------------------------------------------------

func WrapNoSqlCollection(runtime *goja.Runtime, collection nosqldrivers.INoSqlCollection) goja.Value {
	instance := new(JsCollection)
	instance.runtime = runtime
	instance.collection = collection

	object := instance.runtime.NewObject()
	exportCollection(instance, object)

	return object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsCollection) name(call goja.FunctionCall) goja.Value {
	if nil != instance.collection {
		return instance.runtime.ToValue(instance.collection.Name())
	}
	return goja.Undefined()
}

func (instance *JsCollection) ensureIndex(call goja.FunctionCall) goja.Value {
	if nil != instance.collection {
		if len(call.Arguments) == 3 {
			typeName := commons.GetString(call, 0)
			fields := commons.GetArrayOfString(call, 1)
			unique := commons.GetBool(call, 2)
			_, err := instance.collection.EnsureIndex(typeName, fields, unique)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsCollection) removeIndex(call goja.FunctionCall) goja.Value {
	if nil != instance.collection {
		if len(call.Arguments) < 3 {
			typeName := commons.GetString(call, 0)
			fields := qbc.Convert.ToArrayOfString(commons.GetArray(call, 1))
			_, err := instance.collection.RemoveIndex(typeName, fields)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func exportCollection(instance *JsCollection, o *goja.Object) {
	_ = o.Set("name", instance.name)
	_ = o.Set("ensureIndex", instance.ensureIndex)
	_ = o.Set("removeIndex", instance.removeIndex)

}

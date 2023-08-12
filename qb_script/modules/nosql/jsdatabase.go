package nosql

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/nosql/nosqldrivers"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsDatabase struct {
	database nosqldrivers.INoSqlDatabase
	runtime  *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	JsDatabase
//----------------------------------------------------------------------------------------------------------------------

func Wrap(runtime *goja.Runtime, database nosqldrivers.INoSqlDatabase) goja.Value {
	instance := new(JsDatabase)
	instance.runtime = runtime
	instance.database = database

	object := instance.runtime.NewObject()
	export(instance, object)

	return object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDatabase) close(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		err := instance.database.Close()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsDatabase) query(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		query := commons.GetString(call, 0)
		bindVars := commons.GetMap(call, 1)
		if len(query) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.Query(query, bindVars)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) exec(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		query := commons.GetString(call, 0)
		bindVars := commons.GetMap(call, 1)
		if len(query) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.Exec(query, bindVars)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) insert(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		collection := commons.GetString(call, 0)
		item := commons.GetMap(call, 1)
		if len(collection) == 0 || nil == item {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.Insert(collection, item)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) update(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		collection := commons.GetString(call, 0)
		item := commons.GetMap(call, 1)
		if len(collection) == 0 || nil == item {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.Update(collection, item)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) upsert(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		collection := commons.GetString(call, 0)
		item := commons.GetMap(call, 1)
		if len(collection) == 0 || nil == item {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.Upsert(collection, item)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) delete(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) == 2 {
			collection := commons.GetString(call, 0)
			itemOrKey := call.Argument(1).Export()
			response, err := instance.database.Delete(collection, itemOrKey)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
			return instance.runtime.ToValue(response)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsDatabase) count(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		query := commons.GetString(call, 0)
		bindVars := commons.GetMap(call, 1)
		if len(query) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.Count(query, bindVars)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) collection(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		name := commons.GetString(call, 0)
		create := true
		if len(call.Arguments) == 2 {
			create = call.Arguments[1].ToBoolean()
		}
		coll, err := instance.database.Collection(name, create)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return WrapNoSqlCollection(instance.runtime, coll)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func export(instance *JsDatabase, o *goja.Object) {
	_ = o.Set("close", instance.close)
	_ = o.Set("query", instance.query)
	_ = o.Set("exec", instance.exec)
	_ = o.Set("insert", instance.insert)
	_ = o.Set("update", instance.update)
	_ = o.Set("upsert", instance.upsert)
	_ = o.Set("delete", instance.delete)
	_ = o.Set("count", instance.count)
	// collection
	_ = o.Set("collection", instance.collection)

}

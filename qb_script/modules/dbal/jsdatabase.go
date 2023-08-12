package dbal

import (
	"errors"
	"strings"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsDbal struct {
	database drivers.IDatabase
	runtime  *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	JsDatabase
//----------------------------------------------------------------------------------------------------------------------

func WrapDbal(runtime *goja.Runtime, database drivers.IDatabase) goja.Value {
	instance := new(JsDbal)
	instance.runtime = runtime
	instance.database = database

	return instance.export(runtime.NewObject())
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDbal) close(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		err := instance.database.Close()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsDbal) uid(_ goja.FunctionCall) goja.Value {
	if nil != instance.database {
		response := instance.database.Uid()
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbal) name(_ goja.FunctionCall) goja.Value {
	if nil != instance.database {
		response := instance.database.DriverName()
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbal) ensureIndex(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var collection, typeName string
		var fields []string
		var unique bool
		switch len(call.Arguments) {
		case 4:
			collection = commons.GetString(call, 0)
			typeName = commons.GetString(call, 1)
			fields = commons.GetArrayOfString(call, 2)
			unique = commons.GetBool(call, 3)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		response, err := instance.database.EnsureIndex(collection, typeName, fields, unique)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbal) ensureCollection(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var collection string
		switch len(call.Arguments) {
		case 1:
			collection = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		response, err := instance.database.EnsureCollection(collection)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbal) get(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var collection, key string
		switch len(call.Arguments) {
		case 2:
			collection = commons.GetString(call, 0)
			key = commons.GetString(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		response, err := instance.database.Get(collection, key)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbal) upsert(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var collection string
		var doc map[string]interface{}
		switch len(call.Arguments) {
		case 2:
			collection = commons.GetString(call, 0)
			doc = commons.GetMap(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		response, err := instance.database.Upsert(collection, doc)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbal) remove(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var collection, key string
		switch len(call.Arguments) {
		case 2:
			collection = commons.GetString(call, 0)
			key = commons.GetString(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		err := instance.database.Remove(collection, key)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(true)
	}
	return goja.Undefined()
}

func (instance *JsDbal) foreach(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var collection string
		var callback goja.Callable
		switch len(call.Arguments) {
		case 2:
			collection = commons.GetString(call, 0)
			callback = commons.GetCallbackIfAny(call)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		var err error
		err = instance.database.ForEach(collection, func(m map[string]interface{}) bool {
			result, cErr := callback(instance.runtime.ToValue(m))
			if nil != cErr {
				err = cErr
				return true // exit due error
			}
			if result.Equals(goja.Undefined()) || result.Equals(goja.Null()) {
				return false // continue
			}
			//result.ToBoolean()
			return result.ToBoolean()
		})
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsDbal) exec(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		query := commons.GetString(call, 0)
		bindVars := commons.GetMap(call, 1)
		if len(query) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.ExecNative(query, bindVars)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

// query returns a query wrapper
func (instance *JsDbal) query(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		queryOrFilename := commons.GetString(call, 0)
		if len(queryOrFilename) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.buildQuery(queryOrFilename)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return response
	}
	return goja.Undefined()
}

// parseParamNames returns an array of name extracted from query parameters
func (instance *JsDbal) parseParamNames(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		queryOrFilename := commons.GetString(call, 0)
		if len(queryOrFilename) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		query, err := instance.readQuery(queryOrFilename)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		response := instance.database.QueryGetParamNames(query)
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

// matchParams get a query (or path to a query) and some parameters, return parameters used in query
func (instance *JsDbal) matchParams(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		queryOrFilename := commons.GetString(call, 0)
		params := commons.GetMap(call, 0)
		if len(queryOrFilename) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if nil==params{
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		query, err := instance.readQuery(queryOrFilename)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		response := instance.database.QuerySelectParams(query, params)
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDbal) readQuery(queryOrFilename string) (string, error) {
	var query string
	if b, e := qbc.Paths.IsFile(queryOrFilename); b && nil == e {
		text, err := qbc.IO.ReadTextFromFile(queryOrFilename)
		if nil != err {
			return "", err
		}
		query = text
	} else {
		// non an existing file or a query
		if strings.Index(queryOrFilename, "/") == -1 || strings.Index(queryOrFilename, "\\") == -1 {
			// should be a query, not a file path
			query = queryOrFilename
		}
	}
	if len(query) > 0 {
		return query, nil
	}
	return "", errors.New("invalid file or query: " + queryOrFilename)
}

func (instance *JsDbal) buildQuery(queryOrFilename string) (goja.Value, error) {
	query, err := instance.readQuery(queryOrFilename)
	if len(query) > 0 {
		return WrapDbalQuery(instance.runtime, instance.database, query), nil
	}
	return nil, err
}

func (instance *JsDbal) export(o *goja.Object) *goja.Object {
	_ = o.Set("close", instance.close)
	_ = o.Set("name", instance.name)
	_ = o.Set("uid", instance.uid)
	_ = o.Set("ensureIndex", instance.ensureIndex)
	_ = o.Set("ensureCollection", instance.ensureCollection)
	_ = o.Set("get", instance.get)
	_ = o.Set("upsert", instance.upsert)
	_ = o.Set("foreach", instance.foreach)
	_ = o.Set("exec", instance.exec)
	_ = o.Set("query", instance.query)
	_ = o.Set("parseParamNames", instance.parseParamNames)
	_ = o.Set("matchParams", instance.matchParams)

	return o
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

package dbal

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsDbalQuery struct {
	queryCommand string
	database     drivers.IDatabase
	runtime      *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	JsDbalQuery
//----------------------------------------------------------------------------------------------------------------------

func WrapDbalQuery(runtime *goja.Runtime, database drivers.IDatabase, command string) goja.Value {
	instance := new(JsDbalQuery)
	instance.runtime = runtime
	instance.database = database
	instance.queryCommand = command
	return instance.export(instance.runtime.NewObject())
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

/**
 command returns the string of query to execute
 usage: var cmd = query.command();
 */
func (instance *JsDbalQuery) command(_ goja.FunctionCall) goja.Value {
	if nil != instance.database {
		response := instance.queryCommand
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDbalQuery) foreach(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		var callback goja.Callable
		switch len(call.Arguments) {
		case 1:
			callback = commons.GetCallbackIfAny(call)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		var err error
		err = instance.database.ForEach(instance.queryCommand, func(m map[string]interface{}) bool {
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

func (instance *JsDbalQuery) exec(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		bindVars := commons.GetMap(call, 0)
		if len(instance.queryCommand) == 0 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.database.ExecNative(instance.queryCommand, bindVars)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDbalQuery) export(o *goja.Object) *goja.Object {
	_ = o.Set("command", instance.command)
	_ = o.Set("foreach", instance.foreach)
	_ = o.Set("exec", instance.exec)

	return o
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

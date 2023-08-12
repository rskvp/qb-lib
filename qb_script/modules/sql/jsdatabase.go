package sql

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/sql/helpers"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsDatabase struct {
	database *helpers.Database
	runtime  *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	JsDatabase
//----------------------------------------------------------------------------------------------------------------------

func Wrap(runtime *goja.Runtime, database *helpers.Database) goja.Value {
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
		query := call.Argument(0).String()
		args := commons.ToArray(call.Arguments[1:])
		if len(call.Arguments) < 1 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		rows := instance.database.Query(query, args...)
		defer rows.Close()
		err := rows.GetError()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		dataset := make([]map[string]interface{}, 0)
		err = rows.ForEach(func(item map[string]interface{}) bool {
			dataset = append(dataset, item)
			return false
		})
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(dataset)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) exec(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) < 2 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		query := call.Argument(0).String()
		args := commons.ToArray(call.Arguments[1:])
		result := instance.database.Exec(query, args...)
		err := result.GetError()
		if nil != err {
			panic(instance.runtime.NewTypeError(result.Error()))
		}
		response, err := sqlResultToMap(result)
		if nil != err {
			panic(instance.runtime.NewTypeError(result.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsDatabase) insert(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) < 2 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		tableName := commons.GetString(call, 0)
		items := commons.GetArray(call, 1)
		if len(tableName) > 0 && nil != items {
			commands := helpers.BuildInsertCommands(tableName, items)
			var lastInsertId int64
			var rowsAffected int64
			for _, command := range commands {
				if len(command) > 0 {
					result := instance.database.Exec(command)
					err := result.GetError()
					if nil != err {
						panic(instance.runtime.NewTypeError(result.Error()))
					}
					lastInsertId, _ = result.LastInsertId()
					r, _ := result.RowsAffected()
					rowsAffected += r
				}
			}

			return instance.runtime.ToValue(map[string]interface{}{
				"last_id":       lastInsertId,
				"rows_affected": rowsAffected,
			})
		}
	}
	return goja.Undefined()
}

func (instance *JsDatabase) update(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) < 4 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		tableName := commons.GetString(call, 0)
		keyName := commons.GetString(call, 1)
		keyValue := commons.GetExport(call, 2)
		item := commons.GetExport(call, 3)

		if len(tableName) > 0 && nil != item {
			if m, b := item.(map[string]interface{}); b {
				command := helpers.BuildUpdateCommand(tableName, keyName, keyValue, m)
				var lastInsertId int64
				var rowsAffected int64
				if len(command) > 0 {
					result := instance.database.Exec(command)
					err := result.GetError()
					if nil != err {
						panic(instance.runtime.NewTypeError(result.Error()))
					}
					lastInsertId, _ = result.LastInsertId()
					r, _ := result.RowsAffected()
					rowsAffected += r
				}

				return instance.runtime.ToValue(map[string]interface{}{
					"last_id":       lastInsertId,
					"rows_affected": rowsAffected,
				})
			}

		}
	}
	return goja.Undefined()
}

func (instance *JsDatabase) delete(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) < 2 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		tableName := commons.GetString(call, 0)
		filter := commons.GetString(call, 1)

		if len(tableName) > 0 && len(filter) > 0 {
			command := helpers.BuildDeleteCommand(tableName, filter)
			var lastInsertId int64
			var rowsAffected int64
			if len(command) > 0 {
				result := instance.database.Exec(command)
				err := result.GetError()
				if nil != err {
					panic(instance.runtime.NewTypeError(result.Error()))
				}
				lastInsertId, _ = result.LastInsertId()
				r, _ := result.RowsAffected()
				rowsAffected += r
			}

			return instance.runtime.ToValue(map[string]interface{}{
				"last_id":       lastInsertId,
				"rows_affected": rowsAffected,
			})
		}
	}
	return goja.Undefined()
}

func (instance *JsDatabase) count(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) < 1 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		tableName := commons.GetString(call, 0)
		filter := commons.GetString(call, 1)

		if len(tableName) > 0 {
			command := helpers.BuildCountCommand(tableName, filter)
			var response int64
			if len(command) > 0 {
				result := instance.database.QueryRow(command, &response)
				err := result.GetError()
				if nil != err {
					panic(instance.runtime.NewTypeError(result.Error()))
				}
			}
			return instance.runtime.ToValue(response)
		}
	}
	return goja.Undefined()
}

func (instance *JsDatabase) countDistinct(call goja.FunctionCall) goja.Value {
	if nil != instance.database {
		if len(call.Arguments) < 1 {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		tableName := commons.GetString(call, 0)
		fieldName := commons.GetString(call, 1)
		filter := commons.GetString(call, 2)

		if len(tableName) > 0 {
			command := helpers.BuildCountDistinctCommand(tableName, fieldName, filter)
			var response int64
			if len(command) > 0 {
				result := instance.database.QueryRow(command, &response)
				err := result.GetError()
				if nil != err {
					panic(instance.runtime.NewTypeError(result.Error()))
				}
			}
			return instance.runtime.ToValue(response)
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func sqlResultToMap(result *helpers.DatabaseResult) (map[string]interface{}, error) {
	if nil != result.GetError() {
		return nil, result.GetError()
	}
	lastInsertId, err := result.LastInsertId()
	if nil != err {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if nil != err {
		return nil, err
	}
	response := map[string]interface{}{
		"last_id":       lastInsertId,
		"rows_affected": rowsAffected,
	}

	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func export(instance *JsDatabase, o *goja.Object) {
	_ = o.Set("close", instance.close)
	_ = o.Set("query", instance.query)
	_ = o.Set("exec", instance.exec)
	_ = o.Set("insert", instance.insert)
	_ = o.Set("update", instance.update)
	_ = o.Set("delete", instance.delete)
	_ = o.Set("count", instance.count)
	_ = o.Set("countDistinct", instance.countDistinct)

}

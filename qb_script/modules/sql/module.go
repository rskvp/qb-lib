package sql

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
	"github.com/rskvp/qb-lib/qb_script/modules/sql/helpers"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "sql"

type ModuleSql struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// sql.open(driverName, dataSourceName)
func (instance *ModuleSql) open(call goja.FunctionCall) goja.Value {
	driverName := call.Argument(0).String()
	dataSourceName := call.Argument(1).String()
	if len(driverName) > 0 && len(dataSourceName) > 0 {
		db, err := helpers.NewDatabase(driverName, dataSourceName)
		if nil != err {
			// throw back error to javascript
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		if nil != db {
			return instance.runtime.ToValue(Wrap(instance.runtime, db))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &ModuleSql{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("open", instance.open)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

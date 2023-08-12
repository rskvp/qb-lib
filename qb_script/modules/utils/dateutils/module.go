package dateutils

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "date-utils"

type DateUtils struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// wrap Return a date wrapper
// @param [date, string, null]
// @usage date.wrap(new Date())
func (instance *DateUtils) wrap(call goja.FunctionCall) goja.Value {
	source := call.Argument(0).Export()
	return WrapDate(instance.runtime, source)
}

// getLayout Return layout of date string
// @param strDate:string A date in string format
// @usage date.getLayout(datestring)
func (instance *DateUtils) getLayout(call goja.FunctionCall) goja.Value {
	if nil != instance {
		response := "" // empty response
		strDate := qbc.Convert.ToString(call.Argument(0).Export())
		if len(strDate) > 0 {
			v, err := qbc.Dates.ParseFormat(strDate)
			if nil != err {
				panic(instance.runtime.NewTypeError(err.Error()))
			}
			response = v
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}

		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &DateUtils{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	// date utility
	_ = o.Set("wrap", instance.wrap)
	_ = o.Set("getLayout", instance.getLayout)
}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

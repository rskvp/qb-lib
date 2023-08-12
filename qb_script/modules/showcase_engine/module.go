package showcase_engine

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "showcase-qb_sms_engine"

type ModuleShowcaseSearch struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// showcase.newEngine(configuration)
func (instance *ModuleShowcaseSearch) newEngine(call goja.FunctionCall) goja.Value {
	if nil != instance {
		if len(call.Arguments) > 0 {
			settings := commons.GetExport(call, 0)
			return WrapEngine(instance.runtime, settings)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &ModuleShowcaseSearch{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("newEngine", instance.newEngine)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

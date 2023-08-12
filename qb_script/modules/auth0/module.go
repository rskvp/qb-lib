package auth0

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "auth0"

type ModuleHttp struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// auth0.newEngine(config)
func (instance *ModuleHttp) newEngine(call goja.FunctionCall) goja.Value {
	if nil != instance {
		if len(call.Arguments) > 0 {
			config := commons.GetExport(call, 0)
			return instance.runtime.ToValue(WrapAuth0Config(instance.runtime, config).Value())
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
	instance := &ModuleHttp{
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

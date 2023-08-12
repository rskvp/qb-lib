package http

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "http"

type ModuleHttp struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// http.newClient()
func (instance *ModuleHttp) newClient(call goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(WrapClient(instance.runtime))
}

// http.newServer()
func (instance *ModuleHttp) newServer(call goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(WrapServer(instance.runtime))
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &ModuleHttp{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("newClient", instance.newClient)
	_ = o.Set("newServer", instance.newServer) // TODO: DEPRECATE

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

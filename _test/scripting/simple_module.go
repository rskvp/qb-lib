package scripting

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "simple"

type ModuleSimple struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// simple.echo("hello")
func (instance *ModuleSimple) echo(call goja.FunctionCall) goja.Value {
	message := call.Argument(0).String()
	return instance.runtime.ToValue(message)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	instance := &ModuleSimple{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("echo", instance.echo)
}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})

	// ctx.Runtime.Set(NAME, require.Require(ctx.Runtime, NAME))
}

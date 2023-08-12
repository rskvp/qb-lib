package process

import (
	"github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	p r o c e s s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "process"

type Process struct {
	runtime *goja.Runtime
	object  *goja.Object
	env     goja.Value
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *Process) close(_ goja.FunctionCall) goja.Value {
	if nil != instance {

	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &Process{
		runtime: runtime,
		env:     NewEnv(runtime),
	}

	o := module.Get("exports").(*goja.Object)

	// properties
	_ = o.Set("env", instance.env)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})

	// add require to javascript context
	ctx.Runtime.Set(NAME, require.Require(ctx.Runtime, NAME))
}

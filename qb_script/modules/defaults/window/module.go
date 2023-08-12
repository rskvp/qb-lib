package window

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

var FIELDS = []string{"context", "params"}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

const NAME = "window"

type ModuleRuntime struct {
	runtime *goja.Runtime
	object  *goja.Object
	context goja.Value
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *ModuleRuntime) params(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		m := make(map[string]interface{})
		keys := instance.object.Keys()
		for _, k := range keys {
			if qbc.Arrays.IndexOf(k, FIELDS) == -1 {
				m[k] = instance.object.Get(k).Export()
			}
		}
		return instance.runtime.ToValue(m)
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
	instance := &ModuleRuntime{
		runtime: runtime,
		object:  module.Get("exports").(*goja.Object),
		context: NewContext(runtime),
	}

	// properties
	// context: is same as runtime.context or an empty object
	if c := commons.GetRtDeepValue(runtime, "runtime.context"); nil != c {
		_ = instance.object.Set(FIELDS[0], c)
	} else {
		_ = instance.object.Set(FIELDS[0], runtime.NewObject())
	}
	_ = instance.object.Set(FIELDS[1], instance.params) // params

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

package nio

import (
	"strings"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "nio"

type ModuleNio struct {
	runtime *goja.Runtime
}

type ModuleNioMessage struct {
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// nio.newClient(address)
func (instance *ModuleNio) newClient(call goja.FunctionCall) goja.Value {
	address := call.Argument(0).String()
	if len(address) > 0 {
		host := address
		port := 10001
		tokens := strings.Split(address, ":")
		if len(tokens) == 2 {
			host = tokens[0]
			port = qbc.Convert.ToInt(tokens[1])
		}
		return instance.runtime.ToValue(WrapClient(instance.runtime, host, port))
	}
	return goja.Undefined()
}

// nio.newClient(address)
func (instance *ModuleNio) newServer(call goja.FunctionCall) goja.Value {
	port := int(call.Argument(0).ToInteger())
	if port > 0 {
		return instance.runtime.ToValue(WrapServer(instance.runtime, port))
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &ModuleNio{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("newClient", instance.newClient)
	_ = o.Set("newServer", instance.newServer)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

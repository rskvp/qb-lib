package sys

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "sys"

type ModuleSys struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// sys.id()
func (instance *ModuleSys) id(_ goja.FunctionCall) goja.Value {
	// p := call.Argument(0).String()
	id, err := qbc.Sys.ID()
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	if len(id) > 0 {
		return instance.runtime.ToValue(id)
	}
	return goja.Undefined()
}

func (instance *ModuleSys) shutdown(_ goja.FunctionCall) goja.Value {
	err := qbc.Sys.Shutdown()
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return goja.Undefined()
}

func (instance *ModuleSys) getOS(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Sys.GetOS())
}

func (instance *ModuleSys) getOSVersion(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Sys.GetOSVersion())
}

func (instance *ModuleSys) getInfo(_ goja.FunctionCall) goja.Value {
	info := qbc.Sys.GetInfo()
	if nil != info {
		return instance.runtime.ToValue(map[string]interface{}{
			"core":     info.Core,
			"cpu":      info.CPUs,
			"hostname": info.Hostname,
			"kernel":   info.Kernel,
			"os":       info.OS,
			"platform": info.Platform, // x86_64
		})
	}
	return goja.Undefined()
}

func (instance *ModuleSys) isLinux(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Sys.IsLinux())
}

func (instance *ModuleSys) isMac(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Sys.IsMac())
}

func (instance *ModuleSys) isWindows(_ goja.FunctionCall) goja.Value {
	return instance.runtime.ToValue(qbc.Sys.IsWindows())
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &ModuleSys{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("id", instance.id)
	_ = o.Set("shutdown", instance.shutdown)
	_ = o.Set("getOS", instance.getOS)
	_ = o.Set("getOSVersion", instance.getOSVersion)
	_ = o.Set("getInfo", instance.getInfo)
	_ = o.Set("isLinux", instance.isLinux)
	_ = o.Set("isMac", instance.isMac)
	_ = o.Set("isWindows", instance.isWindows)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

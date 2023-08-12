package require

import (
	"errors"
	"fmt"
	"path/filepath"

	js "github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

const NAME = "require"

var (
	InvalidModuleError     = errors.New("Invalid module")
	IllegalModuleNameError = errors.New("Illegal module name")
)

type RequireModule struct {
	ctx *commons.RuntimeContext
	registry *Registry
	modules  map[string]*js.Object
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewRequireModule(ctx *commons.RuntimeContext, registry *Registry) *RequireModule {
	return &RequireModule{
		registry: registry,
		ctx:  ctx,
		modules:  make(map[string]*js.Object),
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Require can be used to import modules from Go source (similar to JS require() function).
func (instance *RequireModule) GoRequire(p string) (ret js.Value, err error) {
	p = filepath.Clean(p)
	if p == "" {
		err = IllegalModuleNameError
		return
	}
	module := instance.modules[p]
	if module == nil {
		module = instance.ctx.Runtime.NewObject()
		_ = module.Set("exports", instance.ctx.Runtime.NewObject())
		instance.modules[p] = module
		err = instance.loadModule(p, module)
		if err != nil {
			delete(instance.modules, p)
			err = fmt.Errorf("Could not load module '%s': %v", p, err)
			return
		}
	}
	ret = module.Get("exports")
	return
}

func (instance *RequireModule) JsRequire(call js.FunctionCall) js.Value {
	ret, err := instance.GoRequire(call.Argument(0).String())
	if err != nil {
		panic(instance.ctx.Runtime.NewGoError(err))
	}
	return ret
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *RequireModule) loadModule(path string, jsModule *js.Object) error {
	contextId := ""
	if nil!=instance.ctx && nil!=instance.ctx.Uid{
		contextId = *instance.ctx.Uid
	}
	if info, exists := GetModuleInfo(path, contextId); exists {
		ldr := info.Loader
		arguments := make([]interface{}, 0)
		if nil != info.Context && len(info.Context.Arguments) > 0 {
			arguments = info.Context.Arguments
		}
		ldr(instance.ctx.Runtime, jsModule, arguments...)
		return nil
	}

	prg, err := instance.registry.GetCompiled(path)

	if err != nil {
		return err
	}

	f, err := instance.ctx.Runtime.RunProgram(prg)
	if err != nil {
		return err
	}

	if call, ok := js.AssertFunction(f); ok {
		jsExports := jsModule.Get("exports")

		// Run the module source, with "jsModule" as the "module" variable, "jsExports" as "this"(Nodejs capable).
		_, err = call(jsExports, jsModule, jsExports)
		if err != nil {
			return err
		}
	} else {
		return InvalidModuleError
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func Require(runtime *js.Runtime, name string) js.Value {
	if r, ok := js.AssertFunction(runtime.Get(NAME)); ok {
		mod, err := r(js.Undefined(), runtime.ToValue(name))
		if err != nil {
			panic(err)
		}
		return mod
	}
	panic(runtime.NewTypeError("Please enable require for this runtime using new(require.Require).Enable(runtime)"))
}

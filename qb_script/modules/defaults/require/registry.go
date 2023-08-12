package require

import (
	"io/ioutil"
	"sync"

	js "github.com/dop251/goja"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type SourceLoader func(path string) ([]byte, error)

// type ModuleLoader func(*js.Runtime, *js.Object)

var native map[string]map[string]*commons.ModuleInfo
var nativeMux sync.Mutex

// Registry contains a cache of compiled modules which can be used by multiple Runtimes
type Registry struct {
	sync.Mutex
	compiled map[string]*js.Program

	srcLoader SourceLoader
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewRegistryWithLoader(srcLoader SourceLoader) *Registry {
	return &Registry{
		srcLoader: srcLoader,
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Enable adds the require() function to the specified runtime.
func (instance *Registry) Enable(ctx *commons.RuntimeContext) *RequireModule {
	// creates module "require"
	rrt := NewRequireModule(ctx, instance)
	_ = ctx.Runtime.Set(NAME, rrt.JsRequire)
	_ = ctx.Runtime.Set("use", rrt.JsRequire)

	// creates "runtime" object with a context and require method
	rto := ctx.Runtime.NewObject()
	_ = rto.Set(NAME, rrt.JsRequire)
	_ = rto.Set("use", rrt.JsRequire)
	_ = rto.Set("context", ctx.Runtime.NewObject())
	ctx.Runtime.Set("runtime", rto)

	return rrt
}

func (instance *Registry) GetCompiled(path string) (prg *js.Program, err error) {
	return instance.getCompiledSource(path)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Registry) getCompiledSource(path string) (prg *js.Program, err error) {
	instance.Lock()
	defer instance.Unlock()

	prg = instance.compiled[path]
	if prg == nil {
		srcLoader := instance.srcLoader
		if srcLoader == nil {
			srcLoader = ioutil.ReadFile
		}
		if s, err1 := srcLoader(path); err1 == nil {
			source := "(function(module, exports) {" + string(s) + "\n})"
			prg, err = js.Compile(path, source, false)
			if err == nil {
				if instance.compiled == nil {
					instance.compiled = make(map[string]*js.Program)
				}
				instance.compiled[path] = prg
			}
		} else {
			err = err1
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func GetModuleInfo(name, contextId string) (*commons.ModuleInfo, bool) {
	nativeMux.Lock()
	defer nativeMux.Unlock()

	if native == nil {
		native = make(map[string]map[string]*commons.ModuleInfo)
	}
	if len(contextId) == 0 {
		contextId = name
	}
	if moduleRuntimes, b := native[name]; b {
		if info, b := moduleRuntimes[contextId]; b {
			return info, true
		}
		if info, b := moduleRuntimes[name]; b {
			return info, true
		}
	}
	return nil, false
}

func RegisterNativeModule(name string, info *commons.ModuleInfo) {
	nativeMux.Lock()
	defer nativeMux.Unlock()

	if native == nil {
		native = make(map[string]map[string]*commons.ModuleInfo)
	}
	contextId := name
	if nil != info.Context && nil != info.Context.Uid && len(*info.Context.Uid) > 0 {
		contextId = *info.Context.Uid
	}
	if _, b := native[name]; !b {
		native[name] = make(map[string]*commons.ModuleInfo)
	}
	moduleRuntimes := native[name]
	moduleRuntimes[contextId] = info
}

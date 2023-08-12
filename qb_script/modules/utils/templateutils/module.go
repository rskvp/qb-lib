package templateutils

import (
	"github.com/cbroglie/mustache"
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "template-utils"

type TemplateUtils struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// template.mergeFile(path, data)
func (instance *TemplateUtils) mergeFile(call goja.FunctionCall) goja.Value {
	if nil != instance {
		var path string
		var context map[string]interface{}
		switch len(call.Arguments) {
		case 2:
			path = commons.GetString(call, 0)
			context = commons.GetMap(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		text, err := qbc.IO.ReadTextFromFile(path)
		if nil != err {
			panic(instance.runtime.NewTypeError(err))
		}
		content, err := mustache.Render(text, context)
		if nil != err {
			panic(instance.runtime.NewTypeError(err))
		}
		return instance.runtime.ToValue(content)
	}
	return goja.Undefined()
}

// template.mergeText(text, data)
func (instance *TemplateUtils) mergeText(call goja.FunctionCall) goja.Value {
	if nil != instance {
		var text string
		var context map[string]interface{}
		switch len(call.Arguments) {
		case 2:
			text = commons.GetString(call, 0)
			context = commons.GetMap(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		content, err := mustache.Render(text, context)
		if nil != err {
			panic(instance.runtime.NewTypeError(err))
		}
		return instance.runtime.ToValue(content)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &TemplateUtils{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	// format
	_ = o.Set("mergeFile", instance.mergeFile)
	_ = o.Set("mergeText", instance.mergeText)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

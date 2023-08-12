package linereader

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

const NAME = "line-reader"

type LineReader struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// lineReader.eachLine((path, callback)
func (instance *LineReader) eachLine(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	callback, _ := goja.AssertFunction(call.Argument(len(call.Arguments) - 1))
	if len(path) > 0 && nil != callback {
		err := qbc.IO.ScanTextFromFile(path, func(text string) bool {
			vtext := instance.runtime.ToValue(text)
			v, err := callback(call.This, vtext)
			stop := nil != err || v.ToBoolean()
			return stop
		})
		if nil != err {
			// throw back error to javascript
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

// lineReader.readLines(path, max_count)
func (instance *LineReader) readLines(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	max := int(call.Argument(1).ToInteger())
	count := 0
	var lines strings.Builder
	err := qbc.IO.ScanTextFromFile(path, func(text string) bool {
		lines.WriteString(text + "\n")
		count++
		stop := count >= max
		return stop
	})
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return instance.runtime.ToValue(lines.String())
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &LineReader{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("eachLine", instance.eachLine)
	_ = o.Set("readLines", instance.readLines)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

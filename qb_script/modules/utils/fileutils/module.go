package fileutils

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "file-utils"

type FileUtils struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// utils.fileReadText(path)
func (instance *FileUtils) fileReadText(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	text, err := qbc.IO.ReadTextFromFile(path)
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return instance.runtime.ToValue(text)
}

// utils.fileReadBytes(path)
func (instance *FileUtils) fileReadBytes(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	bytes, err := qbc.IO.ReadBytesFromFile(path)
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return instance.runtime.ToValue(bytes)
}

// utils.fileWrite(path, data)
func (instance *FileUtils) fileWrite(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	data := call.Argument(1).Export()
	err := WriteDataToFile(path, data)
	if nil != err {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func WriteDataToFile(path string, data interface{}) error {
	if bytes, b := data.([]byte); b {
		_, err := qbc.IO.WriteBytesToFile(bytes, path)
		return err
	} else if text, b := data.(string); b {
		_, err := qbc.IO.WriteTextToFile(text, path)
		return err
	}
	return nil
}

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &FileUtils{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	// file utility
	_ = o.Set("fileReadText", instance.fileReadText)
	_ = o.Set("fileReadBytes", instance.fileReadBytes)
	_ = o.Set("fileWrite", instance.fileWrite)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

package fs

import (
	"os"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/fileutils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "fs"

type FileSystem struct {
	runtime *goja.Runtime
	util    *goja.Object
	root    string
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// fs.readFile(path[, options], callback)
func (instance *FileSystem) readFile(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	callback := commons.GetCallbackIfAny(call)
	if len(path) > 0 && nil != callback {
		data, err := qbc.IO.ReadTextFromFile(path)
		if nil != err {
			_, _ = callback(call.This, instance.runtime.ToValue(err.Error()), goja.Undefined())
		} else {
			_, _ = callback(call.This, goja.Undefined(), instance.runtime.ToValue(data))
		}
	}
	return goja.Undefined()
}

// fs.readFileSync(path[, options])
func (instance *FileSystem) readFileSync(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	if len(path) > 0 {
		data, err := qbc.IO.ReadTextFromFile(path)
		if nil != err {
			// throw back error to javascript
			panic(instance.runtime.NewTypeError(err.Error()))
		} else {
			return instance.runtime.ToValue(data)
		}
	}
	return goja.Undefined()
}

// fs.writeFileSync(file, data[, options])
func (instance *FileSystem) writeFileSync(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	data := call.Argument(1).Export()
	if len(path) > 0 && nil != data {
		err := fileutils.WriteDataToFile(path, data)
		if nil != err {
			// throw back error to javascript
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

// fs.writeFile(file, data[, options], callback)
func (instance *FileSystem) writeFile(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	data := call.Argument(1).Export()
	callback := commons.GetCallbackIfAny(call)
	if len(path) > 0 && nil != callback {
		err := fileutils.WriteDataToFile(path, data)
		if nil != err {
			_, _ = callback(call.This, instance.runtime.ToValue(err.Error()))
		} else {
			_, _ = callback(call.This, goja.Undefined())
		}
	}
	return goja.Undefined()
}

// fs.existsSync(path)
// Returns true if the path exists, false otherwise.
func (instance *FileSystem) existsSync(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	val, err := qbc.Paths.Exists(path)
	if nil != err {
		return instance.runtime.ToValue(false)
	}
	return instance.runtime.ToValue(val)
}

// fs.statSync(path[, options])
// A fs.Stats object provides information about a file.
// Objects returned from fs.stat(), fs.lstat() and fs.fstat() and their synchronous counterparts
// are of this type. If bigint in the options passed to those methods is true, the numeric values
// will be bigint instead of number, and the object will contain additional nanosecond-precision
// properties suffixed with Ns.
func (instance *FileSystem) statSync(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	val, err := instance.stats(path)
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return instance.runtime.ToValue(val)
}

// fs.stat(path[, options], callback)
func (instance *FileSystem) stat(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	callback := commons.GetCallbackIfAny(call)
	if len(path) > 0 && nil != callback {
		val, err := instance.stats(path)
		if nil != err {
			_, _ = callback(call.This, instance.runtime.ToValue(err.Error()), goja.Undefined())
		} else {
			_, _ = callback(call.This, goja.Undefined(), instance.runtime.ToValue(val))
		}
	}
	return goja.Undefined()
}

// fs.unlinkSync(path)
// removes a file or symbolic link.
func (instance *FileSystem) unlinkSync(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	err := qbc.IO.Remove(path)
	if nil != err {
		// throw back error to javascript
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return goja.Undefined()
}

// fs.unlink(path, callback)
// Asynchronously removes a file or symbolic link.
// No arguments other than a possible exception are given to the completion callback.
func (instance *FileSystem) unlink(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	callback := commons.GetCallbackIfAny(call)
	if len(path) > 0 && nil != callback {
		err := qbc.IO.Remove(path)
		if nil != err {
			_, _ = callback(call.This, instance.runtime.ToValue(err.Error()))
		} else {
			_, _ = callback(call.This, goja.Undefined())
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *FileSystem) stats(path string) (map[string]interface{}, error) {
	info, err := os.Stat(path)
	if nil != err {
		return nil, err
	}
	val := map[string]interface{}{
		"size": info.Size(),
		"name": info.Name(),
		"mode": info.Mode(),
		"isDirectory": func() bool {
			return info.IsDir()
		},
		"isSymbolicLink": func() bool {
			return info.Mode()&os.ModeSymlink != 0
		},
		"isFile": func() bool {
			return info.Mode().IsRegular()
		},
	}
	return val, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	fs := &FileSystem{
		runtime: runtime,
	}
	//c.util = require.Require(runtime, "util").(*goja.Object)

	if len(args) > 0 {
		root := qbc.Reflect.ValueOf(args[0]).String()
		if len(root) > 0 {
			fs.root = root
		}
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("readFile", fs.readFile)
	_ = o.Set("readFileSync", fs.readFileSync)
	_ = o.Set("writeFile", fs.writeFile)
	_ = o.Set("writeFileSync", fs.writeFileSync)
	_ = o.Set("existsSync", fs.existsSync)
	_ = o.Set("stat", fs.stat)
	_ = o.Set("statSync", fs.statSync)
	_ = o.Set("unlink", fs.unlink)
	_ = o.Set("unlinkSync", fs.unlinkSync)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})

	// ctx.Runtime.Set(NAME, require.Require(ctx.Runtime, NAME))
}

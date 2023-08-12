package http

import (
	"mime/multipart"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsFile struct {
	runtime *goja.Runtime
	object  *goja.Object
	group   string
	ctx     *fiber.Ctx
	file    *multipart.FileHeader
	mux     sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	JsFile
//----------------------------------------------------------------------------------------------------------------------

func WrapFile(runtime *goja.Runtime, group string, ctx *fiber.Ctx, file *multipart.FileHeader) *JsFile {
	instance := new(JsFile)
	instance.runtime = runtime
	instance.group = group
	instance.ctx = ctx
	instance.file = file

	instance.object = instance.runtime.NewObject()
	instance.export()

	return instance
}

func (instance *JsFile) Value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsFile) contentType(_ goja.FunctionCall) goja.Value {
	if nil != instance && nil != instance.file {
		return instance.runtime.ToValue(instance.file.Header["Content-Type"][0])
	}
	return goja.Undefined()
}

func (instance *JsFile) size(_ goja.FunctionCall) goja.Value {
	if nil != instance && nil != instance.file {
		return instance.runtime.ToValue(instance.file.Size)
	}
	return goja.Undefined()
}

func (instance *JsFile) name(_ goja.FunctionCall) goja.Value {
	if nil != instance && nil != instance.file {
		return instance.runtime.ToValue(instance.file.Filename)
	}
	return goja.Undefined()
}

func (instance *JsFile) save(call goja.FunctionCall) goja.Value {
	if nil != instance && nil != instance.file && nil != instance.ctx {
		var root string
		var level int
		var maxSize int64
		var returnFullPath bool
		switch len(call.Arguments) {
		case 1:
			root = commons.GetString(call, 0)
			level = 3 // yyyy/MM/dd
			maxSize = 0
			returnFullPath = false
		case 2:
			root = commons.GetString(call, 0)
			level = int(commons.GetInt(call, 1))
			maxSize = 0
			returnFullPath = false
		case 3:
			root = commons.GetString(call, 0)
			level = int(commons.GetInt(call, 1))
			maxSize = commons.GetInt(call, 2)
			returnFullPath = false
		case 4:
			root = commons.GetString(call, 0)
			level = int(commons.GetInt(call, 1))
			maxSize = commons.GetInt(call, 2)
			returnFullPath = commons.GetBool(call, 3)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if maxSize > 0 {
			size := instance.file.Size
			if size > maxSize {
				panic(instance.runtime.NewTypeError("file_size_exceed:max=" + qbc.Formatter.FormatBytes(uint64(maxSize)) +
					":actual=" + qbc.Formatter.FormatBytes(uint64(size))))
			}
		}
		filename := instance.file.Filename
		root = qbc.Paths.Absolute(root) // absolute
		ext := qbc.Paths.Extension(filename)
		name := qbc.Strings.Slugify(qbc.Paths.FileName(filename, false))
		dir := qbc.Paths.Dir(filename)
		filename = qbc.Paths.Concat(dir, name+"_"+qbc.Coding.MD5(qbc.Rnd.Uuid())+ext)
		path := strings.ToLower(qbc.Paths.DatePath(root, filename, level, true))
		err := instance.ctx.SaveFile(instance.file, path)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		if returnFullPath {
			return instance.runtime.ToValue(path)
		} else {
			return instance.runtime.ToValue(strings.Replace(path, root, ".", 1))
		}
	}

	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsFile) export() {
	o := instance.object
	// methods
	_ = o.Set("save", instance.save)
	_ = o.Set("size", instance.size)
	_ = o.Set("name", instance.name)
	_ = o.Set("contentType", instance.contentType)

}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

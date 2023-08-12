package commons

import (
	"github.com/dop251/goja"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	// reserved js variables
	ObjClosable = "_closable"

	// errors
	ErrorMissingParam = "missing_param"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type RuntimeContext struct {
	Uid       *string
	Workspace string
	Runtime   *goja.Runtime
	Arguments []interface{}
}

type ModuleLoader func(*goja.Runtime, *goja.Object, ...interface{})

type ModuleInfo struct {
	Context *RuntimeContext
	Loader  ModuleLoader
}

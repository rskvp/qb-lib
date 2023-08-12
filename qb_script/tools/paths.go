package tools

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_PATHS = "$paths"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolPaths struct {
	params  *ScriptingToolParams
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolPaths(params *ScriptingToolParams) *ScriptingToolPaths {
	result := new(ScriptingToolPaths)
	result.params = params
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolPaths) Init(params *ScriptingToolParams) {
	tool.params = params
	tool.runtime = params.Runtime
}

func (tool *ScriptingToolPaths) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Return absolute path
// @param [string] path
// @return string
func (tool *ScriptingToolPaths) GetAbsolutePath(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		path := tool.getArgsString(args)
		if len(path) > 0 {
			absolutePath := qbc.Paths.Absolute(path)
			return tool.runtime.ToValue(absolutePath)
		}
	}
	return tool.runtime.ToValue([]map[string]string{})
}

func (tool *ScriptingToolPaths) GetWorkspacePath(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		path := tool.getArgsString(args)
		if len(path) > 0 {
			absolutePath := qbc.Paths.WorkspacePath(path)
			return tool.runtime.ToValue(absolutePath)
		}
	}
	return tool.runtime.ToValue([]map[string]string{})
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolPaths) getArgsString(args []goja.Value) string {
	arg1 := ""

	if len(args) > 0 {
		switch len(args) {
		case 1:
			arg1 = qbc.Convert.ToString(args[0])
		default:
			arg1 = qbc.Convert.ToString(args[0])
		}
	}

	return arg1
}

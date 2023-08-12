package tools

import (
	"strings"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_STRINGS = "$strings"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolStrings struct {
	params  *ScriptingToolParams
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolStrings(params *ScriptingToolParams) *ScriptingToolStrings {
	result := new(ScriptingToolStrings)
	result.params = params
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolStrings) Init(params *ScriptingToolParams) {
	tool.params = params
	tool.runtime = params.Runtime
}

func (tool *ScriptingToolStrings) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Split a string by sep.
// Support multiple separators
// @param sep string
// @param text string (Optional) CONTEXT is used if not found
// @return []string
func (tool *ScriptingToolStrings) Split(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		sep, text := tool.getArgsStringString(args)
		if len(sep) > 0 && len(text) > 0 {
			tokens := qbc.Strings.Split(text, sep)
			return tool.runtime.ToValue(tokens)
		}
	}

	return tool.runtime.ToValue("")
}

// Get a substring
// @param start int Start index
// @param end int End index
// @param text string (Optional) CONTEXT is used if not found
// @return []string
func (tool *ScriptingToolStrings) Sub(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		start, end, text := tool.getArgsIntIntString(args)
		if len(text) > 0 {
			value := qbc.Strings.Sub(text, start, end)
			return tool.runtime.ToValue(value)
		}
	}

	return tool.runtime.ToValue("")
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m p o u n d
//----------------------------------------------------------------------------------------------------------------------

// Split a string by spaces AND get a word at index
// @param index int
// @param text string (Optional) CONTEXT is used if not found
// @return bool
func (tool *ScriptingToolStrings) SplitBySpaceWordAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		index, text := tool.getArgsIntString(args)
		if index > -1 && len(text) > 0 {
			tokens := qbc.Strings.Split(text, " \n")
			if len(tokens) > index {
				result := tokens[index]
				return tool.runtime.ToValue(result)
			}
		}
	}

	return tool.runtime.ToValue("")
}

// Split a string by "sep" AND get a word at index
// @param index int
// @param sep string
// @param text string (Optional) CONTEXT is used if not found
// @return string
func (tool *ScriptingToolStrings) SplitWordAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		index, sep, text := tool.getArgsIntStringString(args)
		if index > -1 && len(text) > 0 {
			tokens := strings.Split(text, sep)
			if len(tokens) > index {
				result := tokens[index]
				return tool.runtime.ToValue(result)
			}
		}
	}

	return tool.runtime.ToValue("")
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolStrings) getArgsStringString(args []goja.Value) (string, string) {
	var arg1, argCtx string
	arg1 = qbc.Convert.ToString(args[0].Export())
	if len(arg1) > 0 {
		if len(args) == 2 {
			argCtx = qbc.Convert.ToString(args[1].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = qbc.Convert.ToString(tool.context)
		}
	}

	return arg1, argCtx
}

func (tool *ScriptingToolStrings) getArgsIntString(args []goja.Value) (int, string) {
	var arg1 int
	var argCtx string

	arg1 = qbc.Convert.ToInt(args[0].Export())
	if arg1 > -1 {
		if len(args) == 2 {
			argCtx = qbc.Convert.ToString(args[1].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = qbc.Convert.ToString(tool.context)
		}
	}

	return arg1, argCtx
}

func (tool *ScriptingToolStrings) getArgsIntStringString(args []goja.Value) (int, string, string) {
	var arg1 int
	var arg2 string
	var argCtx string

	arg1 = qbc.Convert.ToInt(args[0].Export())
	if arg1 > -1 {

		if len(args) > 1 {
			arg2 = qbc.Convert.ToString(args[1].Export())
		}

		if len(args) == 3 {
			argCtx = qbc.Convert.ToString(args[2].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = qbc.Convert.ToString(tool.context)
		}
	}

	return arg1, arg2, argCtx
}

func (tool *ScriptingToolStrings) getArgsIntIntString(args []goja.Value) (int, int, string) {
	var arg1 int
	var arg2 int
	var argCtx string

	arg1 = qbc.Convert.ToInt(args[0].Export())
	if arg1 > -1 {

		if len(args) > 1 {
			arg2 = qbc.Convert.ToInt(args[1].Export())
		}

		if len(args) == 3 {
			argCtx = qbc.Convert.ToString(args[2].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = qbc.Convert.ToString(tool.context)
		}
	}

	return arg1, arg2, argCtx
}

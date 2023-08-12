package tools

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_num2word"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_CONVERT = "$convert"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolConvert struct {
	params  *ScriptingToolParams
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolConvert(params *ScriptingToolParams) *ScriptingToolConvert {
	result := new(ScriptingToolConvert)
	result.params = params
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolConvert) Init(params *ScriptingToolParams) {
	tool.params = params
	tool.runtime = params.Runtime
}

func (tool *ScriptingToolConvert) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	Number to Word
//----------------------------------------------------------------------------------------------------------------------

// Convert a number into word
// @param num int Number to convert into word
// @param lang string Language for conversion (ex: "it")
// @return string
func (tool *ScriptingToolConvert) Num2Word(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		num, lang := tool.getArgsIntString(args)
		if num > 0 && len(lang) > 0 {
			num2word := qb_num2word.Num2Word2
			num2word.Options.WordSeparator = ""
			value := num2word.Convert(num, lang)

			return tool.runtime.ToValue(value)
		}
	}

	return tool.runtime.ToValue("")
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolConvert) getArgsStringString(args []goja.Value) (string, string) {
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

func (tool *ScriptingToolConvert) getArgsIntString(args []goja.Value) (int, string) {
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

func (tool *ScriptingToolConvert) getArgsIntStringString(args []goja.Value) (int, string, string) {
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

func (tool *ScriptingToolConvert) getArgsIntIntString(args []goja.Value) (int, int, string) {
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

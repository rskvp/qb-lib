package tools

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_CSV = "$csv"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolCSV struct {
	params  *ScriptingToolParams
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolCSV(params *ScriptingToolParams) *ScriptingToolCSV {
	result := new(ScriptingToolCSV)
	result.params = params
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolCSV) Init(params *ScriptingToolParams) {
	tool.params = params
	tool.runtime = params.Runtime
}

func (tool *ScriptingToolCSV) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Load csv file from file and return []map[string]string
// @param [string] filename
// @param [bool] headerOnFirstRow
// @return []map[string]string
func (tool *ScriptingToolCSV) LoadFromFile(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		filename, headerOnFirstRow := tool.getArgsStringBool(args)
		if len(filename) > 0 {
			options := qbc.CSV.NewCsvOptions(";", "#", headerOnFirstRow)
			text, err := qbc.IO.ReadTextFromFile(filename)
			if nil==err{
				data, err := qbc.CSV.ReadAll(text, options)
				if nil == err {
					return tool.runtime.ToValue(data)
				} else {
					return tool.runtime.ToValue(err.Error())
				}
			} else {
				return tool.runtime.ToValue(err.Error())
			}
		}
	}
	return tool.runtime.ToValue([]map[string]string{})
}


//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolCSV) getArgsStringBool(args []goja.Value) (string, bool) {
	arg1 := ""
	arg2 := false

	if len(args) > 0 {
		switch len(args) {
		case 1:
			arg1 = qbc.Convert.ToString(args[0])
		case 2:
			arg1 = qbc.Convert.ToString(args[0])
			arg2 = qbc.Convert.ToBool(args[1])
		default:
			arg1 = qbc.Convert.ToString(args[0])
		}
	}

	return arg1, arg2
}

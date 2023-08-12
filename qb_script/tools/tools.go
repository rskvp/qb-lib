package tools

import "github.com/dop251/goja"

type ScriptingTool interface {
	Init(params *ScriptingToolParams)
	SetContext(context interface{})
}

type ScriptingToolParams struct {
	Root    *string // pointer to external string (the qb_sms_engine Root)
	Name    *string // pointer to external string (the qb_sms_engine Name)
	Runtime *goja.Runtime
}

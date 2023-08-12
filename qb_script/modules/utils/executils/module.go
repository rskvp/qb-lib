package executils

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

const NAME = "exec-utils"

type ExecUtils struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// exec.run(command)
func (instance *ExecUtils) run(call goja.FunctionCall) goja.Value {
	command := call.Argument(0).String()
	if len(command) > 0 {
		data, err := run(command)
		if nil == err {
			return instance.runtime.ToValue(data)
		} else {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	} else {
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return goja.Undefined()
}

func (instance *ExecUtils) open(call goja.FunctionCall) goja.Value {
	command := call.Argument(0).String()
	if len(command) > 0 {
		err := open(command)
		if nil == err {
			return instance.runtime.ToValue(true)
		} else {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	} else {
		panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func run(command string) (map[string]interface{}, error) {
	tokens := strings.Split(command, " ")
	cmd := tokens[0]
	params := make([]string, 0)
	if len(tokens) > 1 {
		params = tokens[1:]
	}
	executor := qbc.Exec.NewExecutor(cmd)
	err := executor.Run(params...)
	if nil != err {
		return nil, err
	}
	err = executor.Wait()
	if nil != err {
		return nil, err
	}
	return map[string]interface{}{
		"pid":        executor.PidLatest(),
		"lines":      executor.StdOutLines(),
		"error":      executor.StdErr(),
		"elapsed_ms": executor.Elapsed(),
	}, nil
}

func open(command string) error {
	tokens := strings.Split(command, " ")
	return qbc.Exec.Open(tokens...)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &ExecUtils{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)
	// file utility
	_ = o.Set("run", instance.run)
	_ = o.Set("open", instance.open)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

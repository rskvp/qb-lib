package qb_script

import (
	"strings"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

var _js *goja.Runtime

type OnReadyCallback func(rtContext *commons.RuntimeContext)
type EnvSettings struct {
	Root           string
	EngineName     string
	FileName       string
	ProgramName    string
	ProgramScript  string
	Context        map[string]interface{}
	Logger         qb_log.ILogger
	LoggerFileName string
	LogReset       bool
	OnReady        OnReadyCallback
}

type ScriptingHelper struct {
	logger qb_log.ILogger
}

var Scripting *ScriptingHelper

func init() {
	Scripting = new(ScriptingHelper)
}

func (instance *ScriptingHelper) SetLogger(logger qb_log.ILogger) {
	instance.logger = logger
}

func (instance *ScriptingHelper) NewEngine(name ...interface{}) *ScriptEngine {
	return NewJsEngine()
}

func (instance *ScriptingHelper) ValidateEnv(env *EnvSettings) error {
	// check env validity
	if len(env.EngineName) == 0 {
		env.EngineName = "js"
	}
	if nil == env.Logger {
		env.Logger = instance.logger
	}

	if len(env.FileName) > 0 {
		env.FileName = qbc.Paths.Absolute(env.FileName)
		text, ioErr := qbc.IO.ReadTextFromFile(env.FileName)
		if nil != ioErr {
			return ioErr
		}
		env.ProgramScript = text
		env.ProgramName = qbc.Paths.FileName(env.FileName, false)
		env.EngineName = qbc.Paths.ExtensionName(env.FileName)
		env.Root = qbc.Paths.Dir(env.FileName)
	}

	if len(env.Root) == 0 {
		env.Root = qbc.Paths.WorkspacePath("./")
	}
	if len(env.ProgramName) == 0 {
		env.ProgramName = qbc.Rnd.RndId()
	}
	return nil
}

func (instance *ScriptingHelper) RunFile(filename string, context map[string]interface{}, args ...interface{}) (response string, err error) {
	env := new(EnvSettings)
	env.FileName = filename
	env.Context = context

	parseArgs(env, args...)

	return instance.Run(env)
}

func (instance *ScriptingHelper) RunProgram(programName string, script string, context map[string]interface{}, args ...interface{}) (response string, err error) {
	env := new(EnvSettings)
	env.EngineName = "js"
	env.ProgramName = programName
	env.ProgramScript = script
	env.Context = context

	parseArgs(env, args...)

	return instance.Run(env)
}

func (instance *ScriptingHelper) Run(env *EnvSettings) (response string, err error) {
	// check env validity
	err = instance.ValidateEnv(env)
	if nil != err {
		return "", err
	}

	// create engine
	vm := instance.NewEngine(env.EngineName)
	vm.Name = env.ProgramName
	// logger
	if nil != env.Logger {
		vm.SetLogger(env.Logger)
		vm.AllowMultipleLoggers = false
	} else {
		if len(env.LoggerFileName) > 0 {
			vm.LogFile = env.LoggerFileName
		} else {
			vm.LogFile = qbc.Paths.Concat(env.Root, env.ProgramName+".log")
		}
	}
	vm.ResetLogOnEachRun = env.LogReset
	// context
	if nil != env.Context {
		for k, v := range env.Context {
			// add prefix
			if strings.Contains(k, "$") == false {
				k = "$" + k
			}
			vm.Set(k, v)
		}
	}

	// registry
	registry := NewModuleRegistry(Loader)
	rtContext := registry.Start(vm)
	if nil != rtContext && nil != env.OnReady {
		// here we can add modules to runtime context
		env.OnReady(rtContext) // usage: "console.Enable(rtContext)"
	}

	v, e := vm.RunProgram(env.ProgramName, env.ProgramScript)
	if nil != e {
		err = e
	} else {
		response = qb_utils.Convert.ToString(v.Export())
	}
	vm.Close()
	return
}

func (instance *ScriptingHelper) EvalJs(expression string, useGlobalRuntime, overwriteContext bool, context ...map[string]interface{}) (interface{}, error) {
	if len(expression) > 0 {
		var runtime *goja.Runtime
		if useGlobalRuntime {
			runtime = js()

		} else {
			runtime = goja.New()
		}
		for _, ctx := range context {
			for k, v := range ctx {
				if overwriteContext || nil == runtime.Get(k) {
					_ = runtime.Set(k, v)
				}
			}
		}
		value, err := runtime.RunString(expression)
		if nil != err {
			return nil, err
		}
		return value.Export(), err
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// Loader default javascript Module registry loader
func Loader(path string) ([]byte, error) {
	path = qbc.Paths.Concat(qbc.Paths.WorkspacePath("modules"), path)
	return qbc.IO.ReadBytesFromFile(path)
}

func parseArgs(env *EnvSettings, args ...interface{}) {
	for _, arg := range args {
		if engine, ok := arg.(string); ok {
			env.EngineName = engine
		} else if callback, ok := arg.(OnReadyCallback); ok {
			env.OnReady = callback
		}
	}
}

func js() *goja.Runtime {
	if nil == _js {
		_js = goja.New()
	}
	return _js
}

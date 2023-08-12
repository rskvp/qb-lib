package qb_script

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/tools"
)

var (
	errorNotInitialized = errors.New("engine_not_initialized")
)

func NewJsEngine() *ScriptEngine {
	instance := new(ScriptEngine)
	instance.Name = qbc.Rnd.Uuid()
	instance.Root = qbc.Paths.Absolute("./")
	instance.LogLevel = qb_log.InfoLevel
	instance.ResetLogOnEachRun = false
	instance.AllowMultipleLoggers = true

	// init runtime
	instance.runtime = goja.New()
	instance.initialized = false

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptEngine struct {
	Root                 string
	Name                 string
	Silent               bool // if enabled do not log at console output but only on files
	LogLevel             qb_log.Level
	LogFile              string
	ResetLogOnEachRun    bool
	AllowMultipleLoggers bool

	//-- private --//
	initialized bool
	logger      qb_log.ILogger
	runtime     *goja.Runtime
	loop        *commons.EventLoop
	tools       []tools.ScriptingTool
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScriptEngine) Open() {
	if nil != instance {
		instance.initialize()
	}
}

func (instance *ScriptEngine) Close() {
	if nil != instance {
		instance.finalize()
	}
}

func (instance *ScriptEngine) Runtime() *goja.Runtime {
	if nil != instance && nil != instance.runtime {
		return instance.runtime
	}
	return nil
}

func (instance *ScriptEngine) GetLogger() qb_log.ILogger {
	if nil != instance {
		return instance.getLogger()
	}
	return nil
}

func (instance *ScriptEngine) SetLogger(logger qb_log.ILogger) {
	if nil != instance {
		instance.logger = logger
	}
}

func (instance *ScriptEngine) Set(name string, value interface{}) {
	if nil != instance && nil != instance.runtime {
		instance.runtime.Set(name, value)
	}
}

func (instance *ScriptEngine) Get(name string) goja.Value {
	if nil != instance && nil != instance.runtime {
		return instance.runtime.Get(name)
	}
	return nil
}

func (instance *ScriptEngine) ToValue(value interface{}) goja.Value {
	if nil != instance && nil != instance.runtime {
		return instance.runtime.ToValue(value)
	}
	return nil
}

func (instance *ScriptEngine) ToObject(value goja.Value) *goja.Object {
	if nil != instance && nil != instance.runtime {
		if nil != value {
			return value.ToObject(instance.runtime)
		}
	}
	return nil
}

func (instance *ScriptEngine) HasMember(value goja.Value, fieldName string) bool {
	if nil != instance && nil != instance.runtime {
		if nil != value {
			obj := value.ToObject(instance.runtime)
			if nil != obj {
				return obj.Get(fieldName) != nil
			}
		}
	}
	return false
}

func (instance *ScriptEngine) GetMember(value goja.Value, fieldName string) goja.Value {
	if nil != instance && nil != instance.runtime {
		if nil != value {
			obj := value.ToObject(instance.runtime)
			if nil != obj {
				return obj.Get(fieldName)
			}
		}
	}
	return nil
}

func (instance *ScriptEngine) GetMemberAsCallable(value goja.Value, fieldName string) goja.Callable {
	if nil != instance && nil != instance.runtime {
		if nil != value && value != goja.Undefined() {
			obj := value.ToObject(instance.runtime)
			if nil != obj {
				value := obj.Get(fieldName)
				if nil != value {
					callable, _ := goja.AssertFunction(value)
					return callable
				}
			}
		}
	}
	return nil
}

func (instance *ScriptEngine) NewObject() *goja.Object {
	if nil != instance && nil != instance.runtime {
		return instance.runtime.NewObject()
	}
	return nil
}

func (instance *ScriptEngine) NewTypeError(args ...interface{}) *goja.Object {
	if nil != instance && nil != instance.runtime {
		return instance.runtime.NewTypeError(args...)
	}
	return nil
}

func (instance *ScriptEngine) NewArrayBuffer(data []byte) goja.ArrayBuffer {
	if nil != instance && nil != instance.runtime {
		return instance.runtime.NewArrayBuffer(data)
	}
	return goja.ArrayBuffer{}
}

func (instance *ScriptEngine) RunProgram(name, program string) (goja.Value, error) {
	if nil != instance && instance.initialized {
		v, err := instance.runtime.RunScript(name, program)
		return processResponse(instance.runtime, v, err)
	}
	return nil, errorNotInitialized
}

func (instance *ScriptEngine) RunScript(program string) (goja.Value, error) {
	if nil != instance && instance.initialized {
		v, err := instance.runtime.RunScript(instance.Name, program)
		return processResponse(instance.runtime, v, err)
	}
	return nil, errorNotInitialized
}

//----------------------------------------------------------------------------------------------------------------------
//	t o o l s
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScriptEngine) AddTool(name string, tool tools.ScriptingTool) error {
	if nil != instance && instance.initialized {
		params := instance.getParams()
		tool.Init(params)
		instance.runtime.Set(name, tool)
		instance.tools = append(instance.tools, tool)
	}
	return errorNotInitialized
}

//----------------------------------------------------------------------------------------------------------------------
//	T O O L S    c o n t e x t
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScriptEngine) SetToolsContext(value interface{}) {
	for _, tool := range instance.tools {
		tool.SetContext(value)
	}
}

func (instance *ScriptEngine) SetToolContext(toolName string, value interface{}) {
	v := instance.runtime.Get(toolName)
	if nil != v {
		tool := v.Export().(tools.ScriptingTool)
		if nil != tool {
			tool.SetContext(value)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	W I N D O W    c o n t e x t
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScriptEngine) SetRootContextValue(name string, value interface{}) error {
	if nil != instance && nil != instance.runtime {
		if c := getRootContext(instance.runtime); nil != c {
			return c.Set(name, value)
		}
	}
	return nil
}

func (instance *ScriptEngine) GetRootContextValue(name string) goja.Value {
	if nil != instance && nil != instance.runtime {
		if c := getRootContext(instance.runtime); nil != c {
			return c.Get(name)
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScriptEngine) initialize() {
	if !instance.initialized {
		instance.initialized = true

		if len(instance.LogFile) == 0 {
			instance.LogFile = qbc.Paths.Concat(instance.Root, "./logging/script-logging.log")
			if len(instance.Name) > 0 {
				name := qbc.Paths.FileName(instance.Name, false)
				instance.LogFile = qbc.Paths.Concat(instance.Root, fmt.Sprintf("./logging/%s.log", name))
			}
		}

		// internal variables (__root)
		instance.initVariables()

		// internal tools initialization
		instance.initTools()

		// add support to event loop
		instance.loop = commons.NewEventLoop(instance.runtime)
		instance.loop.Start()
	}
}

func (instance *ScriptEngine) finalize() {
	if instance.initialized {
		instance.initialized = false
		if nil != instance.loop {
			instance.loop.Stop()
		}
		instance.loop = nil

		// closable objects must be exposed to avoid locks
		commons.CloseObjects(instance.runtime)
		instance.runtime.Interrupt(nil)
		instance.runtime.ClearInterrupt()
	}
}

func (instance *ScriptEngine) getLogger() qb_log.ILogger {
	if nil != instance && !instance.AllowMultipleLoggers {
		if nil == instance.logger {
			logger := qb_log.NewLogger()
			logger.SetFilename(instance.LogFile)
			logger.SetLevel("debug")
			instance.logger = logger
		}
		return instance.logger
	}
	return nil
}

func (instance *ScriptEngine) getParams() *tools.ScriptingToolParams {
	params := new(tools.ScriptingToolParams)
	params.Runtime = instance.runtime
	params.Root = &instance.Root
	params.Name = &instance.Name

	return params
}

func (instance *ScriptEngine) initVariables() {
	//instance.runtime.Set("__root", instance.Root)
	commons.SetRtRoot(instance.runtime, instance.Root)
}

func (instance *ScriptEngine) initTools() {

	params := instance.getParams()

	Tconsole := tools.NewToolConsole(params)
	instance.runtime.Set(tools.TOOL_CONSOLE, Tconsole)
	instance.tools = append(instance.tools, Tconsole)

	Tstrings := tools.NewToolStrings(params)
	instance.runtime.Set(tools.TOOL_STRINGS, Tstrings)
	instance.tools = append(instance.tools, Tstrings)

	Tarrays := tools.NewToolArrays(params)
	instance.runtime.Set(tools.TOOL_ARRAYS, Tarrays)
	instance.tools = append(instance.tools, Tarrays)

	Tregexps := tools.NewToolRegExps(params)
	instance.runtime.Set(tools.TOOL_REGEXPS, Tregexps)
	instance.tools = append(instance.tools, Tregexps)

	Tconvert := tools.NewToolConvert(params)
	instance.runtime.Set(tools.TOOL_CONVERT, Tconvert)
	instance.tools = append(instance.tools, Tconvert)

	Tcsv := tools.NewToolCSV(params)
	instance.runtime.Set(tools.TOOL_CSV, Tcsv)
	instance.tools = append(instance.tools, Tcsv)

	Tpaths := tools.NewToolPaths(params)
	instance.runtime.Set(tools.TOOL_PATHS, Tpaths)
	instance.tools = append(instance.tools, Tcsv)

}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// Root Context is an object contained in "window".
// "window" is an emulation of browser window object
func getRootContextValue(runtime *goja.Runtime, name string) goja.Value {
	return commons.GetRtDeepValue(runtime, "runtime.context."+name)
}

func getRootContext(runtime *goja.Runtime) *goja.Object {
	if window, b := runtime.Get("runtime").(*goja.Object); b {
		if context, b := window.Get("context").(*goja.Object); b {
			return context
		}
	}
	return nil
}

func processResponse(runtime *goja.Runtime, v goja.Value, e error) (goja.Value, error) {
	if nil != e {
		return v, e
	} else {
		if v := getRootContextValue(runtime, "return"); nil != v {
			return v, e
		} else if v := getRootContextValue(runtime, "result"); nil != v {
			return v, e
		} else if v := getRootContextValue(runtime, "response"); nil != v {
			return v, e
		}
	}
	return v, e
}

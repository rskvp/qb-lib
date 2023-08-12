package console

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "console"

type Console struct {
	runtime  *goja.Runtime
	util     *goja.Object
	name     string
	filename string
	logger   qb_log.ILogger
	silent   bool // enable/disable console output
	mux      sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *Console) close(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.logger.Debug("-> javascript console closed.")
	}
	return goja.Undefined()
}

func (instance *Console) reset(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		_ = qbc.IO.Remove(instance.filename)
		instance.logger.Debug("-> javascript console reset.")
	}
	return goja.Undefined()
}

func (instance *Console) log(call goja.FunctionCall) goja.Value {
	if message, err := instance.format(call); nil == err {
		instance.write(qb_log.InfoLevel, message)
	} else {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return nil
}

func (instance *Console) error(call goja.FunctionCall) goja.Value {
	if message, err := instance.format(call); nil == err {
		instance.write(qb_log.ErrorLevel, message)
	} else {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return nil
}

func (instance *Console) warn(call goja.FunctionCall) goja.Value {
	if message, err := instance.format(call); nil == err {
		instance.write(qb_log.WarnLevel, message)
	} else {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return nil
}

func (instance *Console) info(call goja.FunctionCall) goja.Value {
	if message, err := instance.format(call); nil == err {
		instance.write(qb_log.InfoLevel, message)
	} else {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return nil
}

func (instance *Console) debug(call goja.FunctionCall) goja.Value {
	if message, err := instance.format(call); nil == err {
		instance.write(qb_log.DebugLevel, message)
	} else {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Console) format(call goja.FunctionCall) (string, error) {
	if format, ok := goja.AssertFunction(instance.util.Get("format")); ok {
		ret, err := format(instance.util, call.Arguments...)
		if err != nil {
			return "", err
		}
		message := ret.String()

		return message, nil
	} else {
		return "", errors.New("util.format is not a function")
	}
}

func (instance *Console) write(level qb_log.Level, message string) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	// PANIC RECOVERY
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			message := qbc.Strings.Format("[panic] MODULE %s ERROR: %s", NAME, r)
			fmt.Println(message)
		}
	}()

	if nil != instance.logger {
		switch level {
		case qb_log.WarnLevel:
			instance.logger.Warn(message)
		case qb_log.InfoLevel:
			instance.logger.Info(message)
		case qb_log.ErrorLevel:
			instance.logger.Error(message)
		case qb_log.DebugLevel:
			instance.logger.Debug(message)
		case qb_log.TraceLevel:
			instance.logger.Trace(message)
		case qb_log.PanicLevel:
			instance.logger.Panic(message)
		}
	} else {
		if len(instance.filename) > 0 {
			if b, _ := qbc.Paths.Exists(instance.filename); !b {
				_ = qbc.Paths.Mkdir(instance.filename)
			}
			_, _ = qbc.IO.AppendTextToFile(message+"\n", instance.filename)
		}
	}

	if !instance.silent {
		fmt.Println("["+instance.name+"] ", message)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	instance := &Console{
		runtime: runtime,
	}
	instance.util = require.Require(runtime, "util").(*goja.Object)

	if len(args) > 3 {
		root := qbc.Reflect.ValueOf(args[0]).String()
		name := qbc.Reflect.ValueOf(args[1]).String()
		silent := qbc.Reflect.ValueOf(args[2]).Bool()
		level := qbc.Reflect.ValueOf(args[3]).Uint()
		resetLog := qbc.Reflect.ValueOf(args[4]).Bool()
		getEngine := args[5]
		filename := qbc.Reflect.ValueOf(args[6]).String()

		if len(name) > 0 && len(root) > 0 {
			instance.name = name
			instance.silent = silent
			if f, b := getEngine.(func() qb_log.ILogger); b {
				instance.logger = f()
				if l, b := instance.logger.(*qb_log.Logger); b {
					instance.filename = l.GetFilename()
				}
			}
			if nil == instance.logger {
				if len(filename) > 0 {
					instance.filename = filename
				} else {
					instance.filename = qbc.Paths.Concat(qbc.Paths.Absolute(root), "logging", name+".log")
				}
				// reset log
				if resetLog {
					_ = qbc.IO.Remove(instance.filename)
				}
				logger := qb_log.NewLogger()
				logger.SetFilename(instance.filename)
				logger.SetLevel(level)
				instance.logger = logger
			}
		}
	}

	o := module.Get("exports").(*goja.Object)

	//_ = o.Set("open", instance.open)
	_ = o.Set("close", instance.close) // close console log file
	_ = o.Set("reset", instance.reset) // remove console file

	_ = o.Set("log", instance.log)
	_ = o.Set("error", instance.error)
	_ = o.Set("warn", instance.warn)
	_ = o.Set("debug", instance.debug)
	_ = o.Set("info", instance.info)

	commons.AddClosableObject(instance.runtime, o)
}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})

	// add module to javascript context
	_ = ctx.Runtime.Set(NAME, require.Require(ctx.Runtime, NAME))
}

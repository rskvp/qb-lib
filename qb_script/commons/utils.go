package commons

import (
	"strings"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	runtime
//----------------------------------------------------------------------------------------------------------------------

func GetRtDeepValue(runtime *goja.Runtime, name string) goja.Value {
	keys := strings.Split(name, ".")
	size := len(keys)
	if size > 1 {
		var obj *goja.Object
		for i := 0; i < size-1; i++ {
			if nil != obj {
				obj = GetObject(obj, keys[i]) // obj.Get(keys[i]).(*goja.Object)
			} else {
				obj = GetObject(runtime, keys[i]) // runtime.Get(keys[i]).(*goja.Object)
			}
		}
		if nil != obj {
			return obj.Get(keys[size-1])
		}
	} else {
		return runtime.Get(name)
	}
	return nil
}

func GetRtDeepObject(runtime *goja.Runtime, name string) *goja.Object {
	return AsObject(GetRtDeepValue(runtime, name))
}

func GetRtDeepExport(runtime *goja.Runtime, name string) interface{} {
	return AsExport(GetRtDeepValue(runtime, name))
}

func GetRtObject(runtime *goja.Runtime, name string) interface{} {
	if nil != runtime {
		value := runtime.Get(name)
		if nil != value {
			return value.Export()
		}
	}
	return nil
}

func GetRtRoot(runtime *goja.Runtime) string {
	if nil != runtime {
		value := GetRtObject(runtime, "__root")
		if nil != value {
			if v, b := value.(string); b {
				return v
			}
		}
	}
	return qbc.Paths.Absolute(".")
}

func SetRtRoot(runtime *goja.Runtime, value string) {
	if nil != runtime {
		runtime.Set("__root", value)
	}
}

func AddClosableObject(runtime *goja.Runtime, object *goja.Object) {
	if nil != runtime && nil != object {
		array := make([]*goja.Object, 0)
		a := runtime.Get(ObjClosable)
		if nil != a {
			if v, b := a.Export().([]*goja.Object); b {
				array = v
			}
		}
		array = append(array, object)
		runtime.Set(ObjClosable, array)
	}

}

func CloseObjects(runtime *goja.Runtime) {
	if nil != runtime {
		a := runtime.Get(ObjClosable)
		if nil != a {
			if array, b := a.Export().([]*goja.Object); b {
				for _, o := range array {
					if nil != o {
						// retrieve close function from closable object
						function, _ := goja.AssertFunction(o.Get("close"))
						if nil != function {
							_, _ = function(nil)
						}
					}
				}
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	call
//----------------------------------------------------------------------------------------------------------------------

func GetCallbackIfAny(call goja.FunctionCall) goja.Callable {
	index := len(call.Arguments) - 1
	if index > -1 {
		callback, _ := goja.AssertFunction(call.Argument(index))
		return callback
	}
	return nil
}

func GetString(call goja.FunctionCall, index int) string {
	max := len(call.Arguments) - 1
	if index <= max {
		return call.Argument(index).String()
	}
	return ""
}

func GetMap(call goja.FunctionCall, index int) map[string]interface{} {
	max := len(call.Arguments) - 1
	if index <= max {
		exp := call.Argument(index).Export()
		if v, b := exp.(map[string]interface{}); b {
			return v
		}
	}
	return nil
}

func GetExport(call goja.FunctionCall, index int) interface{} {
	max := len(call.Arguments) - 1
	if index <= max {
		return call.Argument(index).Export()
	}
	return nil
}

func GetBool(call goja.FunctionCall, index int) bool {
	max := len(call.Arguments) - 1
	if index <= max {
		return call.Argument(index).ToBoolean()
	}
	return false
}

func GetArray(call goja.FunctionCall, index int) []interface{} {
	max := len(call.Arguments) - 1
	if index <= max {
		return ToArray(call.Argument(index))
	}
	return nil
}

func GetArrayOfString(call goja.FunctionCall, index int) []string {
	response := make([]string, 0)
	intArray := GetArray(call, index)
	for _, item := range intArray {
		response = append(response, qbc.Convert.ToString(item))
	}
	return response
}

func GetInt(call goja.FunctionCall, index int) int64 {
	max := len(call.Arguments) - 1
	if index <= max {
		return call.Argument(index).ToInteger()
	}
	return 0
}

//----------------------------------------------------------------------------------------------------------------------
//	v a l u e s
//----------------------------------------------------------------------------------------------------------------------

func AsExport(value goja.Value) interface{} {
	if nil != value {
		return value.Export()
	}
	return nil
}

func AsObject(value goja.Value) *goja.Object {
	if nil != value {
		if v, b := value.(*goja.Object); b {
			return v
		}
	}
	return nil
}

func GetObject(obj interface{}, name string) *goja.Object {
	if nil != obj {
		if o, b := obj.(*goja.Object); b {
			return AsObject(o.Get(name))
		} else if r, b := obj.(*goja.Runtime); b {
			return AsObject(r.Get(name))
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n v e r t
//----------------------------------------------------------------------------------------------------------------------

func ToArrayOfMap(val interface{}) []map[string]interface{} {
	response := make([]map[string]interface{}, 0)
	if v, b := val.([]interface{}); b {
		for _, item := range v {
			if vv, b := item.(map[string]interface{}); b {
				response = append(response, vv)
			}
		}
	} else if v, b := val.(map[string]interface{}); b {
		response = append(response, v)
	}
	return response
}

func ToArray(v interface{}) []interface{} {
	response := make([]interface{}, 0)
	if args, b := v.([]interface{}); b {
		response = append(response, args...)
	} else if args, b := v.([]goja.Value); b {
		for _, val := range args {
			response = append(response, val.Export())
		}
	} else if args, b := v.([]goja.Object); b {
		for _, val := range args {
			response = append(response, val.Export())
		}
	} else if val, b := v.(goja.Value); b {
		exp := val.Export()
		response = append(response, ToArray(exp)...)
	} else if val, b := v.(goja.Object); b {
		response = append(response, val.Export())
	} else if args, b := v.(map[string]interface{}); b {
		for _, val := range args {
			response = append(response, val)
		}
	}

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	context
//----------------------------------------------------------------------------------------------------------------------

func GetArgsString(context interface{}, args []goja.Value) string {
	arg1 := ""

	switch len(args) {
	case 1:
		arg1 = qbc.Convert.ToString(args[0].Export())
	default:
		if nil != context {
			arg1 = qbc.Convert.ToString(context)
		}
	}
	return arg1
}

func GetArgsStringString(context interface{}, args []goja.Value) (string, string) {
	arg1 := ""
	arg2 := ""

	switch len(args) {
	case 1:
		arg1 = qbc.Convert.ToString(args[0].Export())
		// fallback on context for latest arg
		if nil != context {
			arg2 = qbc.Convert.ToString(context)
		}
	case 2:
		arg1 = qbc.Convert.ToString(args[0].Export())
		arg2 = qbc.Convert.ToString(args[1].Export())
	default:
		if nil != context {
			arg1 = qbc.Convert.ToString(context)
		}
	}

	return arg1, arg2
}

func GetArgsStringStringString(context interface{}, args []goja.Value) (string, string, string) {
	arg1 := ""
	arg2 := ""
	arg3 := ""

	switch len(args) {
	case 1:
		arg1 = qbc.Convert.ToString(args[0].Export())
		// fallback on context for latest arg
		if nil != context {
			arg2 = qbc.Convert.ToString(context)
		}
	case 2:
		arg1 = qbc.Convert.ToString(args[0].Export())
		arg2 = qbc.Convert.ToString(args[1].Export())
		// fallback on context for latest arg
		if nil != context {
			arg3 = qbc.Convert.ToString(context)
		}
	case 3:
		arg1 = qbc.Convert.ToString(args[0].Export())
		arg2 = qbc.Convert.ToString(args[1].Export())
		arg3 = qbc.Convert.ToString(args[2].Export())
	default:
		if nil != context {
			arg1 = qbc.Convert.ToString(context)
		}
	}

	return arg1, arg2, arg3
}

func GetArgsIntStringString(context interface{}, args []goja.Value) (int, string, string) {
	arg1 := 0
	arg2 := ""
	arg3 := ""

	switch len(args) {
	case 1:
		arg1 = qbc.Convert.ToInt(args[0].Export())
		// fallback on context for latest arg
		if nil != context {
			arg2 = qbc.Convert.ToString(context)
		}
	case 2:
		arg1 = qbc.Convert.ToInt(args[0].Export())
		arg2 = qbc.Convert.ToString(args[1].Export())
		// fallback on context for latest arg
		if nil != context {
			arg3 = qbc.Convert.ToString(context)
		}
	case 3:
		arg1 = qbc.Convert.ToInt(args[0].Export())
		arg2 = qbc.Convert.ToString(args[1].Export())
		arg3 = qbc.Convert.ToString(args[2].Export())
	default:
		if nil != context {
			arg1 = qbc.Convert.ToInt(context)
		}
	}

	return arg1, arg2, arg3
}

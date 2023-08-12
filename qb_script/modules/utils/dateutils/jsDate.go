package dateutils

import (
	"time"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsDate struct {
	runtime  *goja.Runtime
	original interface{}
	date     time.Time
}

//----------------------------------------------------------------------------------------------------------------------
//	JsCollection
//----------------------------------------------------------------------------------------------------------------------

func WrapDate(runtime *goja.Runtime, original interface{}) goja.Value {
	instance := new(JsDate)
	instance.runtime = runtime
	instance.original = original

	object := instance.runtime.NewObject()
	exportFields(instance, object)

	instance.init(instance.original)

	return object
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// format Format a date
// @param pattern:string Optional parameter
// @usage dt.format([pattern_or_null])
func (instance *JsDate) format(call goja.FunctionCall) goja.Value {
	if nil != instance {
		response := "" // empty response
		template := qbc.Convert.ToString(call.Argument(0).Export())
		if len(template) == 0 {
			response = instance.date.String()
		} else {
			response = qbc.Dates.FormatDate(instance.date, template)
		}

		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}


//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsDate) init(source interface{}) {
	instance.date = time.Now()
	if nil != source {
		if v, b := source.(time.Time); b {
			instance.date = v
		} else if v, b := source.(string); b {
			// try parse string date
			if t, err := qbc.Dates.TryParseAny(v); nil == err {
				instance.date = t
			}
		}
	}
}


//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func exportFields(instance *JsDate, o *goja.Object) {
	_ = o.Set("format", instance.format)

}

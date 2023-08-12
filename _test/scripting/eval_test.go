package _test

import (
	"fmt"
	"testing"

	ggx "github.com/rskvp/qb-lib"
)

func Test_Eval(t *testing.T) {
	expression := "A=2;A==1"
	context := map[string]interface{}{
		"A": 1,
	}
	response, err := ggx.Scripting.EvalJs(expression, true, false, context)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(expression, response)

	expression = "A==2"
	response, err = ggx.Scripting.EvalJs(expression, true, false, context)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(expression, response)
}

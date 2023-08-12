package dbsql

import (
	"fmt"
	"log"
	"testing"
)

func TestConvert(t *testing.T) {
	values := map[string]interface{}{
		"1":    "float64",
		"0001": "string",
		"":     "string",
		"a1":   "string",
		"1.1":  "float64",
	}
	for k, v := range values {
		c := convert(k)
		got := fmt.Sprintf("%T", c)
		if v != got {
			t.Error(fmt.Sprintf("Expected '%v', got '%v': %v", v, got, k))
			t.FailNow()
		} else {
			log.Println(got, c)
		}
	}
}

func TestPutValue(t *testing.T) {
	m := map[string]interface{}{}
	err := putValue(1, m, "", nil)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	log.Println(m)
}

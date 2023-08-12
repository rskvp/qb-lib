package dbal

import (
	"fmt"
	"testing"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/qb_dbal_drivers"
)

func TestDriverSqlite(t *testing.T) {
	var m map[string]string
	err := qbc.JSON.ReadFromFile("./sql_sqlite.json", &m)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	driver, err := qb_dbal_drivers.NewDatabase(m["driver"], m["dsn"])
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	if nil == driver {
		t.Error("Driver not found")
		t.FailNow()
	}

	item, err := driver.Upsert("test", map[string]interface{}{"name": "Mariolino"})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(item)

	err = driver.ForEach("test", func(doc map[string]interface{}) bool {
		fmt.Println(qbc.JSON.Stringify(doc))
		id := qbc.Convert.ToInt(doc["id"])
		if id > 4 {
			// remove
			fmt.Println("\tREMOVING: ", qbc.JSON.Stringify(doc))
			err = driver.Remove("test", qbc.Convert.ToString(id))
			if nil != err {
				t.Error(err)
				t.FailNow()
			}
		} else {
			doc["surname"] = "Geminiani"
			item, err = driver.Upsert("test", doc)
			if nil != err {
				t.Error(err)
				t.FailNow()
			}
			fmt.Println("\tUPDATED: ", qbc.JSON.Stringify(item))
		}
		return false // continue loop
	})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

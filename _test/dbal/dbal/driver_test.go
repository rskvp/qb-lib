package _test

import (
	"fmt"
	"testing"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/qb_dbal_drivers"
	// _ "github.com/alexbrainman/odbc"
)

func TestDriver(t *testing.T) {
	var m map[string]string
	err := qbc.JSON.ReadFromFile("./dsn.json", &m)
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
	err = driver.ForEach("cache", func(doc map[string]interface{}) bool {
		fmt.Println(doc)
		return false // continue loop
	})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	commands := make([]string, 0)
	bindVars := make([]map[string]interface{}, 0)
	commands = append(commands, "FOR u IN users\n  UPDATE u._key WITH { name: CONCAT(u.firstName, \" \", u.lastName) } IN users RETURN u")
	bindVars = append(bindVars, map[string]interface{}{})
	data, err := driver.ExecMultiple(commands, bindVars,
		map[string]interface{}{"read": []string{}, "write": []string{"users"}, "exclusive": []string{}})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(data)
}

func TestDriverSQL(t *testing.T) {
	var m map[string]string
	err := qbc.JSON.ReadFromFile("./sql_mysql.json", &m)
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
	err = driver.ForEach("table1", func(doc map[string]interface{}) bool {
		fmt.Println(doc)
		return false // continue loop
	})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

func TestDriverSQLGet(t *testing.T) {
	var m map[string]string
	err := qbc.JSON.ReadFromFile("./sql_mysql.json", &m)
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
	data, err := driver.Get("table1", "1")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(qbc.JSON.Stringify(data))
}

// mysql.server start
func TestDriverSQLExec(t *testing.T) {
	var m map[string]string
	err := qbc.JSON.ReadFromFile("./sql_mysql.json", &m)
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
	data, err := driver.ExecNative("SELECT COUNT(DISTINCT id) FROM table1", nil)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(qbc.JSON.Stringify(data))

	data, err = driver.ExecNative("SELECT t.* FROM table1 t WHERE id=@mykey",
		map[string]interface{}{
			"mykey": 1,
		})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(qbc.JSON.Stringify(data))
}

func TestDriverODBCExec(t *testing.T) {
	var m map[string]string
	err := qbc.JSON.ReadFromFile("./sql_odbc.json", &m)
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
	data, err := driver.ExecNative("SELECT COUNT(DISTINCT codice) FROM dbo.banche", nil)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(qbc.JSON.Stringify(data))

	data, err = driver.ExecNative("SELECT t.* FROM dbo.cafliexc t WHERE id=@mykey",
		map[string]interface{}{
			"mykey": 1,
		})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(qbc.JSON.Stringify(data))
}

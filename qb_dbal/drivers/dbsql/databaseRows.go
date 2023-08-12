package dbsql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type DatabaseRows struct {
	err  error
	rows *sql.Rows
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDatabaseRows(rows *sql.Rows, err error) *DatabaseRows {
	instance := new(DatabaseRows)
	instance.err = err
	instance.rows = rows

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DatabaseRows) HasError() bool {
	return nil != instance.err
}

func (instance *DatabaseRows) GetError() error {
	return instance.err
}

func (instance *DatabaseRows) Error() string {
	if nil != instance.err {
		return instance.err.Error()
	}
	return ""
}

func (instance *DatabaseRows) Close() error {
	if nil != instance.rows {
		return instance.rows.Close()
	}
	return instance.err
}

func (instance *DatabaseRows) Columns() ([]string, error) {
	response := make([]string, 0)
	if nil != instance.rows {
		return instance.rows.Columns()
	}
	return response, instance.err
}

func (instance *DatabaseRows) ForEach(callback func(item map[string]interface{}) bool) error {
	if nil != instance.rows && nil != callback {
		columns, err := instance.Columns()
		if err != nil {
			return err
		}
		count := len(columns)
		for instance.rows.Next() {
			values := make([]interface{}, count)
			scanArgs := make([]interface{}, count)
			for i := range values {
				scanArgs[i] = &values[i]
			}

			err = instance.rows.Scan(scanArgs...)
			if nil != err {
				return err
			}

			item, err := toMap(columns, values)
			if nil != err {
				return err
			}

			exit := callback(item)
			if exit {
				break
			}
		}
	}
	return nil
}

func (instance *DatabaseRows) First() (map[string]interface{}, error) {
	if instance.HasError() {
		return nil, instance.GetError()
	}
	columns, err := instance.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	if instance.rows.Next() {
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = instance.rows.Scan(scanArgs...)
		if nil != err {
			return nil, err
		}

		return toMap(columns, values)
	}
	return nil, nil
}

func (instance *DatabaseRows) All() ([]map[string]interface{}, error) {
	if instance.HasError() {
		return nil, instance.GetError()
	}
	columns, err := instance.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	response := make([]map[string]interface{}, 0)
	for instance.rows.Next() {
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = instance.rows.Scan(scanArgs...)
		if nil != err {
			return nil, err
		}
		item, err := toMap(columns, values)
		if nil != err {
			return nil, err
		}
		response = append(response, item)
	}
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toMap(columns []string, values []interface{}) (map[string]interface{}, error) {
	item := map[string]interface{}{}
	for i, bytes := range values {
		if nil != bytes {
			err := putValue(i, item, columns[i], bytes)
			if nil != err {
				return nil, err
			}
		} else {
			item[columns[i]] = nil
		}
	}
	return item, nil
}

func putValue(index int, item map[string]interface{}, name string, value interface{}) error {
	if len(name) == 0 {
		name = fmt.Sprintf("col_%v", index)
	}
	if v, b := value.([]byte); b {
		item[name] = convert(string(v))
	} else if v, b := value.([]int32); b {
		item[name] = convert(string(v))
	} else {
		item[name] = value
	}
	return nil
}

func convert(s string) interface{} {
	isNumber, _ := regexp.MatchString("^[1-9][\\.\\d]*(,\\d+)?$", s)
	//isString := len(s) > 0 && strings.HasPrefix(s, "0")
	if !isNumber {
		return s
	} else if nx, ok := strconv.ParseFloat(s, 64); ok == nil {
		return nx
	} else if b, ok := strconv.ParseBool(s); ok == nil {
		return b
	} else {
		return s
	}
}

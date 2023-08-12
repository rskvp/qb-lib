package helpers

import (
	"database/sql"
	"errors"
	"fmt"
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

			item := map[string]interface{}{}
			for i, v := range values {
				if nil != v {
					x := v.([]byte)
					if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
						item[columns[i]] = nx
					} else if b, ok := strconv.ParseBool(string(x)); ok == nil {
						item[columns[i]] = b
					} else if "string" == fmt.Sprintf("%T", string(x)) {
						item[columns[i]] = string(x)
					} else {
						return errors.New(fmt.Sprintf("Failed on if for type %T of %v\n", x, x))
					}
				} else {
					item[columns[i]] = nil
				}
			}
			exit := callback(item)
			if exit {
				break
			}
		}
	}
	return nil
}

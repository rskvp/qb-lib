package dbsql

import "database/sql"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type DatabaseRow struct {
	err      error
	row      *sql.Row
	response interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDatabaseRow(row *sql.Row, response interface{}) *DatabaseRow {
	instance := new(DatabaseRow)
	instance.row = row
	instance.response = response

	instance.err = row.Scan(instance.response)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DatabaseRow) HasError() bool {
	return nil != instance.err
}

func (instance *DatabaseRow) GetError() error {
	return instance.err
}

func (instance *DatabaseRow) Error() string {
	if nil != instance.err {
		return instance.err.Error()
	}
	return ""
}

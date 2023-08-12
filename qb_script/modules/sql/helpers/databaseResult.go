package helpers

import "database/sql"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type DatabaseResult struct {
	err    error
	result sql.Result
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDatabaseResult(result sql.Result, err error) *DatabaseResult {
	instance := new(DatabaseResult)
	instance.err = err
	instance.result = result

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DatabaseResult) HasError() bool {
	return nil != instance.err
}

func (instance *DatabaseResult) GetError() error {
	return instance.err
}

func (instance *DatabaseResult) Error() string {
	if nil != instance.err {
		return instance.err.Error()
	}
	return ""
}

func (instance *DatabaseResult) LastInsertId() (int64, error) {
	if nil != instance.result {
		return instance.result.LastInsertId()
	}
	return -1, instance.err
}

func (instance *DatabaseResult) RowsAffected() (int64, error) {
	if nil != instance.result {
		return instance.result.RowsAffected()
	}
	return -1, instance.err
}
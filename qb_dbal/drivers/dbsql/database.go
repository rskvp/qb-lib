package dbsql

import (
	"database/sql"
	"errors"

	// _ "github.com/alexbrainman/odbc" REMOVED FROM LIB. USE THIS IN FINAL PROJECT TO AVOID CROSS COMPILE ISSUE!!!!

	_ "github.com/go-sql-driver/mysql"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

var (
	DatabaseNotInitializedError = errors.New("database_not_initialized")
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Database struct {
	driver         string // i.e. "mysql"
	dataSourceName string // i.e. "user:password@/dbname", "admin:admin@tcp(localhost:3306)/test"
	db             *sql.DB
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDatabase(driverName, connectionString string) (*Database, error) {
	instance := new(Database)
	instance.driver = driverName
	instance.dataSourceName = connectionString

	err := instance.init()
	if nil != err {
		return nil, err
	}

	return instance, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Database) Close() error {
	if nil != instance.db {
		err := instance.db.Close()
		instance.db = nil
		return err
	}
	return DatabaseNotInitializedError
}

func (instance *Database) Query(query string, args ...interface{}) *DatabaseRows {
	if nil != instance.db {
		return NewDatabaseRows(instance.db.Query(query, args...))
	}
	return nil
}

func (instance *Database) QueryRow(query string, response interface{}, args ...interface{}) *DatabaseRow {
	if nil != instance.db {
		return NewDatabaseRow(instance.db.QueryRow(query, args...), response)
	}
	return nil
}

func (instance *Database) Exec(query string, args ...interface{}) *DatabaseResult {
	if nil != instance.db {
		return NewDatabaseResult(instance.db.Exec(query, args...))
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Database) init() error {
	db, err := sql.Open(instance.driver, instance.dataSourceName)
	if nil != err {
		return qbc.Errors.Prefix(err, "Error Opening Database")
	}
	err = db.Ping()
	if nil != err {
		return qbc.Errors.Prefix(err, "Error Testing Connection")
	}

	// assign database
	instance.db = db

	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func SqlResultToMap(result *DatabaseResult) (map[string]interface{}, error) {
	if nil != result.GetError() {
		return nil, result.GetError()
	}
	lastInsertId, err := result.LastInsertId()
	if nil != err {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if nil != err {
		return nil, err
	}
	response := map[string]interface{}{
		"last_id":       lastInsertId,
		"rows_affected": rowsAffected,
	}

	return response, nil
}
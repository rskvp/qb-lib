package drivers

import (
	qbc "github.com/rskvp/qb-core"
	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
)

type ForEachCallback func(map[string]interface{}) bool // if return TRUE, exit loop

type IDatabase interface {
	Uid() string
	DriverName() string
	Enabled() bool
	Open() error
	Close() error

	Upsert(collection string, doc map[string]interface{}) (map[string]interface{}, error)
	Remove(collection string, key string) error
	Get(collection string, key string) (map[string]interface{}, error)
	ForEach(collection string, callback ForEachCallback) error
	Find(collection string, fieldName string, fieldValue interface{}) (interface{}, error)

	// optional methods. May be not supported from all databases
	EnsureIndex(collection string, typeName string, fields []string, unique bool) (bool, error)
	EnsureCollection(collection string) (bool, error)

	// native methods does not support cross-database
	ExecNative(command string, bindingVars map[string]interface{}) (interface{}, error)
	ExecMultiple(commands []string, bindVars []map[string]interface{}, options interface{}) ([]interface{}, error)

	// utils
	QueryGetParamNames(query string) []string
	QuerySelectParams(query string, allParams map[string]interface{}) map[string]interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	I N I T
//----------------------------------------------------------------------------------------------------------------------

var cache *CacheManager

func init() {
	cache = NewCache()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// return a thread safe cache manager
func Cache() *CacheManager {
	return cache
}

func NewDatabase(driverName, connectionString string) (driver IDatabase, err error) {
	dsn := dbalcommons.NewDsn(connectionString)
	return NewDatabaseFromDsn(driverName, dsn)
}

func NewDatabaseFromDsn(driverName string, dsn *dbalcommons.Dsn) (driver IDatabase, err error) {
	if dsn.IsValid() {
		switch driverName {
		case NameArango:
			driver = NewDriverArango(dsn)
			err = driver.Open()
		case NameBolt:
			driver = NewDriverBolt(dsn)
			err = driver.Open()
		case NameMsSQL, NameODBC:
			// ODBC, MsSQL
			driver = NewDriverODBC(driverName, dsn)
			err = driver.Open()
		case NameMySQL, NameOracle:
			// SQL database
			driver = NewDriverSQL(driverName, dsn)
			err = driver.Open()
		default:
			driver = NewDriverGorm(driverName, dsn)
			err = driver.Open()
		}

		return driver, err
	}
	return driver, err
}

func OpenDatabase(driver, connectionString string) (IDatabase, error) {
	db, err := NewDatabase(driver, connectionString)
	if nil != err {
		return nil, err
	}
	err = db.Open()
	if nil != err {
		return nil, err
	}
	return db, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func keyFrom(driver, dsn string) string {
	return qbc.Coding.MD5(driver + dsn)
}

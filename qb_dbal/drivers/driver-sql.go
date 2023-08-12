package drivers

import (
	"fmt"

	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
	"github.com/rskvp/qb-lib/qb_dbal/drivers/dbsql"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

// mysql.server start

const NameMySQL = "mysql"
const NameOracle = "oracle"

//----------------------------------------------------------------------------------------------------------------------
//	NewDriverSQL
//----------------------------------------------------------------------------------------------------------------------

type DriverSQL struct {
	uid    string
	driver string
	dsn    *dbalcommons.Dsn
	db     *dbsql.Database
	err    error
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

// NewDriverSQL "admin:admin@tcp(localhost:3306)/test"
func NewDriverSQL(driver string, dsn ...interface{}) *DriverSQL {
	instance := new(DriverSQL)
	instance.driver = driver

	if len(dsn) == 1 {
		if s, b := dsn[0].(string); b {
			instance.dsn = dbalcommons.NewDsn(s)
		} else if d, b := dsn[0].(dbalcommons.Dsn); b {
			instance.dsn = &d
		} else if d, b := dsn[0].(*dbalcommons.Dsn); b {
			instance.dsn = d
		} else {
			instance.err = dbalcommons.ErrorInvalidDsn
		}
	}
	if nil == instance.dsn && nil == instance.err {
		instance.err = dbalcommons.ErrorInvalidDsn
	}
	if nil != instance.dsn {
		instance.uid = keyFrom(driver, instance.dsn.String())
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverSQL) Uid() string {
	return instance.uid
}

func (instance *DriverSQL) DriverName() string {
	return instance.driver
}

func (instance *DriverSQL) Enabled() bool {
	return nil != instance && nil != instance.dsn && nil == instance.err && instance.dsn.IsValid()
}

func (instance *DriverSQL) Open() error {
	if nil != instance {
		if nil == instance.err {
			instance.err = instance.init()
		}
		return instance.err
	}
	return nil
}

func (instance *DriverSQL) Close() error {
	if nil != instance.db {
		return instance.db.Close()
	}
	return nil
}

func (instance *DriverSQL) Remove(collection, key string) error {
	if nil != instance && nil != instance.db {
		query := dbsql.BuildDeleteCommand(collection, "id="+key)
		response := instance.db.Exec(query)
		return response.GetError()
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) Get(collection string, key string) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		query := dbsql.BuildSelect(collection) +
			fmt.Sprintf(" WHERE t.id=%v", key)
		rows := instance.db.Query(query)
		if rows.HasError() {
			return nil, rows.GetError()
		}
		return rows.First()
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) Upsert(collection string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		id := fmt.Sprintf("%v", item["id"])
		query := dbsql.BuildUpdateCommand(collection, "id", id, item)
		response := instance.db.Exec(query)
		if response.HasError() {
			return nil, response.GetError()
		}
		return item, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) ForEach(collection string, callback ForEachCallback) error {
	if nil != instance && nil != instance.db {
		if nil != callback {
			query := dbsql.BuildSelect(collection)
			rows := instance.db.Query(query)
			return rows.ForEach(callback)
		}
		return nil // do nothing
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) ExecNative(command string, bindVars map[string]interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {
		query, args := mergeParams(command, bindVars)
		result := instance.db.Query(query, args...)
		if result.HasError() {
			return nil, result.GetError()
		}
		return result.All()
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) ExecMultiple(commands []string, bindVars []map[string]interface{}, options interface{}) (response []interface{}, err error) {
	if nil != instance && nil != instance.db {
		if len(commands) != len(bindVars) {
			return nil, dbalcommons.ErrorCommandAndParamsDoNotMatch
		}
		for i := 0; i < len(commands); i++ {
			command := commands[i]
			bindVar := bindVars[i]
			data, execErr := instance.ExecNative(command, bindVar)
			if nil != execErr {
				err = execErr
				return
			}
			response = append(response, data)
		}
		return
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) EnsureIndex(collection string, typeName string, fields []string, unique bool) (bool, error) {
	if nil != instance && nil != instance.db {

		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) EnsureCollection(collection string) (bool, error) {
	if nil != instance && nil != instance.db {

		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) Find(collection string, fieldName string, fieldValue interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {

	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverSQL) QueryGetParamNames(query string) []string {
	return QueryGetParamNames(query)
}

func (instance *DriverSQL) QuerySelectParams(query string, allParams map[string]interface{}) map[string]interface{} {
	return QuerySelectParams(query, allParams)
}

//----------------------------------------------------------------------------------------------------------------------
//	e x t e n d
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverSQL) init() error {
	if nil != instance.dsn {
		db, err := dbsql.NewDatabase(instance.driver, instance.dsn.String())
		if nil != err {
			return err
		}
		instance.db = db
	}
	return nil
}

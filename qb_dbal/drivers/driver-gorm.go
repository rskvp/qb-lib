package drivers

import (
	"fmt"

	qbc "github.com/rskvp/qb-core"
	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
	"github.com/rskvp/qb-lib/qb_dbal/drivers/dbsql"
	"github.com/rskvp/qb-lib/qbl_commons"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	NewDriverSQL
//----------------------------------------------------------------------------------------------------------------------

type DriverGorm struct {
	uid    string
	driver string
	dsn    string
	db     *gorm.DB
	err    error
	mode   string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

// NewDriverGorm "admin:admin@tcp(localhost:3306)/test"
func NewDriverGorm(driver string, dsn ...interface{}) *DriverGorm {
	instance := new(DriverGorm)
	instance.driver = driver
	instance.mode = qbc.ModeProduction

	if len(dsn) == 1 {
		if s, b := dsn[0].(string); b {
			instance.dsn = s
		} else if d, b := dsn[0].(dbalcommons.Dsn); b {
			instance.dsn = (&d).String()
		} else if d, b := dsn[0].(*dbalcommons.Dsn); b {
			instance.dsn = d.String()
		} else {
			instance.err = dbalcommons.ErrorInvalidDsn
		}
	}
	if len(instance.dsn) == 0 && nil == instance.err {
		instance.err = dbalcommons.ErrorInvalidDsn
	}
	if len(instance.dsn) > 0 {
		instance.uid = keyFrom(driver, instance.dsn)
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverGorm) SetMode(value string) *DriverGorm {
	instance.mode = value
	return instance
}

func (instance *DriverGorm) GetMode() string {
	return instance.mode
}

func (instance *DriverGorm) Uid() string {
	return instance.uid
}

func (instance *DriverGorm) DriverName() string {
	return instance.driver
}

func (instance *DriverGorm) Enabled() bool {
	return nil != instance && len(instance.dsn) > 0
}

func (instance *DriverGorm) Open() error {
	if nil != instance {
		if nil == instance.err {
			instance.err = instance.init()
		}
		return instance.err
	}
	return nil
}

func (instance *DriverGorm) Close() error {
	if nil != instance.db {
		instance.db = nil
	}
	return nil
}

func (instance *DriverGorm) Remove(collection, key string) (err error) {
	if nil != instance && nil != instance.db {
		query := dbsql.BuildDeleteCommand(collection, "id="+key)
		tx := instance.db.Exec(query)
		if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
			err = tx.Error
		}
	} else {
		err = dbalcommons.ErrorDatabaseDoesNotExists
	}
	return
}

func (instance *DriverGorm) Get(collection string, key string) (response map[string]interface{}, err error) {
	if nil != instance && nil != instance.db && len(key) > 0 && len(collection) > 0 {
		query := dbsql.BuildSelect(collection) +
			fmt.Sprintf(" WHERE t.id=%v", key)

		tx := instance.db.Raw(query).First(&response)
		if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
			err = tx.Error
		}
	} else {
		err = dbalcommons.ErrorDatabaseDoesNotExists
	}
	return
}

func (instance *DriverGorm) Upsert(collection string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db && len(collection) > 0 {
		id := qbc.Convert.ToString(item["id"])
		var tx *gorm.DB
		if ok, _ := instance.Exists(collection, id); ok {
			query := dbsql.BuildUpdateCommand(collection, "id", id, item)
			tx = instance.db.Exec(query)
		} else {
			query := dbsql.BuildInsertCommands(collection, []interface{}{item})[0]
			tx = instance.db.Exec(query)
		}

		if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
			return nil, tx.Error
		}
		return item, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverGorm) ForEach(collection string, callback ForEachCallback) error {
	if nil != instance && nil != instance.db {
		if nil != callback {
			var response []map[string]interface{}
			query := dbsql.BuildSelect(collection)
			tx := instance.db.Raw(query)
			tx.Scan(&response)
			if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
				return tx.Error
			}
			if len(response) > 0 {
				for _, item := range response {
					if callback(item) {
						break
					}
				}
			}
		}
		return nil // do nothing
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverGorm) ExecNative(command string, bindVars map[string]interface{}) (interface{}, error) {
	var result interface{}
	if nil != instance && nil != instance.db {

		query, args := mergeParams(command, bindVars)
		tx := instance.db.Raw(query, args...)
		tx.Scan(&result)

		if nil != tx.Error {
			if IsRecordNotFoundError(tx.Error) {
				return nil, nil
			}
			return nil, tx.Error
		}
	}
	return result, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverGorm) ExecMultiple(commands []string, bindVars []map[string]interface{}, options interface{}) (response []interface{}, err error) {
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

func (instance *DriverGorm) EnsureIndex(collection string, typeName string, fields []string, unique bool) (bool, error) {
	if nil != instance && nil != instance.db {

		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverGorm) EnsureCollection(collection string) (bool, error) {
	if nil != instance && nil != instance.db {

		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverGorm) Find(collection string, fieldName string, fieldValue interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {
		//empty
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverGorm) QueryGetParamNames(query string) []string {
	return QueryGetParamNames(query)
}

func (instance *DriverGorm) QuerySelectParams(query string, allParams map[string]interface{}) map[string]interface{} {
	return QuerySelectParams(query, allParams)
}

//----------------------------------------------------------------------------------------------------------------------
//	e x t e n d
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverGorm) Exists(collection, key string) (bool, interface{}) {
	if len(key) > 0 {
		response, err := instance.Get(collection, key)
		return nil == err && nil != response, response
	}
	return false, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverGorm) init() (err error) {
	if len(instance.dsn) > 0 {
		switch instance.driver {
		case "sqlite":
			filename := qbc.Paths.WorkspacePath(instance.dsn)
			instance.db, err = gorm.Open(sqlite.Open(filename), qbl_commons.GormConfig(instance.mode))
		case "mysql":
			// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
			instance.db, err = gorm.Open(mysql.Open(instance.dsn), &gorm.Config{})
		case "postgres":
			// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
			instance.db, err = gorm.Open(postgres.Open(instance.dsn), &gorm.Config{})
		case "sqlserver":
			// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
			instance.db, err = gorm.Open(sqlserver.Open(instance.dsn), &gorm.Config{})
		default:
			instance.db = nil
			err = qbc.Errors.Prefix(dbalcommons.ErrorDriverNotImplemented,
				fmt.Sprintf("'%s': ", instance.driver))
		}

	}
	return
}

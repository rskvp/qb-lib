package drivers

import (
	"encoding/json"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/bolt"
	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
)

const NameBolt = "bolt"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type DriverBolt struct {
	uid string
	dsn *dbalcommons.Dsn
	db  *bolt.BoltDatabase
	err error
}

func NewDriverBolt(dsn ...interface{}) *DriverBolt {
	instance := new(DriverBolt)
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
		instance.uid = keyFrom(NameBolt, instance.dsn.String())
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverBolt) Uid() string {
	return instance.uid
}

func (instance *DriverBolt) DriverName() string {
	return NameBolt
}

func (instance *DriverBolt) Enabled() bool {
	return nil != instance && nil != instance.dsn && nil == instance.err && instance.dsn.IsValid()
}

func (instance *DriverBolt) Open() error {
	if nil != instance {
		if nil == instance.err {
			filename := qbc.Paths.Absolute(instance.dsn.Database)
			err := qbc.Paths.Mkdir(filename)
			if nil != err {
				instance.err = err
			} else {
				config := bolt.NewBoltConfig()
				config.Name = filename
				instance.db = bolt.NewBoltDatabase(config)
				instance.err = instance.db.Open()
			}
		}
		return instance.err
	}
	return nil
}

func (instance *DriverBolt) Close() error {
	if nil != instance && nil != instance.db {
		return instance.db.Close()
	}
	return nil
}

func (instance *DriverBolt) Remove(collection string, key string) error {
	if nil != instance && nil != instance.db {
		return instance.remove(collection, key)
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) Get(collection string, key string) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		coll, err := instance.db.Collection(collection, true)
		if nil != err {
			return nil, err
		}
		item, err := coll.Get(key)
		if nil != err {
			return nil, err
		}
		if v, b := item.(map[string]interface{}); b {
			return v, nil
		}
		return nil, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) Upsert(collection string, doc map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		coll, err := instance.db.Collection(collection, true)
		if nil != err {
			return nil, err
		}

		if _, b := doc["_key"]; !b {
			doc["_key"] = qbc.Rnd.Uuid()
		}

		err = coll.Upsert(doc)
		if nil != err {
			return nil, err
		}
		return doc, nil
	}
	return nil, nil
}

func (instance *DriverBolt) ForEach(collection string, callback ForEachCallback) error {
	if nil != instance && nil != instance.db {
		if nil != callback {
			coll, err := instance.db.Collection(collection, true)
			if nil != err {
				return err
			}
			var doc map[string]interface{}
			err = coll.ForEach(func(k, v []byte) bool {
				e := json.Unmarshal(v, &doc)
				if nil != e {
					err = e
					return true // exit
				}
				return callback(doc)
			})

			return err
		}
		return nil
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) ExecNative(command string, bindingVars map[string]interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {

		return nil, dbalcommons.ErrorCommandNotSupported
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) ExecMultiple(commands []string, bindVars []map[string]interface{}, options interface{}) ([]interface{}, error) {
	if nil != instance && nil != instance.db {

		return nil, dbalcommons.ErrorCommandNotSupported
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) EnsureIndex(collection string, typeName string, fields []string, unique bool) (bool, error) {
	if nil != instance && nil != instance.db {
		_, err := instance.db.Collection(collection, true)
		if nil != err {
			return false, err
		}
		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) EnsureCollection(collection string) (bool, error) {
	if nil != instance && nil != instance.db {
		_, err := instance.db.Collection(collection, true)
		if nil != err {
			return false, err
		}
		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) Find(collection string, fieldName string, fieldValue interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {
		coll, err := instance.db.Collection(collection, true)
		if nil != err {
			return nil, err
		}
		return coll.GetByFieldValue(fieldName, fieldValue)
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverBolt) QueryGetParamNames(query string) []string {
	return QueryGetParamNames(query)
}

func (instance *DriverBolt) QuerySelectParams(query string, allParams map[string]interface{}) map[string]interface{} {
	return QuerySelectParams(query, allParams)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverBolt) remove(collectionName string, key string) (err error) {
	if nil != instance && nil != instance.db {
		var coll *bolt.BoltCollection
		coll, err = instance.db.Collection(collectionName, true)
		if nil == err {
			err = coll.Remove(key)
		}
	}
	return err
}

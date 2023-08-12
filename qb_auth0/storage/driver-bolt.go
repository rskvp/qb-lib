package storage

import (
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/bolt"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type DriverBolt struct {
	dsn         *Dsn
	enableCache bool
	db          *bolt.BoltDatabase
	err         error
}

func NewDriverBolt(dsn ...interface{}) *DriverBolt {
	instance := new(DriverBolt)
	if len(dsn) == 1 {
		if s, b := dsn[0].(string); b {
			instance.dsn = NewDsn(s)
		} else if d, b := dsn[0].(Dsn); b {
			instance.dsn = &d
		} else if d, b := dsn[0].(*Dsn); b {
			instance.dsn = d
		} else {
			instance.err = ErrorInvalidDsn
		}
	}
	if nil == instance.dsn && nil == instance.err {
		instance.err = ErrorInvalidDsn
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverBolt) EnableCache(value bool) {
	instance.enableCache = value
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

//----------------------------------------------------------------------------------------------------------------------
//	a u t h
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverBolt) AuthRegister(key, payload string) error {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}
		coll, err := instance.db.Collection(CollectionAuth, true)
		if nil != err {
			return err
		}
		i, err := coll.Get(key)
		if nil != err {
			return err
		}
		if nil != i {
			return ErrorEntityAlreadyRegistered
		}
		item := map[string]interface{}{
			"_key":    key,
			"payload": payload,
		}
		return coll.Upsert(item)
	}
	return nil
}

func (instance *DriverBolt) AuthOverwrite(key, payload string) error {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}
		coll, err := instance.db.Collection(CollectionAuth, true)
		if nil != err {
			return err
		}

		item := map[string]interface{}{
			"_key":    key,
			"payload": payload,
		}
		return coll.Upsert(item)
	}
	return nil
}

func (instance *DriverBolt) AuthGet(key string) (payload string, err error) {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return payload, ErrorDatabaseCacheCannotAuthenticate
		}
		var coll *bolt.BoltCollection
		coll, err = instance.db.Collection(CollectionAuth, true)
		if nil == err {
			var i interface{}
			i, err = coll.Get(key)
			if nil == err {
				if nil != i {
					payload = qbc.Reflect.GetString(i, "payload")
				} else {
					err = ErrorEntityDoesNotExists
				}
			}
		}
	}
	return payload, err
}

func (instance *DriverBolt) AuthRemove(key string) (err error) {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}
		err = instance.remove(CollectionAuth, key)
	}
	return err
}

//----------------------------------------------------------------------------------------------------------------------
//	c a c h e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverBolt) CacheGet(key string) (string, error) {
	if nil != instance && nil != instance.db {
		if !instance.enableCache {
			return "", ErrorDatabaseCacheNotEnabled
		}

		coll, err := instance.db.Collection(CollectionCache, true)
		if nil != err {
			return "", err
		}
		item, err := coll.Get(key)
		if nil != err {
			return "", err
		}
		if nil != item {
			if m, b := item.(map[string]interface{}); b {
				var token string
				var err error
				now := int(time.Now().Unix())
				expire := qbc.Convert.ToInt(m[bolt.FieldExpire])
				token = m["token"].(string)
				if now-expire > 0 {
					// expired
					err = ErrorTokenExpired
				}
				return token, err
			}
			// not found or expired
			return "", ErrorTokenDoesNotExists
		} else {
			// not found or expired
			return "", ErrorTokenDoesNotExists
		}
	}
	return "", ErrorTokenDoesNotExists
}

func (instance *DriverBolt) CacheAdd(key, token string, duration time.Duration) error {
	if nil != instance && nil != instance.db {
		if !instance.enableCache {
			return ErrorDatabaseCacheNotEnabled
		}
		coll, err := instance.db.CollectionAutoCreate(CollectionCache)
		if nil != err {
			return err
		}
		item := map[string]interface{}{
			"_key":  key,
			"token": token,
		}
		item[bolt.FieldExpire] = time.Now().Add(duration).Unix()
		return coll.Upsert(item)
	}
	return nil
}

func (instance *DriverBolt) CacheRemove(key string) error {
	if nil != instance && nil != instance.db {
		if !instance.enableCache {
			return ErrorDatabaseCacheNotEnabled
		}
		return instance.remove(CollectionCache, key)
	}
	return nil
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

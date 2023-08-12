package bolt

import (
	"time"

	qbc "github.com/rskvp/qb-core"
	"go.etcd.io/bbolt"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type BoltDatabase struct {

	//-- private --//
	config *BoltConfig
	name   string
	path   string
	db     *bbolt.DB
	expire *ExpireWorker
	err    error
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewBoltDatabase(config *BoltConfig) *BoltDatabase {
	instance := new(BoltDatabase)
	instance.config = config
	instance.expire = NewExpireWorker(instance) // handle items expiration

	if b, _ := qbc.Paths.IsFile(config.Name); b {
		instance.path = qbc.Paths.Absolute(config.Name)
		instance.name = qbc.Paths.FileName(config.Name, false)
	} else {
		instance.name = config.Name
		instance.path = qbc.Paths.Absolute(config.Name + ".dat")
	}

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltDatabase) Name() string {
	if nil != instance {
		return instance.name
	}
	return ""
}

func (instance *BoltDatabase) Error() error {
	if nil != instance {
		return instance.err
	}
	return nil
}

func (instance *BoltDatabase) HasError() bool {
	if nil != instance {
		return nil != instance.err
	}
	return false
}

func (instance *BoltDatabase) Open() error {
	if nil != instance && nil == instance.db {

		db, err := bbolt.Open(instance.path, 0600, &bbolt.Options{Timeout: instance.config.TimeoutMs * time.Millisecond})
		if nil == err && nil != db {
			instance.db = db
		}
		instance.err = err
		return err
	}
	return nil
}

func (instance *BoltDatabase) Close() error {
	if nil != instance && nil != instance.db {
		err := instance.db.Close()
		if nil == err {
			instance.db = nil
		}
		return err
	}
	return nil
}

func (instance *BoltDatabase) Drop() error {
	if nil != instance && nil != instance.db {
		err := instance.Close()
		if nil != err {
			return err
		}
		// remove file
		return qbc.IO.Remove(instance.path)
	}
	return nil
}

func (instance *BoltDatabase) Size() (int64, error) {
	if nil != instance && nil != instance.db {
		path := instance.db.Path()
		if len(path) > 0 {
			return qbc.IO.FileSize(path)
		}
	}
	return 0, ErrDatabaseIsNotConnected
}

func (instance *BoltDatabase) CollectionAutoCreate(name string) (*BoltCollection, error) {
	return instance.Collection(name, true)
}

func (instance *BoltDatabase) Collection(name string, createIfNotExists bool) (*BoltCollection, error) {
	if nil != instance && nil != instance.db {
		err := instance.db.Update(func(tx *bbolt.Tx) error {
			var e error
			b := tx.Bucket([]byte(name))
			if nil == b {
				if createIfNotExists {
					_, e = tx.CreateBucketIfNotExists([]byte(name))
				} else {
					e = ErrCollectionDoesNotExists
				}
			}
			return e
		})
		if nil == err {
			return NewBoltCollection(instance, instance.db, name), nil
		}
		return nil, err
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

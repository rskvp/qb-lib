package bolt

import (
	"encoding/json"
	"errors"

	qbc "github.com/rskvp/qb-core"
	"go.etcd.io/bbolt"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type BoltCollection struct {

	//-- private --//
	name         string
	db           *BoltDatabase
	boltdb       *bbolt.DB
	enableExpire bool
}

type ForEachCallback func(k, v []byte) bool

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewBoltCollection(db *BoltDatabase, boltdb *bbolt.DB, name string) *BoltCollection {
	instance := new(BoltCollection)
	instance.db = db
	instance.boltdb = boltdb
	instance.name = name

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltCollection) EnableExpire(value bool) {
	if nil != instance {
		instance.enableExpire = value
		if value {
			// start internal task to check expiration
			instance.db.expire.Enable(instance.name)
		} else {
			instance.db.expire.Disable(instance.name)
		}
	}
}

func (instance *BoltCollection) Drop() error {
	if nil != instance && nil != instance.boltdb {
		err := instance.boltdb.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				return tx.DeleteBucket([]byte(instance.name))
			}
			return nil
		})
		return err
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Count() (int64, error) {
	var response int64
	response = 0
	if nil != instance && nil != instance.boltdb {
		err := instance.boltdb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					response++
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) CountByFieldValue(fieldName string, fieldValue interface{}) (int64, error) {
	var response int64
	response = 0
	if nil != instance && nil != instance.boltdb {
		err := instance.boltdb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var entity interface{}
					err := json.Unmarshal(v, &entity)
					if nil == err {
						value := qbc.Reflect.Get(entity, fieldName)
						if qbc.Compare.Equals(fieldValue, value) {
							response++
						}
					}
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Get(key string) (interface{}, error) {
	var response interface{}
	if nil != instance && nil != instance.boltdb {
		err := instance.boltdb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				buf := b.Get([]byte(key))
				if nil != buf {
					err := json.Unmarshal(buf, &response)
					return err
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return nil, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Upsert(entity interface{}) error {
	if nil != instance {
		return instance.update(entity, false)
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) UpsertBatch(entity interface{}) error {
	if nil != instance {
		return instance.update(entity, true)
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Remove(key string) error {
	if nil != instance {
		return instance.remove(key, false)
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) RemoveBatch(key string) error {
	if nil != instance {
		return instance.remove(key, true)
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) GetByFieldValue(fieldName string, fieldValue interface{}) ([]interface{}, error) {
	response := make([]interface{}, 0)
	if nil != instance && nil != instance.boltdb {
		err := instance.boltdb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var entity interface{}
					err := json.Unmarshal(v, &entity)
					if nil == err {
						value := qbc.Reflect.Get(entity, fieldName)
						if qbc.Compare.Equals(fieldValue, value) {
							response = append(response, entity)
						}
					}
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Find(query *BoltQuery) ([]interface{}, error) {
	response := make([]interface{}, 0)
	if nil != instance && nil != instance.boltdb {
		err := instance.boltdb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var entity map[string]interface{}
					err := json.Unmarshal(v, &entity)
					if nil == err {
						if query.MatchFilter(entity) {
							response = append(response, entity)
						}
					}
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) ForEach(callback ForEachCallback) error {
	if nil != instance && nil != instance.boltdb {
		if nil != callback {
			err := instance.boltdb.View(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(instance.name))
				if nil != b {
					_ = b.ForEach(func(k, v []byte) error {
						exit := callback(k, v)
						if exit {
							return errors.New("exit")
						}
						return nil
					})
				} else {
					return ErrCollectionDoesNotExists
				}
				return nil
			})
			return err
		} else {
			return nil
		}
	}
	return ErrDatabaseIsNotConnected
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func unmarshal(v []byte, entity interface{}) (interface{}, error) {
	if nil != entity {
		err := json.Unmarshal(v, entity)
		if nil != err {
			return nil, err
		}
		//var resp struct{}
		//qb_reflect.Copy(reflect.ValueOf(entity).Elem().Interface(), &resp)

		return entity, nil
	} else {
		var resp map[string]interface{}
		err := json.Unmarshal(v, &resp)
		if nil != err {
			return nil, err
		}
		return resp, nil
	}
}

func (instance *BoltCollection) remove(key string, batch bool) (err error) {
	if nil != instance && nil != instance.boltdb {
		if batch {
			err = instance.boltdb.Batch(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(instance.name))
				if nil != b {
					err := b.Delete([]byte(key))
					return err
				} else {
					return ErrCollectionDoesNotExists
				}
			})
		} else {
			err = instance.boltdb.Update(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(instance.name))
				if nil != b {
					err := b.Delete([]byte(key))
					return err
				} else {
					return ErrCollectionDoesNotExists
				}
			})
		}
	} else {
		err = ErrDatabaseIsNotConnected
	}
	return err
}

func (instance *BoltCollection) update(entity interface{}, batch bool) (err error) {
	if nil != instance && nil != instance.boltdb {
		// check key
		key := []byte(qbc.Reflect.GetString(entity, "_key"))
		if len(key) == 0 {
			return ErrMissingDocumentKey
		}
		// get array of bytes
		buf, err := json.Marshal(entity)
		if nil != err {
			return err
		}
		if batch {
			err = instance.boltdb.Batch(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(instance.name))
				if nil != b {
					return b.Put(key, buf)
				} else {
					return ErrCollectionDoesNotExists
				}
			})
		} else {
			err = instance.boltdb.Update(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(instance.name))
				if nil != b {
					return b.Put(key, buf)
				} else {
					return ErrCollectionDoesNotExists
				}
			})
		}
	} else {
		err = ErrDatabaseIsNotConnected
	}
	return err
}

package drivers

import (
	"sync"

	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
)

// ---------------------------------------------------------------------------------------------------------------------
//	t y p e
// ---------------------------------------------------------------------------------------------------------------------

type CacheManager struct {
	cache map[string]IDatabase

	mux sync.Mutex
}

// ---------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewCache() *CacheManager {
	instance := new(CacheManager)
	instance.cache = make(map[string]IDatabase)
	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *CacheManager) Open() {
//empty
}

func (instance *CacheManager) Close() {
	instance.Clear()
}

func (instance *CacheManager) Clear() {
	instance.mux.Lock()
	defer instance.mux.Unlock()
	for _, db := range instance.cache {
		_ = db.Close()
	}
	instance.cache = make(map[string]IDatabase)
}

/**
"driver": "arango",
"dsn": "root:1234567890@tcp(localhost:8529)/my-database)"
*/
func (instance *CacheManager) Get(driver, dsn string) (IDatabase, error) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	d := dbalcommons.NewDsn(dsn)
	key := keyFrom(driver, d.String())
	if _, b := instance.cache[key]; !b {
		// get and open database
		db, err := NewDatabase(driver, dsn)
		if nil == err {
			instance.cache[key] = db
		} else {
			return nil, err
		}
	}
	return instance.cache[key], nil
}

func (instance *CacheManager) Remove(uid string) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if db, b := instance.cache[uid]; !b {
		_ = db.Close()
		delete(instance.cache, uid)
	}
}

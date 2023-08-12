package bolt

import (
	"encoding/json"
	"time"

	"github.com/rskvp/qb-core/qb_ticker"
)

type ExpireItem struct {
	Exp int64 `json:"_expire"`
}

//----------------------------------------------------------------------------------------------------------------------
//	ExpireWorkerJob
//----------------------------------------------------------------------------------------------------------------------

type ExpireWorkerJob struct {
	db       *BoltDatabase
	collName string
	ticker   *qb_ticker.Ticker
}

func newExpireWorkerJob(db *BoltDatabase, collName string) *ExpireWorkerJob {
	instance := new(ExpireWorkerJob)
	instance.db = db
	instance.collName = collName

	if len(collName) > 0 {
		instance.ticker = qb_ticker.NewTicker(10*time.Second, instance.onTick)
		instance.ticker.Start()
	}

	return instance
}

func (instance *ExpireWorkerJob) onTick(ticker *qb_ticker.Ticker) {
	ticker.Pause()
	defer ticker.Resume()
	if nil != instance && nil != instance.db {
		coll, err := instance.db.Collection(instance.collName, false)
		if nil == err {
			// fmt.Println("TEST EXPIRING", FieldExpire)
			expired := make([]string, 0)
			_ = coll.ForEach(func(k, v []byte) bool {
				var e ExpireItem
				_ = json.Unmarshal(v, &e)
				var m map[string]interface{}
				_ = json.Unmarshal(v, &m)
				if e.Exp > 0 && time.Now().Unix()-int64(e.Exp) > 0 {
					expired = append(expired, string(k))
				}
				return false // continue
			})
			for _, key := range expired {
				_ = coll.Remove(key)
			}
		}
	}
}

func (instance *ExpireWorkerJob) Stop() {
	if nil != instance && nil != instance.ticker {
		instance.ticker.Stop()
		instance.ticker = nil
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	ExpireWorker
//----------------------------------------------------------------------------------------------------------------------

type ExpireWorker struct {
	db     *BoltDatabase
	tables map[string]*ExpireWorkerJob
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewExpireWorker(db *BoltDatabase) *ExpireWorker {
	instance := new(ExpireWorker)
	instance.db = db
	instance.tables = make(map[string]*ExpireWorkerJob)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ExpireWorker) Enable(collectionName string) {
	if nil != instance {
		if _, b := instance.tables[collectionName]; !b {
			instance.tables[collectionName] = newExpireWorkerJob(instance.db, collectionName)
		}
	}
}

func (instance *ExpireWorker) Disable(collectionName string) {
	if nil != instance {
		if job, b := instance.tables[collectionName]; b {
			job.Stop()
			delete(instance.tables, collectionName)
		}
	}
}

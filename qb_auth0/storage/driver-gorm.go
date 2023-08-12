package storage

import (
	"errors"
	"fmt"
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qbl_commons"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type DriverGorm struct {
	driver      string
	dsn         string
	enableCache bool
	_db         *gorm.DB
	err         error
	mode        string
}

// Auth "auth" table
type Auth struct {
	Key     string `gorm:"index:auth_key,unique"`
	Payload string
}

// Cache "cache" table
type Cache struct {
	Key    string `gorm:"index:cache_key,unique"`
	Token  string
	Expire int
}

func NewDriverGorm(driverName string, dsn ...interface{}) *DriverGorm {
	instance := new(DriverGorm)
	instance.driver = driverName
	instance.mode = qbc.ModeProduction

	if len(dsn) == 1 {
		if s, b := dsn[0].(string); b {
			instance.dsn = s
		} else if d, b := dsn[0].(Dsn); b {
			instance.dsn = d.String()
		} else if d, b := dsn[0].(*Dsn); b {
			instance.dsn = d.String()
		} else {
			instance.err = ErrorInvalidDsn
		}
	}
	if len(instance.dsn) == 0 && nil == instance.err {
		instance.err = ErrorInvalidDsn
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

func (instance *DriverGorm) EnableCache(value bool) {
	instance.enableCache = value
}

func (instance *DriverGorm) Enabled() bool {
	return nil != instance && len(instance.dsn) > 0 && nil == instance.err
}

func (instance *DriverGorm) Open() error {
	if nil != instance {
		if nil == instance.err {
			_, err := instance.connection()
			if nil != err {
				instance.err = err
			}
		}
		return instance.err
	}
	return nil
}

func (instance *DriverGorm) Close() error {
	if nil != instance && nil != instance._db {
		instance._db = nil
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	a u t h
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverGorm) AuthRegister(key, payload string) error {
	if nil != instance {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}

		db, err := instance.connection()
		if nil != err {
			return err
		}

		// expected error "record not found"
		tx := db.First(&Auth{}, "key=?", key)
		if nil == tx.Error {
			return ErrorEntityAlreadyRegistered
		}

		item := &Auth{
			Key:     key,
			Payload: payload,
		}
		tx = db.Create(item)
		if nil != tx.Error {
			return tx.Error
		}
	}
	return nil
}

func (instance *DriverGorm) AuthOverwrite(key, payload string) error {
	if nil != instance {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}

		db, err := instance.connection()
		if nil != err {
			return err
		}

		var tx *gorm.DB

		tx = db.First(&Auth{}, "key = ?", key)
		if nil != tx.Error {
			// not found
			tx = nil
		} else {
			// update
			tx = db.Model(&Auth{}).Where("key = ?", key).Update("payload", payload)
		}

		if nil == tx {
			auth := &Auth{
				Key:     key,
				Payload: payload,
			}
			tx = db.Create(auth)
		}

		if nil != tx.Error {
			return tx.Error
		}
	}
	return nil
}

func (instance *DriverGorm) AuthGet(key string) (payload string, err error) {
	if nil != instance {
		if instance.enableCache {
			err = ErrorDatabaseCacheCannotAuthenticate
			return
		}

		var db *gorm.DB
		db, err = instance.connection()
		if nil != err {
			return
		}

		var auth Auth
		tx := db.First(&auth, "key=?", key)
		if nil != tx.Error {
			err = tx.Error
		} else {
			payload = auth.Payload
		}
	}
	return payload, err
}

func (instance *DriverGorm) AuthRemove(key string) (err error) {
	if nil != instance {
		if instance.enableCache {
			err = ErrorDatabaseCacheCannotAuthenticate
			return
		}

		var db *gorm.DB
		db, err = instance.connection()
		if nil != err {
			return
		}
		tx := db.Where("key = ?", key).Delete(&Auth{})
		if nil != tx.Error {
			err = tx.Error
		}
	}
	return err
}

//----------------------------------------------------------------------------------------------------------------------
//	c a c h e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverGorm) CacheGet(key string) (string, error) {
	if nil != instance {
		if !instance.enableCache {
			return "", ErrorDatabaseCacheNotEnabled
		}

		db, err := instance.connection()
		if nil != err {
			return "", err
		}

		var cache Cache
		tx := db.First(&cache, "key=?", key)
		if nil != tx.Error {
			err = tx.Error // not found
		} else {
			now := int(time.Now().Unix())
			expire := cache.Expire
			token := cache.Token
			if now-expire > 0 {
				// expired
				err = ErrorTokenExpired
			}
			return token, err
		}
	}
	return "", ErrorTokenDoesNotExists
}

func (instance *DriverGorm) CacheAdd(key, token string, duration time.Duration) error {
	if nil != instance {
		if !instance.enableCache {
			return ErrorDatabaseCacheNotEnabled
		}

		db, err := instance.connection()
		if nil != err {
			return err
		}

		var cache *Cache
		tx := db.First(&cache, "key=?", key)
		if nil != tx.Error {
			tx = nil // not found
		} else {
			tx = db.Model(&cache).Where("key=?", key).Updates(&Cache{
				Token:  token,
				Expire: int(time.Now().Add(duration).Unix()),
			})
		}

		if nil == tx {
			cache = &Cache{
				Key:    key,
				Token:  token,
				Expire: int(time.Now().Add(duration).Unix()),
			}
			tx = db.Create(cache)
		}

		if nil != tx.Error {
			return tx.Error
		}
	}
	return nil
}

func (instance *DriverGorm) CacheRemove(key string) error {
	if nil != instance {
		if !instance.enableCache {
			return ErrorDatabaseCacheNotEnabled
		}

		db, err := instance.connection()
		if nil != err {
			return err
		}

		tx := db.Where("key = ?", key).Delete(&Cache{})
		if nil != tx.Error {
			return tx.Error
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// retrieve a connection
func (instance *DriverGorm) connection() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	if nil == instance._db {
		driver := instance.driver
		dsn := instance.dsn
		switch driver {
		case "sqlite":
			filename := qbc.Paths.WorkspacePath(dsn)
			db, err = gorm.Open(sqlite.Open(filename), qbl_commons.GormConfig(instance.mode))
		case "mysql":
			// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
			// db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgres":
			// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
			// db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		case "sqlserver":
			// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
			// db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		default:
			db = nil
		}
		if db == nil {
			err = qbc.Errors.Prefix(errors.New(fmt.Sprintf("database '%s' not supported", driver)),
				fmt.Sprintf("'%s': ", driver))
		}
		if nil == err {
			instance._db = db

			_ = instance.init(db)
		}
	}

	return instance._db, err
}

func (instance *DriverGorm) init(db *gorm.DB) (err error) {
	if instance.enableCache {
		err = db.AutoMigrate(&Cache{})
	} else {
		err = db.AutoMigrate(&Auth{})
	}
	return
}

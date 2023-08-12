package storage

import (
	"fmt"
	"time"

	qbc "github.com/rskvp/qb-core"
)

type IDatabase interface {
	Enabled() bool
	Open() error
	Close() error
	EnableCache(value bool)
	AuthRegister(key, payload string) error
	AuthGet(key string) (string, error)
	AuthRemove(key string) error
	AuthOverwrite(key, payload string) error

	CacheGet(key string) (string, error)
	CacheRemove(key string) error
	CacheAdd(key, token string, duration time.Duration) error // add new or update existing
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func NewDatabase(driverName, connectionString string, isCache bool) (driver IDatabase, err error) {
	dsn := NewDsn(connectionString)
	isValid := dsn.IsValid()
	switch driverName {
	case "arango":
		if isValid {
			driver = NewDriverArango(dsn)
			driver.EnableCache(isCache)
			err = driver.Open()
		} else {
			driver = nil
			err = qbc.Errors.Prefix(ErrorDriverNotImplemented, fmt.Sprintf("%s: ", driverName))
		}
	case "bolt":
		if isValid {
			driver = NewDriverBolt(dsn)
			driver.EnableCache(isCache)
			err = driver.Open()
		} else {
			driver = nil
			err = qbc.Errors.Prefix(ErrorDriverNotImplemented, fmt.Sprintf("%s: ", driverName))
		}
	default:
		driver = NewDriverGorm(driverName, connectionString)
		driver.EnableCache(isCache)
		err = driver.Open()
	}
	return driver, err
}

func BuildKey(username, password string) string {
	return qbc.Coding.MD5(username + password)
}

func EncryptText(key []byte, value string) (string, error) {
	return qbc.Coding.EncryptTextWithPrefix(value, key)
}

func DecryptText(key []byte, value string) (string, error) {
	return qbc.Coding.DecryptTextWithPrefix(value, key)
}

func EncryptPayload(key []byte, value map[string]interface{}) (string, error) {
	json := qbc.JSON.Stringify(value)
	return EncryptText(key, json)
}

func DecryptPayload(key []byte, value string) (map[string]interface{}, error) {
	data, err := DecryptText(key, value)
	if nil != err {
		return nil, err
	}
	var e map[string]interface{}
	err = qbc.JSON.Read(data, &e)
	return e, err
}

package qb_auth0

import (
	"os"

	qbc "github.com/rskvp/qb-core"
)


const(
	AuthSecretName = "auth"			// used to encrypt authentication data into db
	AccessSecretName = "access"
	RefreshSecretName = "refresh"
)

//----------------------------------------------------------------------------------------------------------------------
//	Auth0ConfigStorage
//----------------------------------------------------------------------------------------------------------------------

type Auth0ConfigStorage struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
}

func Auth0ConfigStorageParse(json string) *Auth0ConfigStorage {
	instance := new(Auth0ConfigStorage)
	_ = qbc.JSON.Read(json, instance)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	Auth0ConfigSecrets
//----------------------------------------------------------------------------------------------------------------------

type Auth0ConfigSecrets map[string]string

func (instance Auth0ConfigSecrets) String() string {
	if nil != instance {
		return qbc.JSON.Stringify(instance)
	}
	return ""
}

func (instance Auth0ConfigSecrets) Put(key, value string) {
	instance[key] = value
}

func (instance Auth0ConfigSecrets) Get(key string) string {
	if _, b := instance[key]; !b {
		// lookup env variables
		env := os.Getenv(key)
		if len(env) > 0 {
			instance[key] = env
		}
	}
	return instance[key]
}

func (instance Auth0ConfigSecrets) GetNotEmpty(key string) string {
	value := instance.Get(key)
	if len(value) == 0 {
		value = "not_empty_secret"
	}
	return value
}

func (instance Auth0ConfigSecrets) Remove(key string) (value string) {
	if _, b := instance[key]; b {
		value = instance[key]
		delete(instance, key)
	}
	return value
}

//----------------------------------------------------------------------------------------------------------------------
//	Auth0Config
//----------------------------------------------------------------------------------------------------------------------

type Auth0Config struct {
	Secrets      Auth0ConfigSecrets  `json:"secrets"`
	CacheStorage *Auth0ConfigStorage `json:"cache-storage"`
	AuthStorage  *Auth0ConfigStorage `json:"auth-storage"`
}

func Auth0ConfigLoad(fileName string) (*Auth0Config, error) {
	var instance Auth0Config
	err := qbc.JSON.ReadFromFile(fileName, &instance)
	if nil != err {
		return nil, err
	}
	return &instance, nil
}

func Auth0ConfigNew() *Auth0Config {
	instance := new(Auth0Config)
	instance.CacheStorage = new(Auth0ConfigStorage)
	instance.AuthStorage = new(Auth0ConfigStorage)
	instance.Secrets = Auth0ConfigSecrets{}

	return instance
}

func Auth0ConfigParse(json string) *Auth0Config {
	instance := new(Auth0Config)
	_ = qbc.JSON.Read(json, instance)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Auth0Config) GoString() string {
	if nil != instance {
		return qbc.JSON.Stringify(instance)
	}
	return ""
}

func (instance *Auth0Config) String() string {
	if nil != instance {
		return instance.GoString()
	}
	return ""
}

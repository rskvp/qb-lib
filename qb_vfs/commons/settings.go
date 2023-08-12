package commons

import (
	"strings"

	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type VfsSettings struct {
	Location string           `json:"location"`
	Auth     *VfsSettingsAuth `json:"auth"`
}

type VfsSettingsAuth struct {
	User     string `json:"user"`
	Password string `json:"pass"`
	Key      string `json:"key"`
}

func NewVfsSettings() *VfsSettings {
	instance := new(VfsSettings)
	instance.Location = ""
	instance.Auth = new(VfsSettingsAuth)

	return instance
}

func LoadVfsSettings(filename string) (*VfsSettings, error) {
	instance := new(VfsSettings)
	err := qbc.JSON.ReadFromFile(filename, &instance)
	if nil != err {
		return nil, err
	}
	return instance, nil
}

func ParseVfsSettings(jsonData string) (*VfsSettings, error) {
	instance := new(VfsSettings)
	err := qbc.JSON.Read(jsonData, &instance)
	if nil != err {
		return nil, err
	}
	return instance, nil
}

func InitVfsSettings(location, user, password, key string) *VfsSettings {
	instance := NewVfsSettings()
	instance.Location = location
	instance.Auth.User = user
	instance.Auth.Password = password
	instance.Auth.Key = key

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	VfsSettings
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsSettings) Schema() string {
	schema, _ := instance.SplitLocation()
	return schema
}

func (instance *VfsSettings) SplitLocation() (scheme, host string) {
	tokens := strings.Split(instance.Location, "://")
	if len(tokens) == 2 {
		scheme = tokens[0]
		host = tokens[1]
	}
	return
}

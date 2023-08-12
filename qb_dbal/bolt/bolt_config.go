package bolt

import (
	"encoding/json"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type BoltConfig struct {
	Name      string        `json:"name"`
	TimeoutMs time.Duration `json:"timeout_ms"`
}

//----------------------------------------------------------------------------------------------------------------------
//	BoltConfig
//----------------------------------------------------------------------------------------------------------------------

func NewBoltConfig() *BoltConfig {
	instance := new(BoltConfig)
	instance.TimeoutMs = 3000

	return instance
}

func (instance *BoltConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *BoltConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

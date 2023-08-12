package arango

import "encoding/json"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type ArangoConfig struct {
	Endpoints      []string          `json:"endpoints"`
	Authentication *ArangoConfigAuth `json:"authentication"`
}

type ArangoConfigAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//----------------------------------------------------------------------------------------------------------------------
//	ArangoConfig
//----------------------------------------------------------------------------------------------------------------------

func NewArangoConfig() *ArangoConfig{
	instance:= new(ArangoConfig)
	instance.Authentication = new (ArangoConfigAuth)
	return instance
}

func (instance *ArangoConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *ArangoConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

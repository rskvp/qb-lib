package commons

import qbc "github.com/rskvp/qb-core"

//----------------------------------------------------------------------------------------------------------------------
//	SemanticConfigDb
//----------------------------------------------------------------------------------------------------------------------

type SemanticConfigDb struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
}

func (instance *SemanticConfigDb) IsValid() bool {
	if nil != instance {
		return len(instance.Driver) > 0 && len(instance.Dsn) > 0
	}
	return false
}

func (instance *SemanticConfigDb) Parse(text string) {
	_ = qbc.JSON.Read(text, &instance)
}

//----------------------------------------------------------------------------------------------------------------------
//	SemanticConfig
//----------------------------------------------------------------------------------------------------------------------

type SemanticConfig struct {
	CaseSensitive bool              `json:"case_sensitive"`
	DbInternal    *SemanticConfigDb `json:"db_internal"` // internal storage
	DbExternal    *SemanticConfigDb `json:"db_external"` // db containing indexed data
}

func NewSemanticConfig() *SemanticConfig {
	instance := new(SemanticConfig)
	instance.CaseSensitive = false
	instance.DbInternal = new(SemanticConfigDb)
	instance.DbExternal = new(SemanticConfigDb)

	return instance
}

func (instance *SemanticConfig) String() string {
	return qbc.JSON.Stringify(instance)
}

func (instance *SemanticConfig) IsValid() bool {
	if nil != instance && nil != instance.DbInternal {
		return instance.DbInternal.IsValid()
	}
	return false
}

func (instance *SemanticConfig) Parse(data interface{}) error {
	if v, b := data.(string); b {
		return qbc.JSON.Read(v, &instance)
	} else {
		return qbc.JSON.Read(qbc.JSON.Stringify(data), &instance)
	}
}

package qb_dbal

import (
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/commons"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
	"github.com/rskvp/qb-lib/qb_dbal/semantic_search"
	"github.com/rskvp/qb-lib/qb_dbal/showcase_search"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func NewDatabase(driverName, connectionString string) (driver drivers.IDatabase, err error) {
	return drivers.NewDatabase(driverName, connectionString)
}

func NewDatabaseFromDsn(driverName string, dsn *commons.Dsn) (drivers.IDatabase, error) {
	return drivers.NewDatabaseFromDsn(driverName, dsn)
}

func OpenDatabase(driver, connectionString string) (drivers.IDatabase, error) {
	return drivers.OpenDatabase(driver, connectionString)
}

func NewSemanticEngine(c interface{}) (*semantic_search.SemanticEngine, error) {
	config, err := getConfig(c)
	if nil != err {
		return nil, err
	}
	return semantic_search.NewSemanticEngine(config)
}

func NewShowcaseEngine(c interface{}) (*showcase_search.ShowcaseEngine, error) {
	config, err := getConfigDB(c)
	if nil != err {
		return nil, err
	}
	return showcase_search.NewShowcaseEngine(config)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getConfig(c interface{}) (*commons.SemanticConfig, error) {
	var config *commons.SemanticConfig
	if v, b := c.(*commons.SemanticConfig); b {
		config = v
	} else if s, b := c.(string); b {
		if qbc.Regex.IsValidJsonObject(s) {
			err := qbc.JSON.Read(s, &config)
			if nil != err {
				return nil, err
			}
		} else {
			// file
			err := qbc.JSON.ReadFromFile(s, &config)
			if nil != err {
				return nil, err
			}
		}
	} else if m, b := c.(map[string]interface{}); b {
		err := qbc.JSON.Read(qbc.JSON.Stringify(m), &config)
		if nil != err {
			return nil, err
		}
	} else if mp, b := c.(*map[string]interface{}); b {
		err := qbc.JSON.Read(qbc.JSON.Stringify(mp), &config)
		if nil != err {
			return nil, err
		}
	}
	if nil == config {
		return nil, commons.ErrorMismatchConfiguration
	}
	return config, nil
}

func getConfigDB(c interface{}) (*commons.SemanticConfigDb, error) {
	var config *commons.SemanticConfigDb
	if v, b := c.(*commons.SemanticConfigDb); b {
		config = v
	} else if s, b := c.(string); b {
		if qbc.Regex.IsValidJsonObject(s) {
			err := qbc.JSON.Read(s, &config)
			if nil != err {
				return nil, err
			}
		} else {
			// file
			err := qbc.JSON.ReadFromFile(s, &config)
			if nil != err {
				return nil, err
			}
		}
	} else if m, b := c.(map[string]interface{}); b {
		err := qbc.JSON.Read(qbc.JSON.Stringify(m), &config)
		if nil != err {
			return nil, err
		}
	} else if mp, b := c.(*map[string]interface{}); b {
		err := qbc.JSON.Read(qbc.JSON.Stringify(mp), &config)
		if nil != err {
			return nil, err
		}
	}
	if nil == config {
		return nil, commons.ErrorMismatchConfiguration
	}
	return config, nil
}

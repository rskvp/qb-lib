package _test

import (
	"fmt"
	"testing"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/qb_dbal_commons"
	"github.com/rskvp/qb-lib/qb_dbal/qb_dbal_semantic_search"
)

const (
	COLL = "test"
)

func TestSearch(t *testing.T) {
	cfg, err := config()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	engine, err := qb_dbal_semantic_search.NewSemanticEngine(cfg)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	err = engine.Put(COLL, "001", "Hello, this is a text added to elaStic super elastiC search, very elastic!")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	err = engine.Put(COLL, "002", "Hello, this is something else with different keywords and matching score!!")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	err = engine.Put(COLL, "003", "Another elastiC here")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	cfg.CaseSensitive = false
	fmt.Println("CASE INSENSITIVE")
	data, err := engine.Get(COLL, "give me the Elastic", 0, -1)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Found: ", len(data))
	for _, item := range data {
		fmt.Println(item.Key, qbc.JSON.Stringify(item))
	}

	cfg.CaseSensitive = true
	fmt.Println("CASE SENSITIVE")
	data, err = engine.Get(COLL, "give me the elastiC", 0, -1)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Found: ", len(data))
	for _, item := range data {
		fmt.Println(item.Key, qbc.JSON.Stringify(item))
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func config() (*qb_dbal_commons.SemanticConfig, error) {
	text, err := qbc.IO.ReadTextFromFile("./dsn.json")
	if nil != err {
		return nil, err
	}
	response := qb_dbal_commons.NewSemanticConfig()
	response.DbInternal.Parse(text)
	response.DbExternal.Parse(text)

	return response, err
}

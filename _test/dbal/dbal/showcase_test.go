package _test

import (
	"fmt"
	"testing"
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/qb_dbal_commons"
	"github.com/rskvp/qb-lib/qb_dbal/qb_dbal_showcase_search"
)

func TestShowcase(t *testing.T) {
	cfg, err := loadConfiguration()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("TESTING SHOWCASE:", cfg)

	engine, err := qb_dbal_showcase_search.NewShowcaseEngine(cfg)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("ENGINE: ", engine)

	// put some data
	count := 0
	for _, category := range qb_dbal_commons.CATEGORIES {
		for i := 0; i < 10; i++ {
			collection := fmt.Sprintf("coll_%v", category)
			key := fmt.Sprintf("entity_%v_%v", category, i)
			payload := map[string]interface{}{
				"collection": collection,
				"key":        key,
				"mode":       "test",
				"category":   category,
			}
			_, _ = engine.Put(payload, time.Now().Unix(), category)

			count++
		}
	}
	fmt.Println("ADDED ITEMS:", count)

	// wait a moment
	time.Sleep(3 * time.Second)

	engine.SetCategoryWeight(qb_dbal_commons.CAT_EVENT, false, 10)
	weight := engine.Categories().Get(qb_dbal_commons.CAT_EVENT).WeightOutDate
	if weight != 10 {
		t.Error(fmt.Sprintf("Expected weight %v, got %v", 10, weight))
		t.FailNow()
	}
	engine.Reset()

	// get data for board
	session := "user1234"
	// change event weight for single user
	_ = engine.SetSessionCategoryWeight(session, qb_dbal_commons.CAT_EVENT, false, 10)
	engine.SetAutoResetSession(true)

	fmt.Println("START ----------")
	items := engine.Query(session, 50)
	for _, item := range items {
		fmt.Println(qbc.JSON.Stringify(item))
	}
	fmt.Println("END ----------", len(items))

	fmt.Println("START ----------")
	items = engine.Query(session, 50)
	for _, item := range items {
		fmt.Println(qbc.JSON.Stringify(item))
	}
	fmt.Println("END ----------", len(items))

	fmt.Println("START ----------")
	items = engine.Query(session, 50)
	if len(items) == 0 {
		t.Error("Expected some items")
		t.FailNow()
	}
	for _, item := range items {
		fmt.Println(qbc.JSON.Stringify(item))
	}
	fmt.Println("END ----------", len(items))

	// remove items one by one
	for _, item := range items {
		key := qbc.Reflect.GetString(item, "_key")
		entity, _ := engine.Delete(key)
		payload := qbc.Reflect.Get(entity, "payload")
		fmt.Println(qbc.JSON.Stringify(payload))
	}

	// reset the session
	engine.ResetSession(session)
	engine.SetAutoResetSession(false)
	fmt.Println("START ----------")
	items = engine.Query(session, 50)
	for _, item := range items {
		fmt.Println(qbc.JSON.Stringify(item))
	}
	fmt.Println("END ----------", len(items))
	fmt.Println("START ----------")
	items = engine.Query(session, 50)
	for _, item := range items {
		fmt.Println(qbc.JSON.Stringify(item))
	}
	fmt.Println("END ----------", len(items))

}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func loadConfiguration() (*qb_dbal_commons.SemanticConfigDb, error) {
	text, err := qbc.IO.ReadTextFromFile("./dsn.json")
	if nil != err {
		return nil, err
	}
	var response *qb_dbal_commons.SemanticConfigDb
	qbc.JSON.Read(text, &response)
	return response, err
}

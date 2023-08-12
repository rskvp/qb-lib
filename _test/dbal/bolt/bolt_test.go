package test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_stopwatch"
	"github.com/rskvp/qb-core/qb_utils"
	"github.com/rskvp/qb-lib/qb_dbal/qb_bolt"
)

type Entity struct {
	Key string `json:"_key"`
	Age int    `json:"age"`
}

type Token struct {
	Key string `json:"_key"`
	Exp int64  `json:"exp"`
}

func TestDatabase(t *testing.T) {

	drop_on_exit := false

	config, err := getConfig()
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	db := qb_bolt.NewBoltDatabase(config)
	err = db.Open()
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	defer db.Close()

	coll, err := db.Collection("my-coll", false)
	if nil != coll {
		t.Error("COLLECTION SHOULD BE NULL")
		t.Fail()
	}
	fmt.Println("Test OK:", err)

	coll, err = db.Collection("my-coll", true)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	// insert item
	item := &map[string]interface{}{
		"_key": "1",
		"name": "Mario",
		"age":  22,
	}
	err = coll.Upsert(item)
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	item2 := &map[string]interface{}{
		"_key": "2",
		"name": "Giorgio",
		"age":  22,
	}
	err = coll.Upsert(item2)
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	item3 := &map[string]interface{}{
		"_key": "3",
		"name": "Mirko",
		"age":  45,
	}
	err = coll.Upsert(item3)
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	count, err := coll.CountByFieldValue("age", 22)
	fmt.Println(count)

	item_des, err := coll.Get("1")
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(item_des)

	data, err := coll.GetByFieldValue("age", 22)
	fmt.Println(data)

	// remove collection
	err = coll.Drop()
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	coll, err = db.Collection("my-coll", false)
	if nil != coll {
		t.Error("COLLECTION SHOULD BE NULL")
		t.Fail()
	}

	// remove database
	if drop_on_exit {
		err = db.Drop()
		if nil != err {
			t.Error(err)
			t.Fail()
		}
	}
}

func TestQuery(t *testing.T) {
	item := &map[string]interface{}{
		"_key": "1",
		"name": "Mario",
		"age":  22,
	}

	query, err := qb_bolt.NewQueryFromFile("./query.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	match := query.MatchFilter(item)
	if !match {
		t.Error("Query do not match")
		t.FailNow()
	}

}

func TestFatData(t *testing.T) {

	drop_on_exit := true

	config, err := getConfig()
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	db := qb_bolt.NewBoltDatabase(config)
	err = db.Open()
	size, _ := db.Size()
	fmt.Println("FILE SIZE: ", qbc.Formatter.FormatBytes(uint64(size)))
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	defer db.Close()

	coll, err := db.Collection("big-coll", true)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	arrayData := generateArray()

	watch := qb_stopwatch.New()
	watch.Start()

	//-- START LOOP TO ADD RECORDS --//
	cur, _ := coll.Count()
	cur++
	for i := cur; i < cur+100; i++ {
		item := &map[string]interface{}{
			"_key": qbc.Strings.Format("%s", i),
			"name": "NAME " + qbc.Strings.Format("%s", i),
			"age":  i,
			"x":    arrayData,
			"y":    arrayData,
		}
		err = coll.Upsert(item)
		if nil != err {
			t.Error(err)
			t.Fail()
			break
		}
	}
	watch.Stop()
	fmt.Println("ELAPSED FOR CREATION: ", watch.Seconds(), "seconds")

	watch.Start()
	count, err := coll.Count()
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	watch.Stop()
	fmt.Println("RECORDS: ", count)
	fmt.Println("ELAPSED FOR COUNT: ", watch.Seconds(), "seconds")

	watch.Start()
	query, err := qb_bolt.NewQueryFromFile("./query.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	query.Conditions[0].Filters[0].Value = 1
	//var entity Entity
	//data, err := coll.Find(query, &entity)
	// data, err := coll.Find(query)
	data := make([]map[string]interface{}, 0)
	err = coll.ForEach(func(k, v []byte) bool {
		key := string(k)
		fmt.Println("KEY", key)
		var e Entity
		_ = json.Unmarshal(v, &e)
		if query.MatchFilter(e) {
			var m map[string]interface{}
			_ = json.Unmarshal(v, &m)
			data = append(data, m)
		}
		return false // continue
	})

	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	watch.Stop()
	fmt.Println("QUERY: ", len(data))
	fmt.Println("ELAPSED FOR QUERY: ", watch.Seconds(), "seconds")
	if len(data) > 0 {
		for _, item := range data {
			fmt.Println("AGE:", qbc.Reflect.Get(item, "age"))
		}
	}

	size, _ = db.Size()
	fmt.Println("FILE SIZE: ", qbc.Formatter.FormatBytes(uint64(size)))

	// remove database
	if drop_on_exit {
		err = db.Drop()
		if nil != err {
			t.Error(err)
			t.Fail()
		}
	}
}

// This test is for a JWT like application.
// Is BBolt good to store tokens?
func TestBigThinData(t *testing.T) {

	dropOnExit := false
	maxRecords := 1000

	filename := "./db/big_thin.dat"
	_ = qbc.Paths.Mkdir(filename)
	config := qb_bolt.NewBoltConfig()
	config.Name = filename
	watch := qb_stopwatch.New()
	var wg sync.WaitGroup

	// OPEN
	watch.Start()
	db := qb_bolt.NewBoltDatabase(config)
	err := db.Open()
	size, _ := db.Size()
	fmt.Println("FILE SIZE: ", qbc.Formatter.FormatBytes(uint64(size)))
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	watch.Stop()
	fmt.Println("ELAPSED FOR OPEN: ", watch.Seconds(), "seconds")

	defer db.Close()

	coll, err := db.Collection("big-coll", true)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	if n, _ := coll.Count(); n == 0 {
		//-- START LOOP TO ADD RECORDS --//
		watch.Start()
		for i := 0; i < maxRecords; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				item := &map[string]interface{}{
					"_key": qbc.Rnd.Uuid(),
					"name": "NAME " + qbc.Strings.Format("%s", i),
					"exp":  time.Now().Add(15 * time.Second).Unix(),
				}
				err = coll.Upsert(item)
				if nil != err {
					t.Error(err)
					t.Fail()
				}
			}()
		}
		wg.Wait()
		watch.Stop()
		fmt.Println("ELAPSED FOR CREATION: ", watch.Seconds(), "seconds")
	}

	// COUNT
	watch.Start()
	count, err := coll.Count()
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	watch.Stop()
	fmt.Println("RECORDS: ", count)
	fmt.Println("ELAPSED FOR COUNT: ", watch.Seconds(), "seconds")

	// GET EXPIRED
	watch.Start()
	expired := make([]string, 0)
	err = coll.ForEach(func(k, v []byte) bool {
		key := string(k)
		var e Token
		_ = json.Unmarshal(v, &e)
		if time.Now().Unix()-e.Exp > 0 {
			expired = append(expired, key)
		}
		return false // continue
	})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	watch.Stop()
	fmt.Println("EXPIRED: ", len(expired))
	fmt.Println("ELAPSED FOR QUERY: ", watch.Seconds(), "seconds")

	size, _ = db.Size()
	fmt.Println("FILE SIZE: ", qbc.Formatter.FormatBytes(uint64(size)))

	// remove expired
	fmt.Println("REMOVING:", len(expired))
	if len(expired) > 0 {
		watch.Start()
		for _, key := range expired {
			wg.Add(1)
			go func(key string) {
				err = coll.Remove(key)
				if nil != err {
					fmt.Println("ERROR REMOVING:", err)
				}
				wg.Done()
			}(key)
		}
		wg.Wait()
		watch.Stop()
		fmt.Println("ELAPSED FOR REMOVE: ", watch.Seconds(), "seconds")
	}

	size, _ = db.Size()
	fmt.Println("FILE SIZE: ", qbc.Formatter.FormatBytes(uint64(size)))

	count, _ = coll.Count()
	fmt.Println("RECORDS: ", count)

	// remove database
	if dropOnExit {
		err = db.Drop()
		if nil != err {
			t.Error(err)
			t.Fail()
		}
	}
}

func TestExpireData(t *testing.T) {
	dropOnExit := false
	maxRecords := 1000

	filename := "./db/expiring.dat"
	_ = qbc.Paths.Mkdir(filename)
	config := qb_bolt.NewBoltConfig()
	config.Name = filename
	watch := qb_stopwatch.New()

	// OPEN
	watch.Start()
	db := qb_bolt.NewBoltDatabase(config)
	err := db.Open()
	size, _ := db.Size()
	fmt.Println("FILE SIZE: ", qb_utils.Formatter.FormatBytes(uint64(size)))
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	watch.Stop()
	fmt.Println("ELAPSED FOR OPEN: ", watch.Seconds(), "seconds")

	defer db.Close()

	coll, err := db.Collection("big-coll", true)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	coll.EnableExpire(true)

	// ADD
	if n, _ := coll.Count(); n == 0 {
		//-- START LOOP TO ADD RECORDS --//
		watch.Start()
		for i := 0; i < maxRecords; i++ {
			item := map[string]interface{}{
				"_key": qbc.Rnd.Uuid(),
				"name": "NAME " + qbc.Strings.Format("%s", i),
			}
			item[qb_bolt.FieldExpire] = time.Now().Add(5 * time.Second).Unix()
			err = coll.Upsert(item)
			if nil != err {
				t.Error(err)
				t.Fail()
			}
		}
		watch.Stop()
		fmt.Println("ELAPSED FOR ADD: ", watch.Seconds(), "seconds")
	}

	for {
		time.Sleep(2 * time.Second)
		n, _ := coll.Count()
		if n == 0 {
			break
		} else {
			fmt.Println(n)
		}
	}

	count, _ := coll.Count()
	fmt.Println("RECORDS: ", count)

	// remove database
	if dropOnExit {
		err = db.Drop()
		if nil != err {
			t.Error(err)
			t.Fail()
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getConfig() (*qb_bolt.BoltConfig, error) {
	text_cfg, err := qbc.IO.ReadTextFromFile("./config.json")
	if nil != err {
		return nil, err
	}
	config := qb_bolt.NewBoltConfig()
	err = config.Parse(text_cfg)

	return config, err
}

func generateArray() []float64 {
	size := 250 * 60 * 10 // 10 minutes of 250Hz data sequence
	data := make([]float64, size)
	for i := 0; i < size; i++ {
		val := float64(i * 2)
		data[i] = val
	}
	return data
}

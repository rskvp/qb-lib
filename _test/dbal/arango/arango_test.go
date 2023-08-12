package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/arangodb/go-driver"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/qb_arango"
)

func TestSimple(t *testing.T) {

	ctext, _ := qbc.IO.ReadTextFromFile("config.json")
	psw, _ := qbc.IO.ReadTextFromFile("psw.txt")
	config := qb_arango.NewArangoConfig()
	config.Parse(ctext)
	config.Authentication.Password = psw

	conn := qb_arango.NewArangoConnection(config)
	err := conn.Open()
	if nil != err {
		// fmt.Println(err)
		t.Error(err, ctext)
		t.Fail()
		return
	}
	// print version
	fmt.Println("ARANGO SERVER", conn.Server)
	fmt.Println("ARANGO VERSION", conn.Version)
	fmt.Println("ARANGO LICENSE", conn.License)

	// remove
	conn.DropDatabase("test_sample")

	// create a db
	db, err := conn.Database("test_sample", true)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}

	if nil != db {
		fmt.Println(db.Name())
	}

	coll, err := db.Collection("not_exists", true)
	if nil != err {
		t.Error(err)
	}
	if nil == coll {
		t.Fail()
	}

	// entity
	entity := map[string]interface{}{
		"_key":    "258647",
		"name":    "Angelo",
		"surname": "Geminiani",
	}

	Key := qbc.Reflect.GetString(entity, "Key")
	fmt.Println("KEY", Key)

	qbc.Maps.Set(entity, "Name", "Gian Angelo")

	doc, meta, err := coll.Upsert(entity)
	if nil != err {
		t.Error(err)
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)

	entity = map[string]interface{}{
		"_key":    qbc.Rnd.Uuid(),
		"name":    "Marco",
		"surname": qbc.Strings.Format("%s", time.Now()),
	}
	doc, meta, err = coll.Upsert(entity)
	if nil != err {
		t.Error(err)
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)

	// bew entity that test upsert used for insert
	newEntity := map[string]interface{}{
		"_key":    qbc.Rnd.Uuid(),
		"name":    "I'm new",
		"surname": qbc.Strings.Format("%s", time.Now()),
	}
	doc, meta, err = coll.Upsert(newEntity)
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)

	doc, meta, err = coll.Read(qbc.Convert.ToString(doc["_key"]))
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)

	// remove
	removed, err := conn.DropDatabase("test_sample")
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}
	if removed {
		fmt.Println("REMOVED", "test_sample")
	}
}

func TestInsert(t *testing.T) {
	ctext, _ := qbc.IO.ReadTextFromFile("config.json")
	config := qb_arango.NewArangoConfig()
	config.Parse(ctext)

	conn := qb_arango.NewArangoConnection(config)
	err := conn.Open()
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
		t.Fail()
		return
	}

	// remove
	conn.DropDatabase("test_sample")

	db, err := conn.Database("test_sample", true)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}
	coll, err := db.Collection("coll_insert", true)
	if nil != err {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		// entity
		entity := map[string]interface{}{
			"_key":    qbc.Strings.Format("key_%s", i),
			"name":    qbc.Strings.Format("Name:%s", i),
			"surname": qbc.Strings.Format("Surname:%s", i),
			"address": qbc.Strings.Format("Address:%s", i),
		}

		doc, meta, err := coll.Insert(entity)
		if nil != err {
			t.Error(err)
		}
		fmt.Println("META", meta)
		fmt.Println("DOC", doc)
	}

	updEntity := map[string]interface{}{
		"_key": "key_1",
		"name": "Gian Angelo",
	}
	_, _, err = coll.Update(updEntity)
	if nil != err {
		t.Error(err)
	}

	noKeyEntity := map[string]interface{}{
		"name":    "NO KEY",
		"surname": "ZERO ZERO",
		"address": "",
	}
	doc, meta, err := coll.Insert(noKeyEntity)
	if nil != err {
		t.Error(err)
	}
	fmt.Println("---------------------")
	fmt.Println("META", meta)
	fmt.Println("DOC KEY", doc["_key"])
	fmt.Println("DOC", qbc.Convert.ToString(doc))
	fmt.Println("---------------------")

	query := "FOR d IN coll_insert RETURN d"
	db.Query(query, nil, gotDocument)

	fmt.Println("---------------------")

}

func TestImport(t *testing.T) {
	ctext, _ := qbc.IO.ReadTextFromFile("config.json")
	config := qb_arango.NewArangoConfig()
	config.Parse(ctext)

	conn := qb_arango.NewArangoConnection(config)
	err := conn.Open()
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
		t.Fail()
		return
	}

	// remove
	conn.DropDatabase("test_sample")

	db, err := conn.Database("test_sample", true)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}

	err = db.ImportFile("./toImport.json")
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}

	err = db.ImportFile("./toImport.csv")
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}

	coll, err := db.Collection("toImport", false)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}
	_, err = coll.EnsureIndex([]string{"name"}, false)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func gotDocument(meta driver.DocumentMeta, doc interface{}, err error) bool {
	fmt.Print("META: ", meta, " ENTITY: ", qbc.Convert.ToString(doc), " ERR: ", err)
	m := qbc.Convert.ToMap(doc)
	if b, _ := qbc.Compare.IsMap(m); b {
		fmt.Printf(" IS MAP: %v\n", b)
	}
	return false // continue
}

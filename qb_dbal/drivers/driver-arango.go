package drivers

import (
	"context"
	"fmt"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	qbc "github.com/rskvp/qb-core"
	dbalcommons "github.com/rskvp/qb-lib/qb_dbal/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

var (
	ArangoConst = struct {
		KeyFieldName string
		IndexPersist string
		IndexGeo     string
		IndexGeoJson string
	}{
		KeyFieldName: "_key",
		IndexPersist: "persist",
		IndexGeo:     "geo",
		IndexGeoJson: "geojson",
	}
)

const KeyFieldName = "_key" // all entities should have this field

const NameArango = "arango"

//----------------------------------------------------------------------------------------------------------------------
//	NewDriverArango
//----------------------------------------------------------------------------------------------------------------------

type DriverArango struct {
	uid            string
	dsn            *dbalcommons.Dsn
	enableCache    bool
	connection     driver.Connection
	authentication driver.Authentication
	client         driver.Client
	version        driver.VersionInfo
	db             driver.Database
	err            error
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDriverArango(dsn ...interface{}) *DriverArango {
	instance := new(DriverArango)

	if len(dsn) == 1 {
		if s, b := dsn[0].(string); b {
			instance.dsn = dbalcommons.NewDsn(s)
		} else if d, b := dsn[0].(dbalcommons.Dsn); b {
			instance.dsn = &d
		} else if d, b := dsn[0].(*dbalcommons.Dsn); b {
			instance.dsn = d
		} else {
			instance.err = dbalcommons.ErrorInvalidDsn
		}
	}
	if nil == instance.dsn && nil == instance.err {
		instance.err = dbalcommons.ErrorInvalidDsn
	}
	if nil != instance.dsn {
		instance.uid = keyFrom(NameArango, instance.dsn.String())
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArango) Uid() string {
	return instance.uid
}

func (instance *DriverArango) DriverName() string {
	return NameArango
}

func (instance *DriverArango) Enabled() bool {
	return nil != instance && nil != instance.dsn && nil == instance.err && instance.dsn.IsValid()
}

func (instance *DriverArango) Open() error {
	if nil != instance {
		if nil == instance.err {
			instance.err = instance.init()
		}
		return instance.err
	}
	return nil
}

func (instance *DriverArango) Close() error {
	if nil != instance.db {
		//1empty
	}
	return nil
}

func (instance *DriverArango) Remove(collection, key string) error {
	if nil != instance && nil != instance.db {
		coll, err := instance.Collection(collection, true)
		if nil != err {
			return err
		}
		_, err = coll.Remove(key)
		return err
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Get(collection string, key string) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		coll, err := instance.Collection(collection, true)
		if nil != err {
			return nil, err
		}
		item, _, err := coll.Read(key)
		if nil != err {
			if e, b := err.(driver.ArangoError); b && e.Code == 404 {
				// not found
				return nil, nil
			}
		}
		return item, err
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Upsert(collection string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		key := qbc.Reflect.GetString(item, ArangoConst.KeyFieldName)
		if len(key) > 0 {
			exists, err := instance.Exists(collection, key)
			if nil != err {
				return nil, err
			} else {
				if exists {
					return instance.Update(collection, item)
				} else {
					return instance.Insert(collection, item)
				}
			}
		} else {
			return instance.Insert(collection, item)
		}
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) ForEach(collectionOrQuery string, callback ForEachCallback) error {
	if nil != instance && nil != instance.db {
		if nil != callback {
			query, err := instance.buildQuery(collectionOrQuery)
			if nil != err {
				return err
			}
			ctx := context.Background()
			cursor, err := instance.db.Query(ctx, query, nil)
			if nil != err {
				return err
			}
			defer cursor.Close()

			// run cursor
			for {
				var doc map[string]interface{}
				_, err := cursor.ReadDocument(ctx, &doc)
				if driver.IsNoMoreDocuments(err) {
					break
				} else {
					if nil != err {
						return err
					}
					if nil != doc {
						exit := callback(doc)
						if exit {
							break
						}
					}
				}
			}
		}
		return nil // do nothing
	}
	return dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) ExecNative(command string, bindVars map[string]interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {
		ctx := context.Background()
		return instance.exec(ctx, command, bindVars)
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) ExecMultiple(commands []string, bindVars []map[string]interface{}, options interface{}) ([]interface{}, error) {
	if nil != instance && nil != instance.db {
		if len(commands) != len(bindVars) {
			return nil, dbalcommons.ErrorCommandAndParamsDoNotMatch
		}
		ctx := context.Background()
		transactional := false

		var cols driver.TransactionCollections
		if nil != options {
			err := qbc.JSON.Read(qbc.JSON.Stringify(options), &cols)
			if nil != err {
				return nil, err
			}
			transactional = true
		}

		response := make([]interface{}, 0)
		// BEGIN
		var id driver.TransactionID
		if transactional {
			tid, err := instance.db.BeginTransaction(ctx, cols, nil)
			if nil != err {
				return nil, err
			}
			id = tid
		}
		// execute
		for i := 0; i < len(commands); i++ {
			command := commands[i]
			bindVar := bindVars[i]
			data, err := instance.exec(ctx, command, bindVar)
			if nil != err {
				// ABORT
				if transactional {
					_ = instance.db.AbortTransaction(ctx, id, nil)
				}
				return response, err
			}
			response = append(response, data)
		}
		// COMMIT
		if transactional {
			err := instance.db.CommitTransaction(ctx, id, nil)
			if nil != err {
				return response, err
			}
		}

		return response, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) EnsureIndex(collection string, typeName string, fields []string, unique bool) (bool, error) {
	if nil != instance && nil != instance.db {
		coll, err := instance.Collection(collection, true)
		if nil != err {
			return false, err
		}
		if len(fields) > 0 {
			if len(typeName) == 0 {
				typeName = ArangoConst.IndexPersist
			}
			return coll.EnsureIndex(typeName, fields, unique)
		}
		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) EnsureCollection(collection string) (bool, error) {
	if nil != instance && nil != instance.db {
		_, err := instance.Collection(collection, true)
		if nil != err {
			return false, err
		}
		return true, nil
	}
	return false, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Find(collection string, fieldName string, fieldValue interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {
		query := fmt.Sprintf("FOR doc IN @@collection FILTER doc.%v == @%v RETURN doc", fieldName, fieldName)
		params := map[string]interface{}{
			"@collection": collection,
		}
		params[fieldName] = fieldValue
		return instance.ExecNative(query, params)
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) QueryGetParamNames(query string) []string {
	return QueryGetParamNames(query)
}

func (instance *DriverArango) QuerySelectParams(query string, allParams map[string]interface{}) map[string]interface{} {
	return QuerySelectParams(query, allParams)
}

//----------------------------------------------------------------------------------------------------------------------
//	e x t e n d
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArango) Exec(query string, bindVars map[string]interface{}) (interface{}, error) {
	if nil != instance && nil != instance.db {
		response := make([]interface{}, 0)
		ctx := context.Background()
		cursor, err := instance.db.Query(ctx, query, bindVars)
		if nil != err {
			return nil, err
		}
		defer cursor.Close()

		// run cursor
		for {
			var doc interface{}
			_, err := cursor.ReadDocument(ctx, &doc)
			if driver.IsNoMoreDocuments(err) {
				break
			} else {
				if nil != err {
					return nil, err
				}
				if nil != doc {
					response = append(response, doc)
				}
			}
		}

		return response, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Exists(collection, key string) (bool, error) {
	bindVars := map[string]interface{}{
		"@collection": collection,
		"key":         key,
	}
	query := "FOR doc IN @@collection\n" +
		"FILTER doc.`_key`==@key\n" +
		"RETURN doc._key"
	count, err := instance.Count(query, bindVars)
	if nil != err {
		return false, err
	}
	return count > 0, err
}

func (instance *DriverArango) Count(query string, bindVars map[string]interface{}) (int64, error) {
	if nil != instance && nil != instance.db {
		ctx := context.Background()
		cursor, err := instance.db.Query(ctx, query, bindVars)
		if nil != err {
			return 0, err
		}

		defer cursor.Close()
		var count int64
		for {
			var doc interface{}
			_, err := cursor.ReadDocument(ctx, &doc)
			if driver.IsNoMoreDocuments(err) {
				break
			} else {
				if nil != doc {
					count++
				}
			}
		}
		return count, nil
	}
	return 0, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Insert(collectionName string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		collection, err := instance.collection(collectionName, true)
		if nil != err {
			return nil, err
		}

		ctx := context.Background()
		meta, err := collection.CreateDocument(ctx, item)
		if nil != err {
			return nil, err
		}
		var doc map[string]interface{}
		_, err = collection.ReadDocument(ctx, meta.Key, &doc)
		if nil != err {
			return nil, err
		}
		return doc, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Update(collectionName string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		collection, err := instance.collection(collectionName, true)
		if nil != err {
			return nil, err
		}

		ctx := context.Background()
		key := qbc.Reflect.GetString(item, ArangoConst.KeyFieldName)
		meta, err := collection.UpdateDocument(ctx, key, item)
		if nil != err {
			if e, b := err.(driver.ArangoError); b && e.Code == 404 {
				return instance.Insert(collectionName, item)
			}
			return nil, err
		}
		var doc map[string]interface{}
		_, err = collection.ReadDocument(ctx, meta.Key, &doc)
		if nil != err {
			return nil, err
		}
		return doc, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Delete(collectionName string, item interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		collection, err := instance.collection(collectionName, true)
		if nil != err {
			return nil, err
		}

		ctx := context.Background()
		var key string
		if v, b := item.(string); b {
			key = v
		} else {
			key = qbc.Reflect.GetString(item, ArangoConst.KeyFieldName)
		}

		var doc map[string]interface{}
		_, err = collection.ReadDocument(ctx, key, &doc)
		if nil != err {
			return nil, err
		}
		_, err = collection.RemoveDocument(ctx, key)
		if nil != err {
			return nil, err
		}

		return doc, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Collection(name string, createIfDoesNotExists bool) (*DriverArangoCollection, error) {
	coll, err := instance.collection(name, createIfDoesNotExists)
	if nil != err {
		return nil, err
	}
	if nil != coll {
		return &DriverArangoCollection{
			name:       name,
			collection: coll,
		}, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

// EnsureAnalyzer https://www.arangodb.com/docs/stable/analyzers.html
// https://www.arangodb.com/docs/stable/arangosearch-case-sensitivity-and-diacritics.html
func (instance *DriverArango) EnsureAnalyzer(definition map[string]interface{}) (err error) {
	if nil != instance {
		name := qbc.Reflect.GetString(definition, "name")
		if len(name) > 0 {
			ctx := context.Background()
			// remove existing
			analyzer, _ := instance.db.Analyzer(ctx, name)
			if nil != analyzer {
				_ = analyzer.Remove(ctx, true)
			}
			// ensure new one
			var settings driver.ArangoSearchAnalyzerDefinition
			err = qbc.JSON.Read(qbc.Convert.ToString(definition), &settings)
			if nil != err {
				return
			}

			_, _, err = instance.db.EnsureAnalyzer(ctx, settings)
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArango) init() error {
	// CONNECTION
	// "http://localhost:8529"
	endpoints := []string{fmt.Sprintf("%v://%v:%v", instance.dsn.Protocol, instance.dsn.Host, instance.dsn.Port)}
	connection, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: endpoints,
		TLSConfig: nil,
	})
	if nil != err {
		return err
	}
	instance.connection = connection

	// AUTHENTICATION
	if len(instance.dsn.User) > 0 && len(instance.dsn.Password) > 0 {
		// BASIC
		instance.authentication = driver.BasicAuthentication(instance.dsn.User, instance.dsn.Password)
	}

	// CLIENT
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     instance.connection,
		Authentication: instance.authentication,
	})
	if nil != err {
		return err
	}
	instance.client = c

	// TEST CONNECTION
	ctx := context.Background()
	v, err := c.Version(ctx)
	if nil != err {
		return err
	}
	instance.version = v

	// DATABASE
	if len(instance.dsn.Database) > 0 {
		db, err := instance.database(instance.dsn.Database, true)
		if nil != err {
			return err
		}
		instance.db = db
	}

	return nil
}

func (instance *DriverArango) database(name string, createIfNotExists bool) (driver.Database, error) {
	ctx := context.Background()
	exists, err := instance.client.DatabaseExists(ctx, name)
	if nil != err {
		return nil, err
	}
	if !exists && createIfNotExists {
		return instance.client.CreateDatabase(ctx, name, nil)
	} else {
		return instance.client.Database(ctx, name)
	}
}

func (instance *DriverArango) collection(name string, createIfNotExists bool) (driver.Collection, error) {
	if nil != instance && nil != instance.db {
		ctx := context.Background()
		exists, err := instance.db.CollectionExists(ctx, name)
		if nil != err {
			return nil, err
		}

		if !exists && createIfNotExists {
			_, err := instance.db.CreateCollection(ctx, name, nil)
			if nil != err {
				return nil, err
			}
		}

		collection, err := instance.db.Collection(ctx, name)
		if nil != err {
			return nil, err
		}
		if nil == collection {
			return nil, dbalcommons.ErrorCollectionDoesNotExists
		}
		return collection, nil
	}
	return nil, dbalcommons.ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) buildQuery(collectionOrQuery string) (string, error) {
	if strings.Index(strings.ToLower(collectionOrQuery), "for") > -1 {
		return collectionOrQuery, nil
	} else {
		_, err := instance.Collection(collectionOrQuery, true)
		if nil != err {
			return "", err
		}
		// "FOR d IN " + collection + " RETURN d"
		return fmt.Sprintf("FOR d IN %v RETURN d", collectionOrQuery), nil
	}
}

func (instance *DriverArango) exec(ctx context.Context, command string, bindVars map[string]interface{}) (interface{}, error) {
	cursor, err := instance.db.Query(ctx, command, bindVars)
	if nil != err {
		return nil, err
	}
	defer cursor.Close()

	response := make([]interface{}, 0)
	// run cursor
	for {
		var doc interface{}
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else {
			if nil != err {
				return nil, err
			}
			if nil != doc {
				response = append(response, doc)
			}
		}
	}
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	DriverArangoCollection
//----------------------------------------------------------------------------------------------------------------------

type DriverArangoCollection struct {
	name       string
	collection driver.Collection
}

func (instance *DriverArangoCollection) Name() string {
	return instance.name
}

//-- indexes --//

func (instance *DriverArangoCollection) RemoveIndex(typeName string, fields []string) (bool, error) {
	if nil != instance && nil != instance.collection {
		ctx := context.Background()

		name := instance.getIndexName(typeName, fields)
		index, err := instance.collection.Index(ctx, name)
		if nil != err {
			return false, err
		}
		if nil != index {
			err = index.Remove(ctx)
		}
		return nil != index, err
	}
	return false, dbalcommons.ErrorCollectionDoesNotExists
}

func (instance *DriverArangoCollection) EnsureIndex(typeName string, fields []string, unique bool) (bool, error) {
	if nil != instance && nil != instance.collection {

		// remove existing
		_, _ = instance.RemoveIndex(typeName, fields)

		switch typeName {
		case ArangoConst.IndexPersist:
			ctx := context.Background()
			options := &driver.EnsurePersistentIndexOptions{
				Name:   instance.getIndexName(typeName, fields),
				Unique: unique,
			}
			_, b, err := instance.collection.EnsurePersistentIndex(ctx, fields, options)
			if nil != err {
				return false, err
			}
			return b, err
		case ArangoConst.IndexGeo:
			ctx := context.Background()
			options := &driver.EnsureGeoIndexOptions{
				Name:    instance.getIndexName(typeName, fields),
				GeoJSON: false,
			}
			_, b, err := instance.collection.EnsureGeoIndex(ctx, fields, options)
			if nil != err {
				return false, err
			}
			return b, err
		case ArangoConst.IndexGeoJson:
			ctx := context.Background()
			options := &driver.EnsureGeoIndexOptions{
				Name:    instance.getIndexName(typeName, fields),
				GeoJSON: true,
			}
			_, b, err := instance.collection.EnsureGeoIndex(ctx, fields, options)
			if nil != err {
				return false, err
			}
			return b, err
		}

	}
	return false, dbalcommons.ErrorCollectionDoesNotExists
}

//-- data --//

func (instance *DriverArangoCollection) Count() (int64, error) {
	if nil != instance && nil != instance.collection {
		ctx := context.Background()
		return instance.collection.Count(ctx)
	}
	return -1, dbalcommons.ErrorCollectionDoesNotExists
}

func (instance *DriverArangoCollection) Exists(key string) (bool, error) {
	if nil != instance && nil != instance.collection {
		if len(key) > 0 {
			ctx := context.Background()
			b, err := instance.collection.DocumentExists(ctx, key)
			if nil != err {
				return false, err
			}
			return b, nil
		}
		return false, nil
	}
	return false, dbalcommons.ErrorCollectionDoesNotExists
}

func (instance *DriverArangoCollection) Read(key string) (map[string]interface{}, driver.DocumentMeta, error) {
	if nil != instance && nil != instance.collection {
		ctx := context.Background()
		var doc map[string]interface{}
		meta, err := instance.collection.ReadDocument(ctx, key, &doc)
		if nil != err {
			return nil, driver.DocumentMeta{}, err
		}
		return doc, meta, nil
	}
	return nil, driver.DocumentMeta{}, dbalcommons.ErrorCollectionDoesNotExists
}

func (instance *DriverArangoCollection) Remove(key string) (driver.DocumentMeta, error) {
	if nil != instance && nil != instance.collection {
		ctx := context.Background()
		meta, err := instance.collection.RemoveDocument(ctx, key)
		if nil != err {
			return driver.DocumentMeta{}, err
		}
		return meta, nil
	}
	return driver.DocumentMeta{}, dbalcommons.ErrorCollectionDoesNotExists
}

func (instance *DriverArangoCollection) Upsert(doc map[string]interface{}) (map[string]interface{}, driver.DocumentMeta, error) {
	if nil != instance && nil != instance.collection {
		if nil != doc {

			key := qbc.Reflect.GetString(doc, KeyFieldName)
			exists, err := instance.Exists(key)
			if nil != err {
				return nil, driver.DocumentMeta{}, err
			}

			ctx := context.Background()
			if !exists {
				// create
				meta, err := instance.collection.CreateDocument(ctx, doc)
				if nil != err {
					return nil, driver.DocumentMeta{}, err
				}
				// check if any key was added
				if instance.addKey(&meta, doc) {
					// update and replace doc
					_, err = instance.collection.ReadDocument(ctx, meta.Key, &doc) // uses doc address
					if nil != err {
						return nil, driver.DocumentMeta{}, err
					}
				}
			} else {
				// update
				_, err := instance.collection.UpdateDocument(ctx, key, doc)
				if nil != err {
					return nil, driver.DocumentMeta{}, err
				}
			}

			// read using a pointer to original doc
			meta, err := instance.collection.ReadDocument(ctx, key, &doc)
			return doc, meta, err
		}
		return nil, driver.DocumentMeta{}, nil
	}
	return nil, driver.DocumentMeta{}, dbalcommons.ErrorCollectionDoesNotExists
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArangoCollection) getIndexName(prefix string, fields []string) string {
	a := qbc.Convert.ToArray(fields)
	name := prefix + "_" + qbc.Strings.ConcatSep("_", a...)
	return "idx_" + qbc.Coding.MD5(name)
}

func (instance *DriverArangoCollection) addKey(meta *driver.DocumentMeta, doc map[string]interface{}) bool {
	key := qbc.Reflect.GetString(doc, ArangoConst.KeyFieldName)
	if len(key) == 0 {
		key = meta.Key
		if len(key) > 0 {
			qbc.Maps.Set(doc, ArangoConst.KeyFieldName, key)
			return true
		}
	}
	return false
}

func (instance *DriverArangoCollection) ensureKey(doc map[string]interface{}) bool {
	key := qbc.Reflect.GetString(doc, ArangoConst.KeyFieldName)
	if len(key) == 0 {
		key = qbc.Rnd.Uuid()
		if len(key) > 0 {
			qbc.Maps.Set(doc, ArangoConst.KeyFieldName, key)
			return true
		}
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

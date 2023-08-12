package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/bolt"
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

//----------------------------------------------------------------------------------------------------------------------
//	NewDriverArango
//----------------------------------------------------------------------------------------------------------------------

type DriverArango struct {
	dsn            *Dsn
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

func NewDriverArango(dsn *Dsn) *DriverArango {
	instance := new(DriverArango)
	instance.dsn = dsn

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArango) EnableCache(value bool) {
	instance.enableCache = value
	// TODO: enable cache exp for arango db
	if value {
		fmt.Println("WARNING!!! ARANGO DB DOES NOT SUPPORT CACHE EXPIRE")
	}
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
	}
	return nil
}

func (instance *DriverArango) AuthRegister(key, payload string) error {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}
		coll, err := instance.Collection(CollectionAuth, true)
		if nil != err {
			return err
		}
		if i, _, _ := coll.Read(key); nil != i {
			return ErrorEntityAlreadyRegistered
		}
		item := map[string]interface{}{
			"_key":    key,
			"payload": payload,
		}
		_, _, err = coll.Upsert(item)
		return err
	}
	return nil
}

func (instance *DriverArango) AuthOverwrite(key, payload string) error {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}
		coll, err := instance.Collection(CollectionAuth, true)
		if nil != err {
			return err
		}

		item := map[string]interface{}{
			"_key":    key,
			"payload": payload,
		}
		_, _, err = coll.Upsert(item)
		return err
	}
	return nil
}

func (instance *DriverArango) AuthGet(key string) (payload string, err error) {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return "", ErrorDatabaseCacheCannotAuthenticate
		}
		var coll *DriverArangoCollection
		coll, err = instance.Collection(CollectionAuth, true)
		if nil == err {
			var i interface{}
			i, _, err = coll.Read(key)
			if nil == err {
				if nil != i {
					payload = qbc.Reflect.GetString(i, "payload")
				}
			} else {
				err = ErrorEntityDoesNotExists
			}
		}
	}
	return payload, err
}

func (instance *DriverArango) AuthRemove(key string) error {
	if nil != instance && nil != instance.db {
		if instance.enableCache {
			return ErrorDatabaseCacheCannotAuthenticate
		}
		coll, err := instance.Collection(CollectionAuth, true)
		if nil != err {
			return err
		}
		_, err = coll.Remove(key)
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	c a c h e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArango) CacheGet(key string) (string, error) {
	if nil != instance && nil != instance.db {
		if !instance.enableCache {
			return "", ErrorDatabaseCacheNotEnabled
		}

		coll, err := instance.Collection(CollectionCache, true)
		if nil != err {
			return "", err
		}
		item, _, err := coll.Read(key)
		if nil != err {
			return "", err
		}
		if nil != item {
			var token string
			var err error
			now := int(time.Now().Unix())
			expire := qbc.Convert.ToInt(item[bolt.FieldExpire])
			token = item["token"].(string)
			if now-expire > 0 {
				// expired
				err = ErrorTokenExpired
			}
			return token, err
		} else {
			// not found or expired
			return "", ErrorTokenDoesNotExists
		}
	}
	return "", ErrorTokenDoesNotExists
}

func (instance *DriverArango) CacheAdd(key, token string, duration time.Duration) error {
	if nil != instance && nil != instance.db {
		if !instance.enableCache {
			return ErrorDatabaseCacheNotEnabled
		}
		coll, err := instance.Collection(CollectionCache, true)
		if nil != err {
			return err
		}
		item := map[string]interface{}{
			"_key":  key,
			"token": token,
		}
		item[bolt.FieldExpire] = time.Now().Add(duration).Unix()
		_, _, err = coll.Upsert(item)
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *DriverArango) CacheRemove(key string) error {
	if nil != instance && nil != instance.db {
		if !instance.enableCache {
			return ErrorDatabaseCacheNotEnabled
		}
		coll, err := instance.Collection(CollectionCache, true)
		if nil != err {
			return err
		}
		_, err = coll.Remove(key)
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	e x t e n d
//----------------------------------------------------------------------------------------------------------------------

func (instance *DriverArango) Query(query string, bindVars map[string]interface{}) ([]interface{}, error) {
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
				if nil != doc {
					response = append(response, doc)
				}
			}
		}

		return response, nil
	}
	return nil, ErrorDatabaseDoesNotExists
}

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
				if nil != doc {
					response = append(response, doc)
				}
			}
		}

		return response, nil
	}
	return nil, ErrorDatabaseDoesNotExists
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
	return 0, ErrorDatabaseDoesNotExists
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
	return nil, ErrorDatabaseDoesNotExists
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
			return nil, err
		}
		var doc map[string]interface{}
		_, err = collection.ReadDocument(ctx, meta.Key, &doc)
		if nil != err {
			return nil, err
		}
		return doc, nil
	}
	return nil, ErrorDatabaseDoesNotExists
}

func (instance *DriverArango) Upsert(collectionName string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance && nil != instance.db {
		key := qbc.Reflect.GetString(item, ArangoConst.KeyFieldName)
		if len(key) > 0 {
			return instance.Update(collectionName, item)
		} else {
			return instance.Insert(collectionName, item)
		}
	}
	return nil, ErrorDatabaseDoesNotExists
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
	return nil, ErrorDatabaseDoesNotExists
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
	return nil, ErrorDatabaseDoesNotExists
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
			return nil, ErrorCollectionDoesNotExists
		}
		return collection, nil
	}
	return nil, ErrorDatabaseDoesNotExists
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
	return false, ErrorCollectionDoesNotExists
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
	return false, ErrorCollectionDoesNotExists
}

//-- data --//

func (instance *DriverArangoCollection) Count() (int64, error) {
	if nil != instance && nil != instance.collection {
		ctx := context.Background()
		return instance.collection.Count(ctx)
	}
	return -1, ErrorCollectionDoesNotExists
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
	return false, ErrorCollectionDoesNotExists
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
	return nil, driver.DocumentMeta{}, ErrorCollectionDoesNotExists
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
	return driver.DocumentMeta{}, ErrorCollectionDoesNotExists
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
	return nil, driver.DocumentMeta{}, ErrorCollectionDoesNotExists
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

package arango

import (
	"context"

	"github.com/arangodb/go-driver"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type ArangoCollection struct {
	collection driver.Collection
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoCollection) Native() driver.Collection {
	if nil != instance && instance.IsReady() {
		return instance.collection
	}
	return nil
}

func (instance *ArangoCollection) IsReady() bool {
	return nil != instance && nil != instance.collection
}

func (instance *ArangoCollection) Name() string {
	if nil != instance && instance.IsReady() {
		return instance.collection.Name()
	}
	return ""
}

func (instance *ArangoCollection) Drop() (bool, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		err := instance.collection.Remove(ctx)
		return nil == err, err
	}
	return false, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Count() (int64, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		return instance.collection.Count(ctx)
	}
	return -1, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Exists(key string) (bool, error) {
	if nil != instance && instance.IsReady() {
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
	return false, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Upsert(doc map[string]interface{}) (map[string]interface{}, driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
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
	return nil, driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Insert(doc map[string]interface{}) (map[string]interface{}, driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
		if nil != doc {
			ctx := context.Background()
			meta, err := instance.collection.CreateDocument(ctx, doc)
			if nil != err {
				return nil, driver.DocumentMeta{}, err
			}

			if instance.addKey(&meta, doc) {
				// update
				meta, err = instance.collection.ReadDocument(ctx, meta.Key, &doc) // uses doc address
			}

			return doc, meta, err
		}
		return nil, driver.DocumentMeta{}, nil
	}
	return nil, driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) InsertDocument(doc interface{}) (interface{}, driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
		if nil != doc {
			ctx := context.Background()
			meta, err := instance.collection.CreateDocument(ctx, doc)
			if nil != err {
				return nil, driver.DocumentMeta{}, err
			}

			return doc, meta, err
		}
		return nil, driver.DocumentMeta{}, nil
	}
	return nil, driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Update(doc map[string]interface{}) (map[string]interface{}, driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
		if nil != doc {
			key := qbc.Reflect.GetString(doc, KeyFieldName)
			if len(key) > 0 {
				ctx := context.Background()
				meta, err := instance.collection.UpdateDocument(ctx, key, doc)
				if nil != err {
					return nil, driver.DocumentMeta{}, err
				}
				return doc, meta, err
			}
			return nil, driver.DocumentMeta{}, ErrMissingDocumentKey
		}
		return nil, driver.DocumentMeta{}, nil
	}
	return nil, driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) UpdateDocument(key string, doc interface{}) (interface{}, driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
		if nil != doc {
			ctx := context.Background()
			meta, err := instance.collection.UpdateDocument(ctx, key, doc)
			if nil != err {
				return nil, driver.DocumentMeta{}, err
			}
			return doc, meta, err
		}
		return nil, driver.DocumentMeta{}, nil
	}
	return nil, driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Remove(key string) (driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		meta, err := instance.collection.RemoveDocument(ctx, key)
		if nil != err {
			return driver.DocumentMeta{}, err
		}
		return meta, nil
	}
	return driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) Read(key string) (map[string]interface{}, driver.DocumentMeta, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		var doc map[string]interface{}
		meta, err := instance.collection.ReadDocument(ctx, key, &doc)
		if nil != err {
			return nil, driver.DocumentMeta{}, err
		}
		return doc, meta, nil
	}
	return nil, driver.DocumentMeta{}, ErrCollectionDoesNotExists
}

//----------------------------------------------------------------------------------------------------------------------
//	i n d e x
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoCollection) RemoveIndex(fields []string) (bool, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()

		name := instance.getIndexName("persist", fields)
		index, err := instance.collection.Index(ctx, name)
		if nil != err {
			return false, err
		}
		if nil != index {
			err = index.Remove(ctx)
		}
		return nil != index, err
	}
	return false, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) RemoveGeoIndex(fields []string) (bool, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()

		name := instance.getIndexName("geo", fields)
		index, err := instance.collection.Index(ctx, name)
		if nil != err {
			return false, err
		}
		if nil != index {
			err = index.Remove(ctx)
		}
		return nil != index, err
	}
	return false, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) EnsureIndex(fields []string, unique bool) (bool, error) {
	if nil != instance && instance.IsReady() {

		// remove existing
		instance.RemoveIndex(fields)

		ctx := context.Background()
		options := &driver.EnsurePersistentIndexOptions{
			Name:   instance.getIndexName("persist", fields),
			Unique: unique,
		}

		_, b, err := instance.collection.EnsurePersistentIndex(ctx, fields, options)
		if nil != err {
			return false, err
		}
		return b, err
	}
	return false, ErrCollectionDoesNotExists
}

func (instance *ArangoCollection) EnsureGeoIndex(fields []string, geoJson bool) (bool, error) {
	if nil != instance && instance.IsReady() {

		// remove existing
		instance.RemoveGeoIndex(fields)

		ctx := context.Background()
		options := &driver.EnsureGeoIndexOptions{
			Name:    instance.getIndexName("geo", fields),
			GeoJSON: geoJson,
		}

		_, b, err := instance.collection.EnsureGeoIndex(ctx, fields, options)
		if nil != err {
			return false, err
		}
		return b, err
	}
	return false, ErrCollectionDoesNotExists
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoCollection) getIndexName(prefix string, fields []string) string {
	a := qbc.Convert.ToArray(fields)
	name := prefix + "_" + qbc.Strings.ConcatSep("_", a...)
	return "idx_" + qbc.Coding.MD5(name)
}

func (instance *ArangoCollection) addKey(meta *driver.DocumentMeta, doc map[string]interface{}) bool {
	key := qbc.Reflect.GetString(doc, KeyFieldName)
	if len(key) == 0 {
		key = meta.Key
		if len(key) > 0 {
			qbc.Maps.Set(doc, KeyFieldName, key)
			return true
		}
	}
	return false
}

func (instance *ArangoCollection) ensureKey(doc map[string]interface{}) bool {
	key := qbc.Reflect.GetString(doc, KeyFieldName)
	if len(key) == 0 {
		key = qbc.Rnd.Uuid()
		if len(key) > 0 {
			qbc.Maps.Set(doc, KeyFieldName, key)
			return true
		}
	}
	return false
}

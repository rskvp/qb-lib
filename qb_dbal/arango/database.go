package arango

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/arangodb/go-driver"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type ArangoDatabase struct {
	database driver.Database
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoDatabase) Native() driver.Database {
	if nil != instance && instance.IsReady() {
		return instance.database
	}
	return nil
}

func (instance *ArangoDatabase) ImportFiles(fileNames []string) error {
	for _, fileName := range fileNames {
		err := instance.ImportFile(fileName)
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *ArangoDatabase) ImportFile(fileName string) error {
	ext := qbc.Paths.ExtensionName(fileName)
	if len(ext) == 0 {
		return errors.New("missing_extension")
	}

	name := qbc.Paths.FileName(fileName, false)
	if len(name) == 0 {
		return errors.New("invalid_filename")
	}

	if ext == "json" {
		return instance.importFileJSON(name, fileName)
	} else if ext == "csv" {
		return instance.importFileCSV(name, fileName)
	} else {
		return errors.New("unsupported_filename")
	}

	//return nil
}

func (instance *ArangoDatabase) IsReady() bool {
	return nil != instance && nil != instance.database
}

func (instance *ArangoDatabase) Name() string {
	if nil != instance && instance.IsReady() {
		return instance.database.Name()
	}
	return ""
}

func (instance *ArangoDatabase) Drop() (bool, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		err := instance.database.Remove(ctx)
		return nil == err, err
	}
	return false, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) CollectionNames() ([]string, error) {
	response := make([]string, 0)
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		collections, err := instance.database.Collections(ctx)
		if nil != err {
			return response, err
		}
		for _, coll := range collections {
			response = append(response, coll.Name())
		}
	}
	return response, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) CollectionExists(name string) (bool, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		exists, err := instance.database.CollectionExists(ctx, name)
		if nil != err {
			return false, err
		}
		return exists, nil
	}
	return false, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) CollectionAutoCreate(name string) (*ArangoCollection, error) {
	return instance.Collection(name, true)
}

func (instance *ArangoDatabase) Collection(name string, createIfNotExists bool) (*ArangoCollection, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		exists, err := instance.database.CollectionExists(ctx, name)
		if nil != err {
			return nil, err
		}

		if !exists && createIfNotExists {
			_, err := instance.database.CreateCollection(ctx, name, nil)
			if nil != err {
				return nil, err
			}
		}

		collection, err := instance.database.Collection(ctx, name)
		if nil != err {
			return nil, err
		}
		if nil == collection {
			return nil, ErrCollectionDoesNotExists
		}
		response := new(ArangoCollection)
		response.collection = collection

		return response, nil
	}
	return nil, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) Query(query string, bindVars map[string]interface{}, callback QueryCallback) error {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		cursor, err := instance.database.Query(ctx, query, bindVars)
		if nil != err {
			return err
		}

		defer cursor.Close()
		for {
			var doc interface{}
			meta, err := cursor.ReadDocument(ctx, &doc)
			if driver.IsNoMoreDocuments(err) {
				break
			} else {
				if nil != callback {
					exit := callback(meta, doc, err)
					if exit {
						break
					}
				}
			}
		}
		// no error
		return nil
	}
	return ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) Count(query string, bindVars map[string]interface{}) (int64, error) {
	if nil != instance && instance.IsReady() {
		ctx := context.Background()
		cursor, err := instance.database.Query(ctx, query, bindVars)
		if nil != err {
			return 0, err
		}

		defer cursor.Close()
		return cursor.Count(), nil
	}
	return 0, ErrDatabaseDoesNotExists
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoDatabase) importFileJSON(collName, fileName string) error {
	text, err := qbc.IO.ReadTextFromFile(fileName)
	if nil == err {
		if len(text) > 0 {

			collection, err := instance.Collection(collName, true)
			if nil != err {
				return err
			}

			var a []map[string]interface{}
			err = json.Unmarshal([]byte(text), &a)
			if nil != err {
				return err
			}

			for _, item := range a {
				if nil != item {
					_, _, err = collection.Upsert(item)
					if nil != err {
						// exit loop if has an error
						return err
					}
				}
			}

			return nil
		}
		return errors.New("empty_file")
	}
	return err
}

func (instance *ArangoDatabase) importFileCSV(collName, fileName string) error {
	text, err := qbc.IO.ReadTextFromFile(fileName)
	if nil == err {
		if len(text) > 0 {

			collection, err := instance.Collection(collName, true)
			if nil != err {
				return err
			}

			a, err := qbc.CSV.ReadAll(text, qbc.CSV.NewCsvOptionsDefaults())
			if nil != err {
				return err
			}

			for _, item := range a {
				if nil != item {
					_, _, err = collection.Upsert(qbc.Convert.ToMap(item))
					if nil != err {
						// exit loop if has an error
						return err
					}
				}
			}

			return nil
		}
		return errors.New("empty_file")
	}
	return err
}

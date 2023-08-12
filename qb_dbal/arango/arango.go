package arango

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

//----------------------------------------------------------------------------------------------------------------------
//	e r r o r s
//----------------------------------------------------------------------------------------------------------------------

var (
	ErrMissingConfiguration    = errors.New("missing_configuration")
	ErrConnectionNotReady      = errors.New("connection_not_ready")
	ErrDatabaseDoesNotExists   = errors.New("database_does_not_exists")
	ErrCollectionDoesNotExists = errors.New("collection_does_not_exists")
	ErrMissingDocumentKey      = errors.New("document_missing_key")
)

//----------------------------------------------------------------------------------------------------------------------
//	const
//----------------------------------------------------------------------------------------------------------------------

const KeyFieldName = "_key" // all entities should have this field

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type ArangoConnection struct {
	Config *ArangoConfig

	Version string
	Server  string
	License string

	//-- private --//
	client driver.Client
}

type QueryCallback func(driver.DocumentMeta, interface{}, error) bool

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewArangoConnection(config *ArangoConfig) *ArangoConnection {
	instance := new(ArangoConnection)
	instance.Config = config

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoConnection) Open() (err error) {
	conn, err := instance.getConnection()
	if nil != err {
		return err
	}

	auth, err := instance.getAuthentication()
	if nil != err {
		return err
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: auth,
	})
	if nil != err {
		return err
	}

	ctx := context.Background()
	v, err := c.Version(ctx)
	if nil != err {
		return err
	}

	instance.client = c
	instance.Version = string(v.Version)
	instance.Server = v.Server
	instance.License = v.License

	return err
}

func (instance *ArangoConnection) IsReady() bool {
	if nil != instance {
		return nil != instance.client
	}
	return false
}

func (instance *ArangoConnection) Database(name string, createIfNotExists bool) (response *ArangoDatabase, err error) {
	if nil==instance || !instance.IsReady() {
		return nil, ErrConnectionNotReady
	}
	ctx := context.Background()
	exists, err := instance.client.DatabaseExists(ctx, name)
	if nil != err {
		return nil, err
	}
	if !exists && createIfNotExists {
		db, err := instance.client.CreateDatabase(ctx, name, nil)
		if nil != err {
			return nil, err
		}
		response = new(ArangoDatabase)
		response.database = db

	} else {
		db, err := instance.client.Database(ctx, name)
		if nil != err {
			return nil, err
		}
		response = new(ArangoDatabase)
		response.database = db
	}
	return response, nil
}

func (instance *ArangoConnection) DropDatabase(name string) (success bool, err error) {
	if !instance.IsReady() {
		return false, ErrConnectionNotReady
	}
	ctx := context.Background()
	exists, err := instance.client.DatabaseExists(ctx, name)
	if nil != err {
		return false, err
	}
	if !exists {
		return false, ErrDatabaseDoesNotExists
	}

	db, err := instance.client.Database(ctx, name)
	if nil != err {
		return false, err
	}

	err = db.Remove(ctx)
	return nil == err, err
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoConnection) getAuthentication() (driver.Authentication, error) {
	config := instance.Config
	if nil != config {
		auth := config.Authentication
		if len(auth.Username) > 0 && len(auth.Password) > 0 {
			// BASIC
			return driver.BasicAuthentication(auth.Username, auth.Password), nil
		}
	}
	return nil, ErrMissingConfiguration
}

func (instance *ArangoConnection) getConnection() (driver.Connection, error) {
	config := instance.Config
	if nil != config {
		return http.NewConnection(http.ConnectionConfig{
			Endpoints: config.Endpoints,
			TLSConfig: nil,
		})
	}

	return nil, ErrMissingConfiguration
}

package nosqldrivers

import "errors"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------
var (
	errorDriverNotImplemented  = errors.New("driver_not_implemented")
	errConnectionNotReady      = errors.New("connection_not_ready")
	errDatabaseDoesNotExists   = errors.New("database_does_not_exists")
	errCollectionDoesNotExists = errors.New("collection_does_not_exists")
)

type INoSqlDatabase interface {
	Close() error
	Query(query string, bindVars map[string]interface{}) ([]interface{}, error)
	Exec(query string, bindVars map[string]interface{}) (interface{}, error)
	Count(query string, bindVars map[string]interface{}) (int64, error)
	Insert(collectionName string, item map[string]interface{}) (map[string]interface{}, error)
	Update(collectionName string, item map[string]interface{}) (map[string]interface{}, error)
	Upsert(collectionName string, item map[string]interface{}) (map[string]interface{}, error)
	Delete(collectionName string, item interface{}) (map[string]interface{}, error)

	Collection(name string, createIfNotExists bool) (INoSqlCollection, error)
}

type INoSqlCollection interface {
	Name() string
	RemoveIndex(typeName string, fields []string) (bool, error)
	EnsureIndex(typeName string, fields []string, unique bool) (bool, error)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func NewDatabase(driverName, connectionString string) (driver INoSqlDatabase, err error) {
	dsn := NewNoSqlDsn(connectionString)
	if len(dsn.Host) > 0 {
		switch driverName {
		case "arango":
			driver, err = NewDriverArango(dsn)
		}

		return
	}
	return driver, errorDriverNotImplemented
}

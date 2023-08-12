package storage

import "errors"

//----------------------------------------------------------------------------------------------------------------------
//	e r r o r s
//----------------------------------------------------------------------------------------------------------------------

var (
	ErrorInvalidDsn = errors.New("invalid_dsn")
	ErrorDriverNotImplemented = errors.New("driver_not_implemented")
	ErrorDatabaseDoesNotExists = errors.New("database_does_not_exists")
	ErrorCollectionDoesNotExists = errors.New("collection_does_not_exists")
	ErrorDatabaseCacheCannotAuthenticate = errors.New("database_cache_cannot_authenticate")
	ErrorDatabaseCacheNotEnabled = errors.New("database_cache_not_enabled")

	ErrorEntityAlreadyRegistered = errors.New("entity_already_registered")
	ErrorEntityDoesNotExists = errors.New("entity_does_not_exists")

	ErrorTokenDoesNotExists = errors.New("token_does_not_exists")
	ErrorTokenExpired = errors.New("token_expired")
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t a n t s
//----------------------------------------------------------------------------------------------------------------------

var (
	CollectionAuth = "auth"		// users authentication
	CollectionCache = "cache"	// tokens cache
)

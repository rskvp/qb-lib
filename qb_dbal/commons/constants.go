package commons

import "errors"

//----------------------------------------------------------------------------------------------------------------------
//	e r r o r s
//----------------------------------------------------------------------------------------------------------------------

var (
	ErrorInvalidDsn                    = errors.New("invalid_dsn")
	ErrorDriverNotImplemented          = errors.New("driver_not_implemented")
	ErrorDatabaseDoesNotExists         = errors.New("database_does_not_exists")
	ErrorCollectionDoesNotExists       = errors.New("collection_does_not_exists")
	ErrorMismatchConfiguration         = errors.New("mismatch_configuration")
	ErrorMissingTransactionOptions     = errors.New("missing_transaction_options")
	ErrorMissingTransactionCollections = errors.New("missing_transaction_collections")
	ErrorCommandAndParamsDoNotMatch    = errors.New("commands_and_params_do_not_match")

	ErrorEngineNotReady      = errors.New("engine_not_ready")
	ErrorCommandNotSupported = errors.New("command_not_supported")
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t a n t s
//----------------------------------------------------------------------------------------------------------------------

const (
	CAT_PERSON   = "person"
	CAT_DOCUMENT = "document"
	CAT_ADV      = "adv"
	CAT_POST     = "post"
	CAT_EVENT    = "event"
)

var (
	CATEGORIES = []string{CAT_PERSON, CAT_DOCUMENT, CAT_ADV, CAT_POST, CAT_EVENT}
)

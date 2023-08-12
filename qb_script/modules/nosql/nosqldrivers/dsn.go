package nosqldrivers

import (
	"fmt"

	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

// "admin:xxxxxxxxx@tcp(localhost:3306)/test"
type NoSqlDsn struct {
	User     string
	Password string
	Protocol string
	Host     string
	Port     int
	Database string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNoSqlDsn(dsn ...string) *NoSqlDsn {
	instance := new(NoSqlDsn)
	if len(dsn) > 0 {
		instance.parse(dsn[0])
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// String return dsn string. i.e. "admin:xxxxxxxx@tcp(localhost:3306)/test"
func (instance *NoSqlDsn) String() string {
	// "admin:!qaz2WSX098@tcp(localhost:3306)/test"
	return fmt.Sprintf("%v:%v@%v(%v:%v)/%v", instance.User, instance.Password, instance.Protocol, instance.Host, instance.Port, instance.Database)
}

func (instance *NoSqlDsn) GoString() string {
	return instance.String()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// parse Parse a dsn string. i.e. "admin:xxxxxxxxx@tcp(localhost:3306)/test"
func (instance *NoSqlDsn) parse(dsn string) {
	tokens := qbc.Strings.Split(dsn, ":@()/")
	count := len(tokens)
	if count > 1 {
		// set username and password
		instance.User = tokens[0]
		instance.Password = tokens[1]
	}
	if count > 4 {
		// set username and password
		instance.Protocol = tokens[2]
		instance.Host = tokens[3]
		instance.Port = qbc.Convert.ToInt(tokens[4])
	}
	if count > 5 {
		instance.Database = tokens[5]
	}
}

package storage

import (
	"fmt"
	"os"
	"strings"

	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

// "admin:xxxxxxxxx@tcp(localhost:3306)/test"
type Dsn struct {
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

func NewDsn(dsn ...string) *Dsn {
	instance := new(Dsn)
	if len(dsn) > 0 {
		instance.parse(dsn[0])
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// String return dsn string. i.e. "admin:xxxxxxxx@tcp(localhost:3306)/test"
func (instance *Dsn) String() string {
	// "admin:!qaz2WSX098@tcp(localhost:3306)/test"
	if len(instance.Host) == 0 {
		return fmt.Sprintf("%v:%v@%v:%v", instance.User, instance.Password, instance.Protocol, instance.Database)
	}
	return fmt.Sprintf("%v:%v@%v(%v:%v)/%v", instance.User, instance.Password, instance.Protocol, instance.Host, instance.Port, instance.Database)
}

func (instance *Dsn) GoString() string {
	return instance.String()
}

func (instance *Dsn) IsValid() bool {
	if nil != instance && len(instance.String()) > 0 && len(instance.Protocol) > 0 {
		if len(instance.Host) == 0 {
			return len(instance.Database) > 0
		} else {
			return len(instance.Database) > 0 && instance.Port > 0
		}
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// parse Parse a dsn string. i.e. "admin:xxxxxxxxx@tcp(localhost:3306)/test"
func (instance *Dsn) parse(dsn string) {
	tokens := qbc.Strings.Split(dsn, ":@()/")
	count := len(tokens)
	if count > 1 {
		// set username and password
		instance.User = tokens[0]
		instance.Password = tokens[1]
	}
	if count == 4 {
		instance.Protocol = tokens[2]
		instance.Database = tokens[3]
	} else {
		if count > 4 {
			// set protocol, host and port
			instance.Protocol = tokens[2]
			if instance.Protocol == "file" {
				instance.Database = strings.Join(tokens[3:], string(os.PathSeparator))
			} else {
				instance.Host = tokens[3]
				instance.Port = qbc.Convert.ToInt(tokens[4])
				if count > 5 {
					instance.Database = tokens[5]
				}
			}
		}
	}
}

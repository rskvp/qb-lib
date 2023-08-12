package commons

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------
const (
	DsnODBC = iota
	DsnFile
	DsnSTD
)

// "admin:xxxxxxxxx@tcp(localhost:3306)/test"
type Dsn struct {
	Type     int
	User     string
	Password string
	Protocol string
	Host     string
	Port     int
	Database string
	Driver   string
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
	switch instance.Type {
	case DsnFile:
		return instance.Database
	case DsnODBC:
		// "driver=mysql;server=%s;database=%s;user=%s;password=%s;",
		// "server=%s;database=%s;uid=%s;pwd=%s;port=%s;TDS_Version=8.0"
		if len(instance.Driver) > 0 {
			// generic odbc
			// "driver=mysql;server=%s;database=%s;user=%s;password=%s;",
			return fmt.Sprintf("driver=%v;server=%v;database=%v;user=%v;password=%v;",
				instance.Driver, instance.Host, instance.Database, instance.User, instance.Password)
		} else {
			// mssql
			return fmt.Sprintf("driver=%v;server=%v,%v;database=%v;uid=%v;pwd=%v;TDS_Version=8.0;",
				defaultDriver(), instance.Host, instance.Port, instance.Database, instance.User, instance.Password)
		}
	default:
		// "admin:!qaz2WSX098@tcp(localhost:3306)/test"
		if len(instance.Host) == 0 {
			return fmt.Sprintf("%v:%v@%v:%v", instance.User, instance.Password, instance.Protocol, instance.Database)
		}
		return fmt.Sprintf("%v:%v@%v(%v:%v)/%v", instance.User, instance.Password, instance.Protocol, instance.Host, instance.Port, instance.Database)
	}
}

func (instance *Dsn) GoString() string {
	return instance.String()
}

func (instance *Dsn) IsValid() bool {
	switch instance.Type {
	case DsnFile:
		return len(instance.Database) > 0
	case DsnODBC:
		return len(instance.Host) > 0 && instance.Port > 0
	default:
		if nil != instance && len(instance.String()) > 0 && len(instance.Protocol) > 0 {
			if len(instance.Host) == 0 {
				return len(instance.Database) > 0
			} else {
				return len(instance.Database) > 0 && instance.Port > 0
			}
		}
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// parse Parse a dsn string. i.e. "admin:xxxxxxxxx@tcp(localhost:3306)/test"
// or odbc like string: "driver=mysql;server=%s;database=%s;user=%s;password=%s;",
// "server=%s;database=%s;uid=%s;pwd=%s;port=%s;TDS_Version=8.0"
func (instance *Dsn) parse(dsn string) {
	if isSQLite(dsn) {
		instance.Type = DsnFile
		instance.Database = dsn
	} else if isODBC(dsn) {
		// ODBC dsn
		instance.Type = DsnODBC
		tokens := qbc.Strings.Split(dsn, ";")
		for _, token := range tokens {
			kv := strings.Split(token, "=")
			if len(kv) == 2 {
				k := kv[0]
				v := kv[1]
				switch k {
				case "driver":
					instance.Driver = v
				case "server":
					if strings.Index(v, ",") > -1 {
						vv := strings.Split(v, ",")
						instance.Host = vv[0]
						instance.Port = qbc.Convert.ToInt(vv[1])
					} else {
						instance.Host = v
					}
				case "database":
					instance.Database = v
				case "user", "uid":
					instance.User = v
				case "password", "pwd":
					instance.Password = v
				case "port":
					instance.Port = qbc.Convert.ToInt(v)
					if instance.Port == 0 {
						instance.Port = 1433
					}
				}
			}
		}
	} else {
		// standard dsn
		instance.Type = DsnSTD
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
}

func isODBC(dsn string) bool {
	return strings.Index(dsn, "=") > -1 && strings.Index(dsn, ";") > -1
}

func isSQLite(dsn string) bool {
	return strings.Index(dsn, "/") > -1 && qbc.Paths.ExtensionName(dsn) == "db"
}

func defaultDriver() string {
	if runtime.GOOS == "windows" {
		return "sql server"
	} else {
		return "freetds"
	}
}

package nosqldrivers

import (
	"fmt"
	"testing"
)

func TestDsn(t *testing.T) {
	txt := "admin:xxxxxxxxx@tcp(localhost:3306)/test"
	dsn := NewNoSqlDsn(txt)

	if dsn.String() != txt {
		t.Error("expected '" + txt + "' but got '" + dsn.String() + "'")
		t.FailNow()
	}

	fmt.Println(dsn)
}

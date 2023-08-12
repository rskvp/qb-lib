package commons

import "testing"

func TestDsn(t *testing.T) {

	sDsn := "admin:!qaz2WSX098@tcp(localhost:3306)/test"
	dsn := NewDsn(sDsn)
	if dsn.String()!=sDsn{
		t.Error("Expected '"+sDsn+"', got ", dsn.String())
		t.FailNow()
	}
	if dsn.User!="admin"{
		t.Error("Expected 'admin', got ", dsn.User)
		t.FailNow()
	}
	if dsn.Password!="!qaz2WSX098"{
		t.Error("Expected '!qaz2WSX098', got ", dsn.Password)
		t.FailNow()
	}
	if dsn.Protocol!="tcp"{
		t.Error("Expected 'tcp', got ", dsn.Protocol)
		t.FailNow()
	}
	if dsn.Host!="localhost"{
		t.Error("Expected 'localhost', got ", dsn.Host)
		t.FailNow()
	}
	if dsn.Port!=3306{
		t.Error("Expected 3306, got ", dsn.Port)
		t.FailNow()
	}
	if dsn.Database!="test"{
		t.Error("Expected 'test', got ", dsn.Database)
		t.FailNow()
	}

	sDsn = "admin:!qaz2WSX098@file:./db/test.dat"
	dsn = NewDsn(sDsn)
	if dsn.String()!=sDsn{
		t.Error("Expected '"+sDsn+"', got ", dsn.String())
		t.FailNow()
	}

	odbc := "driver=mysql;server=192.168.1.1;database=Test;user=foo;password=boo;"
	dsn = NewDsn(odbc)
	if dsn.String()!=odbc{
		t.Error("Expected '"+odbc+"', got ", dsn.String())
		t.FailNow()
	}
}

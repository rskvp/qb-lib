package _test

import (
	"fmt"
	"testing"

	"github.com/rskvp/qb-lib/qb_vfs/vfscommons"
)

func TestSettings(t *testing.T) {

	settings, err := vfscommons.LoadVfsSettings("./settings_os.json")
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(settings)
}

package console

import (
	"testing"

	_ "github.com/rskvp/qb-lib/qb_script/modules/defaults/util"
)

func TestConsole(t *testing.T) {
	/**
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	if c := vm.Get("console"); c == nil {
		t.Fatal("console not found")
	}

	if _, err := vm.RunString("console.log('')"); err != nil {
		t.Fatal("console.log() error", err)
	}

	if _, err := vm.RunString("console.error('')"); err != nil {
		t.Fatal("console.error() error", err)
	}

	if _, err := vm.RunString("console.warn('')"); err != nil {
		t.Fatal("console.warn() error", err)
	}
	*/
}

package scripting

import (
	"fmt"
	"testing"
	"time"

	qbc "github.com/rskvp/qb-core"
	ggx "github.com/rskvp/qb-lib"
	"github.com/rskvp/qb-lib/qb_script"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

func Test_Simple(t *testing.T) {
	qbc.Paths.SetWorkspacePath(qbc.Paths.Absolute("./")) // look for modules in this folder

	response, err := ggx.Scripting.Run(&qb_script.EnvSettings{
		FileName: "./simple.js",
		Logger:   nil, //logger,
		LogReset: true,
		OnReady: func(rtContext *commons.RuntimeContext) {
			Enable(rtContext)
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Script returned: ", response)
	time.Sleep(3 * time.Second)
}

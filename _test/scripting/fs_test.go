package _test

import (
	"fmt"
	"testing"
	"time"

	ggx "bitbucket.org/digi-sense/gg-core-x"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-lib/qb_script"
)

func Test_fs(t *testing.T) {
	const SCRIPT = `
	(function(){
		var test = require("./modules/test_fs.js");
		return test.run();
	})();
	
	`

	expected := map[string]string{
		"exists": "true",
	}

	registry := qb_script.NewModuleRegistry(func(path string) ([]byte, error) {
		return qbc.IO.ReadBytesFromFile(path)
	})

	vm1 := ggx.Scripting.NewEngine("javascript")
	vm1.Name = "fs" // creates log file vm1.log
	registry.Start(vm1)

	v, err := vm1.RunScript(SCRIPT)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	response := v.String()
	if qbc.Regex.IsValidJsonObject(response) {
		m := make(map[string]interface{})
		err = qbc.JSON.Read(response, &m)
		if err != nil {
			t.Error(err)
		}

		for k, v := range expected {
			value := qbc.Reflect.GetString(m, k)
			if value != v {
				t.Error("Expected: " + v + " but got " + value)
			}
		}
	} else {
		t.Error(response)
	}

	time.Sleep(5 * time.Second)
}

func Test_fs2(t *testing.T) {
	qbc.Paths.SetWorkspacePath(qbc.Paths.Absolute("./")) // look for modules in this folder

	logger := qb_log.NewLogger()
	logger.SetFilename(qbc.Paths.WorkspacePath("custom_fs.log"))

	response, err := ggx.Scripting.Run(&qb_script.EnvSettings{
		FileName: "./fs.js",
		Logger:   nil, //logger,
		LogReset: true,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(response)
}

func Test_err(t *testing.T) {
	qbc.Paths.SetWorkspacePath(qbc.Paths.Absolute("./")) // look for modules in this folder
	response, err := ggx.Scripting.RunFile("./fs2.js", nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(response)
}

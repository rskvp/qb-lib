package _test

import (
	"fmt"
	"testing"

	ggx "bitbucket.org/digi-sense/gg-core-x"
	qbc "github.com/rskvp/qb-core"
)

func TestVfs(t *testing.T) {

	filesettings := []string{
		"./settings_ftp.json",
		"./settings_sftp.json",
		"./settings_os.json",
	}

	errorList := make([]error, 0)

	for _, settings := range filesettings {
		fmt.Println("----------------------")
		fmt.Println("TESTING: ", settings)
		fmt.Println("----------------------")
		vfs, err := ggx.VFS.New(settings)
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		defer vfs.Close()
		fmt.Println("current path: ", vfs.Path())
		fmt.Println("\t", vfs)

		// LIST & READ
		list, err := vfs.List("./")
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		for _, file := range list {
			fmt.Println("\t", file.AbsolutePath, file.IsDir, file.Mode, file.Size)
			if !file.IsDir && file.Size > 0 {
				data, err := vfs.Read(file.AbsolutePath)
				if nil != err {
					errorList = append(errorList, qbc.Errors.Prefix(err, "["+settings+"] FILE CONTENT READ ERROR: "))
					fmt.Println("\t", "FILE CONTENT READ ERROR: ", err)
				}
				fmt.Println("\t", "FILE CONTENT READ: ", len(data))
			}
		}

		// UPLOAD, EXISTS, REMOVE
		target := "./test__upload/sample.txt"
		fmt.Println("current path: ", vfs.Path())
		fmt.Println("UPLOADING:", target)
		if b, _ := vfs.Exists(target); b {
			// remove
			err = vfs.Remove(target)
			if nil != err {
				t.Error(err)
				t.FailNow()
			}
		} else {
			data, err := qbc.IO.ReadBytesFromFile("./sample.txt")
			if nil != err {
				t.Error(err)
				t.FailNow()
			}
			fmt.Println("current path: ", vfs.Path())
			fmt.Println("\t", "WRITE:", target)
			_, err = vfs.Write(data, target)
			if nil != err {
				t.Error(err)
				t.FailNow()
			}
			if b, err := vfs.Exists(target); b {
				// read
				fmt.Println("current path: ", vfs.Path())
				fmt.Println("\t", "READ:", target)
				data, err := vfs.Read(target)
				if nil != err {
					t.Error(err)
					t.FailNow()
				}
				fmt.Println("\t", "CONTENT READ:", string(data))

				// remove file
				fmt.Println("current path: ", vfs.Path())
				fmt.Println("\t", "REMOVE FILE:", target)
				err = vfs.Remove(target)
				if nil != err {
					t.Error(err)
					t.FailNow()
				}

				// remove dir
				removeDir := "./" + qbc.Paths.Dir(target)
				fmt.Println("current path: ", vfs.Path())
				fmt.Println("\t", "REMOVE:", removeDir)
				err = vfs.Remove(removeDir)
				if nil != err {
					t.Error(err)
					t.FailNow()
				}
				fmt.Println("current path: ", vfs.Path())
				if b, _ := vfs.Exists(removeDir); b {
					t.Error("Unexpected existing dir: ", removeDir)
					t.FailNow()
				}
			} else {
				if nil != err {
					t.Error(err)
					t.FailNow()
				} else {
					t.Error("file not found")
					t.FailNow()
				}
			}
		}

		// CD
		fmt.Println("current path: ", vfs.Path())
		fmt.Println("\t", "cd ..")
		if b, err := vfs.Cd(".."); b {
			fmt.Println("current path: ", vfs.Path())
		} else {
			t.Error(err)
			t.FailNow()
		}
		list, err = vfs.List("./")
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		for _, file := range list {
			fmt.Println("\t", file.AbsolutePath, file.IsDir, file.Mode, file.Size)
		}

	}

	// log errors
	if len(errorList) > 0 {
		fmt.Println("--------------------------------------------")
		fmt.Println("TEST WARNINGS")
		fmt.Println("--------------------------------------------")
		for _, err := range errorList {
			fmt.Println(err)
		}
		fmt.Println("--------------------------------------------")
	}
}

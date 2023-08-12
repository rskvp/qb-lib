package _test

import (
	"errors"
	"fmt"
	"testing"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_vfs/vfsbackends"
	"github.com/rskvp/qb-lib/qb_vfs/vfscommons"
)

func TestOSList(t *testing.T) {
	settings, err := vfscommons.LoadVfsSettings("./settings_os.json")
	assert(t, err, "Load Settings :")

	fmt.Println("----------------------")
	fmt.Println("TESTING: ", settings)
	fmt.Println("----------------------")
	vfs, err := vfsbackends.NewVfsOS(settings)
	assert(t, err, "NewVfs :")

	defer vfs.Close()
	fmt.Println("current path: ", vfs.Path())
	fmt.Println("\t", vfs)

	// list files
	list, err := vfs.List("./")
	assert(t, err, "list files :")
	for _, file:=range list{
		fmt.Println(file)
		if file.IsDir{
			list2, err := vfs.List(file.RelativePath)
			assert(t, err, "list2 files :")
			for _, file2:=range list2{
				fmt.Println("\t", file2)
			}
		}
	}
	fmt.Println("current path: ", vfs.Path())

}


func TestOS(t *testing.T) {
	settings, err := vfscommons.LoadVfsSettings("./settings_os.json")
	assert(t, err, "Load Settings :")

	fmt.Println("----------------------")
	fmt.Println("TESTING: ", settings)
	fmt.Println("----------------------")
	vfs, err := vfsbackends.NewVfsOS(settings)
	assert(t, err, "NewVfs :")

	defer vfs.Close()
	fmt.Println("current path: ", vfs.Path())
	fmt.Println("\t", vfs)

	// mkdir
	workDir := "./test_os_dir"
	fmt.Println("Create dir:", workDir)
	err = vfs.MkDir(workDir)
	assert(t, err, "MKDIR :")
	fmt.Println("current path: ", vfs.Path())

	exists, err :=  vfs.Exists(workDir)
	fmt.Println("Check Exists: ", workDir, exists)
	assert(t, err, "checking dir exists :")
	if !exists{
		assert(t, errors.New(workDir), "expected dir exists :")
	}
	fmt.Println("current path: ", vfs.Path())

	fmt.Println("cd: ", workDir)
	_, err = vfs.Cd(workDir)
	assert(t, err, "CD :")
	fmt.Println("current path: ", vfs.Path())

	// file create
	file := "./sample.txt"
	data, err := qbc.IO.ReadBytesFromFile(file)
	assert(t, err, "reading local file :")
	fmt.Println("Write: ", file)
	size, err := vfs.Write(data, file)
	assert(t, err, "writing remote file :")
	if size!=len(data){
		assert(t, errors.New(fmt.Sprintf("Wrote %v bytes, expected %v", size, len(data))), "expected same byte was written :")
	}

	// check file exists
	exists, err =  vfs.Exists(file)
	assert(t, err, "checking file exists :")
	if !exists{
		assert(t, errors.New(file), "expected file exists :")
	}

	// list files
	fmt.Println("current path: ", vfs.Path())
	list, err := vfs.List("./")
	assert(t, err, "list files :")
	if len(list)==0{
		assert(t, errors.New(workDir), "expected not empty dir :")
	}
	fmt.Println("List: ", list)
	for _, file:=range list{
		fmt.Println(file)
	}

	// remove file
	err = vfs.Remove(file)
	assert(t, err, "remove: ")

	// remove dir
	fmt.Println("current path: ", vfs.Path())
	fmt.Println("cd..: ")
	if b, err := vfs.Cd("..");b{
		fmt.Println("current path: ", vfs.Path())
	} else {
		assert(t, err, "cd.. :")
	}
	fmt.Println("Removing: ", workDir)
	err = vfs.Remove(workDir)
	assert(t, err, "Removing dir :")
	fmt.Println("current path: ", vfs.Path())

	exists, err =  vfs.Exists(workDir)
	fmt.Println("Check Exists: ", workDir, exists)
	assert(t, err, "checking dir exists :")
	if exists{
		assert(t, errors.New(workDir), "expected dir do not exists :")
	}
	fmt.Println("current path: ", vfs.Path())
}

func assert(t *testing.T, err error, prefix string) {
	if nil != err {
		t.Error(qbc.Errors.Prefix(err, prefix))
		t.FailNow()
	}
}

package backends

import (
	"os"
	"strings"

	qbc "github.com/rskvp/qb-core"
	vfscommons "github.com/rskvp/qb-lib/qb_vfs/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	VfsOS
//----------------------------------------------------------------------------------------------------------------------

type VfsOS struct {
	settings *vfscommons.VfsSettings

	user string

	startDir string
	curDir   string
}

func NewVfsOS(settings *vfscommons.VfsSettings) (instance *VfsOS, err error) {
	instance = new(VfsOS)
	instance.settings = settings

	err = instance.init()

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsOS) String() string {
	return qbc.JSON.Stringify(instance.settings)
}

func (instance *VfsOS) Close() {
//empty
}

func (instance *VfsOS) Path() string {
	return instance.curDir
}

func (instance *VfsOS) Cd(path string) (bool, error) {
	instance.curDir = vfscommons.Absolutize(instance.curDir, path)
	return qbc.Paths.Exists(instance.curDir)
}

func (instance *VfsOS) Stat(path string) (*vfscommons.VfsFile, error) {
	fullPath := vfscommons.Absolutize(instance.curDir, path)
	info, err := os.Stat(fullPath)
	if nil != err {
		return nil, err
	}
	return vfscommons.NewVfsFile(fullPath, instance.curDir, info), nil

}

func (instance *VfsOS) Exists(path string) (bool, error) {
	return qbc.Paths.Exists(vfscommons.Absolutize(instance.curDir, path))
}

func (instance *VfsOS) List(dir string) ([]*vfscommons.VfsFile, error) {
	response := make([]*vfscommons.VfsFile, 0)
	dir = vfscommons.Absolutize(instance.curDir, dir)
	list, err := os.ReadDir(dir)
	if nil == err {
		for _, entry := range list {
			fullPath := qbc.Paths.Concat(dir, entry.Name())
			info, err := instance.Stat(fullPath)
			if nil == err {
				response = append(response, info)
			}
		}
	}
	return response, err
}

func (instance *VfsOS) Read(source string) ([]byte, error) {
	return qbc.IO.ReadBytesFromFile(vfscommons.Absolutize(instance.curDir, source))
}

func (instance *VfsOS) Write(data []byte, target string) (int, error) {
	target = vfscommons.Absolutize(instance.curDir, target)
	err := qbc.Paths.Mkdir(target)
	if nil != err {
		return 0, err
	}
	return qbc.IO.WriteBytesToFile(data, target)
}

func (instance *VfsOS) Download(source, target string) ([]byte, error) {
	data, err := instance.Read(vfscommons.Absolutize(instance.curDir, source))
	if nil != err {
		return nil, err
	}
	_, err = qbc.IO.WriteBytesToFile(data, vfscommons.Absolutize(instance.curDir, target))
	if nil != err {
		return nil, err
	}
	return data, nil
}

func (instance *VfsOS) Remove(source string) error {
	return qbc.IO.Remove(vfscommons.Absolutize(instance.curDir, source))
}

func (instance *VfsOS) MkDir(path string) error {
	return qbc.Paths.Mkdir(vfscommons.Absolutize(instance.curDir, path))
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsOS) init() error {
	if nil != instance.settings {
		_, root := instance.settings.SplitLocation()

		if strings.HasPrefix(root, vfscommons.FileUserHomePrefix) {
			userHome, err := qbc.Paths.UserHomeDir()
			if nil != err {
				return err
			}
			instance.startDir = qbc.Paths.Concat(userHome, root)
		} else if strings.HasPrefix(root, ".") {
			instance.startDir = qbc.Paths.Concat(qbc.Paths.GetWorkspacePath(), root)
		} else {
			// absolute path (user dir not used)
			instance.startDir = root
		}
		instance.curDir = instance.startDir
		if b, err := qbc.Paths.Exists(instance.curDir); !b {
			return err
		}
		return nil
	}
	return vfscommons.ErrorMissingConfiguration
}

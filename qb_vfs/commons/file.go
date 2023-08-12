package commons

import (
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	qbc "github.com/rskvp/qb-core"
)

type VfsFile struct {
	AbsolutePath string    `json:"absolute-path"`
	RelativePath string    `json:"relative-path"`
	Root         string    `json:"root"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	IsDir        bool      `json:"is_dir"`
	Mode         string    `json:"mode"`
}

func NewVfsFile(absolutePath, root string, item interface{}) *VfsFile {
	f := new(VfsFile)
	if file, b := item.(os.FileInfo); b {
		f.Name = file.Name()
		f.Root = root
		f.AbsolutePath = absolutePath
		f.RelativePath = Relativize(root, absolutePath)
		f.Size = file.Size()
		f.ModTime = file.ModTime()
		f.IsDir = file.IsDir()
		f.Mode = file.Mode().String()
	} else if entry, b := item.(*ftp.Entry); b {
		f.Name = entry.Name
		f.Root = root
		f.AbsolutePath = qbc.Paths.Concat(root, entry.Name)
		f.RelativePath = Relativize(root, f.AbsolutePath)
		f.Size = int64(entry.Size)
		f.ModTime = entry.Time
		f.IsDir = entry.Type == ftp.EntryTypeFolder

	}

	return f
}

func (instance *VfsFile) String() string {
	return qbc.JSON.Stringify(instance)
}

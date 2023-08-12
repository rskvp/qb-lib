package commons

import "errors"

type IVfs interface {
	Path() string
	Close()
	Cd(path string) (bool, error)
	Stat(path string) (*VfsFile, error)
	List(dir string) ([]*VfsFile, error)
	Read(source string) ([]byte, error)
	Write(data []byte, target string) (int, error)
	Download(source, target string) ([]byte, error)
	Remove(source string) error
	MkDir(path string) error
	Exists(path string) (bool, error)
}

//----------------------------------------------------------------------------------------------------------------------
//	e r r o r s
//----------------------------------------------------------------------------------------------------------------------

var (
	ErrorMissingConfiguration  = errors.New("missing configuration")
	ErrorMismatchConfiguration = errors.New("mismatch configuration")
	ErrorMissingConnection     = errors.New("missing connection")
	ErrorUnsupportedSchema     = errors.New("unsupported schema")
)

//----------------------------------------------------------------------------------------------------------------------
//	s c h e m a s
//----------------------------------------------------------------------------------------------------------------------

const (
	SchemaFTP  = "ftp"
	SchemaSFTP = "sftp"
	SchemaOS   = "file"
)

//----------------------------------------------------------------------------------------------------------------------
//	schema constants
//----------------------------------------------------------------------------------------------------------------------

const (
	FileUserHomePrefix = "~"
)

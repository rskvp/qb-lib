package backends

import (
	"bytes"
	"fmt"

	//"io/ioutil"

	"io"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	qbc "github.com/rskvp/qb-core"
	vfscommons "github.com/rskvp/qb-lib/qb_vfs/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	VfsFtp
//----------------------------------------------------------------------------------------------------------------------

type VfsFtp struct {
	settings *vfscommons.VfsSettings

	host     string
	port     int
	user     string
	password string
	key      []byte

	conn *VfsFtpConnection

	startDir string
	curDir   string
}

func NewVfsFtp(settings *vfscommons.VfsSettings) (instance *VfsFtp, err error) {
	instance = new(VfsFtp)
	instance.settings = settings

	err = instance.init()

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsFtp) String() string {
	return qbc.JSON.Stringify(instance.settings)
}

func (instance *VfsFtp) Close() {
	if nil != instance.conn {
		instance.conn.Close()
	}
}

func (instance *VfsFtp) Path() string {
	if conn, err := instance.connection(); nil != err {
		if len(instance.curDir) == 0 {
			instance.curDir = "/"
		}
	} else {
		dir, _ := conn.CurrentDir()
		if len(dir) > 0 {
			instance.curDir = dir
		}
	}
	return instance.curDir
}

func (instance *VfsFtp) Cd(path string) (bool, error) {
	if client, err := instance.connection(); nil != err {
		return false, err
	} else {
		path = instance.absolutize(path)
		err = client.ChangeDir(path)
		if instance.isValidError(err) {
			return false, err
		}
		instance.curDir, err = client.CurrentDir()
		return true, nil
	}
}

func (instance *VfsFtp) Stat(path string) (*vfscommons.VfsFile, error) {
	if conn, err := instance.connection(); nil != err {
		return nil, err
	} else {
		return instance.stat(conn, path)
	}
}

func (instance *VfsFtp) Exists(path string) (bool, error) {
	if conn, err := instance.connection(); nil != err {
		return false, err
	} else {
		info, err := instance.stat(conn, path)
		return nil != info, err
	}
}

func (instance *VfsFtp) List(dir string) ([]*vfscommons.VfsFile, error) {
	response := make([]*vfscommons.VfsFile, 0)
	if client, err := instance.connection(); nil != err {
		return nil, err
	} else {
		absolutePath := instance.absolutize(dir)
		entries, err := client.List(absolutePath)
		if instance.isValidError(err) {
			return nil, err
		}
		for _, entry := range entries {
			response = append(response, vfscommons.NewVfsFile(absolutePath, absolutePath, entry))
		}
	}
	return response, nil
}

func (instance *VfsFtp) Read(path string) ([]byte, error) {
	if conn, err := instance.connection(); nil != err {
		return nil, err
	} else {
		curDir, _ := conn.CurrentDir()
		defer conn.ChangeDir(curDir)

		_, name := instance.splitPath(conn, path, true)
		r, err := conn.Retr(name)
		if nil == r || instance.isValidError(err) {
			return nil, err
		}
		defer r.Close()

		data, err := io.ReadAll(r)
		if instance.isValidError(err) {
			return nil, err
		}
		return data, nil
	}
}

func (instance *VfsFtp) Write(data []byte, target string) (int, error) {
	if conn, err := instance.connection(); nil != err {
		return 0, err
	} else {
		curDir, _ := conn.CurrentDir()
		defer conn.ChangeDir(curDir)

		absolute := instance.absolutize(target)
		parent := qbc.Paths.Dir(absolute)
		name := qbc.Paths.FileName(absolute, true)
		if parent != curDir {
			err = instance.mkDir(conn, parent)
			_ = conn.ChangeDir(parent)
		}
		if instance.isValidError(err) {
			return 0, err
		}
		buffer := bytes.NewBuffer(data)
		err = conn.Stor(name, buffer)
		if instance.isValidError(err) {
			return 0, err
		}

		return len(data), nil
	}
}

func (instance *VfsFtp) Download(source, target string) ([]byte, error) {
	data, err := instance.Read(source)
	if nil != err {
		return nil, err
	}
	_, err = qbc.IO.WriteBytesToFile(data, target)
	if instance.isValidError(err) {
		return nil, err
	}
	return data, nil
}

func (instance *VfsFtp) Remove(path string) error {
	if conn, err := instance.connection(); nil != err {
		return err
	} else {
		curDir, _ := conn.CurrentDir()
		defer conn.ChangeDir(curDir)

		_, name := instance.splitPath(conn, path, true)
		if len(name) > 0 {
			if vfscommons.IsFile(name) && name != "." {
				err = conn.Delete(name)
			} else {
				if name != "." && name != ".." && name != "/" {
					err = conn.ChangeDir("..")
					if nil == err {
						err = conn.RemoveDir(name)
						if nil != err {
							err = conn.RemoveDirRecur(name)
						}
					}
				}
			}
		}
		if instance.isValidError(err) {
			return err
		}
		return nil
	}
}

func (instance *VfsFtp) MkDir(path string) error {
	if conn, err := instance.connection(); nil != err {
		return err
	} else {
		return instance.mkDir(conn, path)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsFtp) init() error {
	if nil != instance.settings {
		// prepare configuration
		instance.host, instance.port = vfscommons.SplitHost(instance.settings)
		user := instance.settings.Auth.User
		password := instance.settings.Auth.Password

		instance.user = user
		instance.password = password

		_, err := instance.connection()

		return err
	}
	return vfscommons.ErrorMissingConfiguration
}

func (instance *VfsFtp) connection() (*ftp.ServerConn, error) {
	if nil == instance.conn {
		instance.conn = NewVfsFtpConnection(instance.user, instance.password, instance.host, instance.port)
	}

	client, err := instance.conn.Open()
	if nil != err {
		instance.conn.Close()
		return nil, err
	}

	// update current dir
	instance.curDir, err = client.CurrentDir()
	if len(instance.startDir) == 0 {
		instance.startDir = instance.curDir
	}
	return client, nil
}

func (instance *VfsFtp) absolutize(p string) string {
	return vfscommons.Absolutize(instance.Path(), p)
}

func (instance *VfsFtp) isValidError(err error) bool {
	return nil != err //&& !ErrorContains(err, "is your current location") && !ErrorContains(err, "passive mode ok") && !ErrorContains(err, "file exists")
}

func (instance *VfsFtp) splitPath(conn *ftp.ServerConn, p string, relocate bool) (parent, name string) {
	path := instance.absolutize(p)
	isFile := vfscommons.IsFile(path)

	var dir string
	if isFile {
		dir = qbc.Paths.Dir(path)
	} else {
		dir = path
	}
	if dir == "." || dir == ".." {
		dir = ""
	}
	if relocate && len(dir) > 0 {
		_ = conn.ChangeDir(dir)
	}

	parent, _ = conn.CurrentDir()
	if isFile {
		name = qbc.Paths.FileName(strings.ReplaceAll(path, dir, ""), true)
	} else {
		name = qbc.Paths.FileName(path, false)
	}

	return
}

func (instance *VfsFtp) mkDir(conn *ftp.ServerConn, path string) (err error) {
	_, err = conn.NameList(path) // expected error do not exists
	if nil != err {
		err = conn.MakeDir(instance.absolutize(path))
		if instance.isValidError(err) {
			return err
		}
	}
	return nil
}

func (instance *VfsFtp) stat(conn *ftp.ServerConn, path string) (file *vfscommons.VfsFile, err error) {
	curDir, _ := conn.CurrentDir()

	_, name := instance.splitPath(conn, path, true)
	if !vfscommons.IsFile(path) {
		_ = conn.ChangeDir("..")
	}

	var list []*ftp.Entry
	list, err = conn.List(".")
	if nil != err {
		return
	}

	if len(name) > 0 && name != "." {
		for _, entry := range list {
			if nil != entry && entry.Name == name {
				fullPath := instance.Path()
				file = vfscommons.NewVfsFile(fullPath, instance.Path(), entry)
				break
			}
		}
	}

	_ = conn.ChangeDir(curDir)

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	VfsSftpConnection
//----------------------------------------------------------------------------------------------------------------------

type VfsFtpConnection struct {
	user     string
	password string
	host     string
	port     int

	conn *ftp.ServerConn
}

func NewVfsFtpConnection(user, password string, host string, port int) *VfsFtpConnection {
	instance := new(VfsFtpConnection)
	instance.user = user
	instance.password = password
	instance.host = host
	instance.port = port

	return instance
}

func (instance *VfsFtpConnection) Open() (*ftp.ServerConn, error) {
	if nil != instance.conn {
		// test existing connection
		err := instance.conn.NoOp()
		if nil != err {
			instance.conn = nil
		}
	}
	// creates connection if does not exists
	if nil == instance.conn {
		conn, err := ftp.Dial(fmt.Sprintf("%v:%v", instance.host, instance.port), ftp.DialWithTimeout(5*time.Second))
		if nil != err {
			return nil, err
		}
		if len(instance.user) > 0 {
			err = conn.Login(instance.user, instance.password)
			if err != nil {
				return nil, err
			}
		} else {
			err = conn.Login("anonymous", "anonymous")
			if err != nil {
				return nil, err
			}
		}
		instance.conn = conn
	}
	return instance.conn, nil
}

func (instance *VfsFtpConnection) Close() {
	if nil != instance.conn {
		_ = instance.conn.Logout()
		_ = instance.conn.Quit()
		instance.conn = nil
	}
}

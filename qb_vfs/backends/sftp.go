package backends

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/pkg/sftp"
	qbc "github.com/rskvp/qb-core"
	vfscommons "github.com/rskvp/qb-lib/qb_vfs/commons"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// https://sftptogo.com/blog/go-sftp/
// https://stackoverflow.com/questions/24437809/connect-to-a-server-using-ssh-and-a-pem-key-with-golang

//----------------------------------------------------------------------------------------------------------------------
//	VfsSftp
//----------------------------------------------------------------------------------------------------------------------

type VfsSftp struct {
	settings *vfscommons.VfsSettings

	host     string
	port     int
	user     string
	password string
	key      []byte

	conn *VfsSftpConnection

	startDir string
	curDir   string
}

func NewVfsSftp(settings *vfscommons.VfsSettings) (instance *VfsSftp, err error) {
	instance = new(VfsSftp)
	instance.settings = settings

	err = instance.init()

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsSftp) String() string {
	return qbc.JSON.Stringify(instance.settings)
}

func (instance *VfsSftp) Close() {
	if nil != instance.conn {
		instance.conn.Close()
	}
}

func (instance *VfsSftp) Path() string {
	return instance.curDir
}

func (instance *VfsSftp) Cd(path string) (bool, error) {
	if client, err := instance.connect(); nil != err {
		return false, err
	} else {
		path = instance.absolutize(path)
		_, err := client.ReadDir(path)
		if nil != err {
			return false, err
		}
		instance.curDir = path
		return true, nil
	}
}

func (instance *VfsSftp) Stat(path string) (*vfscommons.VfsFile, error) {
	if client, err := instance.connect(); nil != err {
		return nil, err
	} else {
		absolutePath := instance.absolutize(path)
		info, err := client.Stat(absolutePath)
		if nil != err {
			return nil, err
		}
		return vfscommons.NewVfsFile(absolutePath, instance.curDir, info), nil
	}
}

func (instance *VfsSftp) Exists(path string) (bool, error) {
	if client, err := instance.connect(); nil != err {
		return false, err
	} else {
		info, err := client.Stat(instance.absolutize(path))
		if nil != err {
			return false, err
		}
		return nil != info, nil
	}
}

func (instance *VfsSftp) List(dir string) ([]*vfscommons.VfsFile, error) {
	response := make([]*vfscommons.VfsFile, 0)
	if client, err := instance.connect(); nil != err {
		return nil, err
	} else {
		absolutePath := instance.absolutize(dir)
		files, err := client.ReadDir(absolutePath)
		if nil != err {
			return nil, err
		}
		for _, file := range files {
			response = append(response, vfscommons.NewVfsFile(absolutePath, absolutePath, file))
		}
	}
	return response, nil
}

func (instance *VfsSftp) Read(path string) ([]byte, error) {
	if client, err := instance.connect(); nil != err {
		return nil, err
	} else {
		file, err := client.OpenFile(instance.absolutize(path), os.O_RDONLY)
		if nil != err {
			return nil, err
		}
		var buf bytes.Buffer
		_, err = file.WriteTo(&buf)
		if nil != err {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func (instance *VfsSftp) Write(data []byte, target string) (int, error) {
	if client, err := instance.connect(); nil != err {
		return 0, err
	} else {
		target = instance.absolutize(target)
		err = instance.MkDir(qbc.Paths.Dir(target))
		if nil != err {
			return 0, err
		}
		f, err := client.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		if nil != err {
			return 0, err
		}
		defer f.Close()

		return f.Write(data)
	}
}

func (instance *VfsSftp) Download(source, target string) ([]byte, error) {
	data, err := instance.Read(source)
	if nil != err {
		return nil, err
	}
	_, err = qbc.IO.WriteBytesToFile(data, target)
	if nil != err {
		return nil, err
	}
	return data, nil
}

func (instance *VfsSftp) Remove(path string) error {
	if client, err := instance.connect(); nil != err {
		return err
	} else {
		path = instance.absolutize(path)
		info, err := client.Stat(path)
		if nil!=err{
			return err
		}
		if info.IsDir(){
			return client.RemoveDirectory(path)
		} else {
			return client.Remove(path)
		}
	}
}

func (instance *VfsSftp) MkDir(path string) error {
	if client, err := instance.connect(); nil != err {
		return err
	} else {
		return client.MkdirAll(instance.absolutize(path))
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *VfsSftp) init() error {
	if nil != instance.settings {
		// prepare configuration
		instance.host, instance.port = vfscommons.SplitHost(instance.settings)
		user := instance.settings.Auth.User
		password := instance.settings.Auth.Password
		key, err := vfscommons.ReadKey(instance.settings.Auth.Key)
		if nil != err {
			return err
		}
		instance.user = user
		instance.password = password
		instance.key = key

		instance.conn = NewVfsSftpConnection(instance.user, instance.password, instance.key, instance.host, instance.port)

		_, err = instance.connect()

		return err
	}
	return vfscommons.ErrorMissingConfiguration
}

func (instance *VfsSftp) connect() (*sftp.Client, error) {
	if nil != instance.conn {
		client, err := instance.conn.Open()
		if nil != err {
			instance.conn.Close()
			return nil, err
		}
		instance.startDir, err = client.Getwd()
		if len(instance.curDir) == 0 {
			instance.curDir = instance.startDir
		}

		return client, err
	}
	return nil, vfscommons.ErrorMissingConnection
}

func (instance *VfsSftp) absolutize(p string) string {
	return vfscommons.Absolutize(instance.curDir, p)
}

//----------------------------------------------------------------------------------------------------------------------
//	VfsSftpConnection
//----------------------------------------------------------------------------------------------------------------------

type VfsSftpConnection struct {
	user     string
	password string
	key      []byte
	host     string
	port     int

	conn   *ssh.Client
	client *sftp.Client
}

func NewVfsSftpConnection(user, password string, key []byte, host string, port int) *VfsSftpConnection {
	instance := new(VfsSftpConnection)
	instance.user = user
	instance.password = password
	instance.key = key
	instance.host = host
	instance.port = port

	return instance
}

func (instance *VfsSftpConnection) Open() (*sftp.Client, error) {
	if nil == instance.conn && nil == instance.client {
		var auths []ssh.AuthMethod

		// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
		// for communication with other processes.
		if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
			auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
		}

		// Use password authentication if provided
		if len(instance.password) > 0 {
			auths = append(auths, ssh.Password(instance.password))
		}

		// Initialize client configuration
		config := ssh.ClientConfig{
			User:            instance.user,
			Auth:            auths,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			// HostKeyCallback: ssh.FixedHostKey(hostKey),
		}

		// Connect to server
		addr := fmt.Sprintf("%s:%d", instance.host, instance.port)
		conn, err := ssh.Dial("tcp", addr, &config)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to connecto to [%s]: %v\n", addr, err))
		}

		// Create new SFTP client
		client, err := sftp.NewClient(conn)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to start SFTP subsystem: %v\n", err))
		}

		instance.conn = conn
		instance.client = client
	}
	return instance.client, nil
}

func (instance *VfsSftpConnection) Close() {
	if nil != instance.conn && nil != instance.client {
		_ = instance.conn.Close()
		_ = instance.client.Close()
		instance.conn = nil
		instance.client = nil
	}
}


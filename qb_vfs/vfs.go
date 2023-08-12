package qb_vfs

import (
	qbc "github.com/rskvp/qb-core"
	vfsbackends "github.com/rskvp/qb-lib/qb_vfs/backends"
	vfscommons "github.com/rskvp/qb-lib/qb_vfs/commons"
)

type VFSHelper struct {
}

var VFS *VFSHelper

func init() {
	VFS = new(VFSHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *VFSHelper) New(args ...interface{}) (vfscommons.IVfs, error) {
	switch len(args) {
	case 1:
		arg := args[0]
		if s, b := arg.(string); b {
			settings, err := vfscommons.LoadVfsSettings(s)
			if nil != err {
				return nil, err
			}
			return newVfs(settings)
		} else if c, b := arg.(vfscommons.VfsSettings); b {
			return newVfs(&c)
		} else if p, b := arg.(*vfscommons.VfsSettings); b {
			return newVfs(p)
		} else {
			s = qbc.JSON.Stringify(arg)
			settings, err := vfscommons.ParseVfsSettings(s)
			if nil != err {
				return nil, err
			}
			return newVfs(settings)
		}
	default:
		return nil, vfscommons.ErrorMismatchConfiguration
	}
	//return nil, vfscommons.MismatchConfigurationError
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func newVfs(settings *vfscommons.VfsSettings) (vfscommons.IVfs, error) {
	schema := settings.Schema()
	switch schema {
	case vfscommons.SchemaOS:
		return vfsbackends.NewVfsOS(settings)
	case vfscommons.SchemaSFTP:
		return vfsbackends.NewVfsSftp(settings)
	case vfscommons.SchemaFTP:
		return vfsbackends.NewVfsFtp(settings)
	default:
		return nil, qbc.Errors.Prefix(vfscommons.ErrorUnsupportedSchema, schema+": ")
	}
}

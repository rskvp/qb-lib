package backends

import (
	vfscommons "github.com/rskvp/qb-lib/qb_vfs/commons"
)

// GOOGLE DRIVE SUPPORT
// https://developers.google.com/drive/api/v3/quickstart/go

//----------------------------------------------------------------------------------------------------------------------
//	VfsGD
//----------------------------------------------------------------------------------------------------------------------

type VfsGD struct {
	settings *vfscommons.VfsSettings

	user string

	startDir string
	curDir   string
}




package qbl

import (
	_ "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_html"
	"github.com/rskvp/qb-lib/qb_script"

	"github.com/rskvp/qb-lib/qb_http"
	"github.com/rskvp/qb-lib/qb_vfs"
	_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

var VFS *qb_vfs.VFSHelper
var Http *qb_http.HttpHelper
var HTML *qb_html.HTMLHelper
var Scripting *qb_script.ScriptingHelper

func init() {
	VFS = qb_vfs.VFS
	Http = qb_http.Http
	HTML = qb_html.HTML
	Scripting = qb_script.Scripting
}

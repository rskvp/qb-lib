package qb_script

import (
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/auth0"
	"github.com/rskvp/qb-lib/qb_script/modules/dbal"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/console"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/process"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
	_ "github.com/rskvp/qb-lib/qb_script/modules/defaults/util"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/window"
	"github.com/rskvp/qb-lib/qb_script/modules/elasticsearch"
	"github.com/rskvp/qb-lib/qb_script/modules/fs"
	"github.com/rskvp/qb-lib/qb_script/modules/http"
	"github.com/rskvp/qb-lib/qb_script/modules/linereader"

	"github.com/rskvp/qb-lib/qb_script/modules/nio"
	"github.com/rskvp/qb-lib/qb_script/modules/nodemailer"
	"github.com/rskvp/qb-lib/qb_script/modules/nosql"
	"github.com/rskvp/qb-lib/qb_script/modules/path"
	"github.com/rskvp/qb-lib/qb_script/modules/showcase_engine"
	"github.com/rskvp/qb-lib/qb_script/modules/sql"
	"github.com/rskvp/qb-lib/qb_script/modules/sys"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/cryptoutils"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/dateutils"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/executils"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/fileutils"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/rndutils"
	"github.com/rskvp/qb-lib/qb_script/modules/utils/templateutils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ModuleRegistry struct {
	loader   require.SourceLoader
	registry *require.Registry
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewModuleRegistry(loader require.SourceLoader) *ModuleRegistry {
	instance := new(ModuleRegistry)
	instance.loader = loader
	instance.registry = require.NewRegistryWithLoader(loader)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ModuleRegistry) Start(engine *ScriptEngine) *commons.RuntimeContext {
	context := &commons.RuntimeContext{
		Uid:       &engine.Name,
		Workspace: engine.Root,
		Runtime:   engine.runtime,
		Arguments: []interface{}{
			&engine.Root, &engine.Name, &engine.Silent, &engine.LogLevel,
			&engine.ResetLogOnEachRun, engine.GetLogger, engine.LogFile,
		},
	}

	// creates context registry
	instance.registry.Enable(context)

	// start engine if not already started
	engine.Open()

	// add support to console and other defaults
	console.Enable(context)
	process.Enable(context)
	window.Enable(context)

	//-- add modules --//

	// auth0
	auth0.Enable(context)
	// dbal
	dbal.Enable(context)
	// elastic search
	elasticsearch.Enable(context)
	// showcase qb_sms_engine
	showcase_engine.Enable(context)
	// file system (Nodejs clone)
	fs.Enable(context)
	// http
	http.Enable(context)
	// line text reader
	linereader.Enable(context)
	// network io
	nio.Enable(context)
	// email sender
	nodemailer.Enable(context)
	// nosql
	nosql.Enable(context)
	// path (Nodejs clone)
	path.Enable(context)
	// sql layer
	sql.Enable(context)
	// system
	sys.Enable(context)
	// utility
	cryptoutils.Enable(context)
	dateutils.Enable(context)
	executils.Enable(context)
	fileutils.Enable(context)
	rndutils.Enable(context)
	templateutils.Enable(context)

	return context
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

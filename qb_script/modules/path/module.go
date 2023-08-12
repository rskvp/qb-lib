package path

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "path"

type Path struct {
	runtime *goja.Runtime
	util    *goja.Object
	root    string
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// Provides the platform-specific path delimiter:
//; for Windows
//: for POSIX
func (instance *Path) delimiter() goja.Value {
	val := os.PathListSeparator
	return instance.runtime.ToValue(string(val))
}

// Provides the platform-specific path segment separator: \ on Windows,  / on POSIX
func (instance *Path) separator() goja.Value {
	val := os.PathSeparator
	return instance.runtime.ToValue(string(val))
}

// path.dirname(path)
// The path.dirname() method returns the directory name of a path, similar to the Unix dirname command.
// Trailing directory separators are ignored, see path.sep.
func (instance *Path) dirname(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	val := filepath.Dir(path)

	return instance.runtime.ToValue(val)
}

// path.basename(path[, ext])
// The path.basename() method returns the last portion of a path, similar to the Unix basename command.
// Trailing directory separators are ignored, see path.sep.
func (instance *Path) basename(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	ext := ""
	if len(call.Arguments) == 2 {
		ext = call.Argument(1).String()
	}
	val := strings.Replace(filepath.Base(path), ext, "", 1)

	return instance.runtime.ToValue(val)
}

// path.extname(path)
// The path.extname() method returns the extension of the path, from the last occurrence of the . (period)
// character to end of string in the last portion of the path.
// If there is no . in the last portion of the path, or if there are no . characters other than the
// first character of the basename of path (see path.basename()) , an empty string is returned.
func (instance *Path) extname(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	val := qbc.Paths.Extension(path)

	return instance.runtime.ToValue(val)
}

// path.format(pathObject)
// The path.format() method returns a path string from an object. This is the opposite of path.parse().
// When providing properties to the pathObject remember that there are combinations where one property
// has priority over another:
// pathObject.root is ignored if pathObject.dir is provided
// pathObject.ext and pathObject.name are ignored if pathObject.base exists
func (instance *Path) format(call goja.FunctionCall) goja.Value {
	obj := qbc.Convert.ToMap(call.Argument(0).Export())
	if nil != obj {
		dir := qbc.Reflect.GetString(obj, "dir")
		root := qbc.Reflect.GetString(obj, "root")
		base := qbc.Reflect.GetString(obj, "base")
		name := qbc.Reflect.GetString(obj, "name")
		ext := qbc.Reflect.GetString(obj, "ext")
		if len(dir) == 0 {
			dir = root
		}
		if len(name) == 0 {
			name = base
		}
		if e := qbc.Paths.Extension(name); len(e) > 0 {
			ext = ""
		}
		val := qbc.Paths.Concat(dir, name) + ext

		return instance.runtime.ToValue(val)
	}
	return goja.Undefined()
}

// path.isAbsolute(path)
// The path.isAbsolute() method determines if path is an absolute path.
// If the given path is a zero-length string, false will be returned.
func (instance *Path) isAbsolute(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	val := qbc.Paths.IsAbs(path)

	return instance.runtime.ToValue(val)
}

// path.join([...paths])
// The path.join() method joins all given path segments together using the platform-specific
// separator as a delimiter, then normalizes the resulting path.
// Zero-length path segments are ignored. If the joined path string is a zero-length string
// then '.' will be returned, representing the current working directory.
func (instance *Path) join(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) > 0 {
		paths := make([]string, 0)
		for _, v := range call.Arguments {
			paths = append(paths, v.String())
		}
		val := qbc.Paths.Concat(paths...)
		return instance.runtime.ToValue(val)
	}
	return goja.Undefined()
}

// path.normalize(path)
// The path.normalize() method normalizes the given path, resolving '..' and '.' segments.
// When multiple, sequential path segment separation characters are found
// (e.g. / on POSIX and either \ or / on Windows), they are replaced by a single instance of
// the platform-specific path segment separator (/ on POSIX and \ on Windows).
// Trailing separators are preserved.
// If the path is a zero-length string, '.' is returned, representing the current working directory.
func (instance *Path) normalize(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	val := filepath.Clean(path)

	return instance.runtime.ToValue(val)
}

// path.parse(path)
// The path.parse() method returns an object whose properties represent significant elements of the path. Trailing directory separators are ignored, see path.sep.
// The returned object will have the following properties:
// * dir <string>
// * root <string>
// * base <string>
// * name <string>
// * ext <string>
func (instance *Path) parse(call goja.FunctionCall) goja.Value {
	path := call.Argument(0).String()
	dir, base := filepath.Split(path)
	name := qbc.Paths.FileName(base, false)
	ext := qbc.Paths.Extension(base)
	root := filepath.VolumeName(path) + string(os.PathSeparator)

	return instance.runtime.ToValue(map[string]string{
		"root": root,
		"dir":  dir,
		"base": base,
		"name": name,
		"ext":  ext,
	})
}

// path.relative(from, to)
// The path.relative() method returns the relative path from from to to based on the current working
// directory. If from and to each resolve to the same path (after calling path.resolve() on each),
// a zero-length string is returned.
// If a zero-length string is passed as from or to, the current working directory will be used
// instead of the zero-length strings.
func (instance *Path) relative(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 2 {
		from := call.Argument(0).String()
		to := call.Argument(1).String()
		val, _ := filepath.Rel(from, to)

		return instance.runtime.ToValue(val)
	}
	return goja.Undefined()
}

// path.resolve([...paths])
// The path.resolve() method resolves a sequence of paths or path segments into an absolute path.
// The given sequence of paths is processed from right to left, with each subsequent path prepended
// until an absolute path is constructed. For instance, given the sequence of path
// segments: /foo, /bar, baz, calling path.resolve('/foo', '/bar', 'baz') would
// return /bar/baz because 'baz' is not an absolute path but '/bar' + '/' + 'baz' is.
// If after processing all given path segments an absolute path has not yet been generated,
// the current working directory is used.
// The resulting path is normalized and trailing slashes are removed unless the path is resolved
// to the root directory.
// Zero-length path segments are ignored.
// If no path segments are passed, path.resolve() will return the absolute path of the current
// working directory.
func (instance *Path) resolve(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) > 0 {
		paths := make([]string, 0)
		for i, v := range call.Arguments {
			path := v.String()
			if i == 0 && !filepath.IsAbs(path) {
				path, _ = filepath.Abs(path)
			}
			paths = append(paths, path)
		}
		val := filepath.Join(paths...)

		return instance.runtime.ToValue(val)
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	instance := &Path{
		runtime: runtime,
	}

	if len(args) > 0 {
		root := qbc.Reflect.ValueOf(args[0]).String()
		if len(root) > 0 {
			instance.root = root
		}
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("delimiter", instance.delimiter())
	_ = o.Set("sep", instance.separator())
	_ = o.Set("dirname", instance.dirname)
	_ = o.Set("basename", instance.basename)
	_ = o.Set("extname", instance.extname)
	_ = o.Set("format", instance.format)
	_ = o.Set("isAbsolute", instance.isAbsolute)
	_ = o.Set("join", instance.join)
	_ = o.Set("normalize", instance.normalize)
	_ = o.Set("parse", instance.parse)
	_ = o.Set("relative", instance.relative)
	_ = o.Set("resolve", instance.resolve)
}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})

	//ctx.Runtime.Set(NAME, require.Require(ctx.Runtime, NAME))
}

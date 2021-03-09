package file

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/hocv/gin-swagger-gen/lib/common"
)

// File file
type File struct {
	dirty      bool
	src        string                     // file name or source
	file       *dst.File                  // dst file
	pkg        string                     // package name
	globalVars map[string]string          // global vars
	imports    map[string]string          // import
	types      map[string]string          // types
	funcs      map[string]*dst.FuncDecl   // functions
	structs    map[string]*dst.StructType // structs
}

// New file. src : path of go file or source
func New(src string) (*File, error) {
	fileSet := token.NewFileSet()
	var file *dst.File
	var err error

	if filepath.Ext(src) == ".go" {
		file, err = decorator.ParseFile(fileSet, src, nil, parser.ParseComments)
	} else {
		file, err = decorator.ParseFile(fileSet, "", src, parser.ParseComments)
	}

	if err != nil {
		return nil, err
	}

	f := &File{
		src:        src,
		file:       file,
		dirty:      false,
		globalVars: map[string]string{},
		imports:    map[string]string{},
		types:      map[string]string{},
		funcs:      map[string]*dst.FuncDecl{},
		structs:    map[string]*dst.StructType{},
	}
	f.parse()
	return f, nil
}

func (f *File) Dirty() {
	f.dirty = true
}

// Save file
func (f *File) Save() error {
	if !f.dirty || filepath.Ext(f.src) != ".go" {
		return nil
	}
	var buf bytes.Buffer
	if err := decorator.Fprint(&buf, f.file); err != nil {
		return err
	}
	return ioutil.WriteFile(f.src, buf.Bytes(), 0666)
}

// Pkg package name
func (f *File) Pkg() string {
	return f.pkg
}

func (f *File) parse() {
	f.pkg = f.file.Name.String()

	for _, spec := range f.file.Imports {
		f.imports[spec.Path.Value] = spec.Name.String()
	}

	for _, decl := range f.file.Decls {
		switch decl.(type) {
		case *dst.FuncDecl:
			fn := decl.(*dst.FuncDecl)
			f.funcs[fn.Name.String()] = fn
		case *dst.GenDecl:
			gd := decl.(*dst.GenDecl)
			for _, spec := range gd.Specs {
				switch spec.(type) {
				case *dst.TypeSpec:
					ts := spec.(*dst.TypeSpec)
					tn := ts.Name.String()
					switch ts.Type.(type) {
					case *dst.StructType:
						f.structs[tn] = ts.Type.(*dst.StructType)
					case *dst.Ident:
						f.types[tn] = ts.Type.(*dst.Ident).Name
					case *dst.ArrayType:
						at, ok := ts.Type.(*dst.ArrayType).Elt.(*dst.Ident)
						if ok {
							f.types[tn] = fmt.Sprintf("[]%s", at.String())
						}
					case *dst.MapType:
						mt := ts.Type.(*dst.MapType)
						f.types[tn] = fmt.Sprintf("map[%s]%s", mt.Key.(*dst.Ident).String(), mt.Value.(*dst.Ident).String())
					}
				case *dst.ValueSpec:
					for k, v := range common.GetVars(gd) {
						f.globalVars[k] = v
					}
				}
			}
		}
	}
}

// Func search function with name
func (f *File) Func(name string) (*dst.FuncDecl, error) {
	fd, ok := f.funcs[name]
	if !ok {
		return nil, common.ErrNotFind
	}
	return fd, nil
}

// FuncWithParam search functions with param
func (f *File) FuncWithParam(param string) (fds []*dst.FuncDecl) {
	for _, fd := range f.funcs {
		ps := common.GetFuncParamByType(fd, param)
		if len(ps) > 0 {
			fds = append(fds, fd)
			break
		}
	}
	return
}

// FuncWithSelector search functions that contain this expr
func (f *File) FuncWithSelector(expr string) (fds []*dst.FuncDecl) {
	for _, fd := range f.funcs {
		for _, stmt := range fd.Body.List {
			if common.CheckSelectorExpr(stmt, expr) {
				fds = append(fds, fd)
				break
			}
		}
	}
	return
}

func (f *File) Imported(path string) (string, bool) {
	path = fmt.Sprintf("\"%s\"", path)
	alias, ok := f.imports[path]
	return alias, ok
}

func (f *File) DefaultImport(path string, value string) (string, bool) {
	path = fmt.Sprintf("\"%s\"", path)
	alias, ok := f.imports[path]
	if !ok {
		return "", false
	}
	if alias == "<nil>" {
		return value, true
	}

	return alias, true
}

// GlobalVars global vars in file
func (f *File) GlobalVars() map[string]string {
	return f.imports
}

func (f *File) Struct(name string) (*dst.StructType, error) {
	st, ok := f.structs[name]
	if !ok {
		return nil, common.ErrNotFind
	}
	return st, nil
}

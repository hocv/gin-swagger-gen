package ast

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pkg/errors"
)

var ErrNotFind = errors.New("not found")

// Ast ast
type Ast struct {
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

// New ast. src : path of go file or source
func New(src string) (*Ast, error) {
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

	a := &Ast{
		src:        src,
		file:       file,
		dirty:      false,
		globalVars: map[string]string{},
		imports:    map[string]string{},
		types:      map[string]string{},
		funcs:      map[string]*dst.FuncDecl{},
		structs:    map[string]*dst.StructType{},
	}
	a.parse()
	return a, nil
}

func (a *Ast) Dirty() {
	a.dirty = true
}

// Save file
func (a *Ast) Save() error {
	if !a.dirty || filepath.Ext(a.src) != ".go" {
		return nil
	}
	var buf bytes.Buffer
	if err := decorator.Fprint(&buf, a.file); err != nil {
		return err
	}
	return ioutil.WriteFile(a.src, buf.Bytes(), 0666)
}

// Pkg package name
func (a *Ast) Pkg() string {
	return a.pkg
}

func (a *Ast) parse() {
	a.pkg = a.file.Name.String()

	for _, spec := range a.file.Imports {
		a.imports[spec.Path.Value] = spec.Name.String()
	}

	for _, decl := range a.file.Decls {
		switch decl.(type) {
		case *dst.FuncDecl:
			fn := decl.(*dst.FuncDecl)
			a.funcs[fn.Name.String()] = fn
		case *dst.GenDecl:
			gd := decl.(*dst.GenDecl)
			for _, spec := range gd.Specs {
				switch spec.(type) {
				case *dst.TypeSpec:
					ts := spec.(*dst.TypeSpec)
					tn := ts.Name.String()
					switch ts.Type.(type) {
					case *dst.StructType:
						a.structs[tn] = ts.Type.(*dst.StructType)
					case *dst.Ident:
						a.types[tn] = ts.Type.(*dst.Ident).Name
					case *dst.ArrayType:
						at, ok := ts.Type.(*dst.ArrayType).Elt.(*dst.Ident)
						if ok {
							a.types[tn] = fmt.Sprintf("[]%s", at.String())
						}
					case *dst.MapType:
						mt := ts.Type.(*dst.MapType)
						a.types[tn] = fmt.Sprintf("map[%s]%s", mt.Key.(*dst.Ident).String(), mt.Value.(*dst.Ident).String())
					}
				case *dst.ValueSpec:
					for k, v := range GetVars(gd) {
						a.globalVars[k] = v
					}
				}
			}
		}
	}
}

// Func search function with name
func (a *Ast) Func(name string) (*dst.FuncDecl, error) {
	fd, ok := a.funcs[name]
	if !ok {
		return nil, ErrNotFind
	}
	return fd, nil
}

// FuncWithParam search functions with param
func (a *Ast) FuncWithParam(param string) (fds []*dst.FuncDecl) {
	for _, fd := range a.funcs {
		ps := GetFuncParamByType(fd, param)
		if len(ps) > 0 {
			fds = append(fds, fd)
			break
		}
	}
	return
}

// FuncWithSelector search functions that contain this expr
func (a *Ast) FuncWithSelector(expr string) (fds []*dst.FuncDecl) {
	for _, fd := range a.funcs {
		for _, stmt := range fd.Body.List {
			if CheckSelectorExpr(stmt, expr) {
				fds = append(fds, fd)
				break
			}
		}
	}
	return
}

func (a *Ast) Imported(path string) (string, bool) {
	path = fmt.Sprintf("\"%s\"", path)
	alias, ok := a.imports[path]
	return alias, ok
}

func (a *Ast) DefaultImport(path string, value string) string {
	path = fmt.Sprintf("\"%s\"", path)
	alias, ok := a.imports[path]
	if !ok || alias == "<nil>" {
		return value
	}
	return alias
}

// GlobalVars global vars in file
func (a *Ast) GlobalVars() map[string]string {
	return a.imports
}

func (a *Ast) Struct(name string) (*dst.StructType, error) {
	st, ok := a.structs[name]
	if !ok {
		return nil, ErrNotFind
	}
	return st, nil
}

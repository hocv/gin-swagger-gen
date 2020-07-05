package ast

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pkg/errors"
)

var ErrNotFind = errors.New("not find error")

// Ast ast
type Ast struct {
	path string    // file path
	file *dst.File // dst file
}

// New ast. path : path of go file
func New(path string) (*Ast, error) {
	fileSet := token.NewFileSet()
	file, err := decorator.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &Ast{
		path: path,
		file: file,
	}, nil
}

// Save file
func (a *Ast) Save() error {
	var buf bytes.Buffer
	if err := decorator.Fprint(&buf, a.file); err != nil {
		return err
	}
	return ioutil.WriteFile(a.path, buf.Bytes(), 0666)
}

// Pkg package name
func (a *Ast) Pkg() string {
	return a.file.Name.String()
}

// Func search function with name
func (a *Ast) Func(name string) (fds []*dst.FuncDecl) {
	for _, decl := range a.file.Decls {
		fd, ok := decl.(*dst.FuncDecl)
		if !ok || fd.Name.String() != name {
			continue
		}
		fds = append(fds, fd)
	}
	return
}

// FuncWithParam search functions with param
func (a *Ast) FuncWithParam(param string) (fds []*dst.FuncDecl) {
	for _, decl := range a.file.Decls {
		fd, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}
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
	for _, decl := range a.file.Decls {
		fd, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}
		for _, stmt := range fd.Body.List {
			if CheckSelectorExpr(stmt, expr) {
				fds = append(fds, fd)
				break
			}
		}
	}
	return
}

func (a *Ast) Imported(path string) (alias string, exist bool) {
	path = fmt.Sprintf("\"%s\"", path)
	for _, spec := range a.file.Imports {
		if spec.Path.Value != path {
			continue
		}
		if spec.Name != nil {
			alias = spec.Name.String()
		}
		exist = true
		break
	}
	return
}

func (a *Ast) DefaultImport(path string, value string) string {
	path = fmt.Sprintf("\"%s\"", path)
	for _, spec := range a.file.Imports {
		if spec.Path.Value != path {
			continue
		}
		if spec.Name != nil {
			return spec.Name.String()
		}
	}
	return value
}

// GlobalVars global vars in file
func (a *Ast) GlobalVars() map[string]string {
	vars := make(map[string]string)
	for _, decl := range a.file.Decls {
		gd, ok := decl.(*dst.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for k, v := range GetVars(gd) {
			vars[k] = v
		}
	}
	return vars
}

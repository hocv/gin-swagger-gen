package pkg

import (
	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/lib/common"
	"github.com/hocv/gin-swagger-gen/lib/file"
)

type Pkg struct {
	Name  string
	files []*file.File
}

func New(name string) *Pkg {
	return &Pkg{Name: name}
}

func (p *Pkg) AddFile(files ...*file.File) {
	for _, f := range files {
		if f.Pkg() != p.Name {
			continue
		}
		p.files = append(p.files, f)
	}
}

// Func search function by name
func (p *Pkg) GetFunc(name string) map[*file.File]*dst.FuncDecl {
	af := make(map[*file.File]*dst.FuncDecl)
	for _, a := range p.files {
		fds, err := a.Func(name)
		if err != nil {
			continue
		}
		af[a] = fds
	}
	return af
}

func (p *Pkg) GetFuncWithSelector(expr string) map[*file.File][]*dst.FuncDecl {
	af := make(map[*file.File][]*dst.FuncDecl)
	for _, a := range p.files {
		fs := a.FuncWithSelector(expr)
		af[a] = fs
	}
	return af
}

func (p *Pkg) GetFuncWithParam(param string) map[*file.File][]*dst.FuncDecl {
	af := make(map[*file.File][]*dst.FuncDecl)
	for _, a := range p.files {
		fs := a.FuncWithParam(param)
		af[a] = fs
	}
	return af
}

func (p *Pkg) GetAstWithImported(pkg string) []*file.File {
	as := make([]*file.File, 0)
	for _, a := range p.files {
		_, b := a.Imported(pkg)
		if b {
			as = append(as, a)
		}
	}
	return as
}

func (p *Pkg) GetGlobalVar() map[string]string {
	vars := make(map[string]string)
	for _, a := range p.files {
		for k, v := range a.GlobalVars() {
			vars[k] = v
		}
	}
	return vars
}

func (p *Pkg) GetStruct(name string) (*dst.StructType, error) {
	for _, a := range p.files {
		stru, err := a.Struct(name)
		if err == nil {
			return stru, nil
		}
	}
	return nil, common.ErrNotFind
}

func (p *Pkg) GetImported(path string) (string, error) {
	for _, a := range p.files {
		alias, ok := a.Imported(path)
		if ok {
			return alias, nil
		}
	}
	return "", common.ErrNotFind
}

func (p *Pkg) GetDefaultImported(path, value string) (string, error) {
	for _, a := range p.files {
		alias, ok := a.DefaultImport(path, value)
		if ok {
			return alias, nil
		}
	}
	return "", common.ErrNotFind
}

// Save save file
func (p *Pkg) Save() error {
	for _, a := range p.files {
		if err := a.Save(); err != nil {
			return err
		}
	}
	return nil
}

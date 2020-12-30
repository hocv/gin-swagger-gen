package ast

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dave/dst"
)

const routineNum = 20

type Asts struct {
	as map[string][]*Ast
}

func NewAsts(dir string) Asts {
	arr := scanDir(dir)
	as := make(map[string][]*Ast, 0)
	for _, a := range arr {
		as[a.pkg] = append(as[a.pkg], a)
	}
	return Asts{
		as: as,
	}
}

func (a *Asts) Save() error {
	for _, asts := range a.as {
		for _, ast := range asts {
			if err := ast.Save(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Asts) routine(astFn func(ast *Ast, m *sync.Mutex)) {
	m := sync.Mutex{}
	w := sync.WaitGroup{}
	w.Add(routineNum)
	ch := make(chan *Ast, routineNum*routineNum)

	for i := 0; i < routineNum; i++ {
		go func() {
			defer w.Done()
			for ast := range ch {
				astFn(ast, &m)
			}
		}()
	}

	for _, asts := range a.as {
		for _, ast := range asts {
			ch <- ast
		}
	}
	close(ch)
	w.Wait()
}

// Func search function by name
func (a *Asts) Func(name string) map[*Ast]*dst.FuncDecl {
	af := make(map[*Ast]*dst.FuncDecl)
	a.routine(func(ast *Ast, m *sync.Mutex) {
		fds, err := ast.Func(name)
		if err != nil {
			return
		}
		m.Lock()
		af[ast] = fds
		m.Unlock()
	})
	return af
}

func (a *Asts) FuncInPkg(funcName, pkgName string) (*Ast, *dst.FuncDecl, error) {
	asts, ok := a.as[pkgName]
	if !ok {
		return nil, nil, ErrNotFind
	}
	for _, ast := range asts {
		fd, err := ast.Func(funcName)
		if err != nil {
			continue
		}
		return ast, fd, nil
	}
	return nil, nil, ErrNotFind
}

func (a *Asts) FuncWithSelector(expr string) map[*Ast][]*dst.FuncDecl {
	af := make(map[*Ast][]*dst.FuncDecl)
	a.routine(func(ast *Ast, m *sync.Mutex) {
		fs := ast.FuncWithSelector(expr)
		m.Lock()
		af[ast] = fs
		m.Unlock()
	})
	return af
}

func (a *Asts) FuncWithParam(param string) map[*Ast][]*dst.FuncDecl {
	af := make(map[*Ast][]*dst.FuncDecl)
	a.routine(func(ast *Ast, m *sync.Mutex) {
		fs := ast.FuncWithParam(param)
		m.Lock()
		af[ast] = fs
		m.Unlock()
	})
	return af
}

func (a *Asts) Imported(pkg string) []*Ast {
	as := make([]*Ast, 0)
	a.routine(func(ast *Ast, m *sync.Mutex) {
		_, b := ast.Imported(pkg)
		if b {
			m.Lock()
			as = append(as, ast)
			m.Unlock()
		}
	})
	return as
}

func (a *Asts) GlobalVarInPkg(pkg string) map[string]string {
	asts, ok := a.as[pkg]
	if !ok {
		return nil
	}
	vars := make(map[string]string)
	for _, ast := range asts {
		for k, v := range ast.GlobalVars() {
			vars[k] = v
		}
	}
	return vars
}

// scanDir scan go file and parse to ast
func scanDir(dir string) (asts []*Ast) {
	if strings.Contains(dir, "vendor") {
		return
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range fs {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() {
			asts = append(asts, scanDir(path)...)
		}
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		g, err := New(path)
		if err != nil {
			continue
		}
		asts = append(asts, g)
	}
	return
}

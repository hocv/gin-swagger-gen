package ast

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dave/dst"
)

const routineNum = 20

type Asts []*Ast

func NewAsts(dir string) Asts {
	arr := scanDir(dir)
	asts := make([]*Ast, 0)
	for _, as := range arr {
		asts = append(asts, as)
	}
	return asts
}

func (as Asts) Save() error {
	for _, a := range as {
		if err := a.Save(); err != nil {
			return err
		}
	}
	return nil
}

func (as Asts) routine(astFn func(ast *Ast, m *sync.Mutex)) {
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

	for _, ast := range as {
		a := ast
		ch <- a
	}
	close(ch)
	w.Wait()
}

// Func search function by name
func (as Asts) Func(name string) map[*Ast][]*dst.FuncDecl {
	af := make(map[*Ast][]*dst.FuncDecl)
	as.routine(func(ast *Ast, m *sync.Mutex) {
		fds := ast.Func(name)
		m.Lock()
		af[ast] = fds
		m.Unlock()
	})
	return af
}

func (as Asts) FuncWithSelector(expr string) map[*Ast][]*dst.FuncDecl {
	af := make(map[*Ast][]*dst.FuncDecl)
	as.routine(func(ast *Ast, m *sync.Mutex) {
		fs := ast.FuncWithSelector(expr)
		m.Lock()
		af[ast] = fs
		m.Unlock()
	})
	return af
}

func (as Asts) FuncWithParam(param string) map[*Ast][]*dst.FuncDecl {
	af := make(map[*Ast][]*dst.FuncDecl)
	as.routine(func(ast *Ast, m *sync.Mutex) {
		fs := ast.FuncWithParam(param)
		m.Lock()
		af[ast] = fs
		m.Unlock()
	})
	return af
}

func (as Asts) Imported(pkg string) Asts {
	asts := make([]*Ast, 0)
	as.routine(func(ast *Ast, m *sync.Mutex) {
		_, b := ast.Imported(pkg)
		if b {
			m.Lock()
			asts = append(asts, ast)
			m.Unlock()
		}
	})
	return asts
}

// scanDir scan go file and parse to ast
func scanDir(dir string) (asts []*Ast) {
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

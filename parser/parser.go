package parser

import (
	"fmt"

	"github.com/hocv/gin-swagger-gen/lib/pkg"
	"github.com/hocv/gin-swagger-gen/lib/proj"
	"github.com/pkg/errors"
)

var ginPkg = "github.com/gin-gonic/gin"

type Parser struct {
	proj        *proj.Proj
	specifyFunc string
}

func New(specifyFunc string) *Parser {
	return &Parser{
		specifyFunc: specifyFunc,
		proj:        proj.New(),
	}
}

func (parser *Parser) ScanDir(dir string) {
	parser.proj.ScanDir(dir)
}

func (parser *Parser) Parse(justPrint bool) {
	var hdls []*handle
	ginFn := func(p *pkg.Pkg, expr string) {
		fds := p.GetFuncWithSelector(expr)
		for f, decls := range fds {
			for _, decl := range decls {
				rh := newRoute(parser.proj, expr, parser.specifyFunc)
				rh.Parse(f, decl)
				hdls = append(hdls, rh.Handles...)
			}
		}
	}

	ginAsts := parser.proj.GetPkgWithImported(ginPkg)
	for _, a := range ginAsts {
		alias, err := a.GetDefaultImported(ginPkg, "gin")
		if err != nil {
			continue
		}
		ginFn(a, fmt.Sprintf("%s.New", alias))
		ginFn(a, fmt.Sprintf("%s.Default", alias))
	}

	for _, hdl := range hdls {
		hdl.Parse()
		if justPrint {
			hdl.Print()
		} else {
			hdl.Merge()
		}
	}
}

func (parser *Parser) Save() error {
	if err := parser.proj.Save(); err != nil {
		return errors.Wrap(err, "gen save")
	}
	return nil
}

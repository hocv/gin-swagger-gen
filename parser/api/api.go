package api

import (
	"fmt"

	"github.com/hocv/gin-swagger-gen/ast"
)

const ginPkg = "github.com/gin-gonic/gin"

type Api struct {
	asts        *ast.Asts
	specifyFunc string
}

func NewApiParse(specifyFunc string) *Api {
	return &Api{
		specifyFunc: specifyFunc,
	}
}

func (api *Api) Parse(asts ast.Asts) error {
	api.asts = &asts

	ginFn := func(a *ast.Ast, expr string) {
		fds := a.FuncWithSelector(expr)
		for _, decl := range fds {
			parseRoute(api, a, decl, nil, expr)
		}
	}

	ginAsts := asts.Imported(ginPkg)
	for _, a := range ginAsts {
		alias := a.DefaultImport(ginPkg, "gin")
		ginFn(a, fmt.Sprintf("%s.New", alias))
		ginFn(a, fmt.Sprintf("%s.Default", alias))
	}
	return nil
}

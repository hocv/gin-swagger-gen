package api

import (
	"fmt"
	"strings"

	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/ast"
)

// splitDot split string with dot
func splitDot(str string) (string, string) {
	arr := strings.Split(str, ".")
	if len(arr) != 2 {
		return "", str
	}
	return arr[0], arr[1]
}

// fmtRoutePath remove "" and replace : to {}.
// e.g. "/user/:id" => /user/{id}
func fmtRoutePath(r string) string {
	r = strings.Trim(r, "\"")
	arr := strings.Split(r, "/")
	for i, s := range arr {
		if strings.HasPrefix(s, ":") {
			trim := strings.Trim(s, ":")
			arr[i] = fmt.Sprintf("{%s}", trim)
		}
	}
	return strings.Join(arr, "/")
}

// routePathParams params in path. "/user/{id}"
func routePathParams(r string) []string {
	r = strings.Trim(r, "\"")
	arr := strings.Split(r, "/")
	var params []string
	for _, s := range arr {
		if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			trim := strings.Trim(s, "{")
			trim = strings.Trim(trim, "}")
			params = append(params, trim)
		}
	}
	return params
}

func copyMap(m map[string]string) map[string]string {
	cp := make(map[string]string)
	if m == nil {
		return cp
	}
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// searchGinFunc search function with gin params
// ginType: Engine, Context. srv.ginHandel(c), srv is funcCall, ginHandel is funcName
func searchGinFunc(asts *ast.Asts,
	ginType string,
	funcCall string,
	funcName string,
	params []string,
	fn func(da *ast.Ast, fd *dst.FuncDecl, vs map[string]string)) {

	for a, decls := range asts.Func(funcName) {
		alias := a.DefaultImport(ginPkg, "gin")
		ginCtx := fmt.Sprintf("*%s.%s", alias, ginType)
		for _, decl := range decls {
			if len(ast.GetFuncParamByType(decl, ginCtx)) == 0 {
				continue
			}

			vs := make(map[string]string)
			if params != nil {
				for i, s := range ast.GetFuncParamList(decl) {
					vs[s] = params[i]
				}
			}

			// function
			if len(funcCall) == 0 {
				fn(a, decl, vs)
				continue
			}

			if decl.Recv == nil {
				continue
			}

			// method
			for _, field := range decl.Recv.List {
				ident, ok := field.Type.(*dst.Ident)
				if !ok || ident.Name != funcCall {
					continue
				}
				fn(a, decl, vs)
				return
			}
		}
	}
}

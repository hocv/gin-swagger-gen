package api

import (
	"strings"

	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/ast"
)

type routeHandle struct {
	api         *Api
	initExpr    string
	Vars        map[string]string
	RouteMap    map[string]string
	RouteFns    map[string]routeFunc
	specifyFunc string
}

func parseRoute(api *Api, dstAst *ast.Ast, decl *dst.FuncDecl, routeMap map[string]string, initExpr string, specifyFunc string) {
	r := &routeHandle{
		api:         api,
		initExpr:    initExpr,
		Vars:        make(map[string]string),
		RouteMap:    copyMap(routeMap),
		RouteFns:    map[string]routeFunc{},
		specifyFunc: specifyFunc,
	}

	// vars in function param
	for k, v := range ast.GetFuncParams(decl) {
		r.Vars[k] = v
	}

	// global vars
	for k, v := range dstAst.GlobalVars() {
		r.Vars[k] = v
	}

	// g := gin.New()
	_, ginFunc := splitDot(initExpr)
	r.RouteFns[ginFunc] = func(val string, cal string, call *dst.CallExpr) {
		r.RouteMap[val] = ""
	}

	// api := g.Grout("/api")
	r.RouteFns["Group"] = r.parseGroup

	// g.GET("/usr",handelFunc)
	for _, method := range []string{"POST", "GET", "DELETE", "PATCH", "PUT", "OPTIONS", "HEAD", "Any"} {
		r.RouteFns[method] = r.parseMethod
	}

	parseStmtList(decl.Body.List, r.Vars, r)
}

func (rh *routeHandle) Asts() *ast.Asts {
	return rh.api.asts
}

func (rh *routeHandle) Type() string {
	return "Engine"
}

func (rh *routeHandle) Cond(sel string) bool {
	_, ok := rh.RouteFns[sel]
	return ok
}

func (rh *routeHandle) Parser(val string, vat string, call *dst.CallExpr, vs map[string]string) {
	cal, sel := splitDot(vat)
	rh.RouteFns[sel](val, cal, call)
}

func (rh *routeHandle) Inter(a *ast.Ast, decl *dst.FuncDecl, vars map[string]string) {
	rm := copyMap(rh.RouteMap)
	for k, v := range vars {
		if vv, ok := rh.RouteMap[v]; ok {
			rm[k] = vv
			delete(rm, v)
		}
	}
	parseRoute(rh.api, a, decl, rm, rh.initExpr, rh.specifyFunc)
}

func (rh *routeHandle) parseGroup(val string, cal string, call *dst.CallExpr) {
	if len(call.Args) == 0 {
		return
	}
	// e.g. val := cal.Group("/api")
	routeBase, ok := rh.RouteMap[cal]
	if !ok {
		return
	}
	// "/api",fist arg of function is route path
	path := call.Args[0].(*dst.BasicLit).Value
	rh.RouteMap[val] = routeBase + fmtRoutePath(path)
}

// routeFunc handle function
type routeFunc func(val string, cal string, call *dst.CallExpr)

// parseMethod parse route handler, find the function and add common to it
// api.GET("usr",handleFunc)
func (rh *routeHandle) parseMethod(val string, cal string, call *dst.CallExpr) {
	if len(call.Args) < 2 {
		return
	}
	routeBase, ok := rh.RouteMap[cal]
	if !ok {
		return
	}

	// first arg of function is route path
	firstArg := call.Args[0].(*dst.BasicLit).Value
	// just use last handle function, middle functions maybe middleware
	lastArg := ast.ToStr(call.Args[len(call.Args)-1])
	handleCall, handleFn := splitDot(lastArg)
	v, ok := rh.Vars[handleCall]
	if ok {
		handleCall = v
	}

	// if specify the function, ignore others
	if len(rh.specifyFunc) > 0 && rh.specifyFunc != handleFn {
		return
	}

	sel, ok := call.Fun.(*dst.SelectorExpr)
	if !ok {
		return
	}

	routePath := routeBase + fmtRoutePath(firstArg)

	cmt := &comment{
		Summary:     handleFn,
		RoutePath:   routePath,
		PathParams:  routePathParams(routePath),
		RouteMethod: strings.ToLower(sel.Sel.Name),
	}
	searchGinFunc(rh.Asts(), "Context", handleCall, handleFn, nil, func(da *ast.Ast, sd *dst.FuncDecl, vs map[string]string) {
		parseFuncHandle(rh.api, da, nil, sd, cmt)
	})
}

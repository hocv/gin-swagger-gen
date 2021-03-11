package parser

import (
	"github.com/hocv/gin-swagger-gen/parser/comment"

	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/lib/common"
	"github.com/hocv/gin-swagger-gen/lib/file"
	"github.com/hocv/gin-swagger-gen/lib/proj"
)

var routeParsers = map[string]routeParser{
	"Group":   parseRouteGroup,
	"POST":    parseRouteMethod,
	"GET":     parseRouteMethod,
	"DELETE":  parseRouteMethod,
	"PATCH":   parseRouteMethod,
	"PUT":     parseRouteMethod,
	"OPTIONS": parseRouteMethod,
	"HEAD":    parseRouteMethod,
	"Any":     parseRouteMethod,
}

// routeFunc handle function
type routeParser func(rh *route, val string, cal string, call *dst.CallExpr)

// parseRouteGroup api := g.Grout("/api")
func parseRouteGroup(rh *route, val string, cal string, call *dst.CallExpr) {
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

// parseRouteMethod g.GET("/usr",handelFunc)
func parseRouteMethod(rh *route, val string, cal string, call *dst.CallExpr) {
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
	lastArg := common.ToStr(call.Args[len(call.Args)-1])
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

	cmt := comment.New(handleFn, routeBase, routePath, sel.Sel.Name)
	pps := routePathParams(routePath)
	for _, p := range pps {
		if len(p) > 0 {
			cmt.AddParam(comment.NewPathParam(p, "string", p))
		}
	}

	ffs := rh.proj.GetFunc(rh.curPkg, handleFn)
	for f, fnd := range ffs {
		fh := newHandle(rh.proj, f, fnd, nil, cmt)
		rh.Handles = append(rh.Handles, fh)
	}
}

type route struct {
	proj        *proj.Proj
	engineVar   string
	initExpr    string
	curPkg      string
	specifyFunc string
	Vars        map[string]string
	RouteMap    map[string]string
	Handles     []*handle
}

func newRoute(proj *proj.Proj, initExpr, specifyFunc string) *route {
	// split g := gin.New() or g := gin.Default()
	_, ginFunc := splitDot(initExpr)
	return &route{
		proj:        proj,
		initExpr:    ginFunc,
		specifyFunc: specifyFunc,
		Vars:        make(map[string]string),
		RouteMap:    map[string]string{},
	}
}

func (rh *route) Parse(f *file.File, fnd *dst.FuncDecl) {
	rh.curPkg = f.Pkg()
	// global vars
	for k, v := range rh.proj.GetGlobalVar(f.Pkg()) {
		rh.Vars[k] = v
	}
	// vars in function param
	for k, v := range common.GetFuncParams(fnd) {
		rh.Vars[k] = v
	}

	parseStmtList(fnd.Body.List, rh.Vars, rh.parseItem)
}

func (rh *route) Copy(vars map[string]string) *route {
	rm := copyMap(rh.RouteMap)
	for k, v := range vars {
		if vv, ok := rh.RouteMap[v]; ok {
			rm[k] = vv
			delete(rm, v)
		}
	}

	nrh := &route{
		proj:        rh.proj,
		initExpr:    rh.initExpr,
		specifyFunc: rh.specifyFunc,
		Vars:        make(map[string]string),
		RouteMap:    rm,
	}
	return nrh
}

func (rh *route) parseItem(stmt interface{}, vars map[string]string) {
	vs := rh.proj.GetVarsFromStmt(stmt, rh.curPkg, vars)
	for v, t := range vs {
		_, sel := splitDot(t)

		// g := gin.New()
		if sel == rh.initExpr {
			rh.RouteMap[v] = ""
			rh.engineVar = v
			continue
		}

		// g.GET g.POST g.Group
		if fn, ok := routeParsers[sel]; ok {
			call, err := common.GetCallExprByVarName(stmt, v)
			if err != nil {
				continue
			}
			cal, _ := splitDot(t)
			fn(rh, v, cal, call)
			continue
		}

		if v != "_" {
			vars[v] = t
			continue
		}

		// recursive
		call, err := common.GetCallExprByVarName(stmt, v)
		if err != nil {
			continue
		}

		ps, ok := common.CheckCallExprParam(call, rh.engineVar)
		if !ok {
			continue
		}

		ffs := rh.proj.GetFunc(rh.curPkg, t)
		if len(ffs) == 0 {
			continue
		}

		for f, fnd := range ffs {
			nvs := make(map[string]string)
			if len(ps) > 0 {
				fps := common.GetFuncParamList(fnd)
				if len(fps) != len(ps) {
					continue
				}
				for i, s := range fps {
					nvs[s] = ps[i]
				}
			}

			nrh := rh.Copy(nvs)
			nrh.Parse(f, fnd)
			rh.Handles = append(rh.Handles, nrh.Handles...)
		}
	}
}

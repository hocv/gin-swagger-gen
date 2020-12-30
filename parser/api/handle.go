package api

import (
	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/ast"
)

type recursive func(a *ast.Ast, decl *dst.FuncDecl, vs map[string]string)

type Handle interface {
	Asts() *ast.Asts
	Type() string
	Cond(sel string) bool
	Parser(val string, vat string, call *dst.CallExpr, vs map[string]string)
	Recursive(a *ast.Ast, decl *dst.FuncDecl, vs map[string]string)
}

func parseStmtList(stmts []dst.Stmt, vars map[string]string, hdl Handle) {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *dst.IfStmt:
			ifStmt := stmt.(*dst.IfStmt)
			local := copyMap(vars)
			parseStmtItem(ifStmt.Init, local, hdl)
			parseStmtList(ifStmt.Body.List, local, hdl)
		case *dst.BlockStmt:
			local := copyMap(vars)
			parseStmtList(stmt.(*dst.BlockStmt).List, local, hdl)
		default:
			parseStmtItem(stmt, vars, hdl)
		}
	}
}

func parseStmtItem(stmt interface{}, vars map[string]string, hdl Handle) {
	vs := ast.GetVars(stmt)
	for v, t := range vs {
		_, sel := splitDot(t)

		if hdl.Cond(sel) {
			call, err := ast.GetCallExprByVarName(stmt, v)
			if err != nil {
				continue
			}
			hdl.Parser(v, t, call, vars)
			continue
		}

		if v != "_" {
			vars[v] = t
			continue
		}
		hSel, hName := splitDot(sel)
		v, ok := vars[hSel]
		if !ok {
			hSel = v
		}
		call, err := ast.GetCallExprByVarName(stmt, v)
		if err != nil {
			continue
		}
		ps := make([]string, 0)
		for _, arg := range call.Args {
			ps = append(ps, ast.ToStr(arg))
		}
		searchGinFunc(hdl.Asts(), hdl.Type(), hSel, hName, ps, hdl.Recursive)
	}
}

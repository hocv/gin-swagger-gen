package ast

import (
	"fmt"
	"strings"

	"github.com/dave/dst"
)

// CheckSelectorExpr check stmt has SelectorExpr
func CheckSelectorExpr(stmt interface{}, expr string) bool {
	switch stmt.(type) {
	case *dst.SelectorExpr:
		sel := stmt.(*dst.SelectorExpr)
		ident, ok := sel.X.(*dst.Ident)
		if !ok {
			return false
		}
		exp := fmt.Sprintf("%s.%s", ident.Name, sel.Sel.Name)
		return exp == expr
	case *dst.ExprStmt:
		return CheckSelectorExpr(stmt.(*dst.ExprStmt).X, expr)
	case *dst.CallExpr:
		return CheckSelectorExpr(stmt.(*dst.CallExpr).Fun, expr)
	case *dst.AssignStmt:
		for _, rh := range stmt.(*dst.AssignStmt).Rhs {
			if CheckSelectorExpr(rh, expr) {
				return true
			}
		}
		return false
	case *dst.DeclStmt:
		return CheckSelectorExpr(stmt.(*dst.DeclStmt).Decl, expr)
	case *dst.GenDecl:
		specs := stmt.(*dst.GenDecl).Specs
		for _, spec := range specs {
			if CheckSelectorExpr(spec, expr) {
				return true
			}
		}
		return false
	case *dst.ValueSpec:
		for _, value := range stmt.(*dst.ValueSpec).Values {
			if CheckSelectorExpr(value, expr) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func GetFuncParamByType(decl *dst.FuncDecl, argType string) []string {
	var names []string
	for _, field := range decl.Type.Params.List {
		t := filedType(field.Type)
		if t != argType {
			continue
		}
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
	}
	return names
}

func GetFuncParams(decl *dst.FuncDecl) map[string]string {
	vt := make(map[string]string)
	for _, field := range decl.Type.Params.List {
		t := filedType(field.Type)
		if len(t) == 0 {
			continue
		}
		for _, name := range field.Names {
			vt[name.Name] = t
		}
	}
	return vt
}

func GetFuncParamList(decl *dst.FuncDecl) (ps []string) {
	for _, field := range decl.Type.Params.List {
		t := filedType(field.Type)
		if len(t) == 0 {
			continue
		}
		for _, name := range field.Names {
			ps = append(ps, name.Name)
		}
	}
	return
}

func filedType(filed interface{}) string {
	switch filed.(type) {
	case *dst.SelectorExpr:
		sel := filed.(*dst.SelectorExpr)
		return fmt.Sprintf("%s.%s", sel.X.(*dst.Ident).Name, sel.Sel.Name)
	case *dst.StarExpr:
		return fmt.Sprintf("*%s", filedType(filed.(*dst.StarExpr).X))
	case *dst.Ident:
		return filed.(*dst.Ident).Name
	default:
		return ""
	}
}

func GetFuncVars(decl *dst.FuncDecl) map[string]string {
	vars := make(map[string]string)
	for _, stmt := range decl.Body.List {
		vs := GetVars(stmt)
		for k, v := range vs {
			vars[k] = v
		}
	}
	return vars
}

func GetVars(stmt interface{}) map[string]string {
	vars := make(map[string]string)
	switch stmt.(type) {
	case *dst.ExprStmt:
		vars["_"] = ExprToStr(stmt.(*dst.ExprStmt).X)
	case *dst.DeclStmt:
		genDecl, ok := stmt.(*dst.DeclStmt).Decl.(*dst.GenDecl)
		if !ok {
			return vars
		}
		for k, v := range GetVars(genDecl) {
			vars[k] = v
		}
	case *dst.GenDecl:
		genDecl, ok := stmt.(*dst.GenDecl)
		if !ok {
			return vars
		}
		for _, spec := range genDecl.Specs {
			vp, ok := spec.(*dst.ValueSpec)
			if !ok {
				continue
			}
			vpType := ExprToStr(vp.Type)
			for _, name := range vp.Names {
				vars[name.Name] = vpType
			}
			for _, value := range vp.Values {
				valueStr := ExprToStr(value)
				for _, vpName := range vp.Names {
					vars[vpName.Name] = valueStr
				}
			}
		}
	case *dst.AssignStmt:
		assign := stmt.(*dst.AssignStmt)
		for idx, rh := range assign.Rhs {
			lhName := assign.Lhs[idx].(*dst.Ident).Name
			vars[lhName] = ExprToStr(rh)
		}
	}
	return vars
}

func ExprToStr(stmt interface{}) string {
	switch stmt.(type) {
	case *dst.UnaryExpr:
		return ExprToStr(stmt.(*dst.UnaryExpr).X)
	case *dst.CompositeLit:
		return ExprToStr(stmt.(*dst.CompositeLit).Type)
	case *dst.CallExpr:
		return ExprToStr(stmt.(*dst.CallExpr).Fun)
	case *dst.SelectorExpr:
		sel := stmt.(*dst.SelectorExpr)
		return fmt.Sprintf("%s.%s", sel.X.(*dst.Ident).Name, sel.Sel.Name)
	case *dst.BasicLit:
		return stmt.(*dst.BasicLit).Kind.String()
	case *dst.Ident:
		return stmt.(*dst.Ident).Name
	case *dst.StarExpr:
		return fmt.Sprintf("*%s", ExprToStr(stmt.(*dst.StarExpr).X))
	}
	return ""
}

func CallExprByVarName(stmt interface{}, varName string) (*dst.CallExpr, error) {
	switch stmt.(type) {
	case *dst.ExprStmt:
		call, ok := stmt.(*dst.ExprStmt).X.(*dst.CallExpr)
		if !ok {
			return nil, ErrNotFind
		}
		return call, nil
	case *dst.AssignStmt:
		assign := stmt.(*dst.AssignStmt)
		for idx, rh := range assign.Rhs {
			varIdent, ok := assign.Lhs[idx].(*dst.Ident)
			if !ok || varIdent.Name != varName {
				continue
			}
			call, ok := rh.(*dst.CallExpr)
			if !ok {
				continue
			}
			return call, nil
		}
	case *dst.DeclStmt:
		genDecl, ok := stmt.(*dst.DeclStmt).Decl.(*dst.GenDecl)
		if !ok {
			return nil, ErrNotFind
		}
		for _, spec := range genDecl.Specs {
			vp, ok := spec.(*dst.ValueSpec)
			if !ok {
				continue
			}
			for idx, value := range vp.Values {
				if vp.Names[idx].Name != varName {
					continue
				}
				call, ok := value.(*dst.CallExpr)
				if !ok {
					continue
				}
				return call, nil
			}
		}
	}
	return nil, ErrNotFind
}

func BasicLitValue(basic interface{}) string {
	b, ok := basic.(*dst.BasicLit)
	if !ok {
		return ""
	}
	return strings.Trim(b.Value, "\"")
}

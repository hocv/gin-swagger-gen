package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dave/dst"
)

var ErrNotFind = errors.New("not found")

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

// GetFuncParamByType get param name of function by type
// func Test(a, b string) : string is type , a and b is name
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

// GetFuncParams get params of function,
// return map[name]type
// func Test(a, b string) -> {"a":"string","b":"string")
func GetFuncParams(decl *dst.FuncDecl) map[string]string {
	if decl == nil {
		return nil
	}
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

// GetFuncParamList get param name of function
// func Test(a, b string) -> ["a","b"]
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

// GetVars get vars from statement
// return map[name]type
func GetVars(stmt interface{}) map[string]string {
	vars := make(map[string]string)
	switch stmt.(type) {
	case *dst.ExprStmt:
		vars["_"] = ToStr(stmt.(*dst.ExprStmt).X)
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
			for k, v := range GetVars(spec) {
				vars[k] = v
			}
		}
	case *dst.ValueSpec:
		vs := stmt.(*dst.ValueSpec)
		vpType := ToStr(vs.Type)
		for _, name := range vs.Names {
			vars[name.Name] = vpType
		}
		for _, value := range vs.Values {
			valueStr := ToStr(value)
			for _, vpName := range vs.Names {
				vars[vpName.Name] = valueStr
			}
		}
	case *dst.AssignStmt:
		assign := stmt.(*dst.AssignStmt)
		for idx, rh := range assign.Rhs {
			ident, ok := assign.Lhs[idx].(*dst.Ident)
			if !ok {
				continue
			}
			lhName := ident.Name
			vars[lhName] = ToStr(rh)
		}
	}
	return vars
}

// ToStr convert to string
func ToStr(stmt interface{}) string {
	switch stmt.(type) {
	case *dst.UnaryExpr:
		return ToStr(stmt.(*dst.UnaryExpr).X)
	case *dst.CompositeLit:
		return ToStr(stmt.(*dst.CompositeLit).Type)
	case *dst.CallExpr:
		return ToStr(stmt.(*dst.CallExpr).Fun)
	case *dst.SelectorExpr:
		sel := stmt.(*dst.SelectorExpr)
		return fmt.Sprintf("%s.%s", ToStr(sel.X), sel.Sel.Name)
	case *dst.BasicLit:
		return strings.ToLower(stmt.(*dst.BasicLit).Kind.String())
	case *dst.Ident:
		return stmt.(*dst.Ident).Name
	case *dst.StarExpr:
		return fmt.Sprintf("*%s", ToStr(stmt.(*dst.StarExpr).X))
	case *dst.ArrayType:
		return fmt.Sprintf("[]%s", ToStr(stmt.(*dst.ArrayType).Elt))
	case *dst.MapType:
		mt := stmt.(*dst.MapType)
		return fmt.Sprintf("map[%s]%s", mt.Key.(*dst.Ident).String(), mt.Value.(*dst.Ident).String())
	case string:
		return fmt.Sprintf("%v", stmt)
	}
	return ""
}

// GetCallExprByVarName get CallExpr from stmt by var name
func GetCallExprByVarName(stmt interface{}, varName string) (*dst.CallExpr, error) {
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

func CheckCallExprParam(call *dst.CallExpr, p string) ([]string, bool) {
	var (
		ps []string
		ok = false
	)

	for _, arg := range call.Args {
		str := ToStr(arg)
		if str == p {
			ok = true
		}
		ps = append(ps, str)
	}

	return ps, ok
}

func BasicLitValue(basic interface{}) string {
	b, ok := basic.(*dst.BasicLit)
	if !ok {
		return ""
	}
	return strings.Trim(b.Value, "\"")
}

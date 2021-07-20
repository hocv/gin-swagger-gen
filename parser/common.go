package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/dst"
)

func parseStmtList(stmts []dst.Stmt, vars map[string]string, fn func(stmt interface{}, vars map[string]string)) {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *dst.IfStmt:
			ifStmt := stmt.(*dst.IfStmt)
			local := copyMap(vars)
			fn(ifStmt.Init, local)
			parseStmtList(ifStmt.Body.List, local, fn)
		case *dst.BlockStmt:
			local := copyMap(vars)
			parseStmtList(stmt.(*dst.BlockStmt).List, local, fn)
		default:
			fn(stmt, vars)
		}
	}
}

// splitDot split string with dot
func splitDot(str string) (string, string) {
	arr := strings.Split(str, ".")
	if len(arr) != 2 {
		return "", str
	}
	return arr[0], arr[1]
}

var routePathReg = regexp.MustCompile(":\\w+")

// fmtRoutePath remove "" and replace : to {}.
// e.g. "/user/:id" => /user/{id}
func fmtRoutePath(r string) string {
	r = strings.Trim(r, "\"")
	r = routePathReg.ReplaceAllStringFunc(r, func(s string) string {
		if len(s) < 1 {
			return s
		}
		return fmt.Sprintf("{%s}", s[1:])
	})
	return r
}

var routeParamReg = regexp.MustCompile("{\\w+}")

// routePathParams params in path. "/user/{id}"
func routePathParams(r string) (params []string) {
	_ = routeParamReg.ReplaceAllStringFunc(r, func(s string) string {
		if len(s) < 2 {
			return s
		}
		params = append(params, s[1:len(s)-1])
		return s
	})
	return
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

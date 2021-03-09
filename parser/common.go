package parser

import (
	"fmt"
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

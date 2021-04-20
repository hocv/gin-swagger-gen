package proj

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/lib/common"
	"github.com/hocv/gin-swagger-gen/lib/file"
	"github.com/hocv/gin-swagger-gen/lib/pkg"
)

type Proj struct {
	mtx  sync.Mutex
	pkgs map[string]*pkg.Pkg
}

func New() *Proj {
	return &Proj{
		pkgs: map[string]*pkg.Pkg{},
	}
}

func (proj *Proj) ScanDir(dir string) {
	files := scanDir(dir)
	for _, f := range files {
		proj.AddFile(f)
	}
}

func (proj *Proj) AddFile(files ...*file.File) {
	proj.mtx.Lock()
	defer proj.mtx.Unlock()

	for _, f := range files {
		p, ok := proj.pkgs[f.Pkg()]
		if !ok {
			p = pkg.New(f.Pkg())
			proj.pkgs[f.Pkg()] = p
		}
		p.AddFile(f)
	}
}

func (proj *Proj) GetFunc(pkg, name string) map[*file.File]*dst.FuncDecl {
	p, ok := proj.pkgs[pkg]
	if !ok {
		return nil
	}
	return p.GetFunc(name)
}

func (proj *Proj) GetGlobalVar(pkg string) map[string]string {
	p, ok := proj.pkgs[pkg]
	if !ok {
		return nil
	}
	return p.GetGlobalVar()
}

func (proj *Proj) GetPkgWithImported(path string) (pkgs []*pkg.Pkg) {
	for _, p := range proj.pkgs {
		if _, err := p.GetImported(path); err == nil {
			pkgs = append(pkgs, p)
		}
	}
	return
}

func (proj *Proj) GetStruct(pkg, name string) (*dst.StructType, error) {
	arr := strings.Split(name, ".")
	if len(arr) == 2 {
		pkg = arr[0]
	}

	p, ok := proj.pkgs[pkg]
	if !ok {
		return nil, common.ErrNotFind
	}

	return p.GetStruct(name)
}

func (proj *Proj) GetVarsFromStmt(stmt interface{}, curPkg string, outVars map[string]string) map[string]string {
	selFn := func(sel *dst.SelectorExpr, v string) {
		selX, ok := sel.X.(*dst.Ident)
		if !ok {
			return
		}
		ovt, ok := outVars[selX.Name]
		if !ok {
			return
		}

		if !strings.Contains(ovt, "{") {
			ifs := proj.interfaceOfStmt(curPkg, ovt)
			outVars[selX.Name] = fmt.Sprintf("%s{%s=%s}", ovt, ifs[sel.Sel.Name], v)
			return
		}

		novt := ovt[:strings.Index(ovt, "{")]
		ifs := proj.interfaceOfStmt(curPkg, novt)
		k, ok := ifs[sel.Sel.Name]
		if !ok {
			return
		}
		outVars[selX.Name] = replaceValue(ovt, k, v)
	}

	vars := make(map[string]string)

	switch stmt.(type) {
	case *dst.BasicLit:
		vars["_"] = common.ToStr(stmt)
	case *dst.ExprStmt:
		vars["_"] = common.ToStr(stmt.(*dst.ExprStmt).X)
	case *dst.DeclStmt:
		genDecl, ok := stmt.(*dst.DeclStmt).Decl.(*dst.GenDecl)
		if !ok {
			return vars
		}
		for k, v := range common.GetVars(genDecl) {
			vars[k] = v
		}
	case *dst.GenDecl:
		genDecl := stmt.(*dst.GenDecl)
		for _, spec := range genDecl.Specs {
			for k, v := range common.GetVars(spec) {
				vars[k] = v
			}
		}
	case *dst.ValueSpec:
		vs := stmt.(*dst.ValueSpec)
		vpType := common.ToStr(vs.Type)
		for _, name := range vs.Names {
			vars[name.Name] = vpType
		}
		for _, value := range vs.Values {
			valueStr := common.ToStr(value)
			for _, vpName := range vs.Names {
				vars[vpName.Name] = valueStr
			}
		}
	case *dst.CompositeLit:
		vars["_"] = proj.interfaceOfCompositeLit(curPkg, stmt, outVars)
	case *dst.AssignStmt:
		assign := stmt.(*dst.AssignStmt)
		total := 0
		for _, rh := range assign.Rhs {
			switch rh.(type) {
			case *dst.TypeAssertExpr:
				switch assign.Lhs[total].(type) {
				case *dst.Ident:
					lhName := assign.Lhs[total].(*dst.Ident).Name
					vars[lhName] = proj.interfaceOfCompositeLit(curPkg, rh, outVars)
				case *dst.SelectorExpr:
					v := proj.interfaceOfCompositeLit(curPkg, rh, outVars)
					selFn(assign.Lhs[total].(*dst.SelectorExpr), v)
				}
				total++
				if len(assign.Lhs) > total {
					lhName := assign.Lhs[total].(*dst.Ident).Name
					vars[lhName] = "bool"
					total++
				}
			case *dst.Ident, *dst.BasicLit, *dst.CompositeLit:
				switch assign.Lhs[total].(type) {
				case *dst.Ident:
					lhName := assign.Lhs[total].(*dst.Ident).Name
					vars[lhName] = proj.interfaceOfCompositeLit(curPkg, rh, outVars)
				case *dst.SelectorExpr:
					v := proj.interfaceOfCompositeLit(curPkg, rh, outVars)
					selFn(assign.Lhs[total].(*dst.SelectorExpr), v)
				}
				total++
			case *dst.CallExpr:
				vs := proj.getVarFromCallExprResult(curPkg, rh.(*dst.CallExpr), outVars)
				if len(vs) == 0 {
					vars[common.ToStr(assign.Lhs[total])] = common.ToStr(rh)
					total++
					continue
				}
				for i, v := range vs {
					switch assign.Lhs[total+i].(type) {
					case *dst.Ident:
						vars[assign.Lhs[total+i].(*dst.Ident).Name] = v
					case *dst.SelectorExpr:
						selFn(assign.Lhs[total+i].(*dst.SelectorExpr), v)
					}
				}
				total += len(vs)
			}
		}
	}
	return vars
}

func (proj *Proj) isPkg(name string) bool {
	_, ok := proj.pkgs[name]
	return ok
}

func (proj *Proj) isFunc(name string) bool {
	for _, p := range proj.pkgs {
		if len(p.GetFunc(name)) > 0 {
			return true
		}
	}
	return false
}

func (proj *Proj) isStruct(name string) bool {
	for _, p := range proj.pkgs {
		_, err := p.GetStruct(name)
		if err == nil {
			return true
		}
	}
	return false
}

func (proj *Proj) getFuncResult(pkgName, FnName, recvName string) []string {
	p, ok := proj.pkgs[pkgName]
	if !ok {
		return nil
	}
	return p.GetFuncResultByName(FnName, recvName)
}

func (proj *Proj) getStructFieldType(pkgName, structName, fieldName string) (string, bool) {
	p, ok := proj.pkgs[pkgName]
	if !ok {
		return "", false
	}
	return p.GetStructFieldType(structName, fieldName)
}

func (proj *Proj) getVarFromCallExprResult(curPkg string, stmt interface{}, outVars map[string]string) []string {
	fff := func(args ...string) []string {
		var pkgName, struName, funcName, fieldName string
		for _, arg := range args {
			if proj.isPkg(arg) {
				pkgName = arg
				continue
			}
			if proj.isStruct(arg) {
				struName = arg
				continue
			}
			if proj.isFunc(arg) {
				funcName = arg
				continue
			}
			fieldName = arg
		}
		if len(struName) > 0 {
			ft, ok := proj.getStructFieldType(pkgName, struName, fieldName)
			if ok {
				return []string{ft}
			}
		}
		if len(funcName) > 0 {
			return proj.getFuncResult(pkgName, funcName, struName)
		}
		return nil
	}

	switch stmt.(type) {
	case *dst.SelectorExpr:
		sel := stmt.(*dst.SelectorExpr)
		rs := proj.getVarFromCallExprResult(curPkg, sel.X, outVars)
		if len(rs) != 1 {
			return nil
		}
		var v1, v2 string
		if ovt, ok := outVars[rs[0]]; ok {
			v1, v2 = slitDot(curPkg, ovt)
		} else {
			v1, v2 = slitDot(curPkg, rs[0])
		}
		selStr := common.ToStr(sel.Sel)
		return fff(v1, v2, selStr)
	case *dst.CallExpr:
		call := stmt.(*dst.CallExpr)
		return proj.getVarFromCallExprResult(curPkg, call.Fun, outVars)
	case *dst.Ident:
		str := common.ToStr(stmt)
		if _, ok := outVars[str]; ok {
			return []string{str}
		}
		v1, v2 := slitDot(curPkg, str)
		return fff(v1, v2)
	}
	return nil
}

func (proj *Proj) interfaceOfCompositeLit(curPkg string, stmt interface{}, outVars map[string]string) string {
	value := common.ToStr(stmt)
	ifs := proj.interfaceOfStmt(curPkg, stmt)
	if len(ifs) == 0 {
		return value
	}
	switch stmt.(type) {
	case *dst.CompositeLit:
		clit := stmt.(*dst.CompositeLit)
		var arr []string
		for _, elt := range clit.Elts {
			switch elt.(type) {
			case *dst.KeyValueExpr:
				kve := elt.(*dst.KeyValueExpr)
				jsonTag, ok := ifs[kve.Key.(*dst.Ident).Name]
				if !ok {
					continue
				}
				switch kve.Value.(type) {
				case *dst.Ident:
					name := kve.Value.(*dst.Ident).Name
					if name == "nil" {
						continue
					}
					ovt := outVars[name]
					arr = append(arr, fmt.Sprintf("%s=%s", jsonTag, ovt))
				case *dst.CompositeLit:
					v := proj.interfaceOfCompositeLit(curPkg, kve.Value, outVars)
					arr = append(arr, fmt.Sprintf("%s=%s", jsonTag, v))
				}
			}
		}
		if len(arr) == 0 {
			return value
		}
		value = fmt.Sprintf("%s{%s}", value, strings.Join(arr, ","))
	}
	return value
}

func (proj *Proj) interfaceOfStmt(pkg string, stmt interface{}) map[string]string {
	name := common.ToStr(stmt)
	pkg, name = slitDot(pkg, name)

	fn, err := proj.GetStruct(pkg, name)
	if err != nil {
		return nil
	}

	return proj.getInterfaceOfStruct(pkg, fn)
}

func (proj *Proj) getInterfaceOfStruct(curPkg string, stru *dst.StructType) map[string]string {
	vars := make(map[string]string)

	recur := func(field *dst.Field, tag string, stru *dst.StructType) {
		vs := proj.getInterfaceOfStruct(curPkg, stru)
		for _, name := range field.Names {
			n := common.ToStr(name)
			for k, v := range vs {
				vars[fmt.Sprintf("%s.%s", n, k)] = fmt.Sprintf("%s.%s", tag, v)
			}
		}
	}

	for _, field := range stru.Fields.List {
		var tag string
		if field.Tag != nil {
			tag = getJsonTag(field.Tag.Value)
		} else {
			if len(field.Names) == 0 {
				continue
			}
			tag = snakeCase(field.Names[0].Name)
		}

		switch field.Type.(type) {
		case *dst.InterfaceType:
			for _, name := range field.Names {
				vars[common.ToStr(name)] = tag
			}
		case *dst.StructType:
			recur(field, tag, field.Type.(*dst.StructType))
		case *dst.Ident:
			name := field.Type.(*dst.Ident).Name
			fn, err := proj.GetStruct(curPkg, name)
			if err != nil {
				continue
			}
			recur(field, tag, fn)
		}
	}
	return vars
}

func (proj *Proj) Save() error {
	for _, p := range proj.pkgs {
		if err := p.Save(); err != nil {
			return err
		}
	}
	return nil
}

// scanDir scan go file and parse to ast
func scanDir(dir string) (asts []*file.File) {
	if strings.Contains(dir, "vendor") {
		return
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range fs {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() {
			asts = append(asts, scanDir(path)...)
		}
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		g, err := file.New(path)
		if err != nil {
			continue
		}
		asts = append(asts, g)
	}
	return
}

func slitDot(curPkg, str string) (string, string) {
	if arr := strings.Split(str, "."); len(arr) == 2 {
		return arr[0], arr[1]
	}
	return curPkg, str
}

func replaceValue(str, k, v string) string {
	startIdx := strings.Index(str, fmt.Sprintf("%s=", k))
	if startIdx < 0 {
		prefix := str[:strings.Index(str, "{")+1]
		newPre := fmt.Sprintf("%s%s=%s,", prefix, k, v)
		return strings.Replace(str, prefix, newPre, 1)
	}

	startIdx += len(k) + 1
	endIdx := startIdx
	tag := 0
	for ci, c := range str {
		if ci < startIdx {
			continue
		}
		if c == '{' {
			tag++
		} else if c == '}' {
			if tag == 0 {
				endIdx = ci
				break
			}
			tag--
		} else if c == ',' {
			if tag == 0 {
				endIdx = ci
				break
			}
		}
	}
	tv := str[startIdx-1 : endIdx]
	v = fmt.Sprintf("=%s", v)
	return strings.Replace(str, tv, v, 1)
}

func getJsonTag(str string) string {
	idx := strings.Index(str, "json:\"")
	if idx < 0 {
		return ""
	}
	str = str[idx+6:]
	idx = strings.Index(str, "\"")
	str = str[:idx]
	idx = strings.Index(str, ",")
	if idx < 0 {
		return str
	}
	return str[:idx]
}

func snakeCase(s string) string {
	s = strings.TrimSpace(s)
	buffer := make([]rune, 0, len(s)+3)

	delimiter := '_'

	isLower := func(ch rune) bool {
		return ch >= 'a' && ch <= 'z'
	}
	toLower := func(ch rune) rune {
		if ch >= 'A' && ch <= 'Z' {
			return ch + 32
		}
		return ch
	}
	isUpper := func(ch rune) bool {
		return ch >= 'A' && ch <= 'Z'
	}
	isDelimiter := func(ch rune) bool {
		return ch == '-' || ch == '_' || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
	}

	var prev rune
	var curr rune
	for _, next := range s {
		if isDelimiter(curr) {
			if !isDelimiter(prev) {
				buffer = append(buffer, delimiter)
			}
		} else if isUpper(curr) {
			if isLower(prev) || (isUpper(prev) && isLower(next)) {
				buffer = append(buffer, delimiter)
			}
			buffer = append(buffer, toLower(curr))
		} else if curr != 0 {
			buffer = append(buffer, toLower(curr))
		}
		prev = curr
		curr = next
	}

	if len(s) > 0 {
		if isUpper(curr) && isLower(prev) && prev != 0 {
			buffer = append(buffer, delimiter)
		}
		buffer = append(buffer, toLower(curr))
	}

	return string(buffer)
}

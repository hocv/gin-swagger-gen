package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hocv/gin-swagger-gen/parser/comment"

	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/lib/common"
	"github.com/hocv/gin-swagger-gen/lib/file"
	"github.com/hocv/gin-swagger-gen/lib/proj"
)

type handle struct {
	proj        *proj.Proj
	dstFile     *file.File
	curPkg      string
	DstDecl     *dst.FuncDecl
	SrcDecl     *dst.FuncDecl
	Cmt         *comment.Comment
	Vars        map[string]string
	queryParams map[string]string
}

func newHandle(proj *proj.Proj, f *file.File, dstDecl *dst.FuncDecl, decl *dst.FuncDecl, cmt *comment.Comment) *handle {
	if decl == nil {
		decl = dstDecl
	}
	vars := make(map[string]string) // key: var ,value: var type or method

	// global vars
	for k, v := range f.GlobalVars() {
		vars[k] = v
	}
	// vars in function param
	for k, v := range common.GetFuncParams(decl) {
		vars[k] = v
	}

	return &handle{
		proj:        proj,
		dstFile:     f,
		curPkg:      f.Pkg(),
		DstDecl:     dstDecl,
		SrcDecl:     decl,
		Cmt:         cmt,
		Vars:        vars,
		queryParams: map[string]string{},
	}
}

func (hdl *handle) Parse() {
	parseStmtList(hdl.SrcDecl.Body.List, hdl.Vars, hdl.parseIterm)
}

func (hdl *handle) Merge() {
	if hdl.Cmt.Merge(hdl.DstDecl) {
		hdl.DstDecl.Decs.Start.Clear()
		hdl.DstDecl.Decs.Start.Append(hdl.Cmt.Decs()...)
		hdl.dstFile.Dirty()
	}
}

func (hdl *handle) Print() {
	hdl.Cmt.Merge(hdl.DstDecl)
	for _, s := range hdl.Cmt.Decs() {
		fmt.Println(s)
	}
}

func (hdl *handle) parseIterm(stmt interface{}, vars map[string]string) {
	fn := func(parser handleParser, v string) {
		call, err := common.GetCallExprByVarName(stmt, v)
		if err != nil {
			return
		}
		parser(hdl, vars, v, call)
	}

	vs := hdl.proj.GetVarsFromStmt(stmt, hdl.curPkg, vars)
	for v, t := range vs {
		_, sel := splitDot(t)

		parser, ok := handleParsers[sel]
		if !ok {
			parser, ok = handleParsers[t]
		}

		if ok {
			fn(parser, v)
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

		ginImport, ok := hdl.dstFile.DefaultImport(ginPkg, "gin")
		if !ok {
			continue
		}
		ginCtx := fmt.Sprintf("*%s.Context", ginImport)
		ctx := "c"
		for k, v := range vars {
			if v == ginCtx {
				ctx = k
				break
			}
		}

		ps, ok := common.CheckCallExprParam(call, ctx)
		if !ok {
			continue
		}

		ffs := hdl.proj.GetFunc(hdl.curPkg, t)
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
				fh := newHandle(hdl.proj, f, hdl.DstDecl, fnd, hdl.Cmt)
				for nk, nv := range nvs {
					if ov, ok := vars[nv]; ok {
						fh.Vars[nk] = ov
					}
				}
				fh.Parse()
			}
		}
	}
}

type handleParser func(hdl *handle, vars map[string]string, val string, call *dst.CallExpr)

var handleParsers = map[string]handleParser{
	"strconv.Atoi":       parseStrConv("integer"),
	"strconv.ParseInt":   parseStrConv("integer"),
	"strconv.ParseUint":  parseStrConv("integer"),
	"strconv.ParseFloat": parseStrConv("number"),
	"strconv.ParseBool":  parseStrConv("boolean"),
	"BindJSON":           parseBind("json"),
	"ShouldBindJSON":     parseBind("json"),
	"BindXML":            parseBind("xml"),
	"ShouldBindXML":      parseBind("xml"),
	"BindYAML":           parseBind("yaml"),
	"ShouldBindYAML":     parseBind("yaml"),
	"Query":              parseQuery(""),
	"BindQuery":          parseQuery("BindQuery"),
	"ShouldBindQuery":    parseQuery("ShouldBindQuery"),
	"DefaultQuery":       parseQuery("DefaultQuery"),
	"GetQuery":           parseQuery(""),
	"QueryArray":         parseQuery(""),
	"GetQueryArray":      parseQuery(""),
	"QueryMap":           parseQuery(""),
	"GetQueryMap":        parseQuery(""),
	"PostForm":           parseForm(""),
	"DefaultPostForm":    parseForm("DefaultPostForm"),
	"GetPostForm":        parseForm(""),
	"PostFormArray":      parseForm(""),
	"GetPostFormArray":   parseForm(""),
	"PostFormMap":        parseForm(""),
	"GetPostFormMap":     parseForm(""),
	"FormFile":           parseForm(""),
	"HTML":               parseProduce("html"),
	"IndentedJSON":       parseProduce("json"),
	"SecureJSON":         parseProduce("json"),
	"JSONP":              parseProduce("js"),
	"JSON":               parseProduce("json"),
	"AsciiJSON":          parseProduce("json"),
	"PureJSON":           parseProduce("xml"),
	"XML":                parseProduce("xml"),
	"YAML":               parseProduce("yaml"),
	"ProtoBuf":           parseProduce("protobuf"),
	"String":             parseProduce("string"),
}

var stateCode = map[string]int{
	"http.StatusContinue":                      100,
	"http.StatusSwitchingProtocols":            101,
	"http.StatusProcessing":                    102,
	"http.StatusEarlyHints":                    103,
	"http.StatusOK":                            200,
	"http.StatusCreated":                       201,
	"http.StatusAccepted":                      202,
	"http.StatusNonAuthoritativeInfo":          203,
	"http.StatusNoContent":                     204,
	"http.StatusResetContent":                  205,
	"http.StatusPartialContent":                206,
	"http.StatusMultiStatus":                   207,
	"http.StatusAlreadyReported":               208,
	"http.StatusIMUsed":                        226,
	"http.StatusMultipleChoices":               300,
	"http.StatusMovedPermanently":              301,
	"http.StatusFound":                         302,
	"http.StatusSeeOther":                      303,
	"http.StatusNotModified":                   304,
	"http.StatusUseProxy":                      305,
	"http.StatusTemporaryRedirect":             307,
	"http.StatusPermanentRedirect":             308,
	"http.StatusBadRequest":                    400,
	"http.StatusUnauthorized":                  401,
	"http.StatusPaymentRequired":               402,
	"http.StatusForbidden":                     403,
	"http.StatusNotFound":                      404,
	"http.StatusMethodNotAllowed":              405,
	"http.StatusNotAcceptable":                 406,
	"http.StatusProxyAuthRequired":             407,
	"http.StatusRequestTimeout":                408,
	"http.StatusConflict":                      409,
	"http.StatusGone":                          410,
	"http.StatusLengthRequired":                411,
	"http.StatusPreconditionFailed":            412,
	"http.StatusRequestEntityTooLarge":         413,
	"http.StatusRequestURITooLong":             414,
	"http.StatusUnsupportedMediaType":          415,
	"http.StatusRequestedRangeNotSatisfiable":  416,
	"http.StatusExpectationFailed":             417,
	"http.StatusTeapot":                        418,
	"http.StatusMisdirectedRequest":            421,
	"http.StatusUnprocessableEntity":           422,
	"http.StatusLocked":                        423,
	"http.StatusFailedDependency":              424,
	"http.StatusTooEarly":                      425,
	"http.StatusUpgradeRequired":               426,
	"http.StatusPreconditionRequired":          428,
	"http.StatusTooManyRequests":               429,
	"http.StatusRequestHeaderFieldsTooLarge":   431,
	"http.StatusUnavailableForLegalReasons":    451,
	"http.StatusInternalServerError":           500,
	"http.StatusNotImplemented":                501,
	"http.StatusBadGateway":                    502,
	"http.StatusServiceUnavailable":            503,
	"http.StatusGatewayTimeout":                504,
	"http.StatusHTTPVersionNotSupported":       505,
	"http.StatusVariantAlsoNegotiates":         506,
	"http.StatusInsufficientStorage":           507,
	"http.StatusLoopDetected":                  508,
	"http.StatusNotExtended":                   510,
	"http.StatusNetworkAuthenticationRequired": 511,
}

func parseStrConv(dstType string) handleParser {
	return func(hdl *handle, vars map[string]string, val string, call *dst.CallExpr) {
		p := common.ToStr(call.Args[0])
		name, ok := hdl.queryParams[p]
		if !ok {
			return
		}
		hdl.Cmt.SetParamRefType(name, dstType)
	}
}

func parseBind(bindType string) handleParser {
	return func(hdl *handle, vars map[string]string, val string, call *dst.CallExpr) {
		if len(call.Args) == 0 {
			return
		}

		name := common.ToStr(call.Args[0])
		refType, ok := vars[name]
		if !ok {
			return
		}

		param := comment.NewBodyParam(name, refType, "")
		hdl.Cmt.AddParam(param)
		hdl.Cmt.AddAccept(bindType)
	}
}

func parseQuery(queryType string) handleParser {
	return func(hdl *handle, vars map[string]string, val string, call *dst.CallExpr) {
		if len(call.Args) == 0 {
			return
		}

		name, refType, desc, ok := "", "string", "", false
		if strings.Contains(queryType, "Bind") {
			name = common.ToStr(call.Args[0])
			refType, ok = vars[name]
			if !ok {
				return
			}
			pkg, struName := splitDot(refType)
			if len(pkg) == 0 {
				pkg = hdl.curPkg
			}
			stru, err := hdl.proj.GetStruct(pkg, struName)
			if err != nil {
				return
			}
			for _, field := range stru.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				ft := common.ToStr(field.Type)
				tag := ""
				required := false
				if field.Tag != nil {
					tag = common.GetFormTag(field.Tag.Value)
					required = common.GetTagBindingRequired(field.Tag.Value)
				}
				for _, ident := range field.Names {
					fieldName := tag
					if len(name) == 0 {
						fieldName = common.SnakeCase(ident.Name)
					}
					param := comment.NewQueryParam(fieldName, ft, "")
					param.Required = required
					hdl.Cmt.AddParam(param)
				}
			}
			return
		} else {
			name = common.BasicLitValue(call.Args[0])
			if strings.Contains(queryType, "Default") && len(call.Args) > 1 {
				desc = fmt.Sprintf("default %s", common.BasicLitValue(call.Args[len(call.Args)-1]))
			}
			vars[val] = "string"
			hdl.queryParams[val] = name
		}

		param := comment.NewQueryParam(name, refType, desc)
		hdl.Cmt.AddParam(param)
	}
}

func parseForm(formType string) handleParser {
	return func(hdl *handle, vars map[string]string, val string, call *dst.CallExpr) {
		name, ref, desc := common.BasicLitValue(call.Args[0]), "string", ""
		if formType == "FormFile" {
			ref = "file"
		}

		if strings.Contains(formType, "Default") && len(call.Args) > 1 {
			desc = fmt.Sprintf("default %s", common.BasicLitValue(call.Args[len(call.Args)-1]))
		}
		param := comment.NewFormDataParam(name, ref, desc)
		hdl.Cmt.AddParam(param)
		hdl.Cmt.AddAccept("multipart/form-data")
	}
}

func parseProduce(produceType string) handleParser {
	return func(hdl *handle, vars map[string]string, val string, call *dst.CallExpr) {
		if len(call.Args) < 2 {
			return
		}
		hdl.Cmt.AddProduce(produceType)

		state := common.ToStr(call.Args[0])
		code, ok := stateCode[state]
		if !ok {
			bl, ok := call.Args[0].(*dst.BasicLit)
			if !ok {
				return
			}
			str := common.BasicLitValue(bl)
			v, err := strconv.Atoi(str)
			if err != nil {
				return
			}
			code = v
		}

		r := comment.Resp{
			Code: code,
			Type: produceType,
		}

		lastArg := common.ToStr(call.Args[1])
		if val, ok := vars[lastArg]; ok {
			r.Type = val
		} else {
			nv := hdl.proj.GetVarsFromStmt(call.Args[1], hdl.curPkg, vars)
			if t, ok := nv["_"]; ok {
				r.Type = t
			}
		}

		hdl.Cmt.AddResp(r)
	}
}

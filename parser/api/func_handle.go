package api

import (
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/hocv/gin-swagger-gen/ast"
)

type funcHandle struct {
	api     *Api
	DstDecl *dst.FuncDecl
	SrcDecl *dst.FuncDecl
	Cmt     *comment
	Vars    map[string]string
}

func parseFuncHandle(api *Api, dstAst *ast.Ast, dstDecl *dst.FuncDecl, decl *dst.FuncDecl, cmt *comment) {
	if dstDecl == nil {
		dstDecl = decl
	}
	vars := make(map[string]string) // key: var ,value: var type or method

	// global vars
	for k, v := range dstAst.GlobalVars() {
		vars[k] = v
	}

	// vars in function param
	for k, v := range ast.GetFuncParams(decl) {
		vars[k] = v
	}

	fh := &funcHandle{
		api:     api,
		DstDecl: dstDecl,
		SrcDecl: decl,
		Cmt:     cmt,
		Vars:    vars,
	}

	parseStmtList(decl.Body.List, fh.Vars, fh)
	dstAst.Dirty()
	fh.Cmt.AddToFunc(fh.DstDecl)
}

func (fh *funcHandle) Asts() *ast.Asts {
	return fh.api.asts
}

func (fh *funcHandle) Type() string {
	return "Context"
}

func (fh *funcHandle) Cond(sel string) bool {
	_, ok := ginParsers[sel]
	return ok
}

func (fh *funcHandle) Parser(val string, vat string, call *dst.CallExpr, vs map[string]string) {
	_, sel := splitDot(vat)
	ginParsers[sel](fh.Cmt, vs, call)
}

func (fh *funcHandle) Inter(a *ast.Ast, decl *dst.FuncDecl, vs map[string]string) {
	parseFuncHandle(fh.api, a, fh.DstDecl, decl, fh.Cmt)
}

type parser func(cmt *comment, vars map[string]string, call *dst.CallExpr)

var ginParsers = map[string]parser{
	"BindJSON":         parseBind("json"),
	"ShouldBindJSON":   parseBind("json"),
	"BindXML":          parseBind("xml"),
	"ShouldBindXML":    parseBind("xml"),
	"BindYAML":         parseBind("yaml"),
	"ShouldBindYAML":   parseBind("yaml"),
	"Query":            parseQuery(""),
	"DefaultQuery":     parseQuery("DefaultQuery"),
	"GetQuery":         parseQuery(""),
	"QueryArray":       parseQuery(""),
	"GetQueryArray":    parseQuery(""),
	"QueryMap":         parseQuery(""),
	"GetQueryMap":      parseQuery(""),
	"PostForm":         parseForm(""),
	"DefaultPostForm":  parseForm("DefaultPostForm"),
	"GetPostForm":      parseForm(""),
	"PostFormArray":    parseForm(""),
	"GetPostFormArray": parseForm(""),
	"PostFormMap":      parseForm(""),
	"GetPostFormMap":   parseForm(""),
	"FormFile":         parseForm(""),
	"HTML":             parseProduce("html"),
	"IndentedJSON":     parseProduce("json"),
	"SecureJSON":       parseProduce("json"),
	"JSONP":            parseProduce("js"),
	"JSON":             parseProduce("json"),
	"AsciiJSON":        parseProduce("json"),
	"PureJSON":         parseProduce("xml"),
	"XML":              parseProduce("xml"),
	"YAML":             parseProduce("yaml"),
	"ProtoBuf":         parseProduce("protobuf"),
	"String":           parseProduce("string"),
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

func parseBind(bindType string) parser {
	return func(cmt *comment, vars map[string]string, call *dst.CallExpr) {
		if len(call.Args) == 0 {
			return
		}
		cmt.Accept = append(cmt.Accept, bindType)

		arg := ast.ToStr(call.Args[0])
		body, ok := vars[arg]
		if !ok {
			return
		}

		cmt.BodyParams = append(cmt.BodyParams, bodyParam{
			Var:  arg,
			Body: body,
		})
	}
}

func parseQuery(queryType string) parser {
	return func(cmt *comment, vars map[string]string, call *dst.CallExpr) {
		if len(call.Args) == 0 {
			return
		}
		qp := queryParam{}
		if strings.Contains(queryType, "Bind") {
			arg := ast.ToStr(call.Args[0])
			body, ok := vars[arg]
			if !ok {
				return
			}
			qp.Var = body
		} else {
			qp.Var = ast.BasicLitValue(call.Args[0])
			if strings.Contains(queryType, "Default") && len(call.Args) > 1 {
				qp.Default = ast.BasicLitValue(call.Args[len(call.Args)-1])
			}
		}

		cmt.QueryParams = append(cmt.QueryParams, qp)
	}
}

func parseForm(formType string) parser {
	return func(cmt *comment, vars map[string]string, call *dst.CallExpr) {
		cmt.Accept = append(cmt.Accept, "multipart/form-data")
		fp := formParam{
			Var:  ast.BasicLitValue(call.Args[0]),
			Type: "string",
		}
		if formType == "FormFile" {
			fp.Type = "file"
		}

		if strings.Contains(formType, "Default") && len(call.Args) > 1 {
			fp.Default = ast.BasicLitValue(call.Args[len(call.Args)-1])
		}
		cmt.FormParams = append(cmt.FormParams, fp)
	}
}

func parseProduce(produceType string) parser {
	return func(cmt *comment, vars map[string]string, call *dst.CallExpr) {
		if len(call.Args) < 2 {
			return
		}
		cmt.Produce = append(cmt.Produce, produceType)

		state := ast.ToStr(call.Args[0])
		code, ok := stateCode[state]
		if !ok {
			bl, ok := call.Args[0].(*dst.BasicLit)
			if !ok {
				return
			}
			str := ast.BasicLitValue(bl)
			v, err := strconv.Atoi(str)
			if err != nil {
				return
			}
			code = v
		}

		r := resp{
			Code: code,
			Type: produceType,
		}
		if produceType != "string" {
			lastArg := ast.ToStr(call.Args[1])
			val, ok := vars[lastArg]
			if !ok {
				return
			}
			r.Type = val
		}

		cmt.Resp = append(cmt.Resp, r)
	}
}

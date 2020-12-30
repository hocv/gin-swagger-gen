package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dave/dst"
)

type bodyParam struct {
	Var  string
	Body string
	Flag bool
}

func (bp bodyParam) Decs() string {
	return fmt.Sprintf("// @Param %s body %s true \"%s\"", bp.Var, bp.Body, bp.Var)
}

type queryParam struct {
	Var     string
	Default string
	Flag    bool
}

func (q queryParam) Decs() string {
	if len(q.Default) > 0 {
		return fmt.Sprintf("// @Param %s query string false \"%s default %s\"", q.Var, q.Var, q.Default)
	}
	return fmt.Sprintf("// @Param %s query string true \"%s\"", q.Var, q.Var)
}

type formParam struct {
	Var     string
	Type    string
	Default string
	Flag    bool
}

func (f formParam) Decs() string {
	if len(f.Default) > 0 {
		return fmt.Sprintf("// @Param %s formData %s false \"%s default %s\"", f.Var, f.Type, f.Var, f.Default)
	}
	return fmt.Sprintf("// @Param %s formData %s true \"%s\"", f.Var, f.Type, f.Var)
}

type resp struct {
	Code int
	Type string
}

func (r resp) Decs() string {
	v1 := "// @Failure"
	if r.Code == 200 {
		v1 = fmt.Sprintf("// @Success")
	}
	v2 := "object"
	if r.Type == "string" {
		v2 = "string"
	}
	return fmt.Sprintf("%s %d {%s} %s", v1, r.Code, v2, r.Type)
}

// comment
type comment struct {
	Summary     string
	Tags        string
	Description []string
	Accept      []string
	Produce     []string
	RoutePath   string       // route path
	RouteMethod string       // route method: get post put
	PathParams  []string     // params in path
	BodyParams  []bodyParam  // body params
	QueryParams []queryParam // query params
	FormParams  []formParam  // form params
	Resp        []resp       // resp
}

// Decs common
func (c *comment) Decs() []string {
	strs := []string{
		fmt.Sprintf("// @Summary %s", c.Summary),
	}
	if len(c.Tags) > 0 {
		strs = append(strs, fmt.Sprintf("// @Tags %s", c.Tags))
	}
	fmt.Println(len(c.Description))
	for _, desc := range c.Description {
		strs = append(strs, fmt.Sprintf("// @Description %s", desc))
	}
	if len(c.Accept) > 0 {
		strs = append(strs, trimAndJoin("Accept", c.Accept))
	}
	if len(c.Produce) > 0 {
		strs = append(strs, trimAndJoin("Produce", c.Produce))
	}
	for _, param := range c.PathParams {
		strs = append(strs, fmt.Sprintf("// @Param %s path string true \"%s\"", param, param))
	}
	for _, param := range c.BodyParams {
		strs = append(strs, param.Decs())
	}
	for _, param := range c.QueryParams {
		strs = append(strs, param.Decs())
	}
	for _, param := range c.FormParams {
		strs = append(strs, param.Decs())
	}
	codeMap := make(map[int]string)
	for _, r := range c.Resp {
		_, ok := codeMap[r.Code]
		if ok {
			continue
		}
		codeMap[r.Code] = r.Type
		strs = append(strs, r.Decs())
	}

	strs = append(strs, fmt.Sprintf("// @Router %s [%s]", c.RoutePath, c.RouteMethod))
	return strs
}

func (c *comment) AddToFunc(decl *dst.FuncDecl) {
	if decl == nil {
		return
	}

	for _, cmt := range decl.Decs.Start.All() {
		c.parseComment(cmt)
	}
	decl.Decs.Start.Clear()
	decl.Decs.Start.Append(c.Decs()...)
}

func (c *comment) parseComment(cmt string) {
	commentLine := strings.TrimSpace(strings.TrimLeft(cmt, "//"))
	if len(commentLine) == 0 {
		return
	}
	attribute := strings.Fields(commentLine)[0]
	remainder := strings.TrimSpace(commentLine[len(attribute):])
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "@summary":
		c.Summary = remainder
	case "@description":
		c.Description = append(c.Description, remainder)
	case "@accept":
		if len(remainder) > 0 {
			c.Accept = append(c.Accept, strings.Split(remainder, ",")...)
		}
	case "@produce":
		if len(remainder) > 0 {
			c.Produce = append(c.Produce, strings.Split(remainder, ",")...)
		}
	case "@tags":
		c.Tags = remainder
	default:
		if !strings.HasPrefix(commentLine, "@") {
			c.Description = append(c.Description, commentLine)
		}
	}
}

func trimAndJoin(attr string, arr []string) string {
	var trim []string
	sort.Strings(arr)
	for i, s := range arr {
		if i > 0 && arr[i] == arr[i-1] {
			continue
		}
		trim = append(trim, s)
	}
	return fmt.Sprintf("// @%s %s", attr, strings.Join(trim, ","))
}

package comment

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dave/dst"
)

// Comment
type Comment struct {
	summary     string
	tags        string
	id          string
	description []string
	accept      []string
	produce     []string
	route       Route
	params      Params
	resp        map[int]Resp
}

func New(summary, routeBase, routePath, method string) *Comment {
	return &Comment{
		summary: summary,
		tags:    strings.TrimPrefix(routeBase, "/"),
		route: Route{
			RoutePath:   routePath,
			RouteMethod: method,
		},
		resp: map[int]Resp{},
	}
}

// Decs common
func (c *Comment) Decs() []string {
	desc := []string{
		fmt.Sprintf("// @Summary %s", c.summary),
	}

	if len(c.id) > 0 {
		desc = append(desc, fmt.Sprintf("// @ID %s", c.id))
	}

	if len(c.tags) > 0 {
		desc = append(desc, fmt.Sprintf("// @Tags %s", c.tags))
	}

	for _, d := range c.description {
		desc = append(desc, fmt.Sprintf("// @Description %s", d))
	}

	if len(c.accept) > 0 {
		desc = append(desc, trimAndJoin("Accept", c.accept))
	}
	if len(c.produce) > 0 {
		desc = append(desc, trimAndJoin("Produce", c.produce))
	}

	sort.Sort(c.params)
	for _, p := range c.params {
		desc = append(desc, p.Decs())
	}

	for _, r := range c.resp {
		desc = append(desc, r.Decs())
	}

	desc = append(desc, c.route.Decs())
	return desc
}

func (c *Comment) Merge(decl *dst.FuncDecl) bool {
	if decl == nil {
		return false
	}

	old := decl.Decs.Start.All()
	if len(old) == 0 {
		return true
	}
	for _, cmt := range old {
		c.parseComment(cmt)
	}

	cur := c.Decs()
	dic := make(map[string]struct{})
	for _, s := range cur {
		dic[s] = struct{}{}
	}

	for _, s := range old {
		if _, ok := dic[s]; !ok {
			return true
		}
	}

	return false
}

func (c *Comment) SetParamRefType(name, refType string) {
	for i, p := range c.params {
		if p.Name != name {
			continue
		}
		c.params[i].RefType = refType
		break
	}
}

func (c *Comment) AddParam(param Param) {
	for _, p := range c.params {
		if p.paramType != param.paramType ||
			p.RefType != param.RefType ||
			p.Name != param.Name {
			continue
		}
		p.Description = param.Description
		return
	}
	c.params = append(c.params, param)
}

func (c *Comment) AddResp(resp Resp) {
	c.resp[resp.Code] = resp
}

func (c *Comment) AddProduce(prod string) {
	c.produce = append(c.produce, prod)
}

func (c *Comment) AddAccept(accept string) {
	c.accept = append(c.accept, accept)
}

func (c *Comment) parseComment(cmt string) {
	commentLine := strings.TrimSpace(strings.TrimLeft(cmt, "//"))
	if len(commentLine) == 0 {
		return
	}
	attribute := strings.Fields(commentLine)[0]
	remainder := strings.TrimSpace(commentLine[len(attribute):])
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "@summary":
		c.summary = remainder
	case "@description":
		c.description = append(c.description, remainder)
	case "@accept":
		if len(remainder) > 0 {
			c.accept = append(c.accept, strings.Split(remainder, ",")...)
		}
	case "@produce":
		if len(remainder) > 0 {
			c.produce = append(c.produce, strings.Split(remainder, ",")...)
		}
	case "@param":
		if p, err := parseParam(commentLine); err == nil {
			c.AddParam(p)
		}
	case "@tags":
		c.tags = remainder
	case "@id":
		c.id = remainder
	default:
		if !strings.HasPrefix(commentLine, "@") {
			c.description = append(c.description, commentLine)
		}
	}
}

var paramPattern = regexp.MustCompile(`(\S+)[\s]+([\w]+)[\s]+([\S.]+)[\s]+([\w]+)[\s]+"([^"]+)"`)

func parseParam(commentLine string) (Param, error) {
	matches := paramPattern.FindStringSubmatch(commentLine)
	if len(matches) < 6 {
		return Param{}, nil
	}

	name, paramType, refType, desc := matches[1], matches[2], matches[3], matches[5]

	switch paramType {
	case "path":
		return NewPathParam(name, refType, desc), nil
	case "query":
		return NewQueryParam(name, refType, desc), nil
	case "body":
		return NewBodyParam(name, refType, desc), nil
	case "formData":
		return NewFormDataParam(name, refType, desc), nil
	}
	return Param{}, errors.New("not supported type")
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

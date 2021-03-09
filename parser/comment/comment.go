package comment

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/dave/dst"
)

// Comment
type Comment struct {
	Summary     string
	Tags        string
	ID          string
	Description []string
	Accept      []string
	Produce     []string
	Route       Route
	params      Params
	Resp        []Resp
}

// Decs common
func (c *Comment) Decs() []string {
	desc := []string{
		fmt.Sprintf("// @Summary %s", c.Summary),
	}

	if len(c.ID) > 0 {
		desc = append(desc, fmt.Sprintf("// @ID %s", c.ID))
	}

	if len(c.Tags) > 0 {
		desc = append(desc, fmt.Sprintf("// @Tags %s", c.Tags))
	}

	for _, d := range c.Description {
		desc = append(desc, fmt.Sprintf("// @Description %s", d))
	}

	if len(c.Accept) > 0 {
		desc = append(desc, trimAndJoin("Accept", c.Accept))
	}
	if len(c.Produce) > 0 {
		desc = append(desc, trimAndJoin("Produce", c.Produce))
	}

	sort.Sort(c.params)
	for _, p := range c.params {
		desc = append(desc, p.Decs())
	}

	codeMap := make(map[int]string)
	for _, r := range c.Resp {
		_, ok := codeMap[r.Code]
		if ok {
			continue
		}
		codeMap[r.Code] = r.Type
		desc = append(desc, r.Decs())
	}

	desc = append(desc, c.Route.Decs())
	return desc
}

func (c *Comment) Merge(decl *dst.FuncDecl) bool {
	if decl == nil {
		return false
	}

	old := decl.Decs.Start.All()
	for _, cmt := range old {
		c.parseComment(cmt)
	}

	cur := c.Decs()
	dic := make(map[string]struct{})
	for _, s := range cur {
		dic[s] = struct{}{}
	}

	update := false
	for _, s := range old {
		if _, ok := dic[s]; !ok {
			update = true
			break
		}
	}
	if update {

	}
	return update
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
	case "@param":
		if p, err := parseParam(commentLine); err == nil {
			c.AddParam(p)
		}
	case "@tags":
		c.Tags = remainder
	case "@id":
		c.ID = remainder
	default:
		if !strings.HasPrefix(commentLine, "@") {
			c.Description = append(c.Description, commentLine)
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

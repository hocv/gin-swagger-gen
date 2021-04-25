package comment

import "fmt"

type Param struct {
	Name        string
	Required    bool
	paramType   string
	RefType     string
	Description string
}

func (p Param) Decs() string {
	if len(p.Description) == 0 {
		p.Description = p.Name
	}
	if p.paramType == "error" {
		p.paramType = "string"
	}
	return fmt.Sprintf("// @Param %s %s %s %t \"%s\"", p.Name, p.paramType, p.RefType, p.Required, p.Description)
}

type Params []Param

func (ps Params) Len() int {
	return len(ps)
}

func (ps Params) Less(i, j int) bool {
	return ps[i].paramType < ps[j].paramType
}

func (ps Params) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func NewBodyParam(name, refType, desc string) Param {
	return Param{
		Name:        name,
		Required:    true,
		paramType:   "body",
		RefType:     refType,
		Description: desc,
	}
}

func NewPathParam(name, refType, desc string) Param {
	return Param{
		Name:        name,
		Required:    true,
		paramType:   "path",
		RefType:     refType,
		Description: desc,
	}
}

func NewQueryParam(name, refType, desc string) Param {
	return Param{
		Name:        name,
		paramType:   "query",
		RefType:     refType,
		Description: desc,
	}
}

func NewFormDataParam(name, refType, desc string) Param {
	return Param{
		Name:        name,
		Required:    true,
		paramType:   "formData",
		RefType:     refType,
		Description: desc,
	}
}

type Resp struct {
	Code int
	Type string
}

func (r Resp) Decs() string {
	v1 := "// @Failure"
	if r.Code == 200 {
		v1 = fmt.Sprintf("// @Success")
	}
	v2 := "{object}"
	if r.Type == "string" {
		v2 = "string"
	}
	return fmt.Sprintf("%s %d %s %s", v1, r.Code, v2, r.Type)
}

type Route struct {
	RoutePath   string // route path
	RouteMethod string // route method: get post put
}

func (cr Route) Decs() string {
	return fmt.Sprintf("// @Router %s [%s]", cr.RoutePath, cr.RouteMethod)
}

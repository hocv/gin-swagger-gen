package parser

import (
	"fmt"
	"strings"

	"github.com/dave/dst"

	"github.com/hocv/gin-swagger-gen/ast"

	"github.com/pkg/errors"
)

var (
	version      = "version"
	title        = "title"
	description  = "description"
	contactName  = "contact.name"
	contactEmail = "contact.email"
	contactURL   = "contact.url"
	host         = "host"
	basePath     = "basepath"
	baseInfo     = []string{
		title,
		version,
		description,
		contactName,
		contactEmail,
		contactURL,
		host,
		basePath,
	}
	licenseInfo = []string{
		"license.name",
		"license.url",
	}
	tagInfo = []string{
		"tag.name",
		"tag.description",
		"tag.description.markdown",
		"tag.docs.url",
		"tag.docs.description",
	}
	securityInfo = []string{
		"securitydefinitions.basic",
		"securitydefinitions.apikey",
		"securitydefinitions.oauth2.application",
		"securitydefinitions.oauth2.implicit",
		"securitydefinitions.oauth2.password",
		"securitydefinitions.oauth2.accesscode",
	}
	defaultValue = map[string]string{
		title:        "Swagger Example API",
		version:      "1.0",
		description:  "This is a sample server Petstore server.",
		contactName:  "API Support",
		contactEmail: "http://www.swagger.io/support",
		contactURL:   "support@swagger.io",
		host:         "petstore.swagger.io",
		basePath:     "/v2",
	}
)

type InfoOption struct {
	Base     bool
	License  bool
	Tag      bool
	Security bool
}

type InfoParse struct {
	infos []string
}

func NewInfoParse(opt InfoOption) *InfoParse {
	infos := make([]string, 0, len(baseInfo))
	if opt.Base {
		infos = append(infos, baseInfo...)
	}
	if opt.License {
		infos = append(infos, licenseInfo...)
	}
	if opt.Tag {
		infos = append(infos, tagInfo...)
	}
	if opt.Security {
		infos = append(infos, securityInfo...)
	}
	return &InfoParse{infos: infos}
}

func (p *InfoParse) Parse(asts ast.Asts) error {
	var mainPkg *ast.Ast
	var mainFn *dst.FuncDecl
	var err error
	for _, a := range asts {
		if a.Pkg() != "main" {
			continue
		}
		mainPkg = a
		fds := a.Func("main")
		if len(fds) == 0 {
			return errors.New("not find main function")
		}
		mainFn = fds[0]
		break
	}

	if mainFn == nil || mainPkg == nil {
		return errors.Wrap(err, "main func")
	}

	commons := mainFn.Decs.Start.All()
	commonMap := make(map[string]string, len(commons))
	for _, common := range commons {
		trim := strings.TrimLeft(common, "// @")
		split := strings.Split(trim, " ")
		commonMap[split[0]] = common
	}
	mainFn.Decs.Start.Clear()
	for _, info := range p.infos {
		if common, ok := commonMap[info]; ok {
			mainFn.Decs.Start.Append(common)
			continue
		}
		desc := fmt.Sprintf("// @%s", info)
		if v, ok := defaultValue[info]; ok {
			desc = fmt.Sprintf("%s %s", desc, v)
		}
		mainFn.Decs.Start.Append(desc)
	}
	mainPkg.Dirty()
	return nil
}
package main

import (
	"log"

	"github.com/hocv/gin-swagger-gen/gen"
	"github.com/hocv/gin-swagger-gen/parser"
	"github.com/hocv/gin-swagger-gen/parser/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	searchDir       = kingpin.Flag("dir", "Directory you want to pars").Short('d').Default("./").ExistingDir()
	specifyFunc     = kingpin.Flag("func.name", "specify the function to add comment").Short('f').String()
	apiInfoDisable  = kingpin.Flag("info.disable", "pars api info").Default("false").Bool()
	apiInfoAll      = kingpin.Flag("info.all", "api info all").Default("false").Bool()
	apiInfoLicense  = kingpin.Flag("info.license", "api info license").Default("false").Bool()
	apiInfoTags     = kingpin.Flag("info.tags", "api info tags").Default("false").Bool()
	apiInfoSecurity = kingpin.Flag("info.security", "api info security").Default("false").Bool()
)

func main() {
	kingpin.Parse()

	g := gen.New(*searchDir)
	if len(*specifyFunc) > 0 {
		*apiInfoDisable = false
	}
	if *apiInfoDisable {
		if *apiInfoAll {
			*apiInfoLicense = true
			*apiInfoTags = true
			*apiInfoSecurity = true
		}
		opt := parser.InfoOption{
			Base:     true,
			License:  *apiInfoLicense,
			Tag:      *apiInfoTags,
			Security: *apiInfoSecurity,
		}
		infoParser := parser.NewInfoParse(opt)
		g.AddParser(infoParser)
	}
	g.AddParser(api.NewApiParse(*specifyFunc))

	if err := g.Parse(); err != nil {
		log.Println(err)
	}
	if err := g.Save(); err != nil {
		log.Println(err)
	}
}

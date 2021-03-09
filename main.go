package main

import (
	"log"

	"github.com/hocv/gin-swagger-gen/parser"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	searchDir   = kingpin.Flag("dir", "Directory you want to pars").Short('d').Default("./").ExistingDir()
	specifyFunc = kingpin.Flag("func.name", "specify the function to add comment").Short('f').String()
	justPrint   = kingpin.Flag("just.print", "just print, no save to file").Short('p').Bool()
)

func main() {
	kingpin.Parse()

	p := parser.New(*specifyFunc)
	p.ScanDir(*searchDir)
	p.Parse(*justPrint)
	if !*justPrint {
		if err := p.Save(); err != nil {
			log.Println(err)
		}
	}
}

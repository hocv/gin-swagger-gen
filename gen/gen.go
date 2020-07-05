package gen

import (
	"github.com/hocv/gin-swagger-gen/ast"
	"github.com/hocv/gin-swagger-gen/parser"
	"github.com/pkg/errors"
)

type Gen struct {
	parsers []parser.Parser
	asts    ast.Asts
}

func New(dir string, parsers ...parser.Parser) *Gen {
	g := &Gen{
		parsers: make([]parser.Parser, 0, len(parsers)),
		asts:    ast.NewAsts(dir),
	}
	for _, p := range parsers {
		g.AddParser(p)
	}
	return g
}

func (g *Gen) AddParser(p parser.Parser) {
	g.parsers = append(g.parsers, p)
}

func (g *Gen) Parse() error {
	for _, p := range g.parsers {
		if err := p.Parse(g.asts); err != nil {
			return errors.Wrap(err, "gen parse")
		}
	}
	return nil
}

func (g *Gen) Save() error {
	if err := g.asts.Save(); err != nil {
		return errors.Wrap(err, "gen save")
	}
	return nil
}

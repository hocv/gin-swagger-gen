package parser

import "github.com/hocv/gin-swagger-gen/ast"

type Parser interface {
	Parse(ast.Asts) error
}

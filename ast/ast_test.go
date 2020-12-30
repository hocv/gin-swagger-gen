package ast

import (
	"testing"
)

func TestAst_Struct(t *testing.T) {
	code := `
package test

import (
	f "fmt"
)

type (
	Result struct {
		Code int
		Data interface{}
	}

	Data struct {
		Name string
	}

	Re Result
	Res []Result
	Str string
	Strs []string
	StrMap map[string]string
	StrRes map[string]Result
)

var (
	Name string
	Names []string
	VarMap map[string]string
)

func Test() {
	
}
`
	a, err := New(code)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = a.Struct("Result")
	if err != nil {
		t.Fatal(err)
		return
	}
}

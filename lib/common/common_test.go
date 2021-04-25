package common

import (
	"fmt"
	"go/parser"
	"go/token"
	"testing"

	"github.com/dave/dst"

	"github.com/dave/dst/decorator"
)

func TestGetFuncParamByType(t *testing.T) {
	code := `package a

func handle(c,d *gin.Context) {

} 
`
	fileSet := token.NewFileSet()
	file, _ := decorator.ParseFile(fileSet, "", code, parser.ParseComments)
	value := GetFuncParamByType(file.Decls[0].(*dst.FuncDecl), "*gin.Context")
	if len(value) != 2 {
		t.Fatal("get params size wrong")
	}
	if value[0] != "c" {
		t.Fatalf("should be c, cur is %s", value[0])
	}
}

func TestGetFuncParams(t *testing.T) {
	code := `package a

func params(a, b string, c float64, t time.Duration, cg *gin.Context) {

} 
`
	params := map[string]string{
		"a":  "string",
		"b":  "string",
		"c":  "float64",
		"t":  "time.Duration",
		"cg": "*gin.Context",
	}
	fileSet := token.NewFileSet()
	file, _ := decorator.ParseFile(fileSet, "", code, parser.ParseComments)
	ps := GetFuncParams(file.Decls[0].(*dst.FuncDecl))
	if len(ps) != len(params) {
		t.Fatal("get params size wrong")
	}
	for k, v := range ps {
		if params[k] != v {
			t.Fatalf("%s should be %s", k, params[k])
		}
	}
}

func TestGetFuncVars(t *testing.T) {
	code := `package a

func params() {
	var a int
	var b = gin.New()
	c := g.Group("/d")
	var e,f int = 1,2
	g,h := 1, 2
	i := gen.Gen{}
	var j = gen.Gen{}
	k := &Gen{}
	var l = &Gen{}
	var m = &gen.Gen{}
} 
`
	fileSet := token.NewFileSet()
	file, _ := decorator.ParseFile(fileSet, "", code, parser.ParseComments)
	ps := GetFuncVars(file.Decls[0].(*dst.FuncDecl))
	if len(ps) != 12 {
		t.Fatal("get params size wrong")
	}
	fmt.Println(ps)
}

func TestGetJsonTag(t *testing.T) {
	arr := []struct {
		Tag  string
		Name string
	}{
		{Tag: `json:"a"`, Name: "a"},
		{Tag: `json:"b,omitempty"`, Name: "b"},
		{Tag: `json:"c" db:"c"`, Name: "c"},
		{Tag: `db:"d" json:"d"`, Name: "d"},
	}
	for _, jt := range arr {
		if tag := GetJsonTag(jt.Tag); tag != jt.Name {
			t.Fatalf("%s -> %s , but :%s ", jt.Name, jt.Tag, tag)
		}
	}
}

func TestSnakeCase(t *testing.T) {
	arr := []struct {
		Src string
		Dst string
	}{
		{
			Src: "UserName",
			Dst: "user_name",
		},
		{
			Src: "Lv1",
			Dst: "lv_1",
		},
		{
			Src: "AOName",
			Dst: "ao_name",
		},
		{
			Src: "ABC",
			Dst: "abc",
		},
	}
	for _, v := range arr {
		if str := SnakeCase(v.Dst); str != v.Dst {
			t.Fatalf("%s -> %s , but :%s ", v.Dst, v.Src, str)
		}
	}
}

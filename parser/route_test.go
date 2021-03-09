package parser

import (
	"testing"

	"github.com/hocv/gin-swagger-gen/lib/file"

	"github.com/hocv/gin-swagger-gen/lib/proj"
)

func TestParseRoute(t *testing.T) {
	p := proj.New()
	f, err := file.New("./test/route.go")
	if err != nil {
		t.Fatal(f)
		return
	}
	p.AddFile(f)
	dstFn, err := f.Func("routeTest")
	if err != nil {
		t.Fatal(f)
		return
	}

	rh := newRoute(p, "Default", "")
	rh.Parse(f, dstFn)
	if len(rh.Handles) != 7 {
		t.Fatal()
	}
}

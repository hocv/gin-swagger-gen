package parser

import (
	"fmt"
	"testing"

	"github.com/hocv/gin-swagger-gen/lib/file"
	"github.com/hocv/gin-swagger-gen/lib/proj"
)

func TestHandleAccept(t *testing.T) {
	p := proj.New()
	files := []string{
		"./test/handle.go",
		"./test/model/book/book.go",
		"./test/model/price/price.go",
	}

	for _, s := range files {
		f, err := file.New(s)
		if err != nil {
			t.Fatal(f)
			return
		}
		p.AddFile(f)
	}

	ffnd := p.GetFunc("test", "handleTest")

	rh := newRoute(p, "Default", "handleAccept")
	for f, fnd := range ffnd {
		rh.Parse(f, fnd)
		if len(rh.Handles) != 1 {
			t.Fatal()
		}

		rh.Handles[0].Parse()
		for _, s := range rh.Handles[0].Cmt.Decs() {
			fmt.Println(s)
		}
	}
}

func TestHandleProduct(t *testing.T) {
	p := proj.New()
	files := []string{
		"./test/handle.go",
		"./test/model/book/book.go",
		"./test/model/price/price.go",
	}

	for _, s := range files {
		f, err := file.New(s)
		if err != nil {
			t.Fatal(f)
			return
		}
		p.AddFile(f)
	}

	ffnd := p.GetFunc("test", "handleTest")
	if len(ffnd) != 1 {
		t.Fatal()
		return
	}

	rh := newRoute(p, "Default", "handleProduct")

	for f, fnd := range ffnd {
		rh.Parse(f, fnd)
		if len(rh.Handles) != 1 {
			t.Fatal()
		}

		rh.Handles[0].Parse()
		for _, s := range rh.Handles[0].Cmt.Decs() {
			fmt.Println(s)
		}
	}
}

func TestHandleRecursive(t *testing.T) {
	p := proj.New()

	files := []string{
		"./test/recursive.go",
		"./test/handle.go",
		"./test/model/book/book.go",
		"./test/model/price/price.go",
	}

	for _, s := range files {
		f, err := file.New(s)
		if err != nil {
			t.Fatal(f)
			return
		}
		p.AddFile(f)
	}

	ffnd := p.GetFunc("test", "recursiveTest")
	if len(ffnd) != 1 {
		t.Fatal()
		return
	}

	rh := newRoute(p, "Default", "handleRecursive")

	for f, fnd := range ffnd {
		rh.Parse(f, fnd)
		if len(rh.Handles) != 1 {
			t.Fatal()
		}

		rh.Handles[0].Parse()
		for _, s := range rh.Handles[0].Cmt.Decs() {
			fmt.Println(s)
		}
	}
}

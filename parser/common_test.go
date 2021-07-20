package parser

import (
	"fmt"
	"testing"
)

func TestFmtRoutePath(t *testing.T) {
	path := []string{
		"/user/:id/:aaa",
		"/user",
	}
	for _, s := range path {
		fmt.Println(fmtRoutePath(s))
	}
}

func TestRoutePathParams2(t *testing.T) {
	path := []string{
		"/user/{id}/{aaa}",
		"/user",
	}
	for _, s := range path {
		fmt.Println(routePathParams(s))
	}
}

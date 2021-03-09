package parser

import "testing"

func TestParse(t *testing.T) {
	p := New("")
	p.ScanDir("./test/main")
	p.Parse(true)
}

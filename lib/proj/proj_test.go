package proj

import (
	"fmt"
	"testing"
)

func TestFindStructInterface(t *testing.T) {
	a := New()
	a.ScanDir("./test")
	stru, err := a.GetStruct("test", "Resp3")
	if err != nil {
		t.Fatal(err)
		return
	}
	v := a.getInterfaceOfStruct("test", stru)
	fmt.Println(v)
}

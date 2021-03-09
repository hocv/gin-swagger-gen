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
		if tag := getJsonTag(jt.Tag); tag != jt.Name {
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
		if str := snakeCase(v.Dst); str != v.Dst {
			t.Fatalf("%s -> %s , but :%s ", v.Dst, v.Src, str)
		}
	}
}

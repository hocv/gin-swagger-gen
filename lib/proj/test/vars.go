package test

import (
	"fmt"
	"net/http"
)

type login struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func vars1() {
	lg := &login{}
	r := Resp{
		Code: http.StatusOK,
		Msg:  "ok",
		Data: lg,
	}
	fmt.Println(r)
}

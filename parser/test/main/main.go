package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	g := r.Group("/g")
	g.GET("/get", getHandle)
}

// @Summary get Handle sum
// @ID gff
// @Tags gff
// @Success 200 {object} Resp{data=User}
//aaaa
//bbbb
func getHandle(c *gin.Context) {
	id := c.Query("id")
	var user User
	user.Name = id
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Resp{
			Code: 0,
			Msg:  "",
			Data: nil,
		})
		return
	}
	resp := Resp{
		Code: http.StatusOK,
		Msg:  "ok",
		Data: user,
	}
	c.JSON(http.StatusOK, resp)
}

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

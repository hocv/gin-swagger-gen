package test

import (
	"net/http"

	"github.com/hocv/gin-swagger-gen/parser/test/model"

	"github.com/gin-gonic/gin"
)

func recursiveTest() {
	g := gin.Default()

	g.GET("/hdl_recursive", handleRecursive)

	_ = g.Run(":9090")
}

func handleRecursive(c *gin.Context) {
	b := model.Book{}
	recursive(c, b)
}

func recursive(c *gin.Context, data interface{}) {
	res := Resp{
		Code: http.StatusOK,
		Msg:  "ok",
		Data: data,
	}

	if err, ok := data.(error); ok {
		res.Code = http.StatusBadRequest
		res.Msg = err.Error()
		res.Data = err
	}
	c.JSON(http.StatusOK, res)
}

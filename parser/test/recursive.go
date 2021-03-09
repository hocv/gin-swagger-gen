package test

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func recursiveTest() {
	g := gin.Default()

	g.GET("/hdl_recursive", handleRecursive)

	_ = g.Run(":9090")
}

func handleRecursive(c *gin.Context) {
	recursive(c)
}

func recursive(c *gin.Context) {
	r := Resp{}
	c.JSON(http.StatusOK, r)
}

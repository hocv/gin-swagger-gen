package test

import (
	"github.com/gin-gonic/gin"
)

func routeTest() {
	g := gin.Default()

	g.GET("/get1", getHandler)

	group1 := g.Group("/group")
	group1.POST("/group_get1", groupHandler1)

	setRoute(g)

	_ = g.Run(":9090")
}

func setRoute(a *gin.Engine) {
	a.PUT("/api/:id", putHandler)

	var srv Srv
	a.PUT("/methodHandle/:usr", srv.handle)

	var g1, g2 = a.Group("/group1"), a.Group("/group2")
	g1.POST("/name", groupHandler2)
	g2.DELETE("/g1fu/:id", groupHandler3)
	g11 := g1.Group("/gg")
	g11.GET("/aaaa", groupHandler4)
}

type Srv struct{}

func (receiver Srv) handle(c *gin.Context) {}

func getHandler(c *gin.Context)    {}
func putHandler(c *gin.Context)    {}
func groupHandler1(c *gin.Context) {}
func groupHandler2(c *gin.Context) {}
func groupHandler3(c *gin.Context) {}
func groupHandler4(c *gin.Context) {}

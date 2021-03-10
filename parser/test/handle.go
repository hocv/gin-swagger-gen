package test

import (
	"fmt"
	"strconv"

	"github.com/hocv/gin-swagger-gen/parser/test/model"

	"github.com/gin-gonic/gin"
)

func handleTest() {
	g := gin.Default()

	group := g.Group("/group")
	group.GET("/hdl_accept", handleAccept)
	group.GET("/hdl_product", handleProduct)

	_ = g.Run(":9090")
}

func handleAccept(c *gin.Context) {
	q := c.Query("q1")
	i, _ := strconv.Atoi(q)
	fmt.Print(i)
	b := c.DefaultQuery("q2", "0")
	f, _ := c.GetPostForm("f1")

	fmt.Println(q, b, f)

	lg := &login{}
	_ = c.BindJSON(lg)

	var b1 login
	_ = c.BindXML(&b1)

	var b2 model.Book
	_ = c.BindYAML(&b2)
}

func getBook() (model.Book, error) {
	return model.Book{}, nil
}

type recv struct{}

func (r recv) book() model.Book {
	return model.Book{}
}

func (r recv) books() []model.Book {
	return []model.Book{}
}

func (r recv) book2() ([]model.Book, Resp) {
	return []model.Book{}, Resp{}
}

func handleProduct(c *gin.Context) {
	{
		r := Resp{
			Code: 0,
			Msg:  "",
			Data: model.Book{},
		}
		r.Data = 1
		c.JSON(0, r)
	}
	{
		r := recv{}
		price := r.book().Price()
		resp := Resp{
			Code:  0,
			Msg:   "ok",
			Data:  price,
			Data2: nil,
		}
		c.JSON(1, resp)
	}

	c.String(2, "f")
	{
		c.XML(3, Resp{
			Code:  0,
			Msg:   "",
			Data:  model.Book{},
			Data2: nil,
		})
	}
	{
		r := recv{}
		bs := r.books()
		c.JSON(4, bs)
	}
	{
		r := recv{}
		r1, r2 := Resp{
			Code:  0,
			Msg:   "",
			Data:  nil,
			Data2: Resp{},
		}, r.books()
		c.JSON(5, r1)
		c.JSON(6, r2)
	}
	{
		resp := Resp{
			Code:  0,
			Msg:   "",
			Data:  nil,
			Data2: nil,
		}
		r := recv{}
		resp.Data = r.books()
		c.JSON(7, resp)
	}

	{
		resp := Resp{
			Code:  0,
			Msg:   "",
			Data:  nil,
			Data2: nil,
		}
		r := recv{}
		resp.Data, resp.Data2 = r.book2()
		c.JSON(8, resp)
	}

	{
		resp := Resp{
			Code:  0,
			Msg:   "",
			Data:  nil,
			Data2: nil,
		}
		r := recv{}
		rr := Resp{}
		resp.Data, rr = r.book2()
		fmt.Println(rr)
		c.JSON(9, resp)
	}

	{
		resp := Resp{
			Code:  0,
			Msg:   "",
			Data:  Resp{},
			Data2: nil,
		}
		r := recv{}
		var res []model.Book
		res, resp.Data = r.book2()
		c.JSON(10, resp)
		c.JSON(11, res)
	}

	{
		b, err := getBook()
		c.JSON(12, b)
		c.JSON(13, err)
	}

	{
		var bo model.Book
		resp := Resp{
			Code: 0,
			Msg:  "msg",
			Data: bo,
		}
		c.JSON(14, resp)
	}
	{
		resp := Resp{
			Code: 0,
			Msg:  "msg",
			Data: Resp{
				Code: 0,
				Msg:  "",
				Data: nil,
				Data2: model.Book{
					Name: "",
					Auth: model.Book{
						Name: "",
						Auth: nil,
					},
				},
			},
			Data2: model.Book{
				Name: "",
				Auth: "",
			},
		}
		c.JSON(15, resp)
	}
	{
		b := model.Book{
			Name: "",
			Auth: nil,
		}
		b.Auth = Resp{
			Code:  0,
			Msg:   "",
			Data:  nil,
			Data2: nil,
		}
		resp := Resp{
			Code:  0,
			Msg:   "",
			Data:  b,
			Data2: nil,
		}

		resp.Data = model.Book{
			Name: "",
			Auth: Resp{},
		}

		resp.Data2 = Resp{}
		c.JSON(16, resp)
	}
}

type login struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Resp struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Data2 interface{} `json:"data2"`
}

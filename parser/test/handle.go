package test

import (
	"github.com/hocv/gin-swagger-gen/parser/test/model/book"
	"github.com/hocv/gin-swagger-gen/parser/test/model/price"

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
	// q := c.Query("q1")
	// i, _ := strconv.Atoi(q)
	// fmt.Print(i)
	// b := c.DefaultQuery("q2", "0")
	// f, _ := c.GetPostForm("f1")
	//
	// fmt.Println(q, b, f)
	//
	// lg := &login{}
	// _ = c.BindJSON(lg)
	//
	// var b1 login
	// _ = c.BindXML(&b1)
	//
	// var b2 model.Book
	// _ = c.BindYAML(&b2)

	var b3 price.Price
	_ = c.BindQuery(&b3)
}

func getBook() (book.Book, error) {
	return book.Book{}, nil
}

type recv struct {
	B book.Book
}

func (r recv) book() book.Book {
	return book.Book{}
}

func (r recv) books() []book.Book {
	return []book.Book{}
}

func (r recv) book2() ([]book.Book, Resp) {
	return []book.Book{}, Resp{}
}

func handleProduct(c *gin.Context) {
	// {
	// 	r := Resp{
	// 		Code: 0,
	// 		Msg:  "",
	// 		Data: model.Book{},
	// 	}
	// 	r.Data = 1
	// 	c.JSON(0, r) // Resp{data=int}
	// }
	// {
	// 	r := recv{}
	// 	price := r.book().GetPrice()
	// 	resp := Resp{
	// 		Code:  0,
	// 		Msg:   "ok",
	// 		Data:  price,
	// 		Data2: nil,
	// 	}
	// 	c.JSON(1, resp) // Resp{data=[]Price} todo: add pkg
	// }
	//
	// c.String(2, "f") // string
	// {
	// 	c.XML(3, Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  model.Book{},
	// 		Data2: nil,
	// 	}) //  Resp{data=model.Book}
	// }
	// {
	// 	r := recv{}
	// 	bs := r.books()
	// 	c.JSON(4, bs) // []model.Book
	// }
	// {
	// 	r := recv{}
	// 	r1, r2 := Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  nil,
	// 		Data2: Resp{},
	// 	}, r.books()
	// 	c.JSON(5, r1) // Resp{data2=Resp}
	// 	c.JSON(6, r2) // []model.Book
	// }
	// {
	// 	resp := Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  nil,
	// 		Data2: nil,
	// 	}
	// 	r := recv{}
	// 	resp.Data = r.books()
	// 	c.JSON(7, resp) // Resp{data=[]model.Book}
	// }
	//
	// {
	// 	resp := Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  nil,
	// 		Data2: nil,
	// 	}
	// 	r := recv{}
	// 	resp.Data, resp.Data2 = r.book2()
	// 	c.JSON(8, resp) // Resp{data2=Resp,data=[]model.Book}
	// }
	//
	// {
	// 	resp := Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  nil,
	// 		Data2: nil,
	// 	}
	// 	r := recv{}
	// 	rr := Resp{}
	// 	resp.Data, rr = r.book2()
	// 	fmt.Println(rr)
	// 	c.JSON(9, resp) // Resp{data=[]model.Book}
	// }
	//
	// {
	// 	resp := Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  Resp{},
	// 		Data2: nil,
	// 	}
	// 	r := recv{}
	// 	var res []model.Book
	// 	res, resp.Data = r.book2()
	// 	c.JSON(10, resp) // Resp{data=Resp}
	// 	c.JSON(11, res)  // []model.Book
	// }
	//
	// {
	// 	b, err := getBook()
	// 	c.JSON(12, b)   // model.Book
	// 	c.JSON(13, err) // {object} error todo: to string
	// }
	//
	// {
	// 	var bo model.Book
	// 	resp := Resp{
	// 		Code: 0,
	// 		Msg:  "msg",
	// 		Data: bo,
	// 	}
	// 	c.JSON(14, resp) // Resp{data=model.Book}
	// }
	// {
	// 	resp := Resp{
	// 		Code: 0,
	// 		Msg:  "msg",
	// 		Data: Resp{
	// 			Code: 0,
	// 			Msg:  "",
	// 			Data: nil,
	// 			Data2: model.Book{
	// 				Name: "",
	// 				Auth: model.Book{
	// 					Name: "",
	// 					Auth: nil,
	// 				},
	// 			},
	// 		},
	// 		Data2: model.Book{
	// 			Name: "",
	// 			Auth: "",
	// 		},
	// 	}
	// 	c.JSON(15, resp) // Resp{data=Resp{data2=model.Book{auth=model.Book}},data2=model.Book}
	// }
	// {
	// 	b := model.Book{
	// 		Name: "",
	// 		Auth: nil,
	// 	}
	// 	b.Auth = Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  nil,
	// 		Data2: nil,
	// 	}
	// 	resp := Resp{
	// 		Code:  0,
	// 		Msg:   "",
	// 		Data:  b,
	// 		Data2: nil,
	// 	}
	//
	// 	resp.Data = model.Book{
	// 		Name: "",
	// 		Auth: Resp{},
	// 	}
	//
	// 	resp.Data2 = Resp{}
	// 	c.JSON(16, resp) // Resp{data2=Resp,data=model.Book{auth=Resp}}
	// }
	// {
	// 	p := lib.Bk.GetPrice()
	// 	resp := Resp{
	// 		Code: 0,
	// 		Msg:  "",
	// 		Data: p,
	// 	}
	// 	c.JSON(17, resp) // Resp{data=[]Price} todo: add pkg
	// }
	{
		p := rr.B.GetPrice()
		resp := Resp{
			Code:  0,
			Msg:   "ok",
			Data:  p,
			Data2: nil,
		}
		c.JSON(18, resp) // Resp{data=[]Price}
	}
}

var rr = &recv{B: book.Book{}}
var lib = Lib{}

type Lib struct {
	Bk book.Book
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

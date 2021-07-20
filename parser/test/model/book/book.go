package book

import "github.com/hocv/gin-swagger-gen/parser/test/model/price"

type Book struct {
	Name string      `json:"name"`
	Auth interface{} `json:"auth"`
}

func (b Book) GetPrice() []price.Price {
	return []price.Price{}
}

package model

type Book struct {
	Name string      `json:"name"`
	Auth interface{} `json:"auth"`
}

func (b Book) GetPrice() []Price {
	return []Price{}
}

type Price struct {
	Value int    `form:"p_value" binding:"required"`
	Type  string `form:"p_type"`
}

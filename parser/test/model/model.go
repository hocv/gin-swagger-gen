package model

type Book struct {
	Name string      `json:"name"`
	Auth interface{} `json:"auth"`
}

func (b Book) Price() Price {
	return Price{}
}

type Price struct {
	Value int
	Type  string
}

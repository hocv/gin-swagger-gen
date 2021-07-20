package price

type Price struct {
	Value int    `form:"p_value" binding:"required"`
	Type  string `form:"p_type"`
}

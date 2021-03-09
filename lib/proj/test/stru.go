package test

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Resp2 struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		D1 struct {
			D2 interface{} `json:"d_2"`
		}
		Deep interface{} `json:"deep"`
	}
}

type Resp3 struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Res Resp `json:"res"`
	}
}

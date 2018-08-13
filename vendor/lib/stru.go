package lib

// Resp 用于定义返回数据格式(json)
type Resp struct {
	Ret    Code        `json:"ret"`
	Msg    string      `json:"msg,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// Header 用于定义通用头
type Header struct {
	ID    int64  `json:"h_id"`
	Token string `json:"h_token"`
}

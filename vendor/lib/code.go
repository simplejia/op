package lib

// Code 用于定义返回码
type Code int

// 定义各种返回码
const (
	CodeOk   Code = 1
	CodeSrv       = 2
	CodePara      = 3
)

// CodeMap 定义返回码对应的描述
var CodeMap = map[Code]string{
	CodeSrv:  "服务错误",
	CodePara: "参数错误",
}

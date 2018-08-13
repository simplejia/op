/*
Package model 用于模型层定义，所有db及cache对象封装均定义在这里。
只允许在这里添加对外暴露的接口
*/
package model

import (
	"github.com/simplejia/op/model/history"
	"github.com/simplejia/op/model/srv"
)

func NewSrv() *srv.Srv {
	return srv.NewSrv()
}

func NewHistory() *history.History {
	return history.NewHistory()
}

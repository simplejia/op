// 通用op（后台运营）系统
// author: simplejia
// date: 2017/12/8

//go:generate wsp -s -d

package main

import (
	"fmt"
	"lib"
	"net/http"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/utils"

	"github.com/simplejia/op/conf"
)

func init() {
	clog.AddrFunc = func() (string, error) {
		return lib.NameWrap(conf.C.Addrs.Clog)
	}
	clog.Init(conf.C.Clog.Name, "", conf.C.Clog.Level, conf.C.Clog.Mode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			return
		}
		http.Redirect(w, r, "/srv/list", http.StatusFound)
	})
}

func main() {
	fun := "main"
	clog.Info(fun)

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", conf.C.App.Port)
	err := utils.ListenAndServe(addr, nil)
	if err != nil {
		clog.Error("%s err: %v, addr: %v", fun, err, addr)
	}
}

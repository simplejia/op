// generated by wsp, DO NOT EDIT.

package main

import "net/http"
import "time"
import "github.com/simplejia/op/controller/history"
import "github.com/simplejia/op/controller/srv"
import "github.com/simplejia/op/filter"

func init() {
	http.HandleFunc("/history/remove", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(history.History)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/history/remove"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/history/remove"}); !ok {
			return
		}
		c.Remove(w, r)
	})

	http.HandleFunc("/srv/add", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/add"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/add"}); !ok {
			return
		}
		c.Add(w, r)
	})

	http.HandleFunc("/srv/delete", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/delete"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/delete"}); !ok {
			return
		}
		c.Delete(w, r)
	})

	http.HandleFunc("/srv/get", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/get"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/get"}); !ok {
			return
		}
		c.Get(w, r)
	})

	http.HandleFunc("/srv/list", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/list"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/list"}); !ok {
			return
		}
		c.List(w, r)
	})

	http.HandleFunc("/srv/srv_customer_list", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_customer_list"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_customer_list"}); !ok {
			return
		}
		c.SrvCustomerList(w, r)
	})

	http.HandleFunc("/srv/srv_customer_proc", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_customer_proc"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_customer_proc"}); !ok {
			return
		}
		c.SrvCustomerProc(w, r)
	})

	http.HandleFunc("/srv/srv_delete", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_delete"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_delete"}); !ok {
			return
		}
		c.SrvDelete(w, r)
	})

	http.HandleFunc("/srv/srv_get", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_get"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_get"}); !ok {
			return
		}
		c.SrvGet(w, r)
	})

	http.HandleFunc("/srv/srv_list", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_list"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_list"}); !ok {
			return
		}
		c.SrvList(w, r)
	})

	http.HandleFunc("/srv/srv_update", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_update"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/srv_update"}); !ok {
			return
		}
		c.SrvUpdate(w, r)
	})

	http.HandleFunc("/srv/update", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(srv.Srv)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/update"}); !ok {
				return
			}
		}()
		if ok := filter.Auth(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/srv/update"}); !ok {
			return
		}
		c.Update(w, r)
	})

}
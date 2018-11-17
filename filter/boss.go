package filter

import (
	"net/http"
	"runtime/debug"
	"time"

	"lib"

	"github.com/simplejia/clog/api"
)

// Boss 后置过滤器，用于数据上报，比如调用延时，出错等
func Boss(w http.ResponseWriter, r *http.Request, m map[string]interface{}) bool {
	err := m["__E__"]
	path := m["__P__"]
	bt := m["__T__"].(time.Time)

	if r.Form == nil {
		r.ParseForm()
	}

	for _, vs := range r.Form {
		for pos, v := range vs {
			if maxlen := 512; len(v) > maxlen {
				vs[pos] = lib.TruncateWithSuffix(v, maxlen, "...")
			}
		}
	}

	if err != nil {
		clog.Error("Boss() path: %v, body: %v, err: %v, stack: %s", path, r.Form, err, debug.Stack())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		clog.Info("Boss() path: %v, body: %v, elapse: %s", path, r.Form, time.Since(bt))
	}
	return true
}

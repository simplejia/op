package srv

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/simplejia/clog"

	"lib"

	"github.com/simplejia/op/model"
)

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) Delete(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.Delete"

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	srvModel := model.NewSrv()
	srvModel.ID = id
	if err := srvModel.Delete(); err != nil {
		detail := fmt.Sprintf("%s srv.Delete err: %v, req: %v", fun, err, id)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	srv.WriteJson(w, nil)
}

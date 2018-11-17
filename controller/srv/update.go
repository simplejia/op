package srv

import (
	"fmt"
	"net/http"

	"github.com/simplejia/utils"

	"github.com/simplejia/clog/api"

	"lib"

	"github.com/simplejia/op/model"
)

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) Update(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.Update"

	srvModel := model.NewSrv()
	if err := srvModel.ParseFromRequest(r); err != nil {
		detail := fmt.Sprintf("%s param err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	srvModelExist, err := srvModel.GetByName()
	if err != nil {
		detail := fmt.Sprintf("%s srv.GetByName err: %v, req: %v", fun, err, srvModel)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	if srvModelExist != nil && srvModelExist.ID != srvModel.ID {
		detail := fmt.Sprintf("%s name has exist, req: %v", fun, srvModelExist)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if err := srvModel.Update(); err != nil {
		detail := fmt.Sprintf("%s srv.Update err: %v, req: %v", fun, err, utils.Iprint(srvModel))
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	srv.WriteJson(w, nil)
}

package srv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/simplejia/clog"

	"lib"

	"github.com/simplejia/op/model"
	srv_model "github.com/simplejia/op/model/srv"
)

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) SrvDelete(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.SrvDelete"

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	srvModel := model.NewSrv()
	srvModel.ID = id
	srvModel, err := srvModel.Get()
	if err != nil || srvModel == nil {
		detail := fmt.Sprintf("%s srv.Get err: %v, req: %v", fun, err, id)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	var srvActionField *srv_model.SrvActionField
	for _, actionField := range srvModel.ActionFields {
		if actionField.Action.Kind == srv_model.ActionKindDelete {
			srvActionField = actionField
			break
		}
	}

	if srvActionField == nil {
		detail := fmt.Sprintf("%s must set delete action", fun)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	result, err := srv.FormToMap(r, srvActionField)
	if err != nil {
		detail := fmt.Sprintf("%s field err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	req, _ := json.Marshal(result)
	body, err := lib.PostProxy(srvModel.Addr, srvActionField.Action.Path, req)
	if err != nil {
		detail := fmt.Sprintf("%s post err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	srv.WriteJson(w, body)
}

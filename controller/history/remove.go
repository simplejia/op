package history

import (
	"fmt"
	"net/http"
	"strconv"

	"lib"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/op/model"
	"github.com/simplejia/utils"
)

// @prefilter("Auth")
// @postfilter("Boss")
func (history *History) Remove(w http.ResponseWriter, r *http.Request) {
	fun := "history.History.Update"

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	actionPath := r.URL.Query().Get("action_path")
	pos, _ := strconv.Atoi(r.URL.Query().Get("pos"))

	historyModel := model.NewHistory()
	historyModel.SrvID = id
	historyModel.SrvActionPath = actionPath
	headerVal, _ := history.GetParam(lib.KeyHeader)
	header := headerVal.(*lib.Header)
	historyModel.Uid = header.ID
	historyModel, err := historyModel.GetByUidAndSrv()
	if err != nil {
		detail := fmt.Sprintf("%s history.GetByUidAndSrv err: %v, req: %v", fun, err, historyModel)
		clog.Error(detail)
		history.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	if historyModel == nil || pos >= len(historyModel.Details) {
		detail := fmt.Sprintf("%s history invalid", fun)
		clog.Error(detail)
		history.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	historyModel.Details = append(historyModel.Details[:pos], historyModel.Details[pos+1:]...)

	if err := historyModel.Update(); err != nil {
		detail := fmt.Sprintf("%s history.Update err: %v, req: %v", fun, err, utils.Iprint(historyModel))
		clog.Error(detail)
		history.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusFound)
	return
}

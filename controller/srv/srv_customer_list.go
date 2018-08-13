package srv

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/simplejia/clog"

	"lib"

	"github.com/simplejia/op/model"
	srv_model "github.com/simplejia/op/model/srv"
)

var SrvCustomerListTpl = `
<html>
<head>
	<meta charset="UTF-8">
    <style type="text/css">
    input[type="button"]{
        padding: 0px 10px;
        height: 40px;
        width: 100px;
    }
    body{
		text-align: center;
    }
    table{
		border: 1px solid gray;
		border-collapse: collapse;
		width: auto;
        margin: 4px auto;
    }
    thead{
		background: #DDDDDD;
		border: 1px solid gray;
	}
	tr{
		border: 1px solid gray;
	}
	td{
		border: 1px solid gray;
	}
	th{
		border: 1px solid gray;
	}
    input[type="text"]{
        height: 30px;
        width: 200px;
        padding:10px 10px;
        vertical-align: middle;
    }
    button{
        border: 1px solid #DDDDDD;
        padding: 0px 10px;
    }
    textarea{
		width: 100%;
		height: 100px;
    }
	</style>

    <script type="text/javascript">
    </script>
</head>
<body>
	{{$id := .id}}
	<a href="/">返回首页</a>
	{{if .actions}}
	<table>
	<thead>
		<th>PATH</th><th>描述</th>
	</thead>
	<tbody>
	{{range .actions}}
	<tr>
		<td><button onclick="javascript:window.location.href='/srv/srv_customer_proc?id={{$id}}&action_path={{.Path}}'">{{.Path}}</button></td>
		<td>{{.Desc}}</td>
	</tr>
	{{end}}
	</tbody>
	</table>
	{{end}}
</body>
</html>
`

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) SrvCustomerList(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.SrvCustomerList"

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

	var actions []*srv_model.SrvAction
	for _, actionField := range srvModel.ActionFields {
		if actionField.Action.Kind != srv_model.ActionKindCustomer && actionField.Action.Kind != srv_model.ActionKindTransparent {
			continue
		}

		actions = append(actions, actionField.Action)
	}

	data := map[string]interface{}{
		"id":      id,
		"actions": actions,
	}

	tpl := template.Must(template.New("srv_customer_list").Parse(SrvCustomerListTpl))
	if err := tpl.Execute(w, data); err != nil {
		detail := fmt.Sprintf("%s tpl err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	return
}

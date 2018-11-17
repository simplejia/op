package srv

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/simplejia/clog/api"

	"lib"

	"github.com/simplejia/op/model"
	srv_model "github.com/simplejia/op/model/srv"
)

var SrvGetTpl = `
<html>
<head>
	<meta charset="UTF-8">
    <style type="text/css">
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
		width: 600px;
		height: 60px;
    }
	</style>

    <script type="text/javascript">
    function del_row(tr) {
		tr.innerHTML = "";
    }
    </script>
</head>
<body>
	<a href="/">返回首页</a>
	<form method="post" action="">
	<table>
	<thead>
	<tr></tr>
	</thead>
	<tbody>
		{{range $k, $v := .fields}}
		<tr>
			<td>{{$k}}</td>
			<td><textarea name="{{$k}}">{{$v}}</textarea></td>
			<td><button type="button" onclick="return del_row(this.parentElement.parentElement);">删除</button></td>
		</tr>
		{{end}}
	</tbody>
	</table>
	{{if eq .crud "u"}}
	<button type="submit" formaction="/srv/srv_update?id={{.id}}" onclick="javascript:return confirm('确认吗?')">执行</button>
	{{else}}
	<button type="submit" formaction="/srv/srv_delete?id={{.id}}" onclick="javascript:return confirm('确认吗?')">执行</button>
	{{end}}
	</form>
</body>
</html>
`

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) SrvGet(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.SrvGet"

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	crud := r.URL.Query().Get("crud")

	srvModel := model.NewSrv()
	srvModel.ID = id
	srvModel, err := srvModel.Get()
	if err != nil {
		detail := fmt.Sprintf("%s srv.Get err: %v, req: %v", fun, err, id)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	var srvActionField *srv_model.SrvActionField
	for _, actionField := range srvModel.ActionFields {
		if crud == "u" {
			if actionField.Action.Kind == srv_model.ActionKindUpdate {
				srvActionField = actionField
				break
			}
		} else {
			if actionField.Action.Kind == srv_model.ActionKindDelete {
				srvActionField = actionField
				break
			}
		}
	}

	if srvActionField == nil {
		detail := fmt.Sprintf("%s must set proper action", fun)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	fields := map[string]string{}
	for _, field := range srvActionField.Fields {
		fields[field.Name] = field.Param
	}

	r.ParseForm()
	for k, v := range r.PostForm {
		fields[k] = v[0]
	}

	data := map[string]interface{}{
		"id":     id,
		"fields": fields,
		"crud":   crud,
	}

	tpl := template.Must(template.New("srv_get").Parse(SrvGetTpl))
	if err := tpl.Execute(w, data); err != nil {
		detail := fmt.Sprintf("%s tpl err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	return
}

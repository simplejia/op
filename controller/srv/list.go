package srv

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/simplejia/clog/api"

	"lib"

	"github.com/simplejia/op/model"
)

var ListTpl = `
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
		height: 60px;
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
	</style>
	<script type="text/javascript">
	</script>
</head>
<body>
	<a href="/">返回首页</a>
	<form method="post" action="/srv/add">
		<button type="submit">新建</button>
	</form>
	<table>
	<thead>
	<tr>
		<th>ID</th><th>名字</th><th>名字(次要)</th><th>描述</th><th>地址</th><th>操作</th>
	</tr>
	</thead>
	<tbody>
	{{range .list}}
	<tr>
	<form method="post" action="">
		<td>{{.ID}}</td><td>{{.MName}}</td><td>{{.SName}}</td><td>{{.Desc}}</td><td>{{.Addr}}</td>
		<td>
			<button type="submit" formaction="/srv/get?id={{.ID}}">查看</button>
			<button type="submit" formaction="/srv/srv_list?id={{.ID}}">进入</button>
			<button type="submit" formaction="/srv/delete?id={{.ID}}" onclick="javascript:return confirm('确认吗?')">删除</button>
		</td>
	</form>
	</tr>
	{{end}}
	</tbody>
	</table>
	<button style="position:absolute;margin-left:-60px;" onclick="javascript:window.location.href='/srv/list?offset='+({{.offset}}-2*{{.limit}}>=0?{{.offset}}-2*{{.limit}}:0)+'&limit={{.limit}}'">上一页</button>
	<button style="position:absolute;margin-right:-60px;" onclick="javascript:window.location.href='/srv/list?offset={{.offset}}&limit={{.limit}}'">下一页</button>
	<button style="position:absolute;margin-left:60px;border:0;" type="button">总数: {{.total}}</button>
</body>
</html>
`

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) List(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.List"

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit <= 0 {
		limit = 50
	}

	total, err := model.NewSrv().Count()
	if err != nil {
		detail := fmt.Sprintf("%s srv.Count err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	srvs, err := model.NewSrv().List(offset, limit)
	if err != nil {
		detail := fmt.Sprintf("%s srv.List err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	data := map[string]interface{}{
		"list":   srvs,
		"offset": offset + len(srvs),
		"limit":  limit,
		"total":  total,
	}

	tpl := template.Must(template.New("list").Parse(ListTpl))
	if err := tpl.Execute(w, data); err != nil {
		detail := fmt.Sprintf("%s tpl err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	return
}

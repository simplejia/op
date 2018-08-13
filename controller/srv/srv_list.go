package srv

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/simplejia/clog"

	"lib"

	"github.com/simplejia/op/model"
	srv_model "github.com/simplejia/op/model/srv"
)

var SrvListTpl = `
<html>
<head>
	<meta charset="UTF-8">
    <style type="text/css">
    body{
		text-align: center;
    }
	span{
		width: 400px;
		display: block;
		word-wrap: break-word;
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
	</style>

    <script type="text/javascript">
    </script>
</head>
<body>
	{{$id := .id}}
	<a href="/">返回首页</a>
	<a href="/srv/srv_customer_list?id={{$id}}">其它</a>
	<table>
	<thead>
	<tr>
		<th>序号</th>
		<th>内容</th>
		<th>操作</th>
	</tr>
	</thead>
	<tbody>
	{{if .list}}
	{{range $pos, $elem := .list}}
	<tr>
	<form method="post" action="">
		<td>{{$pos}}</td>
		<td>
		{{range $k, $v := $elem}}
		<input type="hidden" name="{{$k}}" value="{{$v}}"/><span>{{$k}}: {{truncate $v 512 "..."}}</span>
		{{end}}
		</td>
		<td>
			<button type="submit" formaction="/srv/srv_get?id={{$id}}&crud=u">更新</button>
			<button type="submit" formaction="/srv/srv_get?id={{$id}}&crud=d">删除</button>
		</td>
	</form>
	</tr>
	{{end}}
	{{end}}
	</tbody>
	</table>
	<form method="post" action="">
		<input type="hidden" name="_" value="_"/>
		<table>
		<tr>
			{{range $k, $v := .fields}}
			<td>{{$k}}: <input style="width: 50px;" name="{{$k}}" value="{{$v}}"/></td>
			{{end}}
		</tr>
		</table>
		<button style="" type="submit">执行</button>
	</form>
</body>
</html>
`

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) SrvList(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.SrvList"

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
		if actionField.Action.Kind == srv_model.ActionKindList {
			srvActionField = actionField
			break
		}
	}

	if srvActionField == nil {
		http.Redirect(w, r, fmt.Sprintf("/srv/srv_customer_list?id=%d", id), http.StatusFound)
		return
	}

	result := map[string]interface{}{}

	if r.PostFormValue("_") == "" {
		needSupply := false
		for _, field := range srvActionField.Fields {
			if field.Required && field.Param == "" {
				needSupply = true
				break
			}
		}

		if needSupply {
			fields := map[string]string{}
			for _, field := range srvActionField.Fields {
				fields[field.Name] = field.Param
			}

			data := map[string]interface{}{
				"id":     id,
				"fields": fields,
			}

			funcMap := template.FuncMap{
				"truncate": lib.TruncateWithSuffix,
			}
			tpl := template.Must(template.New("srv_list").Funcs(funcMap).Parse(SrvListTpl))
			if err := tpl.Execute(w, data); err != nil {
				detail := fmt.Sprintf("%s tpl err: %v", fun, err)
				clog.Error(detail)
				srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
				return
			}
			return
		}

		for _, field := range srvActionField.Fields {
			if field.Param == "" {
				continue
			}

			v, err := srv.FieldValue(field, field.Param)
			if err != nil {
				detail := fmt.Sprintf("%s field err: %v", fun, err)
				clog.Error(detail)
				srv.ReplyFailWithDetail(w, lib.CodePara, detail)
				return
			}
			result[field.Name] = v
		}
	} else {
		_result, err := srv.FormToMap(r, srvActionField)
		if err != nil {
			detail := fmt.Sprintf("%s field err: %v", fun, err)
			clog.Error(detail)
			srv.ReplyFailWithDetail(w, lib.CodePara, detail)
			return
		}

		result = _result
	}

	req, _ := json.Marshal(result)
	body, err := lib.PostProxy(srvModel.Addr, srvActionField.Action.Path, req)
	if err != nil {
		detail := fmt.Sprintf("%s post err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	resp := &struct {
		lib.Resp
		Data map[string]json.RawMessage `json:"data"`
	}{}

	if err := json.Unmarshal(body, resp); err != nil {
		detail := fmt.Sprintf("%s decode err: %v, req: %s", fun, err, body)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	if resp.Ret != lib.CodeOk {
		srv.ReplyOk(w, resp)
		return
	}

	respDataList := []map[string]json.RawMessage{}
	if err := json.Unmarshal(resp.Data["list"], &respDataList); err != nil {
		detail := fmt.Sprintf("%s list element not found in resp data, err: %v, req: %s", fun, err, body)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	list := []map[string]interface{}{}

	for _, srcElem := range respDataList {
		dstElem := map[string]interface{}{}
		for k, srcV := range srcElem {
			dstElem[k] = string(srcV)
		}

		list = append(list, dstElem)
	}

	fields := map[string]string{}
	for _, field := range srvActionField.Fields {
		fields[field.Name] = field.Param
	}

	for k, v := range result {
		s, _ := json.Marshal(v)
		fields[k] = string(s)
	}

	for k, v := range resp.Data {
		if k == "list" {
			continue
		}

		fields[k] = string(v)
	}

	data := map[string]interface{}{
		"id":     id,
		"list":   list,
		"fields": fields,
	}

	funcMap := template.FuncMap{
		"truncate": lib.TruncateWithSuffix,
	}
	tpl := template.Must(template.New("srv_list").Funcs(funcMap).Parse(SrvListTpl))
	if err := tpl.Execute(w, data); err != nil {
		detail := fmt.Sprintf("%s tpl err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	return
}

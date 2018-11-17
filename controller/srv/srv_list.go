package srv

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/simplejia/clog/api"

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
    textarea{
		width: 100%;
		height: 40px;
    }
	</style>

    <script type="text/javascript">
	function toggle_checkbox(obj, name) {
		document.getElementsByName(name).forEach(function(e){
			e.checked = obj.checked;
		})
	}
    function del_row(tr) {
		tr.innerHTML = "";
    }
    </script>
</head>
<body>
	{{$id := .id}}
	{{$select_field_params := .select_field_params}}
	{{$multi_select_field_params := .multi_select_field_params}}
	<a href="/">返回首页</a>
	<a href="/srv/srv_customer_list?id={{$id}}">其它</a>
	<form method="post" action="">
		<input type="hidden" name="_" value="_"/>
		<table>
		{{range $k, $v := .fields}}
		<tr>
			<td>{{$k}}</td>
			{{if ne ((index $select_field_params $k)|len) 0}}
			<td>
				{{range index $select_field_params $k}}
				<input name="{{$k}}" type="radio" value="{{.Value}}" {{if eq .Value $v}}checked="checked"{{end}}/>{{.Desc}}
				{{end}}
			</td>	
			{{else if ne ((index $multi_select_field_params $k)|len) 0}}
			<td>
				<p>
				<input type="checkbox" onclick="toggle_checkbox(this, '{{$k}}');" />*
				</p>
				{{range index $multi_select_field_params $k}}
				<input name="{{$k}}" type="checkbox" value="{{.Value}}" {{if (is_a_in_b .Value $v)}}checked="checked"{{end}}/>{{.Desc}}
				{{end}}
			</td>	
			{{else}}
			<td>
				<textarea name="{{$k}}">{{$v}}</textarea>
			</td>
			{{end}}
			<td><button type="button" onclick="return del_row(this.parentElement.parentElement);">删除</button></td>
		</tr>
		{{end}}
		</table>
		<button style="" type="submit">执行</button>
	</form>
	{{with .list}}
	<table>
	<thead>
	<tr>
		<th>序号</th>
		<th>内容</th>
		<th>操作</th>
	</tr>
	</thead>
	<tbody>
	{{range $pos, $elem := .}}
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
	</tbody>
	</table>
	{{end}}
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

	selectFieldParams := make(map[string][]*FieldValueDesc)
	multiSelectFieldParams := make(map[string][]*FieldValueDesc)
	for _, field := range srvActionField.Fields {
		if field.Source == srv_model.FieldSourceArray {
			options, err := ParseFieldParam(field.Param)
			if err != nil {
				detail := fmt.Sprintf("%s field err: %v, req: %v", fun, err, field.Param)
				clog.Error(detail)
				srv.ReplyFailWithDetail(w, lib.CodePara, detail)
				return
			}
			if field.Kind == srv_model.FieldKindArray {
				multiSelectFieldParams[field.Name] = options
			} else {
				selectFieldParams[field.Name] = options
			}
		} else if field.Source == srv_model.FieldSourceUrl {
			path := field.Param
			body, err := lib.PostProxy(srvModel.Addr, path, nil)
			if err != nil {
				detail := fmt.Sprintf("%s get field param post err: %v, req: %v", fun, err, path)
				clog.Error(detail)
				srv.ReplyFailWithDetail(w, lib.CodePara, detail)
				return
			}

			resp := &struct {
				lib.Resp
				Data struct {
					List json.RawMessage
				} `json:"data"`
			}{}
			if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Ret != lib.CodeOk {
				detail := fmt.Sprintf("%s get field post response err: %v, req: %v, resp: %s", fun, err, path, body)
				clog.Error(detail)
				srv.ReplyFailWithDetail(w, lib.CodePara, detail)
				return
			}

			options, err := ParseFieldParam(string(resp.Data.List))
			if err != nil {
				detail := fmt.Sprintf("%s field err: %v, req: %v", fun, err, field.Param)
				clog.Error(detail)
				srv.ReplyFailWithDetail(w, lib.CodePara, detail)
				return
			}

			if field.Kind == srv_model.FieldKindArray {
				multiSelectFieldParams[field.Name] = options
			} else {
				selectFieldParams[field.Name] = options
			}
		}
	}

	fields := map[string]string{}
	for _, field := range srvActionField.Fields {
		if field.Source == srv_model.FieldSourceUser {
			fields[field.Name] = field.Param
		} else {
			fields[field.Name] = ""
		}
	}

	if r.PostFormValue("_") == "" {
		needSupply := false
		for _, field := range srvActionField.Fields {
			if field.Required &&
				(field.Source != srv_model.FieldSourceUser || field.Param == "") {
				needSupply = true
				break
			}
		}

		if needSupply {
			data := map[string]interface{}{
				"id":                        id,
				"fields":                    fields,
				"select_field_params":       selectFieldParams,
				"multi_select_field_params": multiSelectFieldParams,
			}

			funcMap := template.FuncMap{
				"truncate":  lib.TruncateWithSuffix,
				"is_a_in_b": IsAInB,
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

	list := []map[string]interface{}{}

	if l := resp.Data["list"]; len(l) > 0 {
		respDataList := []map[string]json.RawMessage{}
		if err := json.Unmarshal(l, &respDataList); err != nil {
			detail := fmt.Sprintf("%s list element not found in resp data, err: %v, req: %s", fun, err, body)
			clog.Error(detail)
			srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
			return
		}

		for _, srcElem := range respDataList {
			dstElem := map[string]interface{}{}
			for k, srcV := range srcElem {
				dstElem[k] = string(srcV)
			}

			list = append(list, dstElem)
		}
	}

	for name, v := range result {
		s, _ := json.Marshal(v)
		fields[name] = string(s)
	}

	for k, v := range resp.Data {
		if k == "list" {
			continue
		}

		fields[k] = string(v)
	}

	data := map[string]interface{}{
		"id":                        id,
		"list":                      list,
		"fields":                    fields,
		"select_field_params":       selectFieldParams,
		"multi_select_field_params": multiSelectFieldParams,
	}

	funcMap := template.FuncMap{
		"truncate":  lib.TruncateWithSuffix,
		"is_a_in_b": IsAInB,
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

package srv

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/simplejia/clog"

	"lib"

	"github.com/simplejia/op/model"
	history_model "github.com/simplejia/op/model/history"
	srv_model "github.com/simplejia/op/model/srv"
)

var SrvCustomerProcTpl = `
<html>
<head>
	<meta charset="UTF-8">
    <style type="text/css">
	span{
		width: 400px;
		display: block;
		word-wrap: break-word;
	}
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
	select{
        width: 100%;
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
    function del_row(tr) {
		tr.innerHTML = "";
	}

	function file_change(file_input) {
		if (file_input.files && file_input.files.length > 0 && file_input.files[0].size > 0) {
			var reader = new FileReader();
			reader.onload = function (evt) {
				if (evt.target.readyState == FileReader.DONE) {
					str = evt.target.result.replace(/^data:.*;base64,/, "");
					document.getElementsByName(file_input.id)[0].value = str;
				}
			};
			reader.onprogress = function(p) {}
			reader.readAsDataURL(file_input.files[0]);
		}
	}

	function safe_get(object, key, default_value="") {
		console.log("safe_get",object,key);
		if(object.hasOwnProperty(key)){
			return object[key];
		} else {
			return default_value;
		}
	}

	function get_selected(get, expect) {
		if(get == expect){
			return ' selected="selected" '
		} else {
			return "";
		}
	}

	function toggle_checkbox(obj, name) {
		document.getElementsByName(name).forEach(function(e){
			e.checked = obj.checked;
		})
	}
    </script>
</head>
<body>
	{{$id := .id}}
	{{$action_path := .action_path}}
	{{$file_fields := .file_fields}}
	{{$select_fields := .select_fields}}
	{{$select_field_params := .select_field_params}}
	{{$multi_select_field_params := .multi_select_field_params}}
	<a href="/">返回首页</a>
	<form method="post" action="">
		<input type="hidden" name="_" value="_"/>
		<table>
		<thead>
		<tr></tr>
		</thead>
		<tbody>
			{{range $k, $v := .fields}}
			<tr>
				<td>{{$k}}</td>
				{{if eq (index $file_fields $k) true}}
				<td>
					<input type="file" id="{{$k}}" onchange="return file_change(this);" onclick="this.value=null;"/>
					<input type="hidden" name="{{$k}}"/>
				</td>	
				{{else if ne ((index $select_field_params $k)|len) 0}}
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
				<td>
					<button type="button" onclick="return del_row(this.parentElement.parentElement);">删除</button>
				</td>
			</tr>
			{{end}}
		</tbody>
		</table>
		<button type="submit" onclick="javascript:return confirm('确认吗?')">执行</button>
	</form>
	{{if .history_details}}
	<table>
	<tr>
	<th>参数</th><th>时间</th><th>操作</th>
	</tr>
	{{range $pos, $elem := .history_details}}
	<tr>
		<form method="post" action="">
		<td>
		{{range $k, $v := $elem.M}}
		<input type="hidden" name="{{$k}}" value="{{$v}}"/><span>{{$k}}: {{truncate $v 512 "..."}}</span>
		{{end}}
		</td>
		<td>{{.Ct}}</td>
		<td>
        <p>
			<button type="submit" formaction="">选我</button>
			<button type="button" onclick="javascript:window.location.href='/history/remove?id={{$id}}&action_path={{$action_path}}&pos={{$pos}}'">删除</button>
        </p>
		</td>
		</form>
	</tr>
	{{end}}
	</table>
	{{end}}
</body>
</html>
`

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) SrvCustomerProc(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.SrvCustomerProc"

	r.Body = http.MaxBytesReader(w, r.Body, 1e9)

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	actionPath := r.URL.Query().Get("action_path")

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
		if actionField.Action.Path == actionPath {
			srvActionField = actionField
			break
		}
	}

	if srvActionField == nil {
		detail := fmt.Sprintf("%s actionField empty", fun)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	fileFields := map[string]bool{}
	selectFieldParams := make(map[string][]*FieldValueDesc)
	multiSelectFieldParams := make(map[string][]*FieldValueDesc)
	for _, field := range srvActionField.Fields {
		if field.Kind == srv_model.FieldKindFile {
			fileFields[field.Name] = true
		} else if field.Source == srv_model.FieldSourceArray {
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

	if r.FormValue("_") == "" {
		fields := map[string]string{}
		for _, field := range srvActionField.Fields {
			if field.Source == srv_model.FieldSourceUser {
				fields[field.Name] = field.Param
			} else {
				fields[field.Name] = ""
			}
		}

		r.ParseForm()
		for k, v := range r.PostForm {
			fields[k] = v[0]
		}

		historyModel := model.NewHistory()
		historyModel.SrvID = id
		historyModel.SrvActionPath = actionPath
		headerVal, _ := srv.GetParam(lib.KeyHeader)
		header := headerVal.(*lib.Header)
		historyModel.Uid = header.ID
		historyModel, err := historyModel.GetByUidAndSrv()
		if err != nil {
			detail := fmt.Sprintf("%s history.GetByUidAndSrv err: %v, req: %v", fun, err, historyModel)
			clog.Error(detail)
			srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
			return
		}

		var historyDetails []*history_model.HistoryDetail
		if historyModel != nil {
			historyDetails = historyModel.Details
		}

		data := map[string]interface{}{
			"id":                        id,
			"action_path":               actionPath,
			"fields":                    fields,
			"file_fields":               fileFields,
			"history_details":           historyDetails,
			"select_field_params":       selectFieldParams,
			"multi_select_field_params": multiSelectFieldParams,
		}

		funcMap := template.FuncMap{
			"truncate":  lib.TruncateWithSuffix,
			"is_a_in_b": IsAInB,
		}
		tpl := template.Must(template.New("srv_customer_proc").Funcs(funcMap).Parse(SrvCustomerProcTpl))
		if err := tpl.Execute(w, data); err != nil {
			detail := fmt.Sprintf("%s tpl err: %v", fun, err)
			clog.Error(detail)
			srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
			return
		}
		return
	}

	result, err := srv.FormToMap(r, srvActionField)
	if err != nil {
		detail := fmt.Sprintf("%s field err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	var body []byte
	if srvActionField.Action.Kind == srv_model.ActionKindTransparent {
		req, err := json.Marshal(result)
		_body, header, err := lib.PostProxyReturnHeader(srvModel.Addr, srvActionField.Action.Path, req)
		if err != nil {
			detail := fmt.Sprintf("%s post err: %v", fun, err)
			clog.Error(detail)
			srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
			return
		}
		for k, vs := range header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		body = _body
	} else {
		req, err := json.Marshal(result)
		_body, err := lib.PostProxy(srvModel.Addr, srvActionField.Action.Path, req)
		if err != nil {
			detail := fmt.Sprintf("%s post err: %v", fun, err)
			clog.Error(detail)
			srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
			return
		}
		body = _body
	}

	if srvActionField.Action.Kind == srv_model.ActionKindTransparent {
		w.Write(body)
	} else {
		srv.WriteJson(w, body)
	}

	historyDetail := &history_model.HistoryDetail{
		M:  map[string]string{},
		Ct: time.Now(),
	}
	for k, v := range result {
		if fileFields[k] {
			continue
		}
		s, _ := json.Marshal(v)
		historyDetail.M[k] = string(s)
	}

	historyModel := model.NewHistory()
	historyModel.SrvID = id
	historyModel.SrvActionPath = actionPath
	historyModel.Details = append(historyModel.Details, historyDetail)
	headerVal, _ := srv.GetParam(lib.KeyHeader)
	header := headerVal.(*lib.Header)
	historyModel.Uid = header.ID

	if err := historyModel.Add(); err != nil {
		detail := fmt.Sprintf("%s history.Add err: %v, req: %v", fun, err, historyModel)
		clog.Warn(detail)
	}

	return
}

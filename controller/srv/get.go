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

var GetTpl = `
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
    incell{
		float: left; list-style: none;
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
    function escape_string(str) {
        return str.replace(/[<>&"']/g, function(c){return {'<':'&lt;','>':'&gt;','&':'&amp;','"':'&quot;','\'':'&apos;'}[c];});
    }

    function safe_get(object, key, default_value="") {
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

    function add_field(tbl, action_num, field={}) {
        var tr = tbl.insertRow();
        var td = tr.insertCell();
        td.innerHTML = '<input type="checkbox"/>';
        var td = tr.insertCell();
        td.innerHTML = '<input type="text" name="field_name_'+action_num+'" value="'+safe_get(field,"name")+'"/>';
        var td = tr.insertCell();
        td.innerHTML = '\
        <select name="field_require_'+action_num+'">\
            <option value="true"'+get_selected(safe_get(field,"required",true),true)+'>必填</option>\
            <option value="false"'+get_selected(safe_get(field,"required",true),false)+'>选填</option>\
        </select>';
        var td = tr.insertCell();
        td.innerHTML = '\
        <select name="field_kind_'+action_num+'">\
            <option value="0"'+get_selected(safe_get(field,"kind"),"0")+'>-</option>\
            <option value="1"'+get_selected(safe_get(field,"kind"),"1")+'>string</option>\
            <option value="2"'+get_selected(safe_get(field,"kind"),"2")+'>integer</option>\
            <option value="3"'+get_selected(safe_get(field,"kind"),"3")+'>float</option>\
            <option value="4"'+get_selected(safe_get(field,"kind"),"4")+'>bool</option>\
            <option value="5"'+get_selected(safe_get(field,"kind"),"5")+'>map</option>\
            <option value="6"'+get_selected(safe_get(field,"kind"),"6")+'>array</option>\
            <option value="7"'+get_selected(safe_get(field,"kind"),"7")+'>file</option>\
        </select>';

        var td = tr.insertCell();
        td.innerHTML = '\
        <select name="field_source_'+action_num+'">\
            <option value="1"'+get_selected(safe_get(field,"source"),"1")+'>自定义</option>\
            <option value="2"'+get_selected(safe_get(field,"source"),"2")+'>从数组</option>\
            <option value="3"'+get_selected(safe_get(field,"source"),"3")+'>从URL</option>\
        </select>';

        var td = tr.insertCell();
        td.innerHTML = '<input type="text" name="field_param_'+action_num+'" value="'+escape_string(safe_get(field,"param"))+'"/>';
	}

    function del_field(tbl) {
        for (var i=tbl.rows.length-1; i>=1; i--) {
            if (tbl.rows[i].cells[0].children[0].checked) {
                tbl.deleteRow(i);
            }
        }
    }

    var g_action_num = 0;

    function add_action(af={}) {
		var action = safe_get(af,"action",{})
		var fields = safe_get(af,"fields",[])
        var action_num = g_action_num++
        var div = document.getElementById("actions");
        var tbl = document.createElement("table");
        tbl.style.width = "800px";
        var tbl_h = tbl.createTHead()
        var tr = tbl_h.insertRow();
        tr.innerHTML = '<th></th><th>PATH</th><th>描述</th><th>类型</th><th colspan="2">操作</th>';
        var tr = tbl_h.insertRow();
        var td = tr.insertCell();
        td.innerHTML = '<input type="checkbox"/>';
        var td = tr.insertCell();
        td.innerHTML = '<input type="text" name="action_path_'+action_num+'" value="'+safe_get(action,"path")+'"/>';
        var td = tr.insertCell();
        td.innerHTML = '<input type="text" name="action_desc_'+action_num+'" value="'+safe_get(action,"desc")+'"/>';
        var td = tr.insertCell();
        td.innerHTML = '\
        <select name="action_kind_'+action_num+'">\
            <option value="1"'+get_selected(safe_get(action,"kind"),"1")+'>list</option>\
            <option value="2"'+get_selected(safe_get(action,"kind"),"2")+'>update</option>\
            <option value="3"'+get_selected(safe_get(action,"kind"),"3")+'>delete</option>\
            <option value="4"'+get_selected(safe_get(action,"kind"),"4")+'>customer</option>\
            <option value="5"'+get_selected(safe_get(action,"kind"),"5")+'>transparent</option>\
		</select>';
        var td = tr.insertCell();
        td.colSpan = '2'
		td.innerHTML = '\
		<button type="button">增加FIELD</button>\
        <button type="button">删除FIELD</button>\
		';
        var tbl_b = tbl.createTBody()
        var tr = tbl_b.insertRow();
        tr.innerHTML = '<th></th><th>字段名</th><th>是否必须</th><th>类型</th><th>数据源</th><th>参数</th>';
        td.children[0].onclick=function() {
            add_field(tbl_b, action_num);
		}
        td.children[1].onclick=function() {
            del_field(tbl_b);
        }
		div.appendChild(tbl)

		for (var i=0; i<fields.length; i++) {
			add_field(tbl_b, action_num, fields[i]);
		}
    }

    function del_action() {
        var div = document.getElementById("actions");
        for (var i=div.childElementCount-1; i>=0; i--) {
            if (div.children[i].tHead.children[1].cells[0].children[0].checked) {
                div.removeChild(div.children[i]);
            }
        }
    }
	
	window.onload=function(){
		{{range .ActionFields}}
		add_action({{.}});
        {{end}}
	}
    </script>
</head>
<body>
	<a href="/">返回首页</a>
    <form action="/srv/update" method="post">
        <input type="hidden" name="_" value="_"/>
        <table>
            <tr><td>ID</td><td><input type="text" name="id" value="{{.ID}}" readonly="true"/></td></tr>
            <tr><td>name</td><td><input type="text" name="m_name" value="{{.MName}}"/></td></tr>
            <tr><td>name(次要)</td><td><input type="text" name="s_name" value="{{.SName}}"/></td></tr>
            <tr><td>desc</td><td><input type="text" name="desc" value="{{.Desc}}"/></td></tr>
            <tr><td>addr</td><td><input type="text" name="addr" value="{{.Addr}}"/></td></tr>
        </table>
        <p>
            <button type="button" onclick="add_action()">增加ACTION</button>
            <button type="button" onclick="del_action()">删除ACTION</button>
        </p>
		<p>
			<div id="actions"></div>
		</p>
        <p><button type="submit">执行</button></p>
    </form>
</body>
</html>
`

// @prefilter("Auth")
// @postfilter("Boss")
func (srv *Srv) Get(w http.ResponseWriter, r *http.Request) {
	fun := "srv.Srv.Get"

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

	data := srvModel

	tpl := template.Must(template.New("get").Parse(GetTpl))
	if err := tpl.Execute(w, data); err != nil {
		detail := fmt.Sprintf("%s tpl err: %v", fun, err)
		clog.Error(detail)
		srv.ReplyFailWithDetail(w, lib.CodeSrv, detail)
		return
	}

	return
}

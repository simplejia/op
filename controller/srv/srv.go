package srv

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"

	"lib"

	srv_model "github.com/simplejia/op/model/srv"
)

type Srv struct {
	lib.Base
}

func (srv *Srv) FieldValue(field *srv_model.SrvField, value string) (ret interface{}, err error) {
	switch field.Kind {
	case srv_model.FieldKindInteger:
		if value == "" {
			ret = 0
		} else {
			var v uint64
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, err
			}
			ret = v
		}
	case srv_model.FieldKindString:
		if value == "" {
			ret = ""
		} else {
			if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
				v, err := strconv.Unquote(value)
				if err != nil {
					return nil, err
				}
				ret = v
			} else {
				ret = value
			}
		}
	case srv_model.FieldKindFloat:
		if value == "" {
			ret = 0.0
		} else {
			var v float64
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			ret = v
		}
	case srv_model.FieldKindBool:
		if value == "" {
			ret = false
		} else {
			var v bool
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			ret = v
		}
	case srv_model.FieldKindMap:
		if value == "" {
			ret = nil
		} else {
			var v map[string]json.RawMessage
			err := json.Unmarshal([]byte(value), &v)
			if err != nil {
				return nil, err
			}
			ret = v
		}
	case srv_model.FieldKindArray:
		if value == "" {
			ret = nil
		} else {
			var v []json.RawMessage
			err := json.Unmarshal([]byte(value), &v)
			if err != nil {
				return nil, err
			}
			ret = v
		}
	case srv_model.FieldKindNone:
		if value == "" {
			ret = ""
		} else {
			var v json.RawMessage
			if err := json.Unmarshal([]byte(value), &v); err != nil {
				return nil, err
			}
			ret = v
		}
	case srv_model.FieldKindFile:
		ret = value
	}

	return
}

func (srv *Srv) FormToMap(r *http.Request, srvActionField *srv_model.SrvActionField) (result map[string]interface{}, err error) {
	result = map[string]interface{}{}

	r.ParseForm()

	fields := map[string]*srv_model.SrvField{}
	for _, field := range srvActionField.Fields {
		fields[field.Name] = field
	}

	for name, vs := range r.PostForm {
		if name == "_" {
			continue
		}

		field := fields[name]
		if field == nil {
			value := vs[0]
			if value == "" {
				result[name] = ""
			} else {
				var v json.RawMessage
				if err := json.Unmarshal([]byte(value), &v); err != nil {
					return nil, err
				}
				result[name] = v
			}
			continue
		}

		if field.Kind == srv_model.FieldKindArray &&
			(field.Source == srv_model.FieldSourceUrl || field.Source == srv_model.FieldSourceArray) {
			var r []json.RawMessage
			for _, value := range vs {
				var v json.RawMessage
				if err := json.Unmarshal([]byte(value), &v); err != nil {
					return nil, err
				}
				r = append(r, v)
			}
			result[name] = r
			continue
		}

		value := vs[0]
		ret, err := srv.FieldValue(field, value)
		if err != nil {
			return nil, err
		}

		result[name] = ret
	}

	for _, field := range srvActionField.Fields {
		name := field.Name
		if r.PostFormValue(name) != "" {
			continue
		}

		param := field.Param
		if field.Required &&
			(field.Source != srv_model.FieldSourceUser || param == "") {
			err = fmt.Errorf("%s must supply", name)
			return nil, err
		}

		if field.Source != srv_model.FieldSourceUser {
			continue
		}

		v, err := srv.FieldValue(field, param)
		if err != nil {
			return nil, err
		}
		result[name] = v
	}

	return
}

func (srv *Srv) WriteJson(w http.ResponseWriter, body []byte) {
	t := `
	<html><head><meta charset="UTF-8"></head><body style="text-align:center;">
	<a href="/" style="display:block;">返回首页</a>
	%s
	</body></html>`

	if len(body) > 0 {
		t = fmt.Sprintf(t, fmt.Sprintf(`
		<script type="text/javascript">
			function unescape_string(str) {
				return str.replace(/&(lt|gt|apos|#39|amp|quot|#34);/g, function(c){return {'&lt;':'<','&gt;':'>','&apos;':'\'','&#39;':'\'','&amp;':'&','&quot;':'"','&#34;':'"'}[c];});
			}

			function escape_string(str) {
				return str.replace(/[<>&"']/g, function(c){return {'<':'&lt;','>':'&gt;','&':'&amp;','"':'&quot;','\'':'&apos;'}[c];});
			}

			function output(s) {
				try {
					return escape_string(JSON.stringify(JSON.parse(unescape_string(s)), null, 2));
				} catch(e) {
					return s;
				}
			}
			document.write("<pre style=\"display:inline-block;text-align:left;\">"+output(%s)+"</pre>")
		</script>
		`, strconv.Quote(html.EscapeString(string(body)))))
	} else {
		t = fmt.Sprintf(t, "")
	}

	io.WriteString(w, t)

	return
}

func IsAInB(a string, b string) (ok bool) {
	var vs []json.RawMessage
	if err := json.Unmarshal([]byte(b), &vs); err != nil {
		return
	}

	for _, v := range vs {
		if a == string(v) {
			return true
		}
	}

	return
}

type FieldValueDesc struct {
	Value string
	Desc  string
}

func ParseFieldParam(param string) (valueDescs []*FieldValueDesc, err error) {
	var vs1 []map[string]json.RawMessage
	if err = json.Unmarshal([]byte(param), &vs1); err == nil && len(vs1) > 0 {
		for _, v := range vs1 {
			var desc string
			var value json.RawMessage
			for desc, value = range v {
				break
			}

			if len(value) == 0 || desc == "" {
				break
			}

			valueDescs = append(valueDescs, &FieldValueDesc{
				Value: string(value),
				Desc:  desc,
			})
		}

		if len(vs1) == len(valueDescs) {
			return
		}

		valueDescs = nil
	}

	var vs2 []json.RawMessage
	if err = json.Unmarshal([]byte(param), &vs2); err == nil {
		for _, v := range vs2 {
			valueDescs = append(valueDescs, &FieldValueDesc{
				Value: string(v),
				Desc:  string(v),
			})
		}
		return
	}

	err = errors.New("parse field param err")
	return
}

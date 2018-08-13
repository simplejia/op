package filter

import (
	"io"
	"net/http"
	"strconv"

	"lib"

	"github.com/simplejia/op/conf"
)

// Auth 前置过滤器，用于登陆态校验，权限校验
func Auth(w http.ResponseWriter, r *http.Request, m map[string]interface{}) (ok bool) {
	if conf.Env == lib.DEV {
		c := m["__C__"].(lib.IBase)
		header := &lib.Header{
			ID:    1,
			Token: "",
		}
		c.SetParam(lib.KeyHeader, header)
		return true
	}

	defer func() {
		if ok {
			return
		}

		io.WriteString(w, `
		<html><head><meta charset="UTF-8"></head><body style="text-align:center;">
		<a href="/" style="display:block;">返回首页</a>
		<form action="/auth/sign_in" method="post">
			<table style="display:inline-block;text-align:left;">
				<tr><td>用户: </td><td><input type="text" name="usr"/></td></tr>
				<tr><td>密码: </td><td><input type="text" name="pwd"/></td></tr>
			</table>
			<p><button type="submit">执行</button></p>
		</form>
		</body></html>`)
	}()

	cookieID, _ := r.Cookie("h_op_id")
	if cookieID == nil || cookieID.Value == "" {
		return false
	}
	id, _ := strconv.ParseInt(cookieID.Value, 10, 64)

	cookieToken, _ := r.Cookie("h_op_token")
	if cookieToken == nil || cookieToken.Value == "" {
		return false
	}
	token := cookieToken.Value

	// TODO: check id & token

	c := m["__C__"].(lib.IBase)
	header := &lib.Header{
		ID:    id,
		Token: token,
	}
	c.SetParam(lib.KeyHeader, header)

	return true
}

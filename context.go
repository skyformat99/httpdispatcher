package httpdispatcher

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//上下文
type Content struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	routerParams   []httprouter.Param
}

func (ctx *Content) init() error {
	return ctx.Request.ParseForm()
}

//GET参数值
func (ctx *Content) QueryValue(key string) string {
	return ctx.Request.Form.Get(key)
}

//POST参数值
func (ctx *Content) FormValue(key string) string {
	value := ctx.Request.PostFormValue(key)
	if value == "" {
		if ctx.Request.ParseMultipartForm(32<<20) != nil {
			return ""
		}
		if ctx.Request.MultipartForm != nil {
			values := ctx.Request.MultipartForm.Value[key]
			if len(values) > 0 {
				value = values[0]
			}
		}
	}
	return value
}

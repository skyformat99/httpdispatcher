package httpdispatcher

import (
	"net/http"
)

//Handler 处理器类型
type Handler func(*Context) error

//500(panic)事件处理
func (d *Dispatcher) panicErrorHandle(resp http.ResponseWriter, req *http.Request, message interface{}) {
	if message != nil && d.EventConfig.ServerError == true {
		d.logger(message, req, 6)
	}
	//如果定义了500事件处理器
	if d.Handler.ServerError != nil {
		//初始化ctx
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		ctx.ctxParams = make(map[string]interface{})
		//执行处理器
		d.Handler.ServerError(&ctx)
	}
}

//404事件处理
func (d *Dispatcher) notFoundHandle(resp http.ResponseWriter, req *http.Request) {
	//如果开启了404事件记录
	if d.EventConfig.NotFound == true {
		//记录事件
		d.logger(http.StatusText(404), req, -1)
	}
	//如果定义了404事件处理器
	if d.Handler.NotFound != nil {
		//初始化ctx
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		ctx.ctxParams = make(map[string]interface{})
		//执行处理器
		d.Handler.NotFound(&ctx)
	}
}

//405事件处理
func (d *Dispatcher) methodNotAllowedHandle(resp http.ResponseWriter, req *http.Request) {
	//如果定义了405事件记录
	if d.EventConfig.MethodNotAllowed == true {
		//记录事件
		d.logger(http.StatusText(405), req, -1)
	}
	//如果定义了405事件处理器
	if d.Handler.MethodNotAllowed != nil {
		//初始化ctx
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		ctx.ctxParams = make(map[string]interface{})
		//执行处理器
		d.Handler.MethodNotAllowed(&ctx)
	}
}

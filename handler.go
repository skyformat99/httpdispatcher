package httpdispatcher

import (
	"net/http"
	"reflect"
)

//处理器类型
type Handler func(*Content) error

//500事件处理
func (d *Dispatcher) serverErrorHandle(resp http.ResponseWriter, req *http.Request, message interface{}) {
	if message != nil && d.EventConfig.ServerError == true {
		//判断消息类型并记录事件
		messageType := reflect.TypeOf(message).String()
		switch messageType {
		case "string":
			d.logger(message.(string), req.RequestURI, 6)
		case "*errors.errorString":
			d.logger(message.(error).Error(), req.RequestURI, 6)
		default:
			d.logger("无法转换消息变量的类型("+messageType+")", req.RequestURI, 6)
		}
	}
	//如果定义了500事件处理器
	if d.Handler.ServerError != nil {
		//初始化ctx
		var ctx Content
		ctx.Request = req
		ctx.ResponseWriter = resp
		if err := ctx.init(); err != nil {
			d.logger(err.Error(), req.RequestURI, 6)
			return
		}
		//执行处理器
		d.Handler.ServerError(&ctx)
	}
}

//404事件处理
func (d *Dispatcher) notFoundHandle(resp http.ResponseWriter, req *http.Request) {
	//如果开启了404事件记录
	if d.EventConfig.NotFound == true {
		//记录事件
		d.logger(http.StatusText(404), req.RequestURI, -1)
	}
	//如果定义了404事件处理器
	if d.Handler.NotFound != nil {
		//初始化ctx
		var ctx Content
		ctx.Request = req
		ctx.ResponseWriter = resp
		if err := ctx.init(); err != nil {
			d.logger(err.Error(), req.RequestURI, 6)
			return
		}
		//执行处理器
		d.Handler.NotFound(&ctx)
	}
}

//405事件处理
func (d *Dispatcher) methodNotAllowedHandle(resp http.ResponseWriter, req *http.Request) {
	//如果定义了405事件记录
	if d.EventConfig.MethodNotAllowed == true {
		//记录事件
		d.loggerURL(req.RequestURI, req.Method, http.StatusText(405))
	}
	//如果定义了405事件处理器
	if d.Handler.MethodNotAllowed != nil {
		//初始化ctx
		var ctx Content
		ctx.Request = req
		ctx.ResponseWriter = resp
		if err := ctx.init(); err != nil {
			d.logger(err.Error(), req.RequestURI, 6)
			return
		}

		//执行处理器
		d.Handler.MethodNotAllowed(&ctx)
	}
}

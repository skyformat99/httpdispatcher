package httpdispatcher

import (
	"errors"
	"net/http"
	"runtime"
	"strconv"
)

//EventHandler 处理器类型
type Handler func(*Context) error

//500(panic)事件处理
func (d *Dispatcher) panicErrorHandle(resp http.ResponseWriter, req *http.Request, err interface{}) {
	//如果定义了500事件处理器
	if d.EventHandler.ServerError != nil {
		//初始化ctx
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		ctx.ctxParams = make(map[string]interface{})

		var event Event
		if errStr, ok := err.(string); ok == true {
			event.Message = errors.New(errStr)
		} else if errErr, ok := err.(error); ok == true {
			event.Message = errErr
		} else {
			event.Message = errors.New("未知的错误消息")
		}
		if d.EventConfig.EnableCaller == true {
			var trace []string
			for i := 0; ; i++ {
				var f, l string
				_, file, line, ok := runtime.Caller(i)
				l = strconv.Itoa(line)
				//if d.EventConfig.ShortCaller == true {
				//	short := file
				//	fileLen := len(file)
				//	for i := fileLen - 1; i > 0; i-- {
				//		if file[i] == '/' {
				//			short = file[i+1:]
				//			break
				//		}
				//	}
				//	file = short
				//}
				f = file
				trace = append(trace, f+":"+l)
				if ok == true {
					break
				}
			}
			event.Trace = trace
		}
		//将事件写入到ContextValue中
		ctx.SetContextValue("_event", event)

		//执行处理器
		d.EventHandler.ServerError(&ctx)
	}
}

//404事件处理
func (d *Dispatcher) notFoundHandle(resp http.ResponseWriter, req *http.Request) {
	//如果定义了404事件处理器
	if d.EventHandler.NotFound != nil {
		//初始化ctx
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		ctx.ctxParams = make(map[string]interface{})
		//执行处理器
		d.EventHandler.NotFound(&ctx)
	}
}

//405事件处理
func (d *Dispatcher) methodNotAllowedHandle(resp http.ResponseWriter, req *http.Request) {
	//如果定义了405事件处理器
	if d.EventHandler.MethodNotAllowed != nil {
		//初始化ctx
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		ctx.ctxParams = make(map[string]interface{})
		//执行处理器
		d.EventHandler.MethodNotAllowed(&ctx)
	}
}

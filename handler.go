package httpdispatcher

import (
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

// Handler 处理器类型
type Handler func(*Context) error

// 500(panic)事件处理
func (d *Dispatcher) panicErrorHandle(resp http.ResponseWriter, req *http.Request, err interface{}) {
	if d.EventHandler.ServerError != nil {
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp

		var event Event
		if errStr, ok := err.(string); ok == true {
			event.Message = errors.New(errStr)
		} else if errErr, ok := err.(error); ok == true {
			event.Message = errErr
		} else {
			event.Message = errors.New("未知错误")
		}
		if d.EventConfig.EnableTrace == true {
			goRoot := runtime.GOROOT()
			for skip := 0; ; skip++ {
				_, file, line, ok := runtime.Caller(skip)
				if strings.HasPrefix(file, goRoot) == false {
					event.Trace = append(event.Trace, file+":"+strconv.Itoa(line))
				}
				if ok == false {
					break
				}
			}
		}
		// 将事件写入到ContextValue中
		ctx.SetContextValue("_event", event)

		// 执行处理器
		d.EventHandler.ServerError(&ctx)
	}
}

// 404事件处理
func (d *Dispatcher) notFoundHandle(resp http.ResponseWriter, req *http.Request) {
	if d.EventHandler.NotFound != nil {
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		d.EventHandler.NotFound(&ctx)
	}
}

// 405事件处理
func (d *Dispatcher) methodNotAllowedHandle(resp http.ResponseWriter, req *http.Request) {
	if d.EventHandler.MethodNotAllowed != nil {
		var ctx Context
		ctx.Request = req
		ctx.ResponseWriter = resp
		d.EventHandler.MethodNotAllowed(&ctx)
	}
}

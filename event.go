package httpdispatcher

import (
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

// 事件结构
type Event struct {
	Status  int      // HTTP状态码
	Trace   []string // 触发事件的trace
	Message error    // 事件消息文本
	Resp    http.ResponseWriter
	Req     *http.Request
}

// EventHandler 事件处理器类型
type EventHandler func(*Event)

// 500(panic)事件处理
func (d *Dispatcher) handle500(resp http.ResponseWriter, req *http.Request, err interface{}) {
	if d.Event.Handler != nil {
		var event Event
		event.Req = req
		event.Resp = resp
		event.Status = 500

		if errStr, ok := err.(string); ok == true {
			event.Message = errors.New(errStr)
		} else if errErr, ok := err.(error); ok == true {
			event.Message = errErr
		} else {
			event.Message = errors.New("未知错误")
		}
		if d.Event.EnableTrace == true {
			goRoot := runtime.GOROOT()
			for skip := 0; ; skip++ {
				_, file, line, ok := runtime.Caller(skip)
				// 排除trace中的标准包信息
				if strings.HasPrefix(file, goRoot) == false {
					// if d.Event.ShortCaller == true {
					// 	short := file
					// 	fileLen := len(file)
					// 	for i := fileLen - 1; i > 0; i-- {
					// 		if file[i] == '/' {
					// 			short = file[i+1:]
					// 			break
					// 		}
					// 	}
					// 	file = short
					// }
					event.Trace = append(event.Trace, file+":"+strconv.Itoa(line))
				}
				if ok == false {
					break
				}
			}
		}

		d.Event.Handler(&event)
	}
}

// 404事件处理
func (d *Dispatcher) handle404(resp http.ResponseWriter, req *http.Request) {
	if d.Event.Handler != nil {
		var event Event
		event.Status = 404
		event.Message = errors.New(http.StatusText(404))
		event.Req = req
		event.Resp = resp
		d.Event.Handler(&event)
	}
}

// 405事件处理
func (d *Dispatcher) handle405(resp http.ResponseWriter, req *http.Request) {
	if d.Event.Handler != nil {
		var event Event
		event.Status = 405
		event.Message = errors.New(http.StatusText(405))
		event.Req = req
		event.Resp = resp
		d.Event.Handler(&event)
	}
}

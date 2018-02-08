package httpdispatcher

import (
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

//事件结构
type Event struct {
	Source   string      //源码文件及行号
	URI      string      //客户端请求的URI
	Method   string      //客户端请求的方法
	ClientIP string      //客户端的IP
	Message  interface{} //消息
}

//事件处理器类型
type EventHandler func(*Event)

//事件记录器
func (d *Dispatcher) logger(message interface{}, req *http.Request, skip int) {
	//如果没有指定接收事件的处理器，则直接退出函数
	if d.Handler.Event == nil {
		return
	}

	var event Event
	event.Message = message
	if req != nil {
		event.URI = req.RequestURI
		event.Method = req.Method
		event.ClientIP = req.RemoteAddr
	}

	if skip >= 0 && d.EventConfig.EnableCaller == true {
		var file string
		var line int

		_, file, line, _ = runtime.Caller(skip)
		if strings.HasSuffix(file, "net/http/server.go") == false && d.EventConfig.ShortCaller == true {
			short := file
			fileLen := len(file)
			for i := fileLen - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			event.Source = file + ":" + strconv.Itoa(line)
		}
	}
	d.Handler.Event(&event)
}

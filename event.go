package httpdispatcher

import (
	"runtime"
	"strconv"
)

//事件结构
type Event struct {
	Source  string //源码文件及行号
	URI string	//客户端请求的URI
	Message string //消息
}

//事件处理器类型
type EventHandler func(*Event)

//事件记录器
func (d *dispatcher) logger(message, uri string, skip int) {
	if d.Handler.Event == nil {
		return
	}

	var event Event
	event.Message = message
	event.URI = uri

	if d.EventConfig.EnableCaller == true {
		var file string
		var line int

		if skip >= 0 {
			//短文件名
			_, file, line, _ = runtime.Caller(skip)
			if d.EventConfig.ShortCaller == true {
				short := file
				fileLen := len(file)
				for i := fileLen - 1; i > 0; i-- {
					if file[i] == '/' {
						short = file[i+1:]
						break
					}
				}
				file = short
			}
			event.Source = file + ":" + strconv.Itoa(line)
		}
	}
	d.Handler.Event(&event)
}

//用于记录URL的事件记录器
func (d *dispatcher) loggerURL(uri, method, message string) {
	if d.Handler.Event == nil {
		return
	}

	var event Event
	event.Message = message
	event.Source = uri + ":" + method
	d.Handler.Event(&event)
}

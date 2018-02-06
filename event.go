package httpdispatcher

import (
	"runtime"
	"strconv"
)

//事件结构
type Event struct {
	Source  string      //源码文件及行号
	URI     string      //客户端请求的URI
	Message interface{} //消息
}

//事件处理器类型
type EventHandler func(*Event)

//事件记录器
func (d *Dispatcher) logger(message interface{}, uri string, skip int) {
	//如果没有指定接收事件的处理器，则直接退出函数
	if d.Handler.Event == nil {
		return
	}

	var event Event
	event.Message = message
	event.URI = uri

	if skip >= 0 && d.EventConfig.EnableCaller == true {
		var file string
		var line int

		_, file, line, _ = runtime.Caller(skip)
		//如果要求记录短文件名
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
		if line == 0 {
			event.Source = ""
		} else {
			event.Source = file + ":" + strconv.Itoa(line)
		}
	}
	d.Handler.Event(&event)
}

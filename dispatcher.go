package httpdispatcher

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type dispatcher struct {
	httpRouter *httprouter.Router
	Router     *RouterGroup
	Handler    struct {
		NotFound         Handler      //404错误处理器
		MethodNotAllowed Handler      //405错误处理器
		ServerError      Handler      //500错误处理器
		Event            EventHandler //事件处理器
	}
	EventConfig struct {
		EnableCaller     bool //启用来源记录
		ShortCaller      bool //缩短来源记录(仅记录源码文件名)
		NotFound         bool //记录404错误事件
		MethodNotAllowed bool //记录405错误事件
		ServerError      bool //记录500错误事件
	}
}

func New() *dispatcher {
	var d dispatcher
	//实例化httprouter路由
	d.httpRouter = httprouter.New()
	//指定各种错误处理
	d.httpRouter.PanicHandler = d.serverErrorHandle                            //panic错误处理器
	d.httpRouter.NotFound = http.HandlerFunc(d.notFoundHandle)                 //404错误处理器
	d.httpRouter.MethodNotAllowed = http.HandlerFunc(d.methodNotAllowedHandle) //405错误处理器
	d.Router = &RouterGroup{
		d: &d,
	}

	return &d
}

func (d *dispatcher) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	d.httpRouter.ServeHTTP(resp, req)
}

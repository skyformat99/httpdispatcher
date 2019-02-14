package httpdispatcher

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 事件
type Event struct {
	Trace   []string
	Message error
}

// Dispatcher 调度器结构
type Dispatcher struct {
	// 事件记录配置
	EventConfig struct {
		EnableTrace bool // 启用500事件的跟踪(影响性能)
		ShortCaller bool // 缩短事件触发的源码文件名(仅记录源码文件名，仅对ctx.Return触发的500事件有效)s
	}
	// 事件处理事
	EventHandler struct {
		NotFound         Handler // 404错误处理器
		MethodNotAllowed Handler // 405错误处理器
		ServerError      Handler // 500错误处理器
	}
	// 路由器
	Router *RouterGroup
	// httprouter的路由器
	httpRouter *httprouter.Router
}

// New 返回一个调度器实例
func New() *Dispatcher {
	var dispatcher Dispatcher
	dispatcher.httpRouter = httprouter.New()
	dispatcher.httpRouter.PanicHandler = dispatcher.panicErrorHandle
	dispatcher.httpRouter.NotFound = http.HandlerFunc(dispatcher.notFoundHandle)
	dispatcher.httpRouter.MethodNotAllowed = http.HandlerFunc(dispatcher.methodNotAllowedHandle)

	dispatcher.Router = &RouterGroup{
		dispatcher: &dispatcher,
	}

	return &dispatcher
}

func (d *Dispatcher) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	d.httpRouter.ServeHTTP(resp, req)
}

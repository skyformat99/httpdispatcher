package httpdispatcher

import (
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

//路由组
type RouterGroup struct {
	handlers []Handler   //处理器
	basePath string      //基路径
	d        *Dispatcher //调度器
}

//定义路由到目录，不支持路由组和中间件
func (r *RouterGroup) PATH(url string, local string, list bool) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 6)
			os.Exit(1)
		}
	}()
	if strings.HasPrefix(url, "/") == false {
		url = "/" + url
	}
	if strings.HasSuffix(url, "/") == false {
		url += "/"
	}
	url += "*filepath"

	//使用GET方法模拟httprouter.ServeFiles()，防止其内部直接输出404消息给客户端
	r.d.httpRouter.GET(url, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		//如果请求的是目录，而判断是否允许列出目录
		if params.ByName("filepath") == "" || params.ByName("filepath")[len(params.ByName("filepath"))-1:] == "/" {
			if list == false {
				//如果不允许列出目录，则触发404事件处理
				r.d.notFoundHandle(resp, req)
				return
			}
		}

		//判断请求的文件是否存在
		file := local + params.ByName("filepath")
		_, err := os.Stat(file)
		if err != nil {
			//404事件处理
			r.d.notFoundHandle(resp, req)
			return
		}
		http.ServeFile(resp, req, file)
	})
}

//定义路由到文件，不支持路由组和中间件
func (r *RouterGroup) FILE(url string, local string) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 6)
			os.Exit(1)
		}
	}()
	//使用GET方法模拟httprouter.ServeFiles()，防止其内部直接输出404消息给客户端
	r.d.httpRouter.GET(url, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		_, err := os.Stat(local)
		if err != nil {
			//404事件处理
			r.d.notFoundHandle(resp, req)
			return
		}
		http.ServeFile(resp, req, local)
	})
}

//定义路由组
func (r *RouterGroup) GROUP(path string, handlers ...Handler) *RouterGroup {
	//声明一个新的路由组
	var group RouterGroup
	group.basePath = r.basePath + path  //继承父组的路径
	group.d = r.d                       //传入调度器
	group.handlers = append(r.handlers) //继承父组的钩子
	//加入当前传入的钩子
	for k := range handlers {
		group.handlers = append(group.handlers, handlers[k])
	}
	return &group
}

//定义GET方法的路由
func (r *RouterGroup) GET(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.GET(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//定义POST方法的路由
func (r *RouterGroup) POST(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.POST(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//定义PUT方法的路由
func (r *RouterGroup) PUT(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.PUT(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//定义HEAD方法的路由
func (r *RouterGroup) HEAD(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.HEAD(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//定义PATCH方法的路由
func (r *RouterGroup) PATCH(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.PATCH(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//定义DELETE方法的路由
func (r *RouterGroup) DELETE(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.DELETE(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//定义OPTIONS方法的路由
func (r *RouterGroup) OPTIONS(path string, handler Handler, handlers ...Handler) {
	defer func() {
		if err := recover(); err != nil {
			//记录panic事件，但不执行 ServerError处理器，而是直接退出进程
			r.d.logger(err.(string), "", 8)
			os.Exit(1)
		}
	}()
	r.d.httpRouter.OPTIONS(r.basePath+path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r.execute(resp, req, params, handler, handlers)
	})
}

//执行Handler
func (r *RouterGroup) execute(resp http.ResponseWriter, req *http.Request, params httprouter.Params, handler Handler, handlers []Handler) {
	//初始化ctx
	var ctx Content
	ctx.Request = req
	ctx.ResponseWriter = resp
	ctx.dispatcher = r.d
	if err := ctx.init(); err != nil {
		//触发500事件
		r.d.panicErrorHandle(resp, req, err.Error())
		return
	}
	//遍历父路由的中间件处理器
	for k := range r.handlers {
		//执行父路由的中间件处理器
		err := r.handlers[k](&ctx)
		if err != nil {
			//触发500事件
			r.d.panicErrorHandle(resp, req, err.Error())
			return
		}
		//如果控制器执行完之后ctx的next属性值为false，则不继续循环执行下一个处理器而是退出整个函数
		if ctx.next == false {
			return
		}
	}
	//遍历刚传入的中间件处理器
	for k := range handlers {
		//执行中间件处理器
		err := handlers[k](&ctx)
		if err != nil {
			//触发500事件
			r.d.panicErrorHandle(resp, req, err.Error())
			return
		}
		//如果控制器执行完之后ctx的next属性值为false，则不继续循环执行下一个处理器而是退出整个函数
		if ctx.next == false {
			return
		}
	}

	//执行处理器
	err := handler(&ctx)
	if err != nil {
		//触发500事件
		r.d.panicErrorHandle(resp, req, err.Error())
	}
}

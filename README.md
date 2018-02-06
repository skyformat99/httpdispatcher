# http dispatcher
Go语言基于[http router](https://github.com/julienschmidt/httprouter)包实现的轻量HTTP调度器，没有对http router包做任何修改，仅封装实现了以下功能：
- [x] 无限层级的路由组
- [x] 路由中间件（钩子）
- [x] Context（会话上下文）
- [x] 使用Context在同一会话的处理器间传递变量（`SetContextValue` / `ContextValue`）
- [x] 使用`PATH()`和`FILE()`替代`httprouter.ServeFiles()`方法，改进如下：
    * 可禁止列举出目录下的所有文件
    * 由`Dispacher.Handler.NotFoundHandler`来处理404，调度器本身不会主动输出任何消息给客户端
  

## 基本示例
``` Go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gitee.com/dxvgef/httpdispatcher"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//获得一个调度器实例
	dispacher := httpdispatcher.New()

	//定义处理器
	//404错误处理器
	dispacher.Handler.NotFound = func(ctx *httpdispatcher.Content) error {
		log.Println("404事件后续自行处理")
		return nil
	}
	//405错误处理器
	dispacher.Handler.MethodNotAllowed = func(ctx *httpdispatcher.Content) error {
		log.Println("405事件后续自行处理")
		return nil
	}
	//500错误处理器
	dispacher.Handler.ServerError = func(ctx *httpdispatcher.Content) error {
		log.Println("500事件后续自行处理")
		return nil
	}
	//事件处理器
	dispacher.Handler.Event = func(e *httpdispatcher.Event) {
		log.Println("事件来源:", e.Source)
		log.Println("事件消息:", e.Message)
		log.Println("事件URI:", e.URI)
	}

	//事件配置
	dispacher.EventConfig.EnableCaller = true     //开启记录触发事件的源文件及行号
	dispacher.EventConfig.NotFound = true         //记录404事件
	dispacher.EventConfig.ServerError = true      //记录500事件
	dispacher.EventConfig.MethodNotAllowed = true //记录405事件

	//静态路由
	{
		//定义路由到目录，不支持路由组和中间件
        //如果第三个参数为false，在直接访问目录时会当做404处理，而不是列出目录下的所有文件
		dispacher.Router.PATH("/static", "./static", false)
		//定义静态路由到文件，不支持路由组和中间件
		dispacher.Router.FILE("/logo.png", "./logo.png")
	}

	//普通路由
	{
		//定义GET路由，handler为路由处理器，hookHandler为钩子(中间件)处理器
		dispacher.Router.GET("/", handler, hookHandler)
		//测试处理器中出现panic
		dispacher.Router.GET("/panic", testPanic)
	}

	//路由组
	{
		//定义路由组，并传入中间件
		authRouter := dispacher.Router.GROUP("/secret", hookHandler)
		//在路由组下面定义一个POST路由，匹配/secret/
		authRouter.POST("/", handler)
	}

	//定义HTTP服务
	svr := &http.Server{
		Addr:         ":8080",
		Handler:      dispacher, //传入调度器
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	//在新协程中启动服务，方便实现退出等待
	go func() {
		if err := svr.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	//Graceful Shutdown(退出时等待10秒的时间让已连接的逻辑都处理完成)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //指定退出超时时间
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Fatal(err.Error())
	}
}

//路由处理器
func handler(ctx *httpdispatcher.Content) error {
	//读取上一个处理器存储在ctx里的变量
	log.Println(ctx.ContextValue("ctx"))
	//读取GET参数值
	log.Println(ctx.QueryValue("get").String())
	//读取POST参数值，在转换为int类型时如果出错则用123默认值返回
	log.Println(ctx.FormValue("post").Int(123))
	return nil
}

//测试触发panic的路由处理器
func testPanic(ctx *httpdispatcher.Content) error {
	log.Println("抛出panic前的逻辑")
	panic("panic消息")
	log.Println("抛出panic后的逻辑")
	return nil
}

//钩子(中间件)处理器
func hookHandler(ctx *httpdispatcher.Content) error {
	//return errors.New("如果函数返回值不是nil，不会继续执行后面的handler")

	//在ctx中写入变量传递到下一个处理器
	ctx.SetContextValue("test", "ok")
	return nil
}
```
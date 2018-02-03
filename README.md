# httpdispatcher
基于golang的http router包实现的http调度器

## 基本示例
``` Go
package main

import (
	"log"
	"net/http"

	"gitee.com/dxvgef/httpdispatcher"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//获得一个调度器实例
	dispacher := httpdispatcher.New()

	//定义404错误处理器
	dispacher.Handler.NotFound = func(ctx *httpdispatcher.Content) {
		log.Println("处理404错误")
	}
	//定义405错误处理器
	dispacher.Handler.MethodNotAllowed = func(ctx *httpdispatcher.Content) {
		log.Println("处理405错误")
	}
	//定义500错误处理器
	dispacher.Handler.ServerError = func(ctx *httpdispatcher.Content) {
		log.Println("处理500错误")
	}
	//定义事件处理器
	dispacher.Handler.Event = func(e *httpdispatcher.Event) {
		log.Println("事件来源:", e.Source)
		log.Println("事件消息:", e.Message)
	}

	//事件配置
	dispacher.EventConfig.NotFound = true         //记录404事件
	dispacher.EventConfig.ServerError = true      //记录500事件
	dispacher.EventConfig.MethodNotAllowed = true //记录405事件

	//定义路由
	dispacher.Router.GET("/", func(ctx *httpdispatcher.Content) {
		log.Println(ctx.Request.URL)
	})

	log.Fatal(http.ListenAndServe(":8080", dispacher))
}

```
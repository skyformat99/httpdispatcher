# HTTP Dispatcher
使用golang基于[http router](https://github.com/julienschmidt/httprouter)路由包实现的轻量HTTP调度器，没有对http router包做任何修改，仅封装实现了以下功能：
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
	"time"``

	"github.com/dxvgef/httpdispatcher"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//获得一个调度器实例
	dispatcher := httpdispatcher.New()
	//事件记录配置
	dispatcher.EventConfig.EnableCaller = true     //开启记录触发事件的源文件及行号(Event.Source的值)
	dispatcher.EventConfig.NotFound = true         //记录404事件
	dispatcher.EventConfig.ServerError = true      //记录500事件
	dispatcher.EventConfig.MethodNotAllowed = true //记录405事件
	//定义接收事件的处理器
	dispatcher.Handler.Event = func(e *httpdispatcher.Event) {
		log.Println("事件来源:", e.Source)
		log.Println("事件消息:", e.Message)
		log.Println("事件URI:", e.URI)
	}
	//定义404事件处理器
	dispatcher.Handler.NotFound = func(ctx *httpdispatcher.Context) error {
		log.Println("404事件后续自行处理")
		return nil
	}
	//定义405事件处理器
	dispatcher.Handler.MethodNotAllowed = func(ctx *httpdispatcher.Context) error {
		log.Println("405事件后续自行处理")
		return nil
	}
	//定义500事件处理器
	dispatcher.Handler.ServerError = func(ctx *httpdispatcher.Context) error {
		log.Println("500事件后续自行处理")
		return nil
	}

	//定义静态路由
	{
		//定义路由到目录，不支持路由组和中间件
		//如果第三个参数为false，在直接访问目录时会当做404处理，而不是列出目录下的所有文件
		dispatcher.Router.PATH("/static", "./static", false)
		//定义静态路由到文件，不支持路由组和中间件
		dispatcher.Router.FILE("/logo.png", "./logo.png")
	}

	//普通路由
	{
		//定义GET路由，handler为路由处理器，hookHandler为中间件(钩子)处理器
		dispatcher.Router.GET("/", handler, hookHandler)
		//测试处理器中出现panic
		dispatcher.Router.GET("/panic", testPanic)
		//定义一个重定向路由处理器
		dispatcher.Router.GET("/redir", func(ctx *httpdispatcher.Context) error {
			return ctx.Redirect(302, "http://github.com/dxvgef/httpdispatcher")
		})
	}

	//路由组（可无限嵌套）
	{
		//定义路由组，并定义中间件处理器
		//如果定义了路由组的中间件处理器，则组下的路由处理器执行前都会先执行组的中间件处理器
		authRouter := dispatcher.Router.GROUP("/secret", hookHandler)
		//在路由组下面定义一个POST路由
		//此路由匹配URL：/secret/test
		authRouter.POST("/test", handler)
	}

	//定义HTTP服务
	svr := &http.Server{
		Addr:         ":8080",
		Handler:      dispatcher, //传入调度器
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

//普通的路由处理器，会在中间件处理器全部执行完成之后才最后执行
func handler(ctx *httpdispatcher.Context) error {
	//定义一个结构体用于json序列化后输出给客户端
	var resp struct {
		Ctx  string
		Get  string
		Post int
	}
	//读取上一个处理器存储在ctx里的变量
	resp.Ctx = ctx.ContextValue("ctx").(string)
	//读取GET参数值并转为string类型
	resp.Get = ctx.QueryValue("get").String()
	//读取POST参数值，在转换为int类型时如果出错则用123默认值返回
	resp.Post, _ = ctx.FormValue("post").Int(123)

	//序列化json
	b, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	//输出json
	ctx.ResponseWriter.WriteHeader(200)
	ctx.ResponseWriter.Write(b)

	//返回nil表示处理器内的业务逻辑都成功完成，返回非nil值会触发500事件
	return nil
}

//测试触发panic的路由处理器
func testPanic(ctx *httpdispatcher.Context) error {
	log.Println("抛出panic前的逻辑")
	panic("panic消息")
	log.Println("抛出panic后的逻辑")

	//如果return的值不是nil，会触发500事件
	//return errors.New("出错了")
	return nil
}

//中间件(钩子)的处理器，优先于普通处理器执行
//多个中间件处理器的执行顺序与传入的顺序相同
func hookHandler(ctx *httpdispatcher.Context) error {
	//在ctx中写入变量传递到下一个处理器
	ctx.SetContextValue("ctx", "ok")

	//执行此函数并且入参值为true，才可继续执行下一个中间件或者最终的处理器
	//此函数在最终的处理器中执行没有任何意义，仅在中间件处理器中有效
	ctx.Next(true)
	//return ctx.Next(true)

	//如果return的error不是nil也不会继续往下执行别的中间件或者最终的处理器，还会触发500事件
	//return errors.New("出错了")
	return nil
}
```

## Logger
可与[github.com/uber-go/zap](https://github.com/uber-go/zap)包整合实现日志记录功能

## Session
可与[github.com/dxvgef/sessions](https://github.com/dxvgef/sessions)包整合实现Session功能，请进入该项目查看示例代码

## 模板引擎
可与[github.com/CloudyKit/jet](https://github.com/CloudyKit/jet)包整合实现模板渲染功能
``` Go
package main

import (
	"log"
	"net/http"

	"github.com/dxvgef/httpdispatcher"

	"github.com/CloudyKit/jet"
)

//定义render对象
type render struct {
	jet *jet.Set
}

var Render render

//执行模板字符串渲染
func (r *render) ExecuteString(resp http.ResponseWriter, code int, tmpl string, vars jet.VarMap) error {
	//解释模板字符串
	t, err := r.jet.Parse("c", tmpl)
	if err != nil {
		return err
	}

	//设置http状态码
	resp.WriteHeader(code)

	//执行模板渲染并输出给客户端
	err = t.Execute(resp, vars, nil)
	if err != nil {
		return err
	}

	return nil
}

//执行模板文件渲染
func (r *render) ExecuteFile(resp http.ResponseWriter, code int, tmpl string, vars jet.VarMap) error {
	//加载模板文件
	t, err := global.Render.GetTemplate(tmpl)
	if err != nil {
		return err
	}

    //为防止出现出现http: multiple response.WriteHeader calls错误
    //将渲染结果赋值给w变量，而不是直接使用resp输出给客户端
	w := bytes.NewBuffer([]byte{})
	err = t.Execute(w, data, nil)
	if err != nil {
		return err
	}
	
	//将w输出给客户端
	ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=UTF-8")
	ctx.ResponseWriter.WriteHeader(code)
	ctx.ResponseWriter.Write(w.Bytes())

	//将bytes.Buffer清空减少内存占用
	w = bytes.NewBuffer([]byte{})

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)

	//设置Jet
	//存放模板文件的路径
	Render.jet = jet.NewHTMLSet("./templates")
	//开发模式下每次请求时都重新解析模板，否则直接从缓存读取
	Render.jet.SetDevelopmentMode(true)

	dispatcher := httpdispatcher.New()
	dispatcher.EventConfig.EnableCaller = true
	dispatcher.EventConfig.NotFound = true
	dispatcher.EventConfig.ServerError = true
	dispatcher.Handler.Event = func(e *httpdispatcher.Event) {
		log.Println("事件来源:", e.Source)
		log.Println("事件消息:", e.Message)
		log.Println("事件URI:", e.URI)
	}
	dispatcher.Handler.NotFound = func(ctx *httpdispatcher.Context) error {
		log.Println("404事件后续自行处理")
		return nil
	}
	dispatcher.Handler.ServerError = func(ctx *httpdispatcher.Context) error {
		log.Println("500事件后续自行处理")
		return nil
	}

	dispatcher.Router.GET("/", func(ctx *httpdispatcher.Context) error {
		//声明模板变量
		vars := make(jet.VarMap)
		//设置模板变量
		vars.Set("test", "ok")
		//渲染模板字符串
		return Render.ExecuteString(ctx.ResponseWriter, 200, "<div>{{test}}</div>", vars)
		//渲染模板文件
		//return Render.ExecuteFile(ctx.ResponseWriter, 200, "index.html", vars)
	})

	if err := http.ListenAndServe(":8080", dispatcher); err != nil {
		log.Fatal(err.Error())
	}
}
```

## Benchmark
#### 路由注册Benchmark代码
``` Go
func BenchmarkTest(b *testing.B) {
    b.ResetTimer()
    d := httpdispatcher.New()
    for i := 0; i < b.N; i++ {
        d.Router.GET("/"+strconv.Itoa(i), func(ctx *httpdispatcher.Context) error {
        return nil
    })
    
    //e := echo.New()
    //for i := 0; i < b.N; i++ {
    //  e.GET("/"+strconv.Itoa(i), func(ctx echo.Context) error {
    //  return nil
    //  })
    //}
}
```

### 路由注册Benchmark结果
```
goos: darwin
goarch: amd64
-------------------------
echo
1000000       1976 ns/op
-------------------------
httpdispatcher
2000000        760 ns/op
```

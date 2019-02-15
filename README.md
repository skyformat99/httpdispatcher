# HTTP Dispatcher

使用Go语言开发，基于高性能的[julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)路由包实现的HTTP调度器，与`net/http`标准包配合使用。

没有对`httprouter`包做任何修改，仅轻量封装实现了更多功能，同时保留`net/http`标准包的API访问，高度可定制性，非常适合自行整合第三方包进行二次开发扩展功能。

## 手册
* [基本示例](https://github.com/dxvgef/httpdispatcher/wiki/%E5%9F%BA%E6%9C%AC%E7%A4%BA%E4%BE%8B)
* [事件处理 `httpdispatcher.Event`](https://github.com/dxvgef/httpdispatcher/wiki/%E4%BA%8B%E4%BB%B6%E5%9B%9E%E8%B0%83)
* [路由/路由组 `httpdispacher.Router`](https://github.com/dxvgef/httpdispatcher/wiki/%E8%B7%AF%E7%94%B1)
* [处理器 httpdispatcher.Handler](https://github.com/dxvgef/httpdispatcher/wiki/%E5%A4%84%E7%90%86%E5%99%A8)
* [会话上下文 `httpdispatcher.Context`](https://github.com/dxvgef/httpdispatcher/wiki/%E4%BC%9A%E8%AF%9D%E4%B8%8A%E4%B8%8B%E6%96%87)
* [客户端输出 String/JSON](https://github.com/dxvgef/httpdispatcher/wiki/%E5%AE%A2%E6%88%B7%E7%AB%AF%E8%BE%93%E5%87%BA)
* [整合`CloudyKit/jet`实现HTML模板渲染](https://github.com/dxvgef/httpdispatcher/wiki/HTML%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93)
* [CORS跨域资源共享控制](https://github.com/dxvgef/httpdispatcher/wiki/CORS%E8%B7%A8%E5%9F%9F%E8%B5%84%E6%BA%90%E5%85%B1%E4%BA%AB%E6%8E%A7%E5%88%B6)
* [整合`dxvgef/session`实现Session](https://github.com/dxvgef/httpdispatcher/wiki/Session)
* [整合`dxvgef/token`实现JWT(JSON Web Token)](https://github.com/dxvgef/httpdispatcher/wiki/JSON-Web-Token)
* [整合`uber-go/zap`记录事件日志](https://github.com/dxvgef/httpdispatcher/wiki/%E6%97%A5%E5%BF%97%E8%AE%B0%E5%BD%95%E5%99%A8(Logger))
* [优雅关闭(Graceful Shutdown)](https://github.com/dxvgef/httpdispatcher/wiki/%E4%BC%98%E9%9B%85%E5%85%B3%E9%97%AD(Graceful-Shutdown))

## Benchmark
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

package httpdispatcher

import (
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

//处理器类型
type Handler func(*Content)

//钩子处理器类型
type HookHandler func(*Content) error

//将httprouter处理器包装转换成框架内部的处理器，并且执行处理器

func (d *dispatcher) initContext(resp http.ResponseWriter, req *http.Request, params httprouter.Params) (*Content, error) {
	//声明一个ctx并初始化值
	var ctx Content
	ctx.Request = req
	ctx.ResponseWriter = resp
	ctx.routerParams = params

	//执行ctx的初始化函数
	err := ctx.init()
	if err != nil {
		d.logger(err.Error(), 1)
		return nil, err
	}

	return &ctx, nil
}

//执行处理器
func (d *dispatcher) executeHandler(ctx *Content, handler Handler, i interface{}) {
	handler(ctx)
}

//执行钩子处理器
func (d *dispatcher) executeHookHandler(ctx *Content, hookHandler HookHandler) error {
	return hookHandler(ctx)
}

//500错误处理
func (d *dispatcher) serverErrorHandle(resp http.ResponseWriter, req *http.Request, message interface{}) {
	if message != nil && d.EventConfig.ServerError == true {
		messageType := reflect.TypeOf(message).String()
		switch messageType {
		case "string":
			d.logger(message.(string), 6)
		case "*errors.errorString":
			d.logger(message.(error).Error(), 6)
		default:
			d.logger("无法转换消息变量的类型("+messageType+")", 6)
		}
	}
	if d.Handler.ServerError != nil {
		ctx, err := d.initContext(resp, req, nil)
		if err != nil {
			d.logger(err.Error(), 6)
			return
		}
		d.executeHandler(ctx, d.Handler.ServerError, message)
	}
}

//404错误处理
func (d *dispatcher) notFoundHandle(resp http.ResponseWriter, req *http.Request) {
	if d.EventConfig.NotFound == true {
		d.loggerURL(req.RequestURI, req.Method, http.StatusText(404))
	}
	if d.Handler.NotFound != nil {
		ctx, err := d.initContext(resp, req, nil)
		if err != nil {
			d.logger(err.Error(), 6)
			return
		}
		d.executeHandler(ctx, d.Handler.NotFound, nil)
	}
}

//405错误处理
func (d *dispatcher) methodNotAllowedHandle(resp http.ResponseWriter, req *http.Request) {
	if d.EventConfig.MethodNotAllowed == true {
		d.loggerURL(req.RequestURI, req.Method, http.StatusText(405))
	}
	if d.Handler.MethodNotAllowed != nil {
		ctx, err := d.initContext(resp, req, nil)
		if err != nil {
			d.logger(err.Error(), 6)
			return
		}
		d.executeHandler(ctx, d.Handler.MethodNotAllowed, nil)
	}
}

package httpdispatcher

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

//Context 上下文
type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	routerParams   httprouter.Params      //路由参数
	ctxParams      map[string]interface{} //ctx参数
	//cc             context.Context
	next       bool //继续往下执行处理器的标识
	dispatcher *Dispatcher
	parsed     bool //是否已解析body
}

//BodyValue 请求的参数值
type BodyValue struct {
	Key   string //参数名
	Value string //参数值
	Error error  //错误
}

//Next 设置标识，用于继续执行下一个处理器
func (ctx *Context) Next(flag bool) error {
	ctx.next = flag
	return nil
}

//SetContextValue 在ctx里存储值，如果key存在则替换值
func (ctx *Context) SetContextValue(key string, value interface{}) {
	ctx.ctxParams[key] = value
}

//ContextValue 获取ctx里的值，取出后根据写入的类型自行断言
func (ctx *Context) ContextValue(key string) interface{} {
	return ctx.ctxParams[key]
}

//Redirect 重定向
func (ctx *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("状态码只能是300-308之间的值")
	}
	ctx.ResponseWriter.Header().Set("Location", url)
	ctx.ResponseWriter.WriteHeader(code)
	return nil
}

//Return 控制器return error时使用，用于精准记录源码文件及行号
func (ctx *Context) Return(err error) error {
	if err != nil {
		//记录事件
		ctx.dispatcher.logger(err, ctx.Request, 2)
		//如果定义了500事件处理器
		if ctx.dispatcher.Handler.ServerError != nil {
			//执行500处理器
			ctx.dispatcher.Handler.ServerError(ctx)
		}
	}
	return nil
}

//RealIP 获得客户端真实IP
func (ctx *Context) RealIP() string {
	ra := ctx.Request.RemoteAddr
	if ip := ctx.Request.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := ctx.Request.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}

//解析body数据
func (ctx *Context) parseBody() error {
	//判断是否已经解析过body
	if ctx.parsed == true {
		return nil
	}
	//如果是form-data类型
	if strings.HasPrefix(ctx.Request.Header.Get("Content-Type"), "multipart/form-data") {
		//使用ParseMultipartForm解析数据
		if err := ctx.Request.ParseMultipartForm(http.DefaultMaxHeaderBytes); err != nil {
			return err
		}
	} else {
		//否则按x-www-form-urlencoded类型来解析数据
		if err := ctx.Request.ParseForm(); err != nil {
			return err
		}
	}
	//标记该context中的body已经解析过
	ctx.parsed = true
	return nil
}

//RouteValue 获取路由参数值
func (ctx *Context) RouteValue(key string) *BodyValue {
	return &BodyValue{
		Key:   key,
		Value: ctx.routerParams.ByName("key"),
	}
}

//QueryValue 获取某个GET参数值
func (ctx *Context) QueryValue(key string) *BodyValue {
	err := ctx.parseBody()
	if err != nil {
		return &BodyValue{
			Key:   key,
			Error: err,
		}
	}
	return &BodyValue{
		Key:   key,
		Value: ctx.Request.Form.Get(key),
	}
}

//FormValue 获取某个POST参数值
func (ctx *Context) FormValue(key string) *BodyValue {
	err := ctx.parseBody()
	if err != nil {
		return &BodyValue{
			Key:   key,
			Error: err,
		}
	}
	return &BodyValue{
		Key:   key,
		Value: ctx.Request.FormValue(key),
	}
}

//将参数值转为string
func (bv *BodyValue) String() string {
	if bv.Error != nil {
		return ""
	}
	return bv.Value
}

//Int 将参数值转为int类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Int(def ...int) (int, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.Atoi(bv.Value)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

//Int32 将参数值转为int32类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Int32(def ...int32) (int32, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseInt(bv.Value, 10, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return int32(value), nil
}

//Int64 将参数值转为int64类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Int64(def ...int64) (int64, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseInt(bv.Value, 10, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

//Uint32 将参数值转为uint32类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Uint32(def ...uint32) (uint32, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseUint(bv.Value, 10, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return uint32(value), nil
}

//Uint64 将参数值转为uint64类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Uint64(def ...uint64) (uint64, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseUint(bv.Value, 10, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

//Float32 将参数值转为float32类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Float32(def ...float32) (float32, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseFloat(bv.Value, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return float32(value), nil
}

//Float64 将参数值转为float64类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Float64(def ...float64) (float64, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return 0, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseFloat(bv.Value, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

//Bool 将参数值转为bool类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *BodyValue) Bool(def ...bool) (bool, error) {
	defLen := len(def)
	if bv.Error != nil {
		if defLen == 0 {
			return false, bv.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseBool(bv.Value)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return false, err
	}
	return value, nil
}

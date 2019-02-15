package httpdispatcher

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"

	"runtime"

	"github.com/julienschmidt/httprouter"
)

// Context 上下文
type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	routerParams   httprouter.Params // 路由参数
	next           bool              // 继续往下执行处理器的标识
	dispatcher     *Dispatcher
	parsed         bool // 是否已解析body
}

// ReqValue 请求参数值
type ReqValue struct {
	Key   string // 参数名
	Value string // 参数值
	Error error  // 错误
}

// Next 设置标识，用于继续执行下一个处理器
func (ctx *Context) Next(flag bool) error {
	ctx.next = flag
	return nil
}

// SetContextValue 在ctx里存储值，如果key存在则替换值
func (ctx *Context) SetContextValue(key string, value interface{}) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), key, value))
}

// ContextValue 获取ctx里的值，取出后根据写入的类型自行断言
func (ctx *Context) ContextValue(key string) interface{} {
	return ctx.Request.Context().Value(key)
}

// Redirect 重定向
func (ctx *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("状态码只能是300-308之间的值")
	}
	ctx.ResponseWriter.Header().Set("Location", url)
	ctx.ResponseWriter.WriteHeader(code)
	return nil
}

// Event 在控制器return时使用，用于精准记录源码文件及行号
func (ctx *Context) Event(err error) error {
	if err != nil {
		// 如果定义了500事件处理器
		if ctx.dispatcher.Event.Handler != nil {
			var event Event
			event.Status = 500
			event.Message = err
			event.Resp = ctx.ResponseWriter
			event.Req = ctx.Request
			if ctx.dispatcher.Event.EnableTrace == true {
				_, file, line, _ := runtime.Caller(1)
				l := strconv.Itoa(line)
				if ctx.dispatcher.Event.ShortCaller == true {
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
				event.Trace = append(event.Trace, file+":"+l)
			}
			ctx.dispatcher.Event.Handler(&event)
		}
	}
	// 不再将传入的error返回，避免再触发handle500函数
	return nil
}

// RealIP 获得客户端真实IP
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

// 解析body数据
func (ctx *Context) parseBody() error {
	// 判断是否已经解析过body
	if ctx.parsed == true {
		return nil
	}
	if strings.HasPrefix(ctx.Request.Header.Get("Content-Type"), "multipart/form-data") {
		if err := ctx.Request.ParseMultipartForm(http.DefaultMaxHeaderBytes); err != nil {
			return err
		}
	} else {
		if err := ctx.Request.ParseForm(); err != nil {
			return err
		}
	}
	// 标记该context中的body已经解析过
	ctx.parsed = true
	return nil
}

// RouteValue 获取路由参数值
func (ctx *Context) RouteValue(key string) *ReqValue {
	return &ReqValue{
		Key:   key,
		Value: ctx.routerParams.ByName(key),
	}
}

// QueryValue 获取某个GET参数值
func (ctx *Context) QueryValue(key string) *ReqValue {
	err := ctx.parseBody()
	if err != nil {
		return &ReqValue{
			Key:   key,
			Error: err,
		}
	}
	return &ReqValue{
		Key:   key,
		Value: ctx.Request.Form.Get(key),
	}
}

// FormValue 获取某个POST参数值
func (ctx *Context) FormValue(key string) *ReqValue {
	err := ctx.parseBody()
	if err != nil {
		return &ReqValue{
			Key:   key,
			Error: err,
		}
	}
	return &ReqValue{
		Key:   key,
		Value: ctx.Request.FormValue(key),
	}
}

// 将参数值转为string，入参为true时自动去除前后空格
func (bv *ReqValue) String(trim bool) string {
	if bv.Error != nil {
		return ""
	}
	if trim == true {
		return strings.TrimSpace(bv.Value)
	}
	return bv.Value
}

// Int 将参数值转为int类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Int(def ...int) (int, error) {
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

// Int32 将参数值转为int32类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Int32(def ...int32) (int32, error) {
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

// Int64 将参数值转为int64类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Int64(def ...int64) (int64, error) {
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

// Uint32 将参数值转为uint32类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Uint32(def ...uint32) (uint32, error) {
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

// Uint64 将参数值转为uint64类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Uint64(def ...uint64) (uint64, error) {
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

// Float32 将参数值转为float32类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Float32(def ...float32) (float32, error) {
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

// Float64 将参数值转为float64类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Float64(def ...float64) (float64, error) {
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

// Bool 将参数值转为bool类型
// 如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (bv *ReqValue) Bool(def ...bool) (bool, error) {
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

package httpdispatcher

import (
	"context"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//上下文
type Content struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	routerParams   []httprouter.Param     //路由参数
	params         map[string]interface{} //ctx参数
	c              context.Context
	next           bool //继续往下执行处理器的标识
}

//request中的GET/POST等方法的参数值
type ReqValue struct {
	Key   string //参数名
	Value string //参数值
	Error error  //错误
}

//初始化ctx
func (ctx *Content) init() error {
	ctx.params = make(map[string]interface{})
	return ctx.Request.ParseForm()
}

//设置标识，用于继续执行下一个处理器
func (ctx *Content) Next(flag bool) {
	ctx.next = flag
}

//在ctx里存储值，如果key存在则替换值
func (ctx *Content) SetContextValue(key string, value interface{}) {
	ctx.params[key] = value
}

//获取ctx里的值，取出后根据写入的类型自行断言
func (ctx *Content) ContextValue(key string) interface{} {
	return ctx.params[key]
}

//获取某个GET参数值
func (ctx *Content) QueryValue(key string) *ReqValue {
	return &ReqValue{
		Key:   key,
		Value: ctx.Request.Form.Get(key),
	}
}

//获取某个POST参数值
func (ctx *Content) FormValue(key string) *ReqValue {
	value := ctx.Request.PostFormValue(key)
	if value == "" {
		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
			return &ReqValue{
				Key:   key,
				Error: err,
			}
		}
		if ctx.Request.MultipartForm != nil {
			values := ctx.Request.MultipartForm.Value[key]
			if len(values) > 0 {
				value = values[0]
			}
		}
	}
	return &ReqValue{
		Key:   key,
		Value: value,
	}
}

//将参数值转为string
func (rv *ReqValue) String() string {
	if rv.Error != nil {
		return ""
	}
	return rv.Value
}

//将参数值转为int类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Int(def ...int) (int, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.Atoi(rv.Value)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return value, nil
}

//将参数值转为int32类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Int32(def ...int32) (int32, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseInt(rv.Value, 10, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return int32(value), nil
}

//将参数值转为int64类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Int64(def ...int64) (int64, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseInt(rv.Value, 10, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return value, nil
}

//将参数值转为uint32类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Uint32(def ...uint32) (uint32, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseUint(rv.Value, 10, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return uint32(value), nil
}

//将参数值转为uint64类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Uint64(def ...uint64) (uint64, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseUint(rv.Value, 10, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return value, nil
}

//将参数值转为float32类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Float32(def ...float32) (float32, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseFloat(rv.Value, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return float32(value), nil
}

//将参数值转为float64类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Float64(def ...float64) (float64, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return 0, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseFloat(rv.Value, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}
	return value, nil
}

//将参数值转为bool类型
//如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (rv *ReqValue) Bool(def ...bool) (bool, error) {
	defLen := len(def)
	if rv.Error != nil {
		if defLen == 0 {
			return false, rv.Error
		} else {
			return def[0], nil
		}
	}
	value, err := strconv.ParseBool(rv.Value)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		} else {
			return false, err
		}
	}
	return value, nil
}

package httpdispatcher

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//路由组
type RouterGroup struct {
	hooks    []HookHandler //钩子
	basePath string        //基路径
	d        *dispatcher   //调度器
}

//路由组
func (r *RouterGroup) GROUP(path string, hooks ...HookHandler) *RouterGroup {
	//声明一个新的路由组
	var group RouterGroup
	//继承父组的路径
	group.basePath = r.basePath + path
	//传入调度器
	group.d = r.d
	//继承父组的钩子
	group.hooks = append(r.hooks)
	//加入当前传入的钩子
	for k := range hooks {
		group.hooks = append(group.hooks, hooks[k])
	}
	return &group
}

func (r *RouterGroup) GET(path string, handler Handler, hooks ...HookHandler) {
	//继承父组的路径
	path = r.basePath + path
	//使用httpRouter注册路由
	r.d.httpRouter.GET(path, func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		//初始化ctx
		ctx, err := r.d.initContext(resp, req, params)
		if err != nil {
			return
		}
		//执行路由组的钩子
		for k := range r.hooks {
			err := r.d.executeHookHandler(ctx, r.hooks[k])
			if err != nil {
				r.d.logger(err.Error(), -1)
				return
			}
		}
		//执行刚传入的钩子
		for k := range hooks {
			err := r.d.executeHookHandler(ctx, hooks[k])
			if err != nil {
				r.d.logger(err.Error(), -1)
				return
			}
		}

		//执行处理器函数
		r.d.executeHandler(ctx, handler, nil)
	})
}

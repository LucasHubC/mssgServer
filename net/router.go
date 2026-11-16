package net

import (
	"mssgServer/infra/log"
	"time"
)

// HandlerFunc 业务 handler 签名（不变）
type HandlerFunc func(req *WsMsgReq, rsp *WsMsgResp)

// Middleware 中间件签名：洋葱模型，接收下一个 handler，返回包装后的 handler
type Middleware func(HandlerFunc) HandlerFunc

type group struct {
	prefix      string
	handlerMap  map[string]HandlerFunc
	middlewares []Middleware // 组级中间件
}

// Router 路由表 + 全局中间件
// 执行顺序：全局 middleware → 组 middleware → handler（洋葱模型）
type Router struct {
	groups      map[string]*group
	middlewares []Middleware // 全局中间件
}

func NewRouter() *Router {
	return &Router{
		groups: make(map[string]*group),
	}
}

// Use 注册全局中间件（所有路由都会走）
func (r *Router) Use(mw ...Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

// Group 注册一个路由组
func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.groups[prefix] = g
	return g
}

// Use 组级中间件（只给这个组的路由加）
func (g *group) Use(mw ...Middleware) {
	g.middlewares = append(g.middlewares, mw...)
}

func (g *group) AddRoute(name string, handler HandlerFunc) {
	g.handlerMap[name] = handler
}

// buildChain 组装中间件链：全局 + 组级 + handler
// 洋葱模型：mw1 → mw2 → ... → handler
func (r *Router) buildChain(g *group, handler HandlerFunc) HandlerFunc {
	// 先包组级中间件（逆序，最后加的在最内层）
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i](handler)
	}
	// 再包全局中间件
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}

// Run 路由分发
func (r *Router) Run(req *WsMsgReq, rsp *WsMsgResp) {
	prefix, name, ok := splitRoute(req.Body.Name)
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "invalid route: " + req.Body.Name}
		return
	}
	g, ok := r.groups[prefix]
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "route not found: " + req.Body.Name}
		return
	}
	handler, ok := g.handlerMap[name]
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "route not found: " + req.Body.Name}
		return
	}
	// 组装中间件链后执行
	chain := r.buildChain(g, handler)
	chain(req, rsp)
}

// splitRoute 按 "." 拆分路由名，必须恰好两部分才算合法
func splitRoute(full string) (string, string, bool) {
	for i := 0; i < len(full); i++ {
		if full[i] == '.' {
			prefix := full[:i]
			name := full[i+1:]
			if prefix == "" || name == "" {
				return "", "", false
			}
			return prefix, name, true
		}
	}
	return "", "", false
}

// ────────────────────────────────────────
//  内置中间件
// ────────────────────────────────────────

// RecoverMiddleware 兜底 handler panic
// 某个 handler panic 不会炸掉整个 WS 连接，其他正常请求继续服务
func RecoverMiddleware(next HandlerFunc) HandlerFunc {
	return func(req *WsMsgReq, rsp *WsMsgResp) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("handler panic recovered",
					"route", req.Body.Name,
					"uid", GetUidFromReq(req),
					"err", r)
				rsp.Body.Code = -1
				rsp.Body.Msg = map[string]string{"errMsg": "internal error"}
			}
		}()
		next(req, rsp)
	}
}

// LogMiddleware 打印每条请求的结构化日志（路由 + uid + 业务码 + 耗时）
func LogMiddleware(next HandlerFunc) HandlerFunc {
	return func(req *WsMsgReq, rsp *WsMsgResp) {
		start := time.Now()
		next(req, rsp)
		duration := time.Since(start)
		log.Debug("ws handler",
			"route", req.Body.Name,
			"uid", GetUidFromReq(req),
			"code", rsp.Body.Code,
			"duration_ms", duration.Milliseconds(),
		)
	}
}

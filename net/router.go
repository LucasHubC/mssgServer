package net

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgResp)

type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

// Router 路由表，外层按 prefix 存 group，内层按 name 存 handler
// 路由查找 = 两次 map 取值，零循环
type Router struct {
	groups map[string]*group
}

func NewRouter() *Router {
	return &Router{
		groups: make(map[string]*group),
	}
}

// Group 注册一个路由组（prefix 重复会覆盖旧的）
func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.groups[prefix] = g
	return g
}

func (g *group) AddRoute(name string, handler HandlerFunc) {
	g.handlerMap[name] = handler
}

// Run 路由分发：两次 O(1) map 查找，零循环
// 消息格式必须是 "prefix.name"（如 "account.login"）
func (r *Router) Run(req *WsMsgReq, rsp *WsMsgResp) {
	// 先按 "." 拆 prefix 和 name
	prefix, name, ok := splitRoute(req.Body.Name)
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "invalid route: " + req.Body.Name}
		return
	}

	// 第一次 map 取值：找 group
	g, ok := r.groups[prefix]
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "route not found: " + req.Body.Name}
		return
	}

	// 第二次 map 取值：找 handler
	handler, ok := g.handlerMap[name]
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "route not found: " + req.Body.Name}
		return
	}

	handler(req, rsp)
}

// splitRoute 按 "." 拆分路由名，必须恰好两部分才算合法
// 返回 prefix, name, ok
func splitRoute(full string) (string, string, bool) {
	// 手动 split 比 strings.Split 更快，少一次切片分配
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

package net

import "strings"

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgResp)

type group struct {
	prefix string

	handlerMap map[string]HandlerFunc
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}

	r.group = append(r.group, g)
	return g
}

type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}
}

func (g *group) AddRoute(name string, handler HandlerFunc) {
	g.handlerMap[name] = handler
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgResp) {

	handler := g.handlerMap[name]

	if handler == nil {
		return
	}

	handler(req, rsp)

}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgResp) {
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}

	// 遍历所有组，找匹配的 handler
	found := false
	for _, g := range r.group {
		if g.prefix == prefix {
			for h, handler := range g.handlerMap {
				if h == name {
					handler(req, rsp)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if !found {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "route not found: " + req.Body.Name}
	}
}

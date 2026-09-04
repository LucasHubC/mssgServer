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
	//req.Body.Name  路由key 登录业务  account.login

	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""

	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}

	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}

}

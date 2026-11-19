package controller

import (
	"encoding/json"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/net"
	wspb "mssgServer/proto"
	invDAO "mssgServer/server/inventory/dao"
	invProto "mssgServer/server/inventory/proto"
	"net/http"
	"sync"
	"time"
)

// ItemController 道具控制器（单例）
var DefaultItem = &ItemController{}

type ItemController struct{}

// ────────────────────────────────────────
//  WebSocket 路由
// ────────────────────────────────────────

func (c *ItemController) Router(r *net.Router) {
	g := r.Group("inventory")
	g.AddRoute("add", c.add)
	g.AddRoute("consume", c.consume)
	g.AddRoute("drop", c.drop)
	g.AddRoute("get", c.get)
}

// add 增加道具
func (c *ItemController) add(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	uid, _ := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "未登录"}
		return
	}
	body, ok := req.Body.Msg.(*wspb.InventoryAdd)
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	if body.Count <= 0 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "数量必须大于 0"}
		return
	}
	before, _ := rediscache.GetItemCount(uid, body.ItemId)
	newCount, err := rediscache.AddItem(uid, body.ItemId, body.Count)
	if err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "操作失败"}
		return
	}
	invDAO.ItemOpLogDAO.InsertSafe(uid, body.ItemId, 1, body.Count, before, newCount, "ws_cmd")
	rsp.Body.Code = 0
	rsp.Body.Msg = &wspb.InventoryAddResp{ItemId: body.ItemId, NewCount: newCount}
}

// consume 使用/消耗道具
func (c *ItemController) consume(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	uid, _ := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "未登录"}
		return
	}
	body, ok := req.Body.Msg.(*wspb.InventoryConsume)
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	if body.Count <= 0 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "数量必须大于 0"}
		return
	}
	ok2, countOrBefore, err := rediscache.ConsumeItem(uid, body.ItemId, body.Count)
	if err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "操作失败"}
		return
	}
	if !ok2 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "道具数量不足"}
		return
	}
	before := countOrBefore + body.Count
	invDAO.ItemOpLogDAO.InsertSafe(uid, body.ItemId, 2, body.Count, before, countOrBefore, "ws_cmd")
	rsp.Body.Code = 0
	rsp.Body.Msg = &wspb.InventoryConsumeResp{ItemId: body.ItemId, NewCount: countOrBefore}
}

// drop 丢弃道具
func (c *ItemController) drop(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	uid, _ := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "未登录"}
		return
	}
	body, ok := req.Body.Msg.(*wspb.InventoryDrop)
	if !ok {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	if body.Count <= 0 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "数量必须大于 0"}
		return
	}
	ok2, newCount, err := rediscache.ConsumeItem(uid, body.ItemId, body.Count)
	if err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "操作失败"}
		return
	}
	if !ok2 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "道具数量不足"}
		return
	}
	before := newCount + body.Count
	invDAO.ItemOpLogDAO.InsertSafe(uid, body.ItemId, 3, body.Count, before, newCount, "ws_cmd")
	rsp.Body.Code = 0
	rsp.Body.Msg = &wspb.InventoryDropResp{ItemId: body.ItemId, NewCount: newCount}
}

// get 查背包
func (c *ItemController) get(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	uid, _ := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "未登录"}
		return
	}
	items, err := rediscache.GetAllItems(uid)
	if err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "查询失败"}
		return
	}
	itemList := make([]*wspb.ItemEntry, 0, len(items))
	for id, count := range items {
		itemList = append(itemList, &wspb.ItemEntry{ItemId: id, Count: count})
	}
	rsp.Body.Code = 0
	rsp.Body.Msg = &wspb.InventoryGetResp{Items: itemList, Total: int32(len(itemList))}
}

// ────────────────────────────────────────
//  HTTP REST 路由
// ────────────────────────────────────────

func (c *ItemController) RegisterHttpRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/inventory/add", c.addHttp)
	mux.HandleFunc("/api/inventory/consume", c.consumeHttp)
	mux.HandleFunc("/api/inventory/get", c.getHttp)
}

func (c *ItemController) addHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		net.WriteBadRequest(w, "method not allowed")
		return
	}
	uid, ok := net.AuthUid(r)
	if !ok {
		net.WriteUnauthorized(w, "token 无效")
		return
	}

	body := &invProto.AddItemReq{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		net.WriteBadRequest(w, "参数解析失败")
		return
	}
	// ✅ Bug #3 fix: 校验 count
	if body.Count <= 0 {
		net.WriteBadRequest(w, "数量必须大于 0")
		return
	}
	// 取 before
	before, _ := rediscache.GetItemCount(uid, body.ItemId)
	newCount, err := rediscache.AddItem(uid, body.ItemId, body.Count)
	if err != nil {
		net.WriteError(w, "操作失败")
		return
	}
	invDAO.ItemOpLogDAO.InsertSafe(uid, body.ItemId, 1, body.Count, before, newCount, "http_api")
	net.WriteJSON(w, map[string]any{"code": 0, "data": invProto.AddItemResp{ItemId: body.ItemId, NewCount: newCount}})
}

func (c *ItemController) consumeHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		net.WriteBadRequest(w, "method not allowed")
		return
	}
	uid, ok := net.AuthUid(r)
	if !ok {
		net.WriteUnauthorized(w, "token 无效")
		return
	}

	body := &invProto.ConsumeItemReq{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		net.WriteBadRequest(w, "参数解析失败")
		return
	}
	// ✅ Bug #3 fix: 校验 count
	if body.Count <= 0 {
		net.WriteBadRequest(w, "数量必须大于 0")
		return
	}
	// ✅ Bug #2 fix: 直接用 ConsumeItem 返回的 newCount
	ok2, newCount, err := rediscache.ConsumeItem(uid, body.ItemId, body.Count)
	if err != nil {
		net.WriteError(w, "操作失败")
		return
	}
	if !ok2 {
		net.WriteBadRequest(w, "道具数量不足")
		return
	}
	before := newCount + body.Count
	invDAO.ItemOpLogDAO.InsertSafe(uid, body.ItemId, 2, body.Count, before, newCount, "http_api")
	net.WriteJSON(w, map[string]any{"code": 0, "data": invProto.ConsumeItemResp{ItemId: body.ItemId, NewCount: newCount}})
}

func (c *ItemController) getHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		net.WriteBadRequest(w, "method not allowed")
		return
	}
	uid, ok := net.AuthUid(r)
	if !ok {
		net.WriteUnauthorized(w, "token 无效")
		return
	}

	items, err := rediscache.GetAllItems(uid)
	if err != nil {
		net.WriteError(w, "查询失败")
		return
	}
	itemList := make([]invProto.ItemInfo, 0, len(items))
	for id, count := range items {
		itemList = append(itemList, invProto.ItemInfo{ItemId: id, Count: count})
	}
	net.WriteJSON(w, map[string]any{"code": 0, "data": invProto.GetInventoryResp{Items: itemList, Total: len(itemList)}})
}

// ────────────────────────────────────────
//  自动保存调度器
// ────────────────────────────────────────

const autoSaveInterval = 30 * time.Second
const saveConcurrency = 10 // 定时落盘最大并发数（防 MySQL 连接池耗尽）

// StartAutoSave 启动定时落盘 goroutine（进程启动时调一次）
func (c *ItemController) StartAutoSave() {
	go func() {
		log.Info("[Inventory] auto-save started", "interval", autoSaveInterval, "concurrency", saveConcurrency)
		ticker := time.NewTicker(autoSaveInterval)
		defer ticker.Stop()
		for range ticker.C {
			c.saveDirtyUids()
		}
	}()
}

// SaveAllDirty 立即保存所有脏 uid（优雅关闭时调用）
func (c *ItemController) SaveAllDirty() {
	c.saveDirtyUids()
}

func (c *ItemController) saveDirtyUids() {
	uids := rediscache.GetDirtyUids()
	if len(uids) == 0 {
		return
	}
	log.Info("[Inventory] auto-save tick", "dirty_count", len(uids))

	// 并行保存（带并发上限）
	sem := make(chan struct{}, saveConcurrency)
	var wg sync.WaitGroup
	for _, uid := range uids {
		wg.Add(1)
		sem <- struct{}{}
		go func(u string) {
			defer wg.Done()
			defer func() { <-sem }()
			c.saveOne(u)
		}(uid)
	}
	wg.Wait()
	log.Info("[Inventory] auto-save done")
}

// saveOne 保存单个 uid 的背包到 MySQL
func (c *ItemController) saveOne(uid string) {
	// 从 Redis 读最新数据
	rediscache.RemoveZeroCount(uid)
	items, err := rediscache.GetAllItems(uid)
	if err != nil {
		log.Warn("[Inventory] read from redis failed", "uid", uid, "err", err)
		return
	}
	// 写 MySQL
	if err := invDAO.ItemDAO.SaveAll(uid, items); err != nil {
		log.Warn("[Inventory] save to mysql failed", "uid", uid, "err", err)
		return
	}
	// 清除脏标记
	rediscache.ClearDirty(uid)
}

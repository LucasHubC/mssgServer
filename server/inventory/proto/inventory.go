package proto

// ItemInfo 道具信息
type ItemInfo struct {
	ItemId string `json:"itemId"`         // 道具实例 ID（唯一）或 模板 ID
	Count  int64  `json:"count"`          // 数量
	Name   string `json:"name,omitempty"` // 可选：前端展示用名称
}

// ──── WebSocket 消息（name 约定: inventory.xxx）────

// AddItemReq 增加道具（GM 命令 / 测试用 / 掉落时）
type AddItemReq struct {
	ItemId string `json:"itemId" validate:"required"`
	Count  int64  `json:"count" validate:"required,min=1"`
}
type AddItemResp struct {
	ItemId   string `json:"itemId"`
	NewCount int64  `json:"newCount"` // 增加后的总数量
}

// ConsumeItemReq 使用/消耗道具
type ConsumeItemReq struct {
	ItemId string `json:"itemId" validate:"required"`
	Count  int64  `json:"count" validate:"required,min=1"`
}
type ConsumeItemResp struct {
	ItemId   string `json:"itemId"`
	NewCount int64  `json:"newCount"`
}

// GetInventoryReq 查询背包（空结构体即可）
type GetInventoryReq struct{}

// GetInventoryResp 返回完整背包
type GetInventoryResp struct {
	Items []ItemInfo `json:"items"`
	Total int        `json:"total"` // 道具种类数
}

// DropItemReq 丢弃道具
type DropItemReq struct {
	ItemId string `json:"itemId" validate:"required"`
	Count  int64  `json:"count" validate:"required,min=1"`
}
type DropItemResp struct {
	ItemId   string `json:"itemId"`
	NewCount int64  `json:"newCount"`
}

// ──── HTTP REST（/api/inventory/...）────

// —— 复用上面的 AddItemReq / ConsumeItemReq / GetInventoryReq ——

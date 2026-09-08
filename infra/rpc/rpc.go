package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mssgServer/infra/log"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 3 * time.Second,
}

// NodeRegisterReq 后端节点启动时向 Gateway 注册自己
type NodeRegisterReq struct {
	Type     string `json:"type"`      // 节点类型: "login" / "game" / "chat"
	Addr     string `json:"addr"`      // 自己的 ws 地址: "ws://127.0.0.1:8001"
	HTTPAddr string `json:"http_addr"` // 自己的 HTTP 地址（Gateway 反向健康检查用）
	MaxConns int64  `json:"max_conns"` // 最大容量
}

// NodeHeartbeatReq 后端节点定期上报心跳
type NodeHeartbeatReq struct {
	Type   string `json:"type"`
	Addr   string `json:"addr"`
	Online int64  `json:"online"` // 当前在线人数
}

// RouteReq 客户端向 Gateway 请求路由
type RouteReq struct {
	Type string `json:"type"` // 要什么服务: "login" / "game" / "chat"
	UId  string `json:"uid"`  // 用户 id（用于 hash 路由，可选）
}

// RouteResp Gateway 返回给客户端
type RouteResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Addr string `json:"addr,omitempty"` // 目标 ws 地址
}

// Register 后端节点启动时注册
func Register(gatewayURL string, req NodeRegisterReq) error {
	body, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/gateway/register", gatewayURL)
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("register request failed: %w", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("register failed status=%d body=%s", resp.StatusCode, string(data))
	}
	log.Info("node registered to gateway", "type", req.Type, "addr", req.Addr)
	return nil
}

// Heartbeat 后端节点定期上报心跳
func Heartbeat(gatewayURL string, req NodeHeartbeatReq) {
	body, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/gateway/heartbeat", gatewayURL)
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Warn("heartbeat failed", "err", err)
		return
	}
	resp.Body.Close()
}

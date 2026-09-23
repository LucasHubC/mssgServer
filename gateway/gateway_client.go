package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"mssgServer/infra/rpc"
	"net/http"
	"net/url"
	"time"
)

// HTTPGatewayClient 实现 NodeProvider 接口
// Broker 进程用这个通过 HTTP 调 Gateway 拿后端节点地址
type HTTPGatewayClient struct {
	baseURL    string // e.g. "http://127.0.0.1:9000"
	httpClient *http.Client
}

func NewHTTPGatewayClient(baseURL string) *HTTPGatewayClient {
	// 🔴 关键:自定义 Transport 加连接池
	// 默认 Transport MaxIdleConnsPerHost=2,3000 并发 Route 会打爆 Gateway
	transport := &http.Transport{
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 500, // Broker→Gateway 全打到同一个 host
		IdleConnTimeout:     30 * time.Second,
	}
	return &HTTPGatewayClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   2 * time.Second, // 快速失败,不要等 3s
		},
	}
}

// Route 调 Gateway /gateway/route?type=xxx&uid=xxx 拿后端节点地址
func (c *HTTPGatewayClient) Route(req rpc.RouteReq) *rpc.RouteResp {
	u := fmt.Sprintf("%s/gateway/route?type=%s", c.baseURL, url.QueryEscape(req.Type))
	if req.UId != "" {
		u += "&uid=" + url.QueryEscape(req.UId)
	}

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return &rpc.RouteResp{Code: -1, Msg: "gateway unreachable: " + err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var routeResp rpc.RouteResp
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return &rpc.RouteResp{Code: -1, Msg: "invalid gateway response"}
	}
	return &routeResp
}

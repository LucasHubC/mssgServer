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
	return &HTTPGatewayClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
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

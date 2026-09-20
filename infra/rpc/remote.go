// Package rpc 提供轻量级跨服 RPC（请求-响应）能力
//
// 实现方式：Redis Stream（请求队列）+ Redis Pub/Sub（响应回传）
//   - 请求方：往目标节点的 Stream 里 XADD 一条请求，然后订阅自己的 Pub/Sub 频道等响应
//   - 响应方：XREAD 监听自己的 Stream，收到请求 → handler 处理 → PUBLISH 回请求方
//
// 为什么不用纯 Pub/Sub？
//
//	Pub/Sub 发了就忘，如果响应方刚好在 GC 停顿或网络抖动期间收到消息，
//	请求方等超时后会认为 RPC 失败，但响应方可能已经处理了业务 —— 有状态不一致风险
//
// 为什么响应还用 Pub/Sub 不用 Stream？
//
//	响应方处理完才发响应，发起方肯定在 Call() 里同步等待（持有 pending map entry），
//	用 Pub/Sub 最快速，且响应丢了发起方只是超时重试，不会有状态丢失
//
// 用法示例：
//
//	// 启动时初始化一次
//	rpc.Init("ws://loginServer:8003")
//
//	// 响应方：注册方法
//	rpc.Handle("get_online_count", func(params map[string]any) (any, error) {
//	    return net.ConnManager.Total(), nil
//	})
//
//	// 发起方：调用
//	result, err := rpc.Call("ws://chatServer:8004", "get_online_count", nil, 3*time.Second)
package rpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	rpcTimeout   = 2 * time.Second // Redis 操作超时
	requestTTL   = 5 * time.Minute // 请求 Stream 里的消息自动过期（防止堆积）
	pendingTTL   = 3 * time.Second // Call() 默认等待超时
	xreadBlockMs = 2000            // XREAD 阻塞轮询间隔
)

var (
	myNodeAddr string // 当前节点地址（启动时 Init 设）

	// handlers 已注册的 RPC 方法表
	handlers = make(map[string]func(params map[string]any) (any, error))
	hmu      sync.RWMutex

	// pending 等待响应的 Call()
	pending = make(map[string]chan *RPCResponse)
	pmu     sync.Mutex

	inited    bool
	initedMu  sync.Mutex
	closeChan = make(chan struct{}) // Init 之后可以 Close()
)

// ────────────────────────────────────────
//  消息体
// ────────────────────────────────────────

// RPCRequest 跨服 RPC 请求体
type RPCRequest struct {
	ID       string         `json:"id"`
	Method   string         `json:"method"`
	Params   map[string]any `json:"params,omitempty"`
	ReplyTo  string         `json:"reply_to"` // 响应发往哪个 Pub/Sub 频道（= 发起方节点地址）
	FromNode string         `json:"from_node"`
}

// RPCResponse 跨服 RPC 响应体
type RPCResponse struct {
	ID     string `json:"id"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"` // 非空=失败
}

// ────────────────────────────────────────
//  生命周期
// ────────────────────────────────────────

// Init 初始化 RPC 模块，传入当前节点地址
// 启动后自动开启两个后台协程：
//  1. listenRequests  —— 监听发往本节点的 RPC 请求 Stream
//  2. listenResponses —— 监听发回本节点的 RPC 响应 Pub/Sub
func Init(nodeAddr string) {
	initedMu.Lock()
	defer initedMu.Unlock()

	if inited {
		log.Warn("[RPC] Init 已经调用过了，跳过")
		return
	}
	if db.RDB == nil {
		log.Warn("[RPC] Redis 未初始化，RPC 模块不可用")
		return
	}

	myNodeAddr = nodeAddr
	inited = true

	go listenRequests()
	go listenResponses()
	log.Info("[RPC] 已初始化，node=", nodeAddr)
}

// Close 停止 RPC 模块（优雅退出用）
func Close() {
	close(closeChan)
}

// ────────────────────────────────────────
//  响应方：注册方法
// ────────────────────────────────────────

// Handle 注册一个 RPC 方法
// 同名方法重复注册会覆盖（启动时注册一般不会重复）
func Handle(method string, fn func(params map[string]any) (any, error)) {
	hmu.Lock()
	defer hmu.Unlock()
	handlers[method] = fn
}

// ────────────────────────────────────────
//  发起方：调用
// ────────────────────────────────────────

// Call 发起跨服 RPC 调用
//
//	targetNode —— 目标节点的地址（如 "ws://chatServer:8004"）
//	method     —— 方法名（必须和目标节点 Handle 注册的一致）
//	params     —— 参数（nil 或 map）
//	timeout    —— 等待响应的超时
//
// 返回 (result, nil) 成功；(nil, err) 失败
// 超时错误：errors.New("rpc timeout")
// 远端 handler 返回错误：err.Error() = 远端 error 字符串
func Call(targetNode, method string, params map[string]any, timeout time.Duration) (any, error) {
	if !inited {
		return nil, errors.New("rpc not initialized")
	}
	if targetNode == "" || method == "" {
		return nil, errors.New("targetNode and method required")
	}

	reqID := genReqID()
	ch := make(chan *RPCResponse, 1) // buffer 1，防止发送方阻塞

	// 注册到 pending
	pmu.Lock()
	pending[reqID] = ch
	pmu.Unlock()
	defer func() {
		pmu.Lock()
		delete(pending, reqID)
		pmu.Unlock()
	}()

	req := RPCRequest{
		ID:       reqID,
		Method:   method,
		Params:   params,
		ReplyTo:  myNodeAddr,
		FromNode: myNodeAddr,
	}

	// XADD 到目标节点的请求 Stream
	streamKey := reqStreamKey(targetNode)
	payload, _ := json.Marshal(req)

	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	id, err := db.RDB.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: 1000, // 防止 Stream 无限增长
		Values: map[string]any{
			"payload": string(payload),
		},
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("rpc xadd failed: %w", err)
	}

	log.Debug("[RPC] call sent", "to", targetNode, "method", method, "reqID", id)

	// 等待响应
	if timeout <= 0 {
		timeout = pendingTTL
	}
	select {
	case rsp := <-ch:
		if rsp.Error != "" {
			return nil, errors.New(rsp.Error)
		}
		return rsp.Result, nil
	case <-time.After(timeout):
		return nil, errors.New("rpc timeout")
	case <-closeChan:
		return nil, errors.New("rpc closed")
	}
}

// ────────────────────────────────────────
//  后台监听
// ────────────────────────────────────────

// listenRequests 监听发往本节点的 RPC 请求 Stream
// 用 XREAD BLOCK 阻塞等待，每隔 xreadBlockMs 唤醒一次（顺便检查 closeChan）
func listenRequests() {
	streamKey := reqStreamKey(myNodeAddr)
	lastID := "0"

	for {
		select {
		case <-closeChan:
			return
		default:
		}

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout+time.Second)
		res, err := db.RDB.XRead(ctx, &redis.XReadArgs{
			Streams: []string{streamKey, lastID},
			Block:   xreadBlockMs * time.Millisecond,
			Count:   10,
		}).Result()
		cancel()

		if err == redis.Nil {
			continue // Block 超时，正常
		}
		if err != nil {
			log.Warn("[RPC] XRead failed, retry...", "err", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		for _, stream := range res {
			for _, msg := range stream.Messages {
				lastID = msg.ID // 继续读下一条

				payloadStr, _ := msg.Values["payload"].(string)
				var req RPCRequest
				if err := json.Unmarshal([]byte(payloadStr), &req); err != nil {
					log.Warn("[RPC] invalid request payload", "err", err)
					continue
				}

				go handleRequest(req) // 每条请求独立 goroutine 处理，不阻塞下一条
			}
		}
	}
}

// handleRequest 处理一条 RPC 请求，然后 PUBLISH 响应回去
func handleRequest(req RPCRequest) {
	var result any
	var rpcErr error

	hmu.RLock()
	fn, ok := handlers[req.Method]
	hmu.RUnlock()

	if !ok {
		rpcErr = fmt.Errorf("method not found: %s", req.Method)
	} else {
		result, rpcErr = safeCall(fn, req.Params) // recover panic
	}

	rsp := RPCResponse{ID: req.ID, Result: result}
	if rpcErr != nil {
		rsp.Error = rpcErr.Error()
	}

	// 发布响应到发起方的 Pub/Sub 频道
	channel := rspChannelKey(req.ReplyTo)
	data, _ := json.Marshal(rsp)

	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()
	if err := db.RDB.Publish(ctx, channel, string(data)).Err(); err != nil {
		log.Warn("[RPC] publish response failed", "err", err, "reqID", req.ID)
	}
}

// listenResponses 订阅发回本节点的 RPC 响应频道
func listenResponses() {
	channel := rspChannelKey(myNodeAddr)
	sub := db.RDB.Subscribe(context.Background(), channel)
	defer sub.Close()

	ch := sub.Channel()
	for {
		select {
		case <-closeChan:
			return
		case msg, ok := <-ch:
			if !ok {
				log.Warn("[RPC] response channel closed, reconnect...")
				// 重订阅
				sub = db.RDB.Subscribe(context.Background(), channel)
				ch = sub.Channel()
				continue
			}
			var rsp RPCResponse
			if err := json.Unmarshal([]byte(msg.Payload), &rsp); err != nil {
				log.Warn("[RPC] invalid response payload", "err", err)
				continue
			}
			pmu.Lock()
			pendingCh, exists := pending[rsp.ID]
			pmu.Unlock()
			if exists {
				select {
				case pendingCh <- &rsp:
				default: // 超时已经移除了，忽略
				}
			}
		}
	}
}

// ────────────────────────────────────────
//  工具函数
// ────────────────────────────────────────

func reqStreamKey(nodeAddr string) string {
	return "rpc:req:" + nodeAddr
}

func rspChannelKey(nodeAddr string) string {
	return "rpc:rsp:" + nodeAddr
}

func genReqID() string {
	var b [8]byte
	rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// safeCall 执行 handler，捕获 panic 转为 error
func safeCall(fn func(map[string]any) (any, error), params map[string]any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("rpc handler panic: %v", r)
			log.Warn("[RPC] handler panicked", "err", r)
		}
	}()
	return fn(params)
}

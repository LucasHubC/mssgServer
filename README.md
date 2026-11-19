# mssgServer — 分布式实时消息服务

基于 Go 的分布式 WebSocket 消息服务，支持登录认证、世界聊天、私聊、聊天室、群组、背包管理。
Gateway 统一路由 + Redis Pub/Sub 跨服广播 + 节点自动摘除以实现弹性伸缩。

---

## 架构

```
客户端
  │
  ├─ 阶段1: HTTP GET /gateway/route?type=xxx&uid=xxx  ← Gateway 路由发现
  │          返回 ws://127.0.0.1:8003 (后端节点地址)
  │
  └─ 阶段2: WebSocket 直连后端节点
              ├─ Worker Pool (2048) 异步执行 handler
              │     └─ Router (两次 O(1) map 取值)
              │           └─ controller → dao
              │                 ├─ Redis 缓存 (先查)
              │                 └─ MySQL 持久化 (miss 回源)
              │
              ├─ Redis Pub/Sub (跨节点广播 / 踢人 / 消息路由)
              ├─ Redis Stream  (跨服 RPC)
              └─ Gateway 注册 + 5s 心跳上报
```


### 节点类型

| 类型 | 二进制 | 默认端口 | 依赖 | 职责 |
|---|---|---|---|---|
| **gateway** | `cmd/gateway` | 9000 | Redis | 路由发现、节点注册/心跳、管理后台 |
| **login** | `cmd/login` | 8003 | MySQL + Redis | 登录认证、背包管理 |
| **chat** | `cmd/chat` | 8004 | Redis | 世界聊天、私聊、群组聊天 |
| **room** | `cmd/room` | 8005 | Redis | 聊天室 |
| **broker** | `cmd/broker` | 9001 | Gateway | 两阶段连接第二阶段：WS 智能路由代理 |

### 路由策略与负载均衡

Gateway 根据 `conf.ini` 的 `[gateway] policy` 配置，从多个同类型节点里选一个返回给 Broker：

| 策略 | 行为 | 适用场景 |
|---|---|---|
| **least_conn** ⭐ | 选 `Online`（心跳上报的在线人数）最小的节点 | 多 login/chat/room 节点自动负载均衡 |
| round_robin | 每次被选中 score+1，选 score 最小的 | 简单均匀分布，不关心节点实际忙闲 |
| uid_hash | `md5(uid)` 模节点数，同一用户固定落到同一节点 | 高缓存命中率场景（玩家数据缓在本机）|

#### 多节点示例

假设开了两个 login + 一个 broker：

```
login-A (8003)  心跳 Online=400  → Gateway 存 login→[A(400), B(150)]
login-B (8004)  心跳 Online=150  →

Broker 收到 WS 连接 → dialAll() 调 Gateway /gateway/route?type=login
                         ↓
Gateway.PickNode(policy="least_conn") → 选 login-B（更空）→ 返回给 Broker
                         ↓
Broker dial ws://127.0.0.1:8004 → 客户端的登录/背包消息全部路由到 login-B
```

- **策略可热切换**：改 `conf.ini` 的 `policy`，**重启 Gateway 即生效**（login/chat/room/broker 不需要重启）
- **自动摘除**：节点 15 秒无心跳 → Gateway CleanupLoop（10s 一次）自动剔除 → 新 route 不再分配已死节点

### 跨服通信

- **广播**：Redis Pub/Sub `broadcast:*` 频道 — 任节点 Publish，全节点 Subscribe 并在本机 uidMap 查找推送
- **RPC**：Redis Stream `rpc:*` — 支持 `ping`、`get_online_count`、`get_inventory` 等跨节点调用
- **用户定位**：Redis Key `online:user:{uid}` = `nodeAddr` — 查"某玩家在哪个节点"只需 O(1)
- **踢人/重连**：BindUser 时 SETNX 抢锁，抢到锁的节点向旧节点发 kick 广播

### 节点生命周期

```
后端启动 → Gateway 注册 (写 lastHeartbeat=now)
    ↓
后端每 5s → Gateway 心跳 (更新 lastHeartbeat + online_count)
    ↓
后端崩了 → 15s 无心跳 → Gateway CleanupLoop (10s 一次) 判定超时
    ↓
新客户端 route → Gateway 不再分配已死节点 ✅
```


---

## 快速开始

### 1. 前置依赖

```bash
# Go 1.21+（go.mod 指定版本）
go version

# MySQL 8.0+（可选，loginServer 需要）
docker run -d -p 3306:3306 mysql:8.0

# Redis 6.0+（强烈推荐，跨节点功能依赖）
docker run -d -p 6379:6379 redis:7
```

### 2. 安装 Go 依赖

```bash
cd mssgServer
go mod tidy
```

### 3. 配置

编辑 `conf/conf.ini`：

```ini
[gateway]
host = 0.0.0.0
port = 9000
policy = least_conn         ; round_robin | least_conn ⭐ | uid_hash
                            ; 多节点生产环境用 least_conn 自动均衡负载

[mysql]
host = 127.0.0.1
port = 3306
user = root
password = 123456
dbname = mssg

[redis]
host = 127.0.0.1
port = 6379
```

> ⚠️ **broker.host 不能用 0.0.0.0**。因为 broker 注册到 Gateway 后，Gateway 会把这个地址返回给客户端。0.0.0.0 是监听地址，不是客户端可连接的地址。生产环境改成具体 IP（如 192.168.x.x）。

### 4. 启动顺序

```bash
# 终端 1 — Gateway（最先启动，其他节点依赖它注册）
go run ./cmd/gateway

# 终端 2 — loginServer
go run ./cmd/login

# 终端 3 — chatServer
go run ./cmd/chat

# 终端 4 — roomServer（可选）
go run ./cmd/room

# 终端 5 — broker（可选，两阶段连接第二阶段）
go run ./cmd/broker
```

### 5. 验证

```bash
# 路由发现
curl "http://127.0.0.1:9000/gateway/route?type=login&uid=user_001"
# → {"code":0,"addr":"ws://127.0.0.1:8003"}

# 节点状态
curl "http://127.0.0.1:9000/gateway/status"

# 管理后台广播（推送给所有在线节点的所有玩家）
curl -X POST http://127.0.0.1:9000/admin/broadcast \
  -H "Content-Type: application/json" \
  -d '{"name":"world.chat","data":"管理员：欢迎新玩家！"}'
```


### 6. 优雅关闭

`Ctrl+C` 或 `SIGTERM` 信号触发后，各节点按以下顺序关闭：
1. 停止接受新连接（`http.Server.Shutdown`）
2. 从 Gateway 注销（立即摘除，新 route 不再分配本节点）
3. 广播 `server_shutdown` 通知在线玩家
4. 主动关闭所有 WebSocket 连接
5. 等待所有 handler 执行完毕（`net.WaitHandlers`）
6. 强制落盘脏数据（`inventory.Shutdown`）
7. `os.Exit(0)`

---

## 协议

使用 Protobuf 序列化，定义在 `proto/ws.proto`。

### 常用消息类型

| 名称 | 方向 | 说明 |
|---|---|---|
| `login.auth` | C→S | 登录认证 |
| `login.auth_ack` | S→C | 登录结果 |
| `chat.world.send` | C→S | 世界聊天 |
| `chat.world.recv` | S→C | 世界聊天广播 |
| `chat.private.send` | C→S | 私聊 |
| `chat.room.send` | C→S | 聊天室消息 |
| `chat.group.send` | C→S | 群组消息 |
| `inventory.add/consume/drop/get` | C→S | 背包操作 |

### WebSocket 连接示例

```javascript
// 1. 先从 Gateway 拿后端地址
const route = await fetch(
  "http://127.0.0.1:9000/gateway/route?type=login&uid=user_001"
).then(r => r.json());
// route.addr = "ws://127.0.0.1:8003"

// 2. 直连后端
const ws = new WebSocket(route.addr);
ws.onopen = () => {
  // 3. 发登录包（protobuf 编码后的二进制）
  ws.send(protobuf.encode(loginAuthReq).finish());
};
```


---

## 项目结构

```
mssgServer/
├── cmd/
│   ├── gateway/    # 路由发现 + 管理后台
│   ├── login/      # 登录 + 背包
│   ├── chat/       # 世界/私聊/群组
│   ├── room/       # 聊天室
│   └── broker/     # 两阶段连接第二阶段
├── config/
│   ├── config.go   # INI 文件加载
│   └── db/         # MySQL + Redis 初始化
├── conf/
│   └── conf.ini    # 配置文件
├── infra/
│   ├── log/        # 结构化日志（轮转 + dual 输出）
│   ├── redis/      # Redis 操作 + Pub/Sub 广播 + Stream RPC
│   ├── rpc/        # 跨服 RPC 封装（Call / Handle / Register）
│   ├── metrics/    # Prometheus 指标注册
│   └── util/       # ListenTCP 平台分文件实现
│       ├── listener_windows.go   # Windows: net.Listen（规避 GC finalizer 关 fd）
│       └── listener_unix.go      # Linux/macOS: syscall 可控 backlog + SO_REUSEADDR
├── gateway/        # 路由策略 + WS 代理 + 节点注册/心跳/清理
├── net/            # WebSocket Server / Router / ConnManager / 中间件
├── server/
│   ├── login/      # 登录 handler + DAO
│   ├── inventory/  # 背包 handler（30s 定时落盘）
│   ├── chat/       # 聊天 handler（world / private / group）
│   └── room/       # 聊天室 handler
├── proto/
│   ├── ws.proto    # Protobuf 定义
│   └── ws.pb.go    # protoc 生成
└── go.mod
```


---

## 测试

```bash
# 全部测试 + 覆盖率
go test ./... -v -cover

# 单元测试（不需要外部依赖）
go test ./gateway -v -run TestPickNode
go test ./net -v

# 需要 MySQL + Redis 连通
go test ./server/login/dao -v -run TestUserDAO
go test ./infra/redis -v -run TestBindUser
```


### 压测结果

测试条件：Windows 开发环境，1000 并发 WebSocket 连接，mixed 模式（world/private/room/group/inventory）

| 指标 | 数值 |
|---|---|
| 连接建立耗时 | 0.86s |
| 最终活跃连接 | 1000 / 1000 |
| 断线数 | 0 |
| 发送 QPS | ~6100/s |
| 接收 QPS | ~10500/s（recv/send ≈ 1.5–1.7x，广播扇出特性）|

---

## 已知问题与修复记录

### Windows 上 Broker GC 提前关闭 Listener fd（已修复）

**现象**：压测 1000 并发连接时，Broker 进程崩 → 108 ECONNREFUSED + 108 Close Code 1006 + 连接建立耗时 15.56s。

**根因**：手动 `syscall.Socket → os.NewFile → net.FileListener` 在 Windows 上，`os.NewFile` 注册的 GC finalizer 在高并发下有 race，会提前关闭底层 fd，导致 `AcceptEx` 返回 "I/O operation has been aborted"。

**修复**：按平台分文件实现 `ListenTCP`：
- `listener_windows.go` — 直接 `net.Listen`，规避 finalizer 问题
- `listener_unix.go` — syscall 精确控制 backlog + SO_REUSEADDR（Linux 生产环境）

**验证**：修复后压测零断线、零 ECONNREFUSED、连接建立耗时 < 1s。

### net.ListenConfig.Backlog 字段不存在（已修复）

不同 Go 版本对 `net.ListenConfig` 的 backlog 字段支持不一致（有的版本叫 `SocketListenBacklog`，有的版本完全没有）。改用手动 syscall 链路跨所有 Go 版本精确控制 backlog。

---

## License

MIT
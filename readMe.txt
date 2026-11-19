
# 装依赖
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
go mod tidy




# 检查 protoc 是否装了
protoc --version

# 如果没装，Windows 用 winget：
winget install Google.Protobuf

# 或手动下载：https://github.com/protocolbuffers/protobuf/releases

# 下载 protoc-xx.x-win64.zip，解压到 C:\tools\protoc\bin，加 PATH

# 安装 Go 代码生成器
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 验证
protoc-gen-go --version


cd e:\mssgServer

# 用 protoc 生成 Go 代码  https://github.com/protocolbuffers/protobuf/releases/latest

# 假设 protoc 装在 C:\tools\protoc\bin
protoc --go_out=. --go_opt=paths=source_relative proto/ws.proto

# 安装 protobuf runtime
go get google.golang.org/protobuf

# go mod tidy
go mod tidy



#客户端 testBench.js / testWS.js 改 Protobuf
# 用 protobufjs 生成 JS（node 侧用 protobufjs）
cd e:\mssgServer
npm install protobufjs
npx pbjs -t static-module -w commonjs -o proto/ws.js proto/ws.proto







cd e:\mssgServer

# 跑所有测试 + 输出覆盖率
go test ./... -v -cover

# 只跑 DAO 测试（必须 MySQL + Redis 都连通）
go test ./server/login/dao -v -run TestUserDAO

# 只跑 Gateway 测试（纯单元测试，不需要外部依赖）
go test ./gateway -v -run TestPickNode

# 只跑 Redis 测试
go test ./infra/redis -v -run TestBindUser

# 想看 goconvey 的彩色 BDD 风格输出
go test -v ./gateway  # 会显示 Given-When-Then 风格




客户端
  │
  ├─→ Gateway (HTTP 路由发现) ──→ 返回 ws://xxx:port
  │
  └─→ loginServer (WebSocket 直连)
        │
        ├─ Worker Pool (2048) 异步执行 handler
        │     └─ Router (两次 O(1) map 取值)
        │           └─ controller → dao
        │                 ├─ Redis 缓存 (先查)
        │                 └─ MySQL 持久化 (miss 回源)
        │
        ├─ Redis Pub/Sub (跨节点广播 / 踢人 / 消息路由)
        └─ Gateway 注册 + 5s 心跳上报






启动顺序：
  ① gatewayServer.exe  → 监听 :9000，启动 CleanupLoop
  ② loginServer.exe    → 注册到 Gateway + 5s 心跳 + WebSocket :8003

客户端流程：
  GET  http://9000/gateway/route?type=login&uid=user_001  ← Gateway 返回 ws://127.0.0.1:8003
  WS   ws://127.0.0.1:8003                                ← 客户端直接连后端

节点生命周期：
  loginServer 启动 → Gateway 注册（写 lastHeartbeat=now）
  loginServer 每 5s → Gateway 心跳（更新 lastHeartbeat + Online）
  loginServer 崩了 → 15s 没心跳 → Gateway 判定超时
  Gateway CleanupLoop（10s 一次）→ 删掉超时节点
  新客户端 route → Gateway 不再分配已死节点




启动顺序：
  ① 本地起 Redis（如果没装的话：docker run -p 6379:6379 redis）
  ② 起 Gateway
  ③ 起 loginServer（现在会 Init MySQL + InitRedis，然后 SubscribeBroadcast）

两个节点场景：
  ③-1 loginServer-A :8003  Register → Redis
  ③-2 loginServer-B :8004  Register → Redis
  
  客户端 A 登录 → Gateway 分到 loginServer-A
    BindUid("user_A", "loginServer-A") → Redis SET online:user:user_A = "loginServer-A"
  
  客户端 B 登录 → Gateway 分到 loginServer-B
    BindUid("user_B", "loginServer-B") → Redis SET online:user:user_B = "loginServer-B"
  
  loginServer-A 想给 user_B 发消息？
    → GetUserNode("user_B") → 返回 "loginServer-B"
    → PublishBroadcast("push", "user_B", msg)
    → loginServer-B 的 SubscribeBroadcast 收到
    → loginServer-B 从本机 uidMap 找到 user_B 的连接
    → Push 给他 ✅

  客户端 A 刷新页面重新登录 → Gateway 分到 loginServer-B
    → BindUser("user_A") SETNX 拿到锁
    → loginServer-B BindUid → 踢掉 loginServer-A 上的旧连接
    → 旧连接 Close → Unregister → UnbindUser ✅

用户缓存：
  第一次登录（缓存 miss）: DB 查 → 回写 Redis → 返回  ~10ms
  第二次登录（缓存 hit）:  Redis 直接取                  ~0.5ms
  ✅ QPS 提升 10~20x




# 终端 1：起 Gateway
.\gatewayServer.exe
# [Gateway] starting on 0.0.0.0:9000 policy=round_robin

# 终端 2：起 loginServer
.\loginServer.exe
# [Gateway] node registered    type=login addr=ws://127.0.0.1:8003

# 终端 3：请求路由
curl "http://127.0.0.1:9000/gateway/route?type=login&uid=user_001"
# {"code":0,"addr":"ws://127.0.0.1:8003"}

# 终端 3：看节点状态
curl "http://127.0.0.1:9000/gateway/status"

# 关掉 loginServer，等 15 秒后再请求 route
# {"code":-1,"msg":"no available node for type: login"}  ← 已自动摘除 ✅


通过网关 发送消息给所有人

cd e:\mssgServer

# 1. 先跑 Gateway（它需要 Redis，但不需要 MySQL）
go build -o gatewayServer.exe gatewayServer.go
.\gatewayServer.exe
# 看到: [Gateway] starting on 0.0.0.0:9000 policy=round_robin
#       [Gateway] admin APIs: /admin/broadcast, /admin/push, /admin/kick

# 2. 再跑 loginServer（多开几个终端窗口跑多个实例模拟集群）
go build -o loginServer.exe loginServer.go
.\loginServer.exe

#3
最简单的公告
$body = @{
    name = "system.notice"
    data = @{
        title   = "🎉 新版本上线啦"
        content = "今晚 22:00 停服维护 30 分钟"
        time    = "2026-09-11 21:00:00"
    }
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Uri "http://127.0.0.1:9000/admin/broadcast" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body



curl -X POST http://127.0.0.1:9000/admin/broadcast `
  -H "Content-Type: application/json" `
  -d '{
    "name": "system.notice",
    "data": {
      "title": "🎉 新版本上线啦",
      "content": "今晚 22:00 停服维护 30 分钟，感谢大家支持！",
      "time": "2026-09-11 21:00:00"
    }
  }'

纯文本广播
curl -X POST http://127.0.0.1:9000/admin/broadcast `
  -H "Content-Type: application/json" `
  -d '{
    "name": "world.chat",
    "data": "管理员：欢迎新来的玩家！"
  }'

活动推送
curl -X POST http://127.0.0.1:9000/admin/broadcast `
  -H "Content-Type: application/json" `
  -d '{
    "name": "activity.lucky_draw",
    "data": {
      "activityId": "act_001",
      "prize": "iPhone 17 Pro",
      "endTime": 1750000000
    }
  }'



  //工具说明 图集拆解
  texture-unpacker 是一个 Node.js 工具，专门用来拆 TexturePacker 打包的图集，支持 Cocos2d-x、Cocos2d、PixiJS 等多种 plist 格式。
  安装：

bash
npm install -g texture-unpacker
使用：

bash
texture-unpacker -s 你的大图.png -f plist -d 你的文件.plist -o 输出目录
参数说明：-s 指定合图 png，-f plist 声明数据格式，-d 指定 plist 文件，-o 指定输出目录。加上 -c 可以在拆分前清空输出目录。


texture-unpacker -s starKeyLockRes.png -f plist -d starKeyLockRes.plist -o output
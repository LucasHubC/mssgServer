package gateway

import (
	"mssgServer/infra/rpc"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func makeTestNode(typ, addr string, online int64, alive bool) *Node {
	n := &Node{
		Type:   typ,
		Addr:   addr,
		Online: online,
	}
	if alive {
		n.lastHeartbeat = time.Now()
	} else {
		n.lastHeartbeat = time.Now().Add(-1 * time.Minute) // 1 分钟没心跳 → 死
	}
	return n
}

// ────────────────────────────────────────
//  路由策略：三种 PickNode 行为
// ────────────────────────────────────────

func TestPickNode(t *testing.T) {
	// PickNode 的职责：从候选节点里挑一个，不负责 IsAlive 过滤
	// IsAlive 过滤是 Gateway.Route 做的，这里只传活着的节点
	nodes := []*Node{
		makeTestNode("login", "ws://A:8001", 100, true),
		makeTestNode("login", "ws://B:8002", 50, true),
		makeTestNode("login", "ws://C:8003", 200, true),
	}

	convey.Convey("RoundRobin 轮询", t, func() {
		// 每次选 score 最小的，被选中后 score++
		// 初始 score 都是 0，第一次选到 A，然后 A.score=1
		// 第二次 B.score=0 → 选 B ...
		selected := make(map[string]int)
		for i := 0; i < 6; i++ {
			n := PickNode(nodes, PolicyRoundRobin, "")
			convey.So(n, convey.ShouldNotBeNil)
			selected[n.Addr]++
		}
		// 3 个节点各被选 2 次（轮询公平性）
		convey.So(selected["ws://A:8001"], convey.ShouldEqual, 2)
		convey.So(selected["ws://B:8002"], convey.ShouldEqual, 2)
		convey.So(selected["ws://C:8003"], convey.ShouldEqual, 2)
	})

	convey.Convey("LeastConn 最少连接", t, func() {
		n := PickNode(nodes, PolicyLeastConn, "")
		convey.So(n, convey.ShouldNotBeNil)
		convey.So(n.Online, convey.ShouldEqual, 50) // B 在线最少
	})

	convey.Convey("UIDHash UID 哈希稳定", t, func() {
		// 同一个 UID 应该稳定地选同一个节点
		n1 := PickNode(nodes, PolicyUIDHash, "user_001")
		n2 := PickNode(nodes, PolicyUIDHash, "user_001")
		convey.So(n1.Addr, convey.ShouldEqual, n2.Addr)

		// 不同 UID 可能选不同节点（不能全选同一个）
		n3 := PickNode(nodes, PolicyUIDHash, "user_999")
		// 不强制不等于，但至少多次 hash 是稳定的
		n4 := PickNode(nodes, PolicyUIDHash, "user_999")
		convey.So(n3.Addr, convey.ShouldEqual, n4.Addr)
	})

	convey.Convey("空切片返回 nil", t, func() {
		n := PickNode(nil, PolicyRoundRobin, "")
		convey.So(n, convey.ShouldBeNil)
	})
}

// ────────────────────────────────────────
//  Gateway 注册 + 心跳 + 路由
// ────────────────────────────────────────

func TestGateway_RegisterAndRoute(t *testing.T) {
	convey.Convey("完整流程：注册→心跳→路由→摘除死节点", t, func() {
		g := NewGateway(PolicyRoundRobin)

		// ① 注册 3 个 login 节点
		err := g.Register(rpc.NodeRegisterReq{Type: "login", Addr: "ws://A:8001", HTTPAddr: "http://A:9001"})
		convey.So(err, convey.ShouldBeNil)
		err = g.Register(rpc.NodeRegisterReq{Type: "login", Addr: "ws://B:8002", HTTPAddr: "http://B:9002"})
		convey.So(err, convey.ShouldBeNil)
		err = g.Register(rpc.NodeRegisterReq{Type: "chat", Addr: "ws://C:8003", HTTPAddr: "http://C:9003"})
		convey.So(err, convey.ShouldBeNil)

		// ② 路由：login 类型有 2 个
		resp := g.Route(rpc.RouteReq{Type: "login", UId: "user_1"})
		convey.So(resp.Code, convey.ShouldEqual, 0)
		convey.So(resp.Addr, convey.ShouldNotBeEmpty)

		// ③ chat 类型只有 1 个
		resp2 := g.Route(rpc.RouteReq{Type: "chat", UId: "user_2"})
		convey.So(resp2.Addr, convey.ShouldEqual, "ws://C:8003")

		// ④ 不存在的类型
		resp3 := g.Route(rpc.RouteReq{Type: "game", UId: "user_3"})
		convey.So(resp3.Code, convey.ShouldNotEqual, 0)
	})
}

func TestGateway_Heartbeat(t *testing.T) {
	convey.Convey("心跳超时后节点会被路由过滤掉", t, func() {
		g := NewGateway(PolicyRoundRobin)
		g.Register(rpc.NodeRegisterReq{Type: "login", Addr: "ws://A:8001"})
		g.Register(rpc.NodeRegisterReq{Type: "login", Addr: "ws://B:8002"})

		// B 节点心跳停了（模拟超时）
		g.Heartbeat(rpc.NodeHeartbeatReq{Type: "login", Addr: "ws://B:8002", Online: 50})
		// 手动把 B 的 lastHeartbeat 推远
		g.mu.Lock()
		for _, n := range g.nodes["login"] {
			if n.Addr == "ws://B:8002" {
				n.lastHeartbeat = time.Now().Add(-20 * time.Second)
			}
		}
		g.mu.Unlock()

		// 路由只能选到 A
		resp := g.Route(rpc.RouteReq{Type: "login", UId: "x"})
		convey.So(resp.Addr, convey.ShouldEqual, "ws://A:8001")
	})
}

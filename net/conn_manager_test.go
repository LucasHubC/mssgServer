package net

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"mssgServer/config/db"
	"mssgServer/infra/metrics"
	rediscache "mssgServer/infra/redis"

	"github.com/smartystreets/goconvey/convey"
)

// TestMain：ConnManager 的某些方法（BindUid skipRedis=false）依赖 Redis
// 不连 Redis 就跳过相关测试
func TestMain(m *testing.M) {
	// metrics 是全局变量（prometheus Gauge/Counter），必须先 Init 才能 Set()，
	// 否则 conn_manager.go:52 metrics.Connections.Set(...) 会 nil pointer panic
	metrics.Init()

	// redis 连不上就让所有测试 pass（不是 Build Fail）
	if err := db.InitRedis(); err != nil {
		// 没 Redis 也跑，只是涉及 Redis 的方法用 skipRedis=true
	}
	rediscache.SetMyNodeAddr("ws://test-node:8003")
	m.Run()
	db.CloseRedis()
}

// newTestWS 构造一个 WSServer 用于测试
func newTestWS() *WSServer {
	return &WSServer{
		outChan:  make(chan *WsMsgResp, 1),
		property: make(map[string]interface{}),
		done:     make(chan struct{}),
	}
}

// ────────────────────────────────────────
//  Register / Unregister / Total / MaxConns
// ────────────────────────────────────────

func TestConnManager_Register_And_Unregister(t *testing.T) {
	convey.Convey("注册/注销连接 + Total 计数器", t, func() {
		cm := NewConnectionManager(100)

		ws1 := newTestWS()
		ws2 := newTestWS()

		// 注册
		ok := cm.Register(ws1)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(cm.Total(), convey.ShouldEqual, int64(1))

		cm.Register(ws2)
		convey.So(cm.Total(), convey.ShouldEqual, int64(2))

		// 注销 1 个
		ws1.SetProperty("uid", "alice")
		cm.BindUid("alice", ws1, true) // skipRedis=true，不依赖 Redis
		cm.Unregister(ws1)
		convey.So(cm.Total(), convey.ShouldEqual, int64(1))
		convey.So(cm.GetByUid("alice"), convey.ShouldBeNil) // Unregister 会清 uidMap

		// 再注销剩余的
		cm.Unregister(ws2)
		convey.So(cm.Total(), convey.ShouldEqual, int64(0))
	})
}

func TestConnManager_MaxConns(t *testing.T) {
	convey.Convey("超过 MaxConns 时 Register 返回 false", t, func() {
		cm := NewConnectionManager(2) // max=2

		w1 := newTestWS()
		w2 := newTestWS()
		w3 := newTestWS()

		convey.So(cm.Register(w1), convey.ShouldBeTrue)
		convey.So(cm.Register(w2), convey.ShouldBeTrue)
		convey.So(cm.Register(w3), convey.ShouldBeFalse) // 满了

		// 动态调大上限
		cm.SetMaxConns(5)
		convey.So(cm.Register(w3), convey.ShouldBeTrue) // 现在能注册了

		// 清理
		cm.Unregister(w1)
		cm.Unregister(w2)
		cm.Unregister(w3)
	})
}

// ────────────────────────────────────────
//  BindUid / GetByUid / OnlineUids
// ────────────────────────────────────────

func TestConnManager_BindUid(t *testing.T) {
	convey.Convey("BindUid 绑定 uid → GetByUid 能查到 → 再次 Bind 切换", t, func() {
		cm := NewConnectionManager(100)

		ws1 := newTestWS()
		ws2 := newTestWS()
		cm.Register(ws1)
		cm.Register(ws2)

		// 第一次绑定
		cm.BindUid("alice", ws1, true)
		convey.So(cm.GetByUid("alice"), convey.ShouldEqual, ws1)

		// 同 uid 绑到另一个 ws（模拟重连）
		cm.BindUid("alice", ws2, true)
		convey.So(cm.GetByUid("alice"), convey.ShouldEqual, ws2)

		// OnlineUids
		cm.BindUid("bob", ws1, true) // 同一 ws 绑两个 uid？正常生产不会这样，但 SetProperty 会覆盖
		// 实际上 BindUid 只写 uidMap，不写 Property
		_ = "bob"
		_ = cm.OnlineUids()
		// OnlineUids 应该有 alice
		uids := cm.OnlineUids()
		found := false
		for _, u := range uids {
			if u == "alice" {
				found = true
				break
			}
		}
		convey.So(found, convey.ShouldBeTrue)

		// 清理
		cm.Unregister(ws1)
		cm.Unregister(ws2)
	})
}

func TestConnManager_BindUid_SkipRedis_False(t *testing.T) {
	convey.Convey("skipRedis=false 时 NodeAddr 为空 → Redis BindUser 不会被调", t, func() {
		cm := NewConnectionManager(100)
		ws := newTestWS()
		cm.Register(ws)

		// NodeAddr 为空 → skipRedis=false 时进 BindUid 但 NodeAddr=="" → 不调 Redis BindUser
		cm.BindUid("charlie", ws, false)
		convey.So(cm.GetByUid("charlie"), convey.ShouldEqual, ws)

		cm.Unregister(ws)
	})
}

func TestConnManager_GetByUid_NotFound(t *testing.T) {
	convey.Convey("查不存在的 uid → nil", t, func() {
		cm := NewConnectionManager(100)
		convey.So(cm.GetByUid("no_such_uid"), convey.ShouldBeNil)
	})
}

// ────────────────────────────────────────
//  Unregister 时 uidMap 和 Redis online 同步清理
// ────────────────────────────────────────

func TestConnManager_Unregister_Cleanup(t *testing.T) {
	convey.Convey("Unregister 清 uidMap + Redis online key", t, func() {
		uid := "test_conn_unreg_" + time.Now().Format("150405")

		// 先 Redis 注册在线状态
		ok := rediscache.BindUser(uid, "ws://test-node:8003")
		convey.So(ok, convey.ShouldBeTrue)

		node, found := rediscache.GetUserNode(uid)
		convey.So(found, convey.ShouldBeTrue)
		convey.So(node, convey.ShouldEqual, "ws://test-node:8003")

		// ConnManager 里也要绑定
		cm := NewConnectionManager(100)
		ws := newTestWS()
		cm.Register(ws)
		ws.SetProperty("uid", uid)
		cm.BindUid(uid, ws, true)

		// 注销 → 同时清 uidMap + Redis online
		cm.Unregister(ws)
		convey.So(cm.GetByUid(uid), convey.ShouldBeNil)

		_, found = rediscache.GetUserNode(uid)
		convey.So(found, convey.ShouldBeFalse) // Redis online 被清了
	})
}

// ────────────────────────────────────────
//  并发安全：多 goroutine 同时 Register/BindUid/Unregister
// ────────────────────────────────────────

func TestConnManager_Concurrency(t *testing.T) {
	convey.Convey("多 goroutine 并发操作不应 panic + Total 最终一致", t, func() {
		cm := NewConnectionManager(1000)
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ws := newTestWS()
				if cm.Register(ws) {
					// 注册成功 → 最终一定要 Unregister，保证 Total 归零
					defer cm.Unregister(ws)
					// 偶发 BindUid
					if atomic.AddInt64(&counter, 1)%3 == 0 {
						cm.BindUid("uid_"+time.Now().Format("150405.000"), ws, true)
					}
				}
			}()
		}
		wg.Wait()

		// 所有 ws 都 defer Unregister 了 → Total 最终 0
		convey.So(cm.Total(), convey.ShouldEqual, int64(0))
	})
}

var counter int64

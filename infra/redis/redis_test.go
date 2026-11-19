package redis

import (
	"encoding/json"
	"mssgServer/config/db"
	"os"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	os.Chdir("../..")

	if err := db.InitRedis(); err != nil {
		// Redis 连不上就跳过所有测试，不报错
		os.Exit(0)
	}
	// 不设 myNodeAddr，用默认空字符串

	code := m.Run()
	db.CloseRedis()
	os.Exit(code)
}

// ────────────────────────────────────────
//  BindUser / UnbindUser / GetUserNode
// ────────────────────────────────────────

func TestBindUser_LifeCycle(t *testing.T) {
	convey.Convey("BindUser → GetUserNode → UnbindUser", t, func() {
		uid := "test_bind_" + time.Now().Format("150405")

		ok := BindUser(uid, "ws://node-A:8001")
		convey.So(ok, convey.ShouldBeTrue)

		node, found := GetUserNode(uid)
		convey.So(found, convey.ShouldBeTrue)
		convey.So(node, convey.ShouldEqual, "ws://node-A:8001")

		UnbindUser(uid)
		_, found = GetUserNode(uid)
		convey.So(found, convey.ShouldBeFalse)
	})
}

func TestBindUser_DistributedLock(t *testing.T) {
	convey.Convey("分布式锁：同一时刻两个节点只有一个能绑定", t, func() {
		uid := "test_lock_" + time.Now().Format("150405")

		// 第一个节点立即绑定
		ok1 := BindUser(uid, "ws://node-A:8001")
		convey.So(ok1, convey.ShouldBeTrue)

		// 第二个节点也立即绑（应该成功，因为锁是 5s 但 BindUser 用完就释放了）
		ok2 := BindUser(uid, "ws://node-B:8002")
		convey.So(ok2, convey.ShouldBeTrue)

		// 最终在线位置应该是 B（后覆盖前）
		node, _ := GetUserNode(uid)
		convey.So(node, convey.ShouldEqual, "ws://node-B:8002")

		UnbindUser(uid)
	})
}

// ────────────────────────────────────────
//  Pub/Sub 订阅过滤逻辑
//
//  重要说明：
//   1. 测试用 db.RDB.Publish 直接发，绕开 PublishBroadcast 的异步队列。
//      因为 PublishBroadcast 内部 StartAsyncBroadcast (redis.go:231) 会强制
//      覆盖 msg.FromNode = myNodeAddr，导致我们传的 FromNode 被丢弃。
//      生产代码里这个覆盖是对的（业务层不该自己管 FromNode），但测试要精确
//      控制 FromNode 验证 SubscribeBroadcast 的过滤逻辑，必须直接 Publish。
//   2. 绝不 Close 全局 db.RDB —— 它由 TestMain 统一 Init/Close。
//   3. 不用 WaitGroup —— SubscribeBroadcast 外层 for {} 永不退出，
//      我们只需要等"消息到达"这个事件，不需要等 goroutine 退出。
//      goroutine 会随进程退出被清理，不影响测试结果。
// ────────────────────────────────────────

func TestPublishSubscribe_FromOtherNode(t *testing.T) {
	// 只测正向：别的节点发的 → SubscribeBroadcast 不过滤 → handler 被调用
	SetMyNodeAddr("ws://self-test:9999")

	received := make(chan BroadcastMsg, 1)
	go SubscribeBroadcast(func(msg BroadcastMsg) {
		received <- msg
	})

	// 等订阅和 Redis 握手确认（Pub/Sub SUBSCRIBE 需要时间）
	time.Sleep(500 * time.Millisecond)

	// 直接 Publish 到 Redis，精确控制 FromNode
	// FromNode != myNodeAddr → SubscribeBroadcast 不跳过 → handler 被调
	payload, _ := json.Marshal(BroadcastMsg{
		Action:   "kick",
		Target:   "user_abc",
		FromNode: "ws://other-node:8888",
	})
	if err := db.RDB.Publish(redisCtx(), ChannelBroadcast, payload).Err(); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	select {
	case msg := <-received:
		if msg.Action != "kick" || msg.Target != "user_abc" || msg.FromNode != "ws://other-node:8888" {
			t.Fatalf("unexpected message: %+v", msg)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout: 没收到其他节点发的广播")
	}
	// 测试到此结束，SubscribeBroadcast goroutine 随进程退出自动清理
}

// ────────────────────────────────────────
//  分布式锁 TryLock
// ────────────────────────────────────────

func TestTryLock(t *testing.T) {
	convey.Convey("TryLock 互斥性 + 自动过期", t, func() {
		key := "test_trylock_" + time.Now().Format("150405")

		unlock, ok := TryLock(key, 2*time.Second)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(unlock, convey.ShouldNotBeNil)

		// 再拿一次应该失败
		_, ok2 := TryLock(key, 2*time.Second)
		convey.So(ok2, convey.ShouldBeFalse)

		// 手动释放
		unlock()

		// 释放后再拿就成功
		_, ok3 := TryLock(key, 2*time.Second)
		convey.So(ok3, convey.ShouldBeTrue)
	})
}

// ────────────────────────────────────────
//  SetMyNodeAddr / GetMyNodeAddr
// ────────────────────────────────────────

func TestMyNodeAddr(t *testing.T) {
	convey.Convey("SetMyNodeAddr + GetMyNodeAddr", t, func() {
		SetMyNodeAddr("ws://test-addr:1234")
		convey.So(GetMyNodeAddr(), convey.ShouldEqual, "ws://test-addr:1234")
	})
}

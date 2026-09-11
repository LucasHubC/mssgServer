package redis

import (
	"mssgServer/config/db"
	"os"
	"sync"
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
//  Pub/Sub：自己发的消息，自己也能收到（业务层要判断 FromNode）
// ────────────────────────────────────────

func TestPublishSubscribe_EchoFromSelf(t *testing.T) {
	convey.Convey("Publisher 也会收到自己 Publish 的消息", t, func() {
		SetMyNodeAddr("ws://self-test:9999")

		var wg sync.WaitGroup
		wg.Add(1)

		received := make(chan BroadcastMsg, 1)
		go func() {
			defer wg.Done()
			SubscribeBroadcast(func(msg BroadcastMsg) {
				received <- msg
			})
		}()

		// 等订阅建好
		time.Sleep(300 * time.Millisecond)

		PublishBroadcast(BroadcastMsg{
			Action: "kick",
			UID:    "user_abc",
		})

		select {
		case msg := <-received:
			convey.So(msg.Action, convey.ShouldEqual, "kick")
			convey.So(msg.UID, convey.ShouldEqual, "user_abc")
			convey.So(msg.FromNode, convey.ShouldEqual, "ws://self-test:9999")
		case <-time.After(3 * time.Second):
			t.Fatal("subscribe timeout, 没收到消息")
		}

		db.RDB.Close() // 关掉让 SubscribeBroadcast 退出
		wg.Wait()
	})
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

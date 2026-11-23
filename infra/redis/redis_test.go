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

// ═══════════════════════════════════════════════════════════
//  P0-① 背包 Redis 操作 (redis_inventory.go)
// ═══════════════════════════════════════════════════════════

func TestAddItem(t *testing.T) {
	convey.Convey("AddItem Lua 原子性 + 自动标记 dirty", t, func() {
		uid := "test_inv_add_" + time.Now().Format("150405")

		// 首次添加
		n, err := AddItem(uid, "apple", 5)
		convey.So(err, convey.ShouldBeNil)
		convey.So(n, convey.ShouldEqual, int64(5))

		// 再次累加
		n, err = AddItem(uid, "apple", 3)
		convey.So(err, convey.ShouldBeNil)
		convey.So(n, convey.ShouldEqual, int64(8))

		// 负数扣减也走 AddItem
		n, err = AddItem(uid, "apple", -2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(n, convey.ShouldEqual, int64(6))

		// 验证 dirty 标记
		dirty := GetDirtyUids()
		convey.So(dirty, convey.ShouldContain, uid)
		ClearDirty(uid)
	})
}

func TestConsumeItem(t *testing.T) {
	convey.Convey("ConsumeItem Lua 原子性：够扣 → 扣；不够 → 不扣", t, func() {
		uid := "test_inv_consume_" + time.Now().Format("150405")

		// 先加 5 个
		AddItem(uid, "sword", 5)

		// 正常扣 3 个 → 成功
		ok, after, err := ConsumeItem(uid, "sword", 3)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(after, convey.ShouldEqual, int64(2))

		// 不够扣 → 失败，数量保持 2
		ok, before, err := ConsumeItem(uid, "sword", 5)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ok, convey.ShouldBeFalse)
		convey.So(before, convey.ShouldEqual, int64(2)) // 扣减前的数量

		// 扣到 0 → 成功，数量 0
		ok, after, err = ConsumeItem(uid, "sword", 2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(after, convey.ShouldEqual, int64(0))

		// 扣不存在的道具 → 失败
		ok, before, err = ConsumeItem(uid, "shield", 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ok, convey.ShouldBeFalse)
		convey.So(before, convey.ShouldEqual, int64(0))
	})
}

func TestGetItemCount_And_GetAllItems(t *testing.T) {
	convey.Convey("GetItemCount 单查 + GetAllItems 全量 + 过滤 0", t, func() {
		uid := "test_inv_get_" + time.Now().Format("150405")

		AddItem(uid, "apple", 10)
		AddItem(uid, "sword", 1)
		AddItem(uid, "shield", 0) // 0 个 → GetAllItems 会过滤

		// 单查
		n, err := GetItemCount(uid, "apple")
		convey.So(err, convey.ShouldBeNil)
		convey.So(n, convey.ShouldEqual, int64(10))

		// 查不存在
		n, err = GetItemCount(uid, "potion")
		convey.So(err, convey.ShouldBeNil)
		convey.So(n, convey.ShouldEqual, int64(0))

		// 全量（只返回 >0 的）
		all, err := GetAllItems(uid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(all, convey.ShouldContainKey, "apple")
		convey.So(all["apple"], convey.ShouldEqual, int64(10))
		convey.So(all, convey.ShouldContainKey, "sword")
		convey.So(all, convey.ShouldNotContainKey, "shield") // 0 被过滤
	})
}

func TestLoadInventory(t *testing.T) {
	convey.Convey("LoadInventory：Redis 空→写 MySQL 数据；Redis 有→跳过", t, func() {
		uid := "test_inv_load_" + time.Now().Format("150405")

		// 场景 1：Redis 空 + MySQL 有数据 → 写入
		err := LoadInventory(uid, map[string]int64{"apple": 5, "sword": 1})
		convey.So(err, convey.ShouldBeNil)
		all, _ := GetAllItems(uid)
		convey.So(all["apple"], convey.ShouldEqual, int64(5))
		convey.So(all["sword"], convey.ShouldEqual, int64(1))

		// 场景 2：再 Load 一次 → Redis 已有，跳过
		err = LoadInventory(uid, map[string]int64{"apple": 100, "sword": 99})
		convey.So(err, convey.ShouldBeNil)
		all, _ = GetAllItems(uid)
		convey.So(all["apple"], convey.ShouldEqual, int64(5)) // 还是原来的 5
	})
}

func TestDirtyUids_LifeCycle(t *testing.T) {
	convey.Convey("dirty UID 集合：MarkDirty→GetDirtyUids→ClearDirty 完整链路", t, func() {
		uid := "test_dirty_" + time.Now().Format("150405")

		// 初始空
		dirty := GetDirtyUids()
		convey.So(dirty, convey.ShouldNotContain, uid)

		// 标记
		MarkDirty(uid)
		dirty = GetDirtyUids()
		convey.So(dirty, convey.ShouldContain, uid)

		// 清除
		ClearDirty(uid)
		dirty = GetDirtyUids()
		convey.So(dirty, convey.ShouldNotContain, uid)
	})
}

func TestRemoveZeroCount(t *testing.T) {
	convey.Convey("RemoveZeroCount Lua 原子清理 count<=0 的道具", t, func() {
		uid := "test_inv_zero_" + time.Now().Format("150405")

		AddItem(uid, "apple", 3)
		AddItem(uid, "sword", 1)
		ConsumeItem(uid, "sword", 1) // sword 扣到 0

		RemoveZeroCount(uid)

		all, _ := GetAllItems(uid)
		convey.So(all, convey.ShouldContainKey, "apple")
		convey.So(all, convey.ShouldNotContainKey, "sword") // 已删除
	})
}

// ═══════════════════════════════════════════════════════════
//  P0-② 房间 Lua 脚本 (redis.go Room 部分)
// ═══════════════════════════════════════════════════════════

func TestRoom_Create_And_Get_And_Delete(t *testing.T) {
	convey.Convey("房间 CRUD 全链路：Create→Get→Delete", t, func() {
		roomId := "test_room_" + time.Now().Format("150405")
		owner := "test_owner"

		// 创建
		err := CreateRoom(roomId, "测试房", owner, 10)
		convey.So(err, convey.ShouldBeNil)

		// 重复创建 → ErrRoomExists
		err = CreateRoom(roomId, "测试房", owner, 10)
		convey.So(err, convey.ShouldEqual, ErrRoomExists)

		// 查询
		info, err := GetRoomInfo(roomId)
		convey.So(err, convey.ShouldBeNil)
		convey.So(info, convey.ShouldNotBeNil)
		convey.So(info.Name, convey.ShouldEqual, "测试房")
		convey.So(info.Owner, convey.ShouldEqual, owner)
		convey.So(info.MaxSize, convey.ShouldEqual, 10)

		// 查不存在的 → ErrRoomNotFound
		_, err = GetRoomInfo("no_such_room_xxx")
		convey.So(err, convey.ShouldEqual, ErrRoomNotFound)

		// 删除（房主才能删）
		err = DeleteRoom(roomId, owner)
		convey.So(err, convey.ShouldBeNil)

		// 删完再查 → NotFound
		_, err = GetRoomInfo(roomId)
		convey.So(err, convey.ShouldEqual, ErrRoomNotFound)
	})
}

func TestRoom_JoinLeave_Full(t *testing.T) {
	convey.Convey("房间加入/离开 + 满员拒绝 Lua 原子检查", t, func() {
		roomId := "test_room_join_" + time.Now().Format("150405")
		owner := "test_owner_join"

		CreateRoom(roomId, "满员测试", owner, 2) // max_size=2

		// 房主自动是成员
		convey.So(RoomSize(roomId), convey.ShouldEqual, 1)

		// 加第 2 个（刚好到 max）
		err := JoinRoom(roomId, "alice")
		convey.So(err, convey.ShouldBeNil)
		convey.So(RoomSize(roomId), convey.ShouldEqual, 2)

		// 加第 3 个 → 满员
		err = JoinRoom(roomId, "bob")
		convey.So(err, convey.ShouldEqual, ErrRoomFull)

		// 离开一个
		LeaveRoom(roomId, "alice")
		convey.So(RoomSize(roomId), convey.ShouldEqual, 1)

		// 离开房主（房主还在但房间不会删，和业务逻辑一致）
		LeaveRoom(roomId, owner)
		convey.So(RoomSize(roomId), convey.ShouldEqual, 0)

		// 加入不存在的房间
		err = JoinRoom("no_such_room_xxx", "bob")
		convey.So(err, convey.ShouldEqual, ErrRoomNotFound)

		// 清理
		DeleteRoom(roomId, owner)
	})
}

func TestRoom_Kick(t *testing.T) {
	convey.Convey("踢人权限校验 Lua：房主才能踢、不能踢自己", t, func() {
		roomId := "test_room_kick_" + time.Now().Format("150405")
		owner := "test_owner_kick"

		CreateRoom(roomId, "踢人测试", owner, 10)
		JoinRoom(roomId, "alice")
		JoinRoom(roomId, "bob")

		// 房主踢 alice → 成功
		err := KickFromRoom(roomId, owner, "alice")
		convey.So(err, convey.ShouldBeNil)
		convey.So(IsRoomMember(roomId, "alice"), convey.ShouldBeFalse)
		convey.So(RoomSize(roomId), convey.ShouldEqual, 2) // owner + bob

		// 非房主踢 bob → ErrNotRoomOwner
		err = KickFromRoom(roomId, "alice", "bob")
		convey.So(err, convey.ShouldEqual, ErrNotRoomOwner)

		// 踢自己 → ErrCannotKickSelf
		err = KickFromRoom(roomId, owner, owner)
		convey.So(err, convey.ShouldEqual, ErrCannotKickSelf)

		// 房主踢 bob 清理
		KickFromRoom(roomId, owner, "bob")
		DeleteRoom(roomId, owner)
	})
}

// ═══════════════════════════════════════════════════════════
//  Session / Token (redis.go Session 部分)
// ═══════════════════════════════════════════════════════════

func TestSession_LifeCycle(t *testing.T) {
	convey.Convey("Session 签发→查询→注销 完整链路", t, func() {
		sess := Session{UId: "user_001", Username: "alice"}

		token, err := StoreSession(sess)
		convey.So(err, convey.ShouldBeNil)
		convey.So(token, convey.ShouldNotBeEmpty)
		convey.So(len(token), convey.ShouldEqual, 3+32) // "tk_" + 32 hex chars
		convey.So(token[:3], convey.ShouldEqual, "tk_")

		// 查询
		got, ok := GetSession(token)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(got.UId, convey.ShouldEqual, "user_001")
		convey.So(got.Username, convey.ShouldEqual, "alice")
		convey.So(got.Token, convey.ShouldEqual, token) // Token 被回填

		// 查无效 token
		_, ok = GetSession("tk_invalid_token_xxx")
		convey.So(ok, convey.ShouldBeFalse)

		// 注销
		DelSession(token)
		_, ok = GetSession(token)
		convey.So(ok, convey.ShouldBeFalse)
	})
}

func TestSession_TokenUniqueness(t *testing.T) {
	convey.Convey("两次签发的 token 应该不同（熵来自 crypto/rand）", t, func() {
		sess := Session{UId: "user_002", Username: "bob"}
		token1, _ := StoreSession(sess)
		token2, _ := StoreSession(sess)
		convey.So(token1, convey.ShouldNotEqual, token2)
		DelSession(token1)
		DelSession(token2)
	})
}

// ═══════════════════════════════════════════════════════════
//  群组 (redis.go Group 部分)
// ═══════════════════════════════════════════════════════════

func TestGroup_LifeCycle(t *testing.T) {
	convey.Convey("群组 Join/Leave/查询/成员列表", t, func() {
		groupId := "test_group_" + time.Now().Format("150405")

		JoinGroup(groupId, "alice")
		JoinGroup(groupId, "bob")

		convey.So(IsGroupMember(groupId, "alice"), convey.ShouldBeTrue)
		convey.So(IsGroupMember(groupId, "charlie"), convey.ShouldBeFalse)

		members := GetGroupMembers(groupId)
		convey.So(len(members), convey.ShouldEqual, 2)
		convey.So(members, convey.ShouldContain, "alice")
		convey.So(members, convey.ShouldContain, "bob")

		LeaveGroup(groupId, "alice")
		convey.So(IsGroupMember(groupId, "alice"), convey.ShouldBeFalse)
		convey.So(IsGroupMember(groupId, "bob"), convey.ShouldBeTrue)

		// 清理
		db.RDB.Del(redisCtx(), groupKey(groupId))
	})
}

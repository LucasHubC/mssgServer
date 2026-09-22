package dao

import (
	"mssgServer/config/db"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

// TestMain 所有 DAO 测试的入口：初始化 DB + Redis，测试结束后清理
func TestMain(m *testing.M) {
	// 切到项目根目录（保证能找到 conf/conf.ini）
	os.Chdir("../../..")

	// ① 先初始化日志（db.Init 里会调 log.Info）
	log.Init("")

	// ② 再初始化 DB + Redis
	if err := db.Init(); err != nil {
		panic("db.Init failed: " + err.Error())
	}
	if err := db.InitRedis(); err != nil {
		panic("db.InitRedis failed: " + err.Error())
	}
	rediscache.SetMyNodeAddr("ws://test-node:9999")

	code := m.Run()

	db.Close()
	db.CloseRedis()
	os.Exit(code)
}

// newTestUser 生成唯一测试用户，避免冲突
func newTestUser(suffix string) *User {
	ts := time.Now().UnixNano()
	return &User{
		UId:      "test_uid_" + suffix + "_" + itoa(ts),
		Username: "test_" + suffix + "_" + itoa(ts),
		Password: "TestPwd@" + suffix,
	}
}

func itoa(i int64) string { return strconv.FormatInt(i, 10) }

// ────────────────────────────────────────
//  Create + GetByUid + GetByUsername（核心读写路径）
// ────────────────────────────────────────

func TestUserDAO_CreateAndGet(t *testing.T) {
	convey.Convey("Create 后能通过 GetByUid / GetByUsername 查到", t, func() {
		u := newTestUser("create")

		// ① Create
		id, err := UserDAO.Create(u)
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldBeGreaterThan, 0)
		convey.So(u.Id, convey.ShouldEqual, id)

		// ② GetByUid
		got, err := UserDAO.GetByUid(u.UId)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got, convey.ShouldNotBeNil)
		convey.So(got.Username, convey.ShouldEqual, u.Username)
		convey.So(got.Password, convey.ShouldEqual, u.Password)

		// ③ GetByUsername
		got2, err := UserDAO.GetByUsername(u.Username)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got2, convey.ShouldNotBeNil)
		convey.So(got2.UId, convey.ShouldEqual, u.UId)

		// ④ GetById
		got3, err := UserDAO.GetById(id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got3, convey.ShouldNotBeNil)
		convey.So(got3.Id, convey.ShouldEqual, id)

		// ⑤ 清理
		_ = UserDAO.DeleteByUid(u.UId)
	})
}

// ────────────────────────────────────────
//  Redis 缓存命中率（第一次走 DB，第二次走缓存）
// ────────────────────────────────────────

func TestUserDAO_CacheHit(t *testing.T) {
	convey.Convey("GetByUsername 第二次查询应命中缓存", t, func() {
		u := newTestUser("cache")
		_, err := UserDAO.Create(u)
		convey.So(err, convey.ShouldBeNil)
		defer UserDAO.DeleteByUid(u.UId)

		// 第一次：DB 回源 → 自动回写缓存
		_, err = UserDAO.GetByUsername(u.Username)
		convey.So(err, convey.ShouldBeNil)

		// 手动清一下 DB 连接池统计不是关键，直接查缓存 key 是否存在
		key := rediscache.UserCacheKeyByUsername(u.Username)
		cached, ok := rediscache.GetUserCache(key)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(len(cached), convey.ShouldBeGreaterThan, 0)

		// 第二次：直接从缓存取（这里只是验证函数不报错 + 数据一致）
		got, err := UserDAO.GetByUsername(u.Username)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got.Username, convey.ShouldEqual, u.Username)
	})
}

// ────────────────────────────────────────
//  UpdatePassword（验证缓存被正确删除）
// ────────────────────────────────────────

func TestUserDAO_UpdatePassword(t *testing.T) {
	convey.Convey("UpdatePassword 后缓存被清，新密码生效", t, func() {
		u := newTestUser("pwd")
		_, err := UserDAO.Create(u)
		convey.So(err, convey.ShouldBeNil)
		defer UserDAO.DeleteByUid(u.UId)

		// 先查一次，建立缓存
		_, _ = UserDAO.GetByUid(u.UId)

		// 改密码
		newPwd := "NewPwd@999"
		err = UserDAO.UpdatePassword(u.UId, newPwd)
		convey.So(err, convey.ShouldBeNil)

		// 缓存 key 应该被删了
		cacheKey := rediscache.UserCacheKeyByUid(u.UId)
		_, ok := rediscache.GetUserCache(cacheKey)
		convey.So(ok, convey.ShouldBeFalse)

		// 重新查 DB，密码应该是新的
		got, err := UserDAO.GetByUid(u.UId)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got.Password, convey.ShouldEqual, newPwd)
	})
}

// ────────────────────────────────────────
//  UpdateUsername（验证新旧缓存处理）
// ────────────────────────────────────────

func TestUserDAO_UpdateUsername(t *testing.T) {
	convey.Convey("UpdateUsername 后旧缓存清掉，新缓存写入", t, func() {
		u := newTestUser("name")
		_, err := UserDAO.Create(u)
		convey.So(err, convey.ShouldBeNil)
		defer UserDAO.DeleteByUid(u.UId)

		_, _ = UserDAO.GetByUsername(u.Username)

		newName := "newname_" + itoa(time.Now().UnixNano())
		err = UserDAO.UpdateUsername(u.UId, newName)
		convey.So(err, convey.ShouldBeNil)

		// 旧 username 缓存不存在了
		_, ok := rediscache.GetUserCache(rediscache.UserCacheKeyByUsername(u.Username))
		convey.So(ok, convey.ShouldBeFalse)

		// 新 username 缓存存在
		_, ok = rediscache.GetUserCache(rediscache.UserCacheKeyByUsername(newName))
		convey.So(ok, convey.ShouldBeTrue)

		// 用新 username 能查到
		got, err := UserDAO.GetByUsername(newName)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got, convey.ShouldNotBeNil)
		convey.So(got.UId, convey.ShouldEqual, u.UId)
	})
}

// ────────────────────────────────────────
//  Delete（两种方式 + 缓存清理）
// ────────────────────────────────────────

func TestUserDAO_Delete(t *testing.T) {
	convey.Convey("DeleteByUid 后 DB 和缓存都没了", t, func() {
		u := newTestUser("del")
		_, err := UserDAO.Create(u)
		convey.So(err, convey.ShouldBeNil)

		// 建缓存
		_, _ = UserDAO.GetByUid(u.UId)

		err = UserDAO.DeleteByUid(u.UId)
		convey.So(err, convey.ShouldBeNil)

		// DB 中没了
		got, err := UserDAO.GetByUid(u.UId)
		convey.So(err, convey.ShouldBeNil)
		convey.So(got, convey.ShouldBeNil)

		// 缓存也没了
		_, ok := rediscache.GetUserCache(rediscache.UserCacheKeyByUid(u.UId))
		convey.So(ok, convey.ShouldBeFalse)
		_, ok = rediscache.GetUserCache(rediscache.UserCacheKeyByUsername(u.Username))
		convey.So(ok, convey.ShouldBeFalse)
	})
}

// ────────────────────────────────────────
//  ExistsByUid / ExistsByUsername
// ────────────────────────────────────────

func TestUserDAO_Exists(t *testing.T) {
	convey.Convey("Exists 系列函数", t, func() {
		u := newTestUser("exists")
		_, err := UserDAO.Create(u)
		convey.So(err, convey.ShouldBeNil)
		defer UserDAO.DeleteByUid(u.UId)

		exists, err := UserDAO.ExistsByUid(u.UId)
		convey.So(err, convey.ShouldBeNil)
		convey.So(exists, convey.ShouldBeTrue)

		exists, err = UserDAO.ExistsByUsername(u.Username)
		convey.So(err, convey.ShouldBeNil)
		convey.So(exists, convey.ShouldBeTrue)

		exists, err = UserDAO.ExistsByUid("not_exist_uid")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exists, convey.ShouldBeFalse)
	})
}

// ────────────────────────────────────────
//  不存在的用户查询返回 nil, nil（不是报错）
// ────────────────────────────────────────

func TestUserDAO_NotFound(t *testing.T) {
	convey.Convey("查不存在的用户返回 nil, nil", t, func() {
		got, err := UserDAO.GetByUid("definitely_not_exist_12345")
		convey.So(err, convey.ShouldBeNil)
		convey.So(got, convey.ShouldBeNil)

		got, err = UserDAO.GetByUsername("no_such_user_67890")
		convey.So(err, convey.ShouldBeNil)
		convey.So(got, convey.ShouldBeNil)
	})
}

// ────────────────────────────────────────
//  List + Count + Search
// ────────────────────────────────────────

func TestUserDAO_ListAndSearch(t *testing.T) {
	convey.Convey("分页和搜索", t, func() {
		// 先造 3 个用户
		for i := 0; i < 3; i++ {
			u := newTestUser("list")
			_, _ = UserDAO.Create(u)
			defer UserDAO.DeleteByUid(u.UId)
		}

		count, err := UserDAO.Count()
		convey.So(err, convey.ShouldBeNil)
		convey.So(count, convey.ShouldBeGreaterThanOrEqualTo, 3)

		list, total, err := UserDAO.List(1, 2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(total, convey.ShouldBeGreaterThanOrEqualTo, 3)
		convey.So(len(list), convey.ShouldBeLessThanOrEqualTo, 2)

		results, err := UserDAO.Search("test_list")
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(results), convey.ShouldBeGreaterThanOrEqualTo, 3)
	})
}

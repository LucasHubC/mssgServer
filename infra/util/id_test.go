package util

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

// TestRandomHex_Length 基础：长度正确 + 只包含 hex 字符
func TestRandomHex_Length(t *testing.T) {
	convey.Convey("RandomHex(n) 返回 2n 个 hex 字符", t, func() {
		for _, n := range []int{1, 8, 16, 32} {
			s := RandomHex(n)
			convey.So(len(s), convey.ShouldEqual, 2*n)
			for i := 0; i < len(s); i++ {
				c := s[i]
				ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
				convey.So(ok, convey.ShouldBeTrue)
			}
		}
	})
}

// TestRandomHex_Uniqueness 基本唯一性（连续 1000 次调用不应重复）
func TestRandomHex_Uniqueness(t *testing.T) {
	convey.Convey("连续 1000 次调用不应重复（crypto/rand 熵足够）", t, func() {
		seen := make(map[string]bool, 1000)
		for i := 0; i < 1000; i++ {
			s := RandomHex(16)
			convey.So(seen[s], convey.ShouldBeFalse) // 如果重复了就 fail
			seen[s] = true
		}
	})
}

// TestRandomHex_N_Zero 边界：n=0 返回空字符串
func TestRandomHex_N_Zero(t *testing.T) {
	convey.Convey("n=0 返回空字符串", t, func() {
		s := RandomHex(0)
		convey.So(s, convey.ShouldEqual, "")
	})
}

// TestRandomHex_Fallback 降级路径：熵池耗尽时 fallback 也能返回有效 hex
func TestRandomHex_Fallback(t *testing.T) {
	convey.Convey("randomHexFallback 也返回正确格式", t, func() {
		s := randomHexFallback(16)
		convey.So(len(s), convey.ShouldEqual, 32)
		for i := 0; i < len(s); i++ {
			c := s[i]
			ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
			convey.So(ok, convey.ShouldBeTrue)
		}

		// 两次调用应不同（seed 用 time.Now().UnixNano，LFSR 每次 seed 不同）
		time.Sleep(10 * time.Millisecond)
		s2 := randomHexFallback(16)
		convey.So(s, convey.ShouldNotEqual, s2)
	})
}
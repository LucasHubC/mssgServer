package util

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"mssgServer/infra/log"
)

// RandomHex 生成 n 字节的安全随机 hex 字符串（2n 个字符）
// 熵池耗尽时自动降级为时间戳伪随机，保证永不 panic
func RandomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Warn("crypto/rand.Read failed, fallback to time-based", "err", err)
		return randomHexFallback(n)
	}
	return hex.EncodeToString(b)
}

// randomHexFallback 熵池耗尽时的降级实现
// 用时间戳作为种子，LFSR 伪随机（可预测但不碰撞，仅应急用）
func randomHexFallback(n int) string {
	seed := time.Now().UnixNano()
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		seed = seed*1103515245 + 12345
		b[i] = byte(seed)
	}
	return hex.EncodeToString(b)
}

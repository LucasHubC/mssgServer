//go:build windows

package util

import "net"

// ListenTCP 在 Windows 上直接用标准 net.Listen。
//
// 绝对不能用 syscall.Socket → os.NewFile → net.FileListener 这条链路！
// os.NewFile 会注册 GC finalizer，Windows Go runtime 在高并发场景下
// finalizer 有 race，会提前关闭底层 fd → AcceptEx 被中断 → Serve 返回错误
// → log.Fatal → broker 进程崩溃。
//
// 验证过的崩溃样本（压测 1000 并发连接时复现）：
//   accept tcp 127.0.0.1:9001: acceptex: The I/O operation has been aborted
//   Close Code 1006 × 108, ECONNREFUSED 127.0.0.1:9001 × 108
//
// Windows 作为开发/测试环境，net.Listen 默认 backlog (~200) 足够。
// SO_REUSEADDR 在 Windows 上语义不同于 Linux（会允许两个进程绑同一端口），
// 故也不设置。生产 Linux 环境请用 listener_unix.go。
func ListenTCP(addr string, backlog int) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
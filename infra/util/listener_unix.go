//go:build linux || darwin

package util

import (
	"net"
	"os"
	"strconv"
	"syscall"
)

// ListenTCP 在 Linux / macOS 上手动调 syscall，精确控制 backlog + SO_REUSEADDR。
//
// 为什么不直接用 net.ListenConfig？因为不同 Go 版本里：
//   - 有的版本有 SocketListenBacklog 字段
//   - 有的版本字段改名或被移除
//   - 你的 Go 版本两个都不认（编译报错 unknown field Backlog / SocketListenBacklog）
//
// 直接 syscall.Listen(fd, backlog) 能跨所有 Go 版本精确控制 listen(2) 的 backlog。
//
// os.NewFile 在 Linux / macOS 上 GC finalizer 行为稳定（不会被提前触发），
// 所以这条链路在 Unix 环境下安全可靠。Windows 请用 listener_windows.go。
func ListenTCP(addr string, backlog int) (net.Listener, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		ip = net.ParseIP("0.0.0.0")
	}
	isIPv6 := ip.To4() == nil

	var family int
	var sa syscall.Sockaddr
	if isIPv6 {
		family = syscall.AF_INET6
		var in6 syscall.SockaddrInet6
		in6.Port = port
		copy(in6.Addr[:], ip.To16())
		sa = &in6
	} else {
		family = syscall.AF_INET
		var in4 syscall.SockaddrInet4
		in4.Port = port
		copy(in4.Addr[:], ip.To4())
		sa = &in4
	}

	fd, err := syscall.Socket(family, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	// SO_REUSEADDR：进程重启时避免 "address already in use"
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		_ = syscall.Close(fd)
		return nil, err
	}
	if err = syscall.Bind(fd, sa); err != nil {
		_ = syscall.Close(fd)
		return nil, err
	}
	if err = syscall.Listen(fd, backlog); err != nil {
		_ = syscall.Close(fd)
		return nil, err
	}
	file := os.NewFile(uintptr(fd), "")
	return net.FileListener(file)
}
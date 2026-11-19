package util

// ListenTCP 按平台分文件实现：
//   - listener_unix.go  (go:build linux || darwin)  → syscall 可控 backlog + SO_REUSEADDR
//   - listener_windows.go (go:build windows)         → 标准 net.Listen（Windows 上 syscall + os.NewFile 会被 GC 提前关 fd）

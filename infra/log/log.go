package log

import (
	"fmt"
	"io"
	"log/slog"
	"mssgServer/config"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 全局 Logger，业务代码用 log.Info / log.Error 即可
var L *slog.Logger

// Init 从配置读取日志级别并初始化全局 Logger
// 支持两种输出格式：text（控制台友好）和 json（结构化收集）
func Init() {
	levelStr := config.File.MustValue("log", "level", "debug")
	format := config.File.MustValue("log", "format", "text")
	output := config.File.MustValue("log", "output", "console")

	level := parseLevel(levelStr)

	// 解析轮转配置
	rotateOpts := RotateOptions()

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(writer(output, rotateOpts...), &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(writer(output, rotateOpts...), &slog.HandlerOptions{
			Level: level,
		})
	}

	L = slog.New(handler)
	slog.SetDefault(L)
	L.Info("logger initialized", "level", levelStr, "format", format, "output", output)
}

// parseLevel 字符串 → slog.Level
func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// writer 根据 output 决定输出目标
// console → 只输出到控制台
// 文件路径（如 logs/app.log）→ 返回 RotatingWriter（自动按日期+大小切分）
// dual → 同时写控制台和轮转文件
// dual:自定义路径 → 同时写控制台和指定路径的轮转文件
func writer(output string, opts ...RotateOption) io.Writer {
	if output == "console" || output == "" {
		return os.Stdout
	}

	// dual:xxx 同时写控制台 + 指定文件
	if strings.HasPrefix(output, "dual:") {
		path := strings.TrimPrefix(output, "dual:")
		if rw := newRotatingWriter(path, opts...); rw != nil {
			return io.MultiWriter(os.Stdout, rw)
		}
		return os.Stdout
	}

	// dual 同时写控制台 + 默认文件
	if output == "dual" {
		if rw := newRotatingWriter("logs/app.log", opts...); rw != nil {
			return io.MultiWriter(os.Stdout, rw)
		}
		return os.Stdout
	}

	// 文件路径 → 轮转 writer
	if rw := newRotatingWriter(output, opts...); rw != nil {
		return rw
	}
	return os.Stdout
}

// ────────────────────────────────────────────────────────────────
//  RotatingWriter  按日期 + 按大小轮转的 io.Writer
//  写入路径: logs/app.log → 实际写入 logs/app-2026-09-04.log
//  文件超过 maxBytes 时，后续写入 -1, -2 ... 序号文件
//  自动删除超过 maxBackups 天的旧文件
// ────────────────────────────────────────────────────────────────

// RotateOption 轮转参数
var (
	optMaxBytes   int64
	optMaxBackups int
)

// RotateOptions 从 config 读取轮转配置
func RotateOptions() (opts []RotateOption) {
	maxBytes := parseSize(config.File.MustValue("log", "max_size", "100mb"))
	maxBackups := config.File.MustInt("log", "max_backups", 7)
	return []RotateOption{
		func(r *RotatingWriter) { r.maxBytes = maxBytes },
		func(r *RotatingWriter) { r.maxBackups = maxBackups },
	}
}

type RotateOption func(*RotatingWriter)

type RotatingWriter struct {
	basePath   string // 基础路径，如 logs/app.log
	dir        string // 目录 logs/
	baseName   string // 基础名 app
	ext        string // 扩展名 .log
	maxBytes   int64  // 单文件最大字节数
	maxBackups int    // 保留最近 N 天的旧文件

	mu       sync.Mutex
	file     *os.File // 当前打开的文件
	fileSize int64    // 当前文件已写字节数
	curDate  string   // 当前文件的日期后缀，如 2026-09-04
	curSeq   int      // 当前日期下的序号（0 = 当天第一个）
}

func newRotatingWriter(basePath string, opts ...RotateOption) *RotatingWriter {
	rw := &RotatingWriter{
		maxBytes:   100 * 1024 * 1024, // 默认 100MB
		maxBackups: 7,                 // 默认保留 7 天
	}
	// 解析 basePath
	rw.dir = filepath.Dir(basePath)
	ext := filepath.Ext(basePath)
	name := strings.TrimSuffix(filepath.Base(basePath), ext)
	rw.baseName = name
	rw.ext = ext
	rw.basePath = basePath

	for _, o := range opts {
		o(rw)
	}

	// 确保目录存在
	if err := os.MkdirAll(rw.dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "log: create dir %s failed: %v\n", rw.dir, err)
		return nil
	}

	// 打开/创建今天的文件
	rw.rotate()
	return rw
}

// Write 实现 io.Writer
func (rw *RotatingWriter) Write(p []byte) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	today := time.Now().Format("2006-01-02")

	// 日期变了，或当前文件超限，切新文件
	if rw.curDate != today || rw.fileSize+int64(len(p)) > rw.maxBytes {
		rw.rotate()
	}

	n, err := rw.file.Write(p)
	rw.fileSize += int64(n)
	return n, err
}

// rotate 切换到下一个文件
func (rw *RotatingWriter) rotate() {
	// 关闭旧文件
	if rw.file != nil {
		rw.file.Close()
	}

	today := time.Now().Format("2006-01-02")

	// 找今天的下一个序号
	seq := 0
	for {
		fname := rw.fileName(today, seq)
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			break
		}
		seq++
	}

	path := rw.fileName(today, seq)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: open rotating file %s failed: %v\n", path, err)
		return
	}

	rw.file = f
	rw.curDate = today
	rw.curSeq = seq
	rw.fileSize = fileSize(f)

	// 清理过期旧文件
	rw.cleanup()
}

// fileName 生成带日期+序号的文件名
func (rw *RotatingWriter) fileName(date string, seq int) string {
	if seq == 0 {
		return filepath.Join(rw.dir, fmt.Sprintf("%s-%s%s", rw.baseName, date, rw.ext))
	}
	return filepath.Join(rw.dir, fmt.Sprintf("%s-%s-%d%s", rw.baseName, date, seq, rw.ext))
}

// cleanup 删除超过 maxBackups 天的旧日志文件
func (rw *RotatingWriter) cleanup() {
	if rw.maxBackups <= 0 {
		return
	}

	entries, err := os.ReadDir(rw.dir)
	if err != nil {
		return
	}

	// 匹配 app-YYYY-MM-DD.log 或 app-YYYY-MM-DD-N.log
	re := regexp.MustCompile(fmt.Sprintf(`^%s-(\d{4}-\d{2}-\d{2})(?:-\d+)?%s$`, regexp.QuoteMeta(rw.baseName), regexp.QuoteMeta(rw.ext)))

	type oldFile struct {
		path string
		date time.Time
	}
	var olds []oldFile

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := re.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		t, err := time.Parse("2006-01-02", m[1])
		if err != nil {
			continue
		}
		olds = append(olds, oldFile{path: filepath.Join(rw.dir, e.Name()), date: t})
	}

	// 按日期升序（最旧的在前面）
	sort.Slice(olds, func(i, j int) bool { return olds[i].date.Before(olds[j].date) })

	// 删除最旧的，只保留最近 maxBackups 天
	cutoff := time.Now().AddDate(0, 0, -rw.maxBackups)
	for _, f := range olds {
		if f.date.Before(cutoff) {
			os.Remove(f.path)
		}
	}
}

// Close 关闭当前文件
func (rw *RotatingWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

// fileSize 获取文件当前大小
func fileSize(f *os.File) int64 {
	info, err := f.Stat()
	if err != nil {
		return 0
	}
	return info.Size()
}

// parseSize 解析大小字符串 "100mb" / "1gb" / "1024kb" → 字节数
func parseSize(s string) int64 {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 100 * 1024 * 1024 // 默认 100MB
	}

	var unit int64 = 1
	var numStr string
	switch {
	case strings.HasSuffix(s, "gb"):
		unit = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(s, "gb")
	case strings.HasSuffix(s, "mb"):
		unit = 1024 * 1024
		numStr = strings.TrimSuffix(s, "mb")
	case strings.HasSuffix(s, "kb"):
		unit = 1024
		numStr = strings.TrimSuffix(s, "kb")
	default:
		// 纯数字，当作字节
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 100 * 1024 * 1024
		}
		return n
	}

	n, err := strconv.ParseFloat(strings.TrimSpace(numStr), 64)
	if err != nil {
		return 100 * 1024 * 1024
	}
	return int64(n * float64(unit))
}

// ────────────────────────────────────────────────────────────────
//  快捷函数（带 format 支持，比 slog 的 With 更顺手）
//  用法：
//    log.Info("server start", "addr", addr)
//    log.Error("db query failed", "err", err, "sql", sql)
//  如果 L 没初始化，自动降级用 slog.Default()（不会 panic）
// ────────────────────────────────────────────────────────────────

func defaultLogger() *slog.Logger {
	if L != nil {
		return L
	}
	return slog.Default()
}

func Debug(msg string, args ...any) { defaultLogger().Debug(msg, args...) }
func Info(msg string, args ...any)  { defaultLogger().Info(msg, args...) }
func Warn(msg string, args ...any)  { defaultLogger().Warn(msg, args...) }
func Error(msg string, args ...any) { defaultLogger().Error(msg, args...) }

// Fatal 输出 error 日志后退出进程
func Fatal(msg string, args ...any) {
	L.Error(msg, args...)
	os.Exit(1)
}

// With 返回带预置 attr 的子 logger
// logger := log.With("module", "shop", "uid", "1001")
// logger.Info("buy item", "itemId", 1)
func With(args ...any) *slog.Logger {
	return L.With(args...)
}

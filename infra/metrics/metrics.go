package metrics

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ────────────────────────────────────────────────────
// Prometheus Metrics — mssgServer 全链路埋点  监控体系
// 如果你以后要部署 Prometheus Server,请注释掉这行
//只需要下载一个独立二进制,写个 prometheus.yml:
/**
	global:
	scrape_interval: 15s      # 全局默认 15s 拉一次,够用了
	evaluation_interval: 15s

	scrape_configs:
	- job_name: 'broker'
		static_configs:
		- targets: ['127.0.0.1:9001']

	- job_name: 'gateway'
		static_configs:
		- targets: ['127.0.0.1:9000']

	- job_name: 'chat'
		static_configs:
		- targets: ['127.0.0.1:8001']


下载 + 启动(Prometheus 官方 GitHub Releases,一个 exe 就搞定):

# 下载(选 prometheus-windows-amd64.zip)
# https://github.com/prometheus/prometheus/releases

# 启动
.\prometheus.exe --config.file=prometheus.yml

# 浏览器打开 UI
curl http://127.0.0.1:9090

**/

// ────────────────────────────────────────────────────

var (
	// ── Gauge（可升可降）──

	// Connections 当前在线连接数
	Connections prometheus.Gauge

	// ActiveHandlers 正在执行的 handler 数
	ActiveHandlers prometheus.Gauge

	// UIDMapSize 当前绑定的 uid→ws 映射数
	UIDMapSize prometheus.Gauge

	// Goroutines 当前 goroutine 数（定时刷新）
	Goroutines prometheus.Gauge

	// HeapAlloc 当前堆内存分配（定时刷新，字节）
	HeapAlloc prometheus.Gauge

	// ── Counter（只增不减）──

	// MessagesIn 收到的消息总数
	MessagesIn prometheus.Counter

	// MessagesOut 发出的消息总数
	MessagesOut prometheus.Counter

	// ErrorsTotal 错误总数（按 type 分类）
	ErrorsTotal *prometheus.CounterVec

	// DroppedMessages 丢包数（outChan 满时）
	DroppedMessages *prometheus.CounterVec

	// BytesIn 入站字节总数
	BytesIn prometheus.Counter

	// BytesOut 出站字节总数
	BytesOut prometheus.Counter

	// HandlersPanic handler panic 次数
	HandlersPanic prometheus.Counter

	// ── Histogram（分布）──

	// HandlerDuration handler 执行耗时
	HandlerDuration prometheus.Histogram

	// WriteDuration WriteMessage 耗时
	WriteDuration prometheus.Histogram

	// BatchSize 每次 batch 合并条数
	BatchSize prometheus.Histogram

	// ConnectionLifetime 连接存活时长
	ConnectionLifetime prometheus.Histogram

	// LatencyEndToEnd 端到端延迟（客户端 send → 广播落地）
	LatencyEndToEnd prometheus.Histogram
)

var once sync.Once

// Init 所有指标（幂等，只调一次）
func Init() {
	once.Do(func() {
		Connections = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ws_connections_total",
			Help: "Current number of WebSocket connections",
		})

		ActiveHandlers = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ws_handlers_active",
			Help: "Current number of handlers executing",
		})

		UIDMapSize = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "ws_uid_map_size",
			Help: "Current number of bound uid→ws mappings",
		})

		Goroutines = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "process_goroutines",
			Help: "Current number of goroutines",
		})

		HeapAlloc = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "process_heap_alloc_bytes",
			Help: "Current heap memory allocation in bytes",
		})

		MessagesIn = promauto.NewCounter(prometheus.CounterOpts{
			Name: "ws_messages_in_total",
			Help: "Total messages received",
		})

		MessagesOut = promauto.NewCounter(prometheus.CounterOpts{
			Name: "ws_messages_out_total",
			Help: "Total messages sent",
		})

		ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "ws_errors_total",
			Help: "Total errors by type",
		}, []string{"type"})

		DroppedMessages = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "ws_dropped_messages_total",
			Help: "Total dropped messages by reason",
		}, []string{"reason"})

		BytesIn = promauto.NewCounter(prometheus.CounterOpts{
			Name: "ws_bytes_in_total",
			Help: "Total bytes received",
		})

		BytesOut = promauto.NewCounter(prometheus.CounterOpts{
			Name: "ws_bytes_out_total",
			Help: "Total bytes sent",
		})

		HandlersPanic = promauto.NewCounter(prometheus.CounterOpts{
			Name: "ws_handlers_panic_total",
			Help: "Total handler panic recoveries",
		})

		HandlerDuration = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "ws_handler_duration_seconds",
			Help:    "Handler execution duration",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.002, 0.005, 0.01, 0.02, 0.05, 0.1, 0.25, 0.5, 1},
		})

		WriteDuration = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "ws_write_duration_seconds",
			Help:    "WriteMessage duration (Marshal + Write)",
			Buckets: []float64{0.00001, 0.00005, 0.0001, 0.0005, 0.001, 0.002, 0.005, 0.01, 0.02, 0.05, 0.1},
		})

		BatchSize = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "ws_batch_size",
			Help:    "Number of messages per batch flush",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128},
		})

		ConnectionLifetime = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "ws_connection_lifetime_seconds",
			Help:    "Connection lifetime from connect to disconnect",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
		})

		LatencyEndToEnd = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "ws_latency_e2e_seconds",
			Help:    "End-to-end latency (client send → broadcast delivered)",
			Buckets: []float64{0.001, 0.002, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1},
		})

		// 后台刷新 runtime 指标（每秒一次）
		go runtimeCollector()
	})
}

// runtimeCollector 定时刷新 Goroutines / HeapAlloc
func runtimeCollector() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		Goroutines.Set(float64(runtime.NumGoroutine()))
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		HeapAlloc.Set(float64(m.HeapAlloc))
	}
}

// RegisterHTTP 注册 Prometheus metrics handler + pprof + health
// 在 http.DefaultServeMux 上注册，服务端直接 http.ListenAndServe 就能用
func RegisterHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		_, _ = w.Write([]byte(`{"status":"ok","goroutines":` + itoa(runtime.NumGoroutine()) +
			`,"heap_alloc":` + itoa(int(m.HeapAlloc)) +
			`,"timestamp":` + itoa(int(time.Now().Unix())) + `}`))
	})
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

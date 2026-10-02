package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. 当前在线WS连接总数（仪表盘，可增减）
var OnlineWSConn = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "im_ws_online_connections",
		Help: "Current online websocket client number",
	},
)

// 2. 消息发送成功总计数
var MsgSendTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "im_msg_send_total",
		Help: "Total successful push message count",
	},
)

// 3. 消息推送失败总计数
var MsgSendFailTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "im_msg_send_fail_total",
		Help: "Total failed push message count",
	},
)

// 4. HTTP接口耗时直方图（区分路径+请求方法）
var HttpReqDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "im_http_request_duration_seconds",
		Help:    "Http request latency distribution",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
	},
	[]string{"path", "method"},
)

// MetricsHandler 对外暴露 /metrics 接口
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// HttpMetricMiddleware Gin全局监控中间件，自动统计接口耗时
func HttpMetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start).Seconds()
		path := c.FullPath()
		method := c.Request.Method
		HttpReqDuration.WithLabelValues(path, method).Observe(dur)
	}
}

package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// BucketConfig is used to initialise bucket
type BucketConfig struct {
	Start float64
	Width float64
	Count int
}

// timers
var (
	dbQueryTimer             *prometheus.HistogramVec
	httpRequestTimer         *prometheus.HistogramVec
	externalHTTPRequestTimer *prometheus.HistogramVec
)

// counters
var (
	dbQueryCounter             *prometheus.CounterVec
	httpRequestCounter         *prometheus.CounterVec
	httpTotalRequestCounter    *prometheus.CounterVec
	httpResponseStatusCounter  *prometheus.CounterVec
	externalHTTPRequestCounter *prometheus.CounterVec
)

// Init is used to initialise metrics
func Init(cfg BucketConfig) {
	httpRequestTimer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "httpResponseTimeSeconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.LinearBuckets(cfg.Start, cfg.Width, cfg.Count),
	}, []string{"path"})

	externalHTTPRequestTimer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "externalHttpRequestTimer",
		Help:    "Duration of External HTTP requests.",
		Buckets: prometheus.LinearBuckets(cfg.Start, cfg.Width, cfg.Count),
	}, []string{"name"})

	dbQueryTimer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "DAOTime",
		Help:    "DAO time",
		Buckets: prometheus.LinearBuckets(cfg.Start, cfg.Width, cfg.Count),
	}, []string{"query"})

	httpTotalRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "httpRequestsTotal",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	dbQueryCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dbRequestsSents",
			Help: "How many DB requests processed, partitioned by query function used",
		},
		[]string{"function"},
	)

	httpRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "routerRequestsTotal",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"status", "host", "path", "method"},
	)

	httpResponseStatusCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "responseStatus",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	externalHTTPRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "externalCount",
			Help: "How many external HTTP requests are processed.",
		},
		[]string{"status", "url", "method"},
	)
}

// GetMetricsMiddleware is to add prometheus timer and counter stats for requests
func GetMetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler := GetHandlerName(ctx)
		timer := prometheus.NewTimer(httpRequestTimer.WithLabelValues(handler))
		httpTotalRequestCounter.WithLabelValues(handler).Inc()
		ctx.Next()
		httpResponseStatusCounter.WithLabelValues(strconv.Itoa(ctx.Writer.Status())).Inc()
		httpRequestCounter.WithLabelValues(strconv.Itoa(ctx.Writer.Status()), ctx.Request.URL.RequestURI(), handler,
			ctx.Request.Method).Inc()
		timer.ObserveDuration()
	}
}

// GetDBQueryTimer is to log the time taken to run a database request
func GetDBQueryTimer(name string) *prometheus.Timer {
	dbQueryCounter.WithLabelValues(name).Inc()
	return prometheus.NewTimer(dbQueryTimer.WithLabelValues(name))
}

// GetExternalHTTPRequestTimer is to log the time taken to run a database request
func GetExternalHTTPRequestTimer(url string) *prometheus.Timer {
	return prometheus.NewTimer(externalHTTPRequestTimer.WithLabelValues(url))
}

// GetExternalHTTPRequestTimer is to log the time taken to run a database request
func SetExternalHTTPRequestCounter(status, url, method string) {
	externalHTTPRequestCounter.WithLabelValues(status, url, method).Inc()
}

// HTTPMetrics is the wrapper to add metrics to HTTP requests
func HTTPMetrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func GetHandlerName(ctx *gin.Context) string {
	handlers := ctx.HandlerNames()
	l := len(handlers)
	if l == 0 {
		return ctx.Request.URL.RequestURI()
	}
	return handlers[l-1]
}

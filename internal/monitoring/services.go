package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

// мидлвар для сбора метрик
func (m *Metrics) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.RequestsInFlight.Inc()
		defer m.RequestsInFlight.Dec()

		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	}
}

// запись операции
func (m *Metrics) RecordOperation(operationType string) {
	m.OperationsTotal.WithLabelValues(operationType).Inc()
}

// запись ошибки
func (m *Metrics) RecordError(errorType, handler string) {
	m.ErrorsTotal.WithLabelValues(errorType, handler).Inc()
}

// регистриция эндпоинта для Прометея
func (m *Metrics) RegisterMetricsHandler(router *gin.Engine, path string) {
	router.GET(path, func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
}

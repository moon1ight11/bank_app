package middleware

import (
	"bank_app/internal/monitoring"
	"github.com/gin-gonic/gin"
)

// мидлвар для сбора метрик
func MetricsMiddleware(metrics *monitoring.Metrics) gin.HandlerFunc {
	return metrics.GinMiddleware()
}

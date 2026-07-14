package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimiter(rps float64, burst int) gin.HandlerFunc {
	limiters := make(map[string]*limiterEntry)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			for ip, entry := range limiters {
				if time.Since(entry.lastSeen) > 30*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		entry, exists := limiters[ip]
		if !exists {
			entry = &limiterEntry{
				limiter:  rate.NewLimiter(rate.Limit(rps), burst),
				lastSeen: time.Now(),
			}
			limiters[ip] = entry
		} else {
			entry.lastSeen = time.Now()
		}
		mu.Unlock()

		if !entry.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

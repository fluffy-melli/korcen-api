// internal/middleware/middleware.go

package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	capacity   int
	tokens     float64
	refillRate float64
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) AllowRequest() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	return false
}

var tokenBuckets sync.Map

func getOrCreateTokenBucket(ip string, capacity int, refillRate float64) *TokenBucket {
	if v, ok := tokenBuckets.Load(ip); ok {
		return v.(*TokenBucket)
	}
	newBucket := NewTokenBucket(capacity, refillRate)
	actual, _ := tokenBuckets.LoadOrStore(ip, newBucket)
	return actual.(*TokenBucket)
}

func TokenBucketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		tb := getOrCreateTokenBucket(ip, 100, 100.0/60.0) // 1분당 100개 (초당 약 1.67개 리필)

		if !tb.AllowRequest() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

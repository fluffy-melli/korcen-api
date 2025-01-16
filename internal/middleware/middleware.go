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

var tokenBuckets = struct {
	sync.Mutex
	buckets map[string]*TokenBucket
}{buckets: make(map[string]*TokenBucket)}

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

func TokenBucketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		tokenBuckets.Lock()
		if _, exists := tokenBuckets.buckets[ip]; !exists {
			tokenBuckets.buckets[ip] = NewTokenBucket(100, 100.0/60.0) // 1분당 100개 (초당 약 1.67개 리필)
		}
		tb := tokenBuckets.buckets[ip]
		tokenBuckets.Unlock()

		if !tb.AllowRequest() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

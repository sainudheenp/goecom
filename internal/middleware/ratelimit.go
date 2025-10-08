package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	requests      int
	windowMinutes int
	clients       map[string]*clientBucket
	mu            sync.RWMutex
}

type clientBucket struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests, windowMinutes int) *RateLimiter {
	limiter := &RateLimiter{
		requests:      requests,
		windowMinutes: windowMinutes,
		clients:       make(map[string]*clientBucket),
	}

	// Clean up old entries periodically
	go limiter.cleanup()

	return limiter
}

// Middleware returns a Gin middleware function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !rl.allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"details": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// allow checks if a request is allowed
func (rl *RateLimiter) allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[clientIP]

	if !exists {
		rl.clients[clientIP] = &clientBucket{
			tokens:    rl.requests - 1,
			lastReset: now,
		}
		return true
	}

	// Reset bucket if window has passed
	if now.Sub(bucket.lastReset) >= time.Duration(rl.windowMinutes)*time.Minute {
		bucket.tokens = rl.requests - 1
		bucket.lastReset = now
		return true
	}

	// Check if tokens available
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup removes old client entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientIP, bucket := range rl.clients {
			if now.Sub(bucket.lastReset) >= time.Duration(rl.windowMinutes*2)*time.Minute {
				delete(rl.clients, clientIP)
			}
		}
		rl.mu.Unlock()
	}
}

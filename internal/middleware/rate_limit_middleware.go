package middleware

import (
	"sync"
	"time"
	"venturo-core/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors    map[string]*visitor
	mu          sync.RWMutex
	rate        int           // requests per window
	window      time.Duration // time window
	cleanupTick time.Duration // cleanup interval
}

type visitor struct {
	requests  int
	lastSeen  time.Time
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:    make(map[string]*visitor),
		rate:        rate,
		window:      window,
		cleanupTick: time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()
	return rl
}

// cleanup removes old entries periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupTick)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.lastSeen) > rl.window*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]

	if !exists {
		rl.visitors[ip] = &visitor{
			requests:  1,
			lastSeen:  now,
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	v.lastSeen = now

	// Reset window if expired
	if now.After(v.windowEnd) {
		v.requests = 1
		v.windowEnd = now.Add(rl.window)
		return true
	}

	// Check if within rate limit
	if v.requests >= rl.rate {
		return false
	}

	v.requests++
	return true
}

// NewRateLimitMiddleware creates a rate limiting middleware
func NewRateLimitMiddleware(requestsPerMinute int) fiber.Handler {
	limiter := NewRateLimiter(requestsPerMinute, time.Minute)

	return func(c *fiber.Ctx) error {
		// Get client IP (consider X-Forwarded-For in production)
		ip := c.IP()
		
		if !limiter.Allow(ip) {
			return response.Error(c, fiber.StatusTooManyRequests, 
				fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded. Please try again later."))
		}

		return c.Next()
	}
}
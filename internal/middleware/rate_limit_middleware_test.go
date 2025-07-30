package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Create a test app with rate limit of 2 requests per minute
	app := fiber.New()
	app.Use(NewRateLimitMiddleware(2))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name           string
		requestCount   int
		expectedStatus int
		description    string
	}{
		{
			name:           "First request should pass",
			requestCount:   1,
			expectedStatus: 200,
			description:    "Within rate limit",
		},
		{
			name:           "Second request should pass",
			requestCount:   2,
			expectedStatus: 200,
			description:    "Still within rate limit",
		},
		{
			name:           "Third request should fail",
			requestCount:   3,
			expectedStatus: 429,
			description:    "Exceeds rate limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make the specified number of requests
			var resp *http.Response
			var err error

			for i := 0; i < tt.requestCount; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Forwarded-For", "192.168.1.1") // Simulate same IP
				resp, err = app.Test(req)
				if err != nil {
					t.Fatalf("Request failed: %v", err)
				}

				// Only check the last response
				if i == tt.requestCount-1 {
					if resp.StatusCode != tt.expectedStatus {
						t.Errorf("Expected status %d, got %d (%s)", tt.expectedStatus, resp.StatusCode, tt.description)
					}
				}
				resp.Body.Close()
			}
		})
	}
}

func TestRateLimitMiddleware_DifferentIPs(t *testing.T) {
	// Create a test app with rate limit of 1 request per minute
	app := fiber.New()
	app.Use(NewRateLimitMiddleware(1))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test that different IPs have separate rate limits
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	for _, ip := range ips {
		t.Run("IP_"+ip, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Forwarded-For", ip)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// All first requests from different IPs should succeed
			if resp.StatusCode != 200 {
				t.Errorf("Request from IP %s should succeed, got status %d", ip, resp.StatusCode)
			}
		})
	}
}

func TestRateLimitMiddleware_WindowReset(t *testing.T) {
	// Create a rate limiter with very short window for testing
	limiter := NewRateLimiter(1, 100*time.Millisecond) // 1 request per 100ms

	// Simulate first request
	if !limiter.Allow("192.168.1.1") {
		t.Error("First request should be allowed")
	}

	// Immediate second request should be blocked
	if limiter.Allow("192.168.1.1") {
		t.Error("Second immediate request should be blocked")
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Request after window reset should be allowed
	if !limiter.Allow("192.168.1.1") {
		t.Error("Request after window reset should be allowed")
	}
}

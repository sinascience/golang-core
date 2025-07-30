package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// SecurityEvent defines different types of security events
type SecurityEventType string

const (
	EventUserRegistration   SecurityEventType = "user_registration"
	EventUserLogin          SecurityEventType = "user_login"
	EventLoginFailed        SecurityEventType = "login_failed"
	EventTokenRefresh       SecurityEventType = "token_refresh"
	EventTokenRefreshFailed SecurityEventType = "token_refresh_failed"
	EventRateLimitExceeded  SecurityEventType = "rate_limit_exceeded"
	EventUnauthorizedAccess SecurityEventType = "unauthorized_access"
)

// SecurityLogger provides structured logging for security events
type SecurityLogger struct{}

// NewSecurityLogger creates a new security logger
func NewSecurityLogger() *SecurityLogger {
	return &SecurityLogger{}
}

// LogSecurityEvent logs a security-related event with structured data
func (sl *SecurityLogger) LogSecurityEvent(ctx context.Context, eventType SecurityEventType, userID *uuid.UUID, ip string, details map[string]interface{}) {
	// Safely get timestamp from context or use current time
	var timestamp time.Time
	if reqTime := ctx.Value("request_time"); reqTime != nil {
		timestamp = reqTime.(time.Time)
	} else {
		timestamp = time.Now()
	}
	
	attrs := []slog.Attr{
		slog.String("event_type", string(eventType)),
		slog.String("ip", ip),
		slog.Time("timestamp", timestamp),
	}
	
	if userID != nil {
		attrs = append(attrs, slog.String("user_id", userID.String()))
	}
	
	// Add additional details
	for key, value := range details {
		attrs = append(attrs, slog.Any(key, value))
	}
	
	// Determine log level based on event type
	var level slog.Level
	switch eventType {
	case EventLoginFailed, EventTokenRefreshFailed, EventRateLimitExceeded, EventUnauthorizedAccess:
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}
	
	slog.LogAttrs(ctx, level, "Security Event", attrs...)
}

// LogUserRegistration logs a user registration event
func (sl *SecurityLogger) LogUserRegistration(ctx context.Context, userID uuid.UUID, email, ip string) {
	sl.LogSecurityEvent(ctx, EventUserRegistration, &userID, ip, map[string]interface{}{
		"email": email,
	})
}

// LogSuccessfulLogin logs a successful login event
func (sl *SecurityLogger) LogSuccessfulLogin(ctx context.Context, userID uuid.UUID, email, ip string) {
	sl.LogSecurityEvent(ctx, EventUserLogin, &userID, ip, map[string]interface{}{
		"email": email,
	})
}

// LogFailedLogin logs a failed login attempt
func (sl *SecurityLogger) LogFailedLogin(ctx context.Context, email, ip, reason string) {
	sl.LogSecurityEvent(ctx, EventLoginFailed, nil, ip, map[string]interface{}{
		"email":  email,
		"reason": reason,
	})
}

// LogTokenRefresh logs a token refresh event
func (sl *SecurityLogger) LogTokenRefresh(ctx context.Context, userID uuid.UUID, ip string, success bool) {
	eventType := EventTokenRefresh
	if !success {
		eventType = EventTokenRefreshFailed
	}
	
	sl.LogSecurityEvent(ctx, eventType, &userID, ip, map[string]interface{}{
		"success": success,
	})
}

// LogRateLimitExceeded logs when rate limit is exceeded
func (sl *SecurityLogger) LogRateLimitExceeded(ctx context.Context, ip, endpoint string) {
	sl.LogSecurityEvent(ctx, EventRateLimitExceeded, nil, ip, map[string]interface{}{
		"endpoint": endpoint,
	})
}

// LogUnauthorizedAccess logs unauthorized access attempts
func (sl *SecurityLogger) LogUnauthorizedAccess(ctx context.Context, ip, endpoint, reason string) {
	sl.LogSecurityEvent(ctx, EventUnauthorizedAccess, nil, ip, map[string]interface{}{
		"endpoint": endpoint,
		"reason":   reason,
	})
}
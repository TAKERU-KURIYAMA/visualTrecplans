package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/visualtrecplans/backend/pkg/logger"
)

// EventType represents different types of audit events
type EventType string

const (
	EventLogin          EventType = "login"
	EventLoginFailed    EventType = "login_failed"
	EventLogout         EventType = "logout"
	EventRegister       EventType = "register"
	EventRegisterFailed EventType = "register_failed"
	EventPasswordChange EventType = "password_change"
	EventProfileUpdate  EventType = "profile_update"
	EventTokenRefresh   EventType = "token_refresh"
	EventAccountLock    EventType = "account_lock"
	EventAccountUnlock  EventType = "account_unlock"
	EventEmailVerify    EventType = "email_verify"
	EventPasswordReset  EventType = "password_reset"
)

// AuditEvent represents an audit log entry
type AuditEvent struct {
	ID        uuid.UUID              `json:"id"`
	EventType EventType              `json:"event_type"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Email     string                 `json:"email,omitempty"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Success   bool                   `json:"success"`
	Reason    string                 `json:"reason,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// AuditLogger handles audit logging
type AuditLogger struct {
	// In a real implementation, this might write to a database,
	// external logging service, or file
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{}
}

// LogEvent logs an audit event
func (al *AuditLogger) LogEvent(ctx context.Context, event *AuditEvent) {
	// Set ID and timestamp if not provided
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Convert to JSON for structured logging
	eventJSON, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal audit event", 
			logger.Error(err),
			logger.String("event_type", string(event.EventType)),
		)
		return
	}

	// Log the event using structured logging
	logger.Info("Audit event", 
		logger.String("audit_event", string(eventJSON)),
		logger.String("event_type", string(event.EventType)),
		logger.String("user_id", getUserIDString(event.UserID)),
		logger.String("email", event.Email),
		logger.String("ip", event.IP),
		logger.Bool("success", event.Success),
	)

	// In production, you might also want to:
	// 1. Write to a separate audit database
	// 2. Send to a SIEM system
	// 3. Write to a secure log file
	// 4. Send to an external audit service
}

// LogLoginSuccess logs a successful login
func (al *AuditLogger) LogLoginSuccess(ctx context.Context, userID uuid.UUID, email, ip, userAgent string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventLogin,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})
}

// LogLoginFailed logs a failed login attempt
func (al *AuditLogger) LogLoginFailed(ctx context.Context, email, ip, userAgent, reason string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventLoginFailed,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   false,
		Reason:    reason,
	})
}

// LogLogout logs a logout event
func (al *AuditLogger) LogLogout(ctx context.Context, userID uuid.UUID, email, ip, userAgent string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventLogout,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})
}

// LogRegistration logs a successful registration
func (al *AuditLogger) LogRegistration(ctx context.Context, userID uuid.UUID, email, ip, userAgent string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventRegister,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})
}

// LogRegistrationFailed logs a failed registration attempt
func (al *AuditLogger) LogRegistrationFailed(ctx context.Context, email, ip, userAgent, reason string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventRegisterFailed,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   false,
		Reason:    reason,
	})
}

// LogPasswordChange logs a password change
func (al *AuditLogger) LogPasswordChange(ctx context.Context, userID uuid.UUID, email, ip, userAgent string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventPasswordChange,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})
}

// LogProfileUpdate logs a profile update
func (al *AuditLogger) LogProfileUpdate(ctx context.Context, userID uuid.UUID, email, ip, userAgent string, changes map[string]interface{}) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventProfileUpdate,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
		Metadata:  map[string]interface{}{"changes": changes},
	})
}

// LogTokenRefresh logs a token refresh
func (al *AuditLogger) LogTokenRefresh(ctx context.Context, userID uuid.UUID, email, ip, userAgent string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventTokenRefresh,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})
}

// LogAccountLock logs an account being locked
func (al *AuditLogger) LogAccountLock(ctx context.Context, userID uuid.UUID, email, ip, userAgent, reason string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventAccountLock,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
		Reason:    reason,
	})
}

// LogEmailVerification logs email verification
func (al *AuditLogger) LogEmailVerification(ctx context.Context, userID uuid.UUID, email, ip, userAgent string) {
	al.LogEvent(ctx, &AuditEvent{
		EventType: EventEmailVerify,
		UserID:    &userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   true,
	})
}

// getUserIDString safely converts UUID pointer to string
func getUserIDString(userID *uuid.UUID) string {
	if userID == nil {
		return ""
	}
	return userID.String()
}

// Global audit logger instance
var globalAuditLogger = NewAuditLogger()

// GetAuditLogger returns the global audit logger instance
func GetAuditLogger() *AuditLogger {
	return globalAuditLogger
}
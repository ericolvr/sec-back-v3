package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	NotificationTypeThresholdReached  = "threshold_reached"
	NotificationTypeCanClose          = "can_close"
	NotificationTypeRiskDetected      = "risk_detected"
	NotificationTypeSnapshotCreated   = "snapshot_created"
	NotificationTypeActionPlanCreated = "action_plan_generated"
)

const (
	NotificationSeverityInfo     = "info"
	NotificationSeverityWarning  = "warning"
	NotificationSeverityCritical = "critical"
)

type Notification struct {
	ID              int64      `json:"id"`
	PartnerID       int64      `json:"partner_id"`
	CompanyID       *int64     `json:"company_id,omitempty"`
	DepartmentID    *int64     `json:"department_id,omitempty"`
	TemplateID *int64     `json:"template_id,omitempty"`
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Message         string     `json:"message"`
	Severity        string     `json:"severity"`
	Metadata        string     `json:"metadata"` // JSON string
	Read            bool       `json:"read"`
	ReadAt          *time.Time `json:"read_at,omitempty"`
	ReadBy          *int64     `json:"read_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, partnerID, id int64) (*Notification, error)
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*Notification, error)
	ListUnread(ctx context.Context, partnerID int64, limit, offset int64) ([]*Notification, error)
	ListByType(ctx context.Context, partnerID int64, notifType string, limit, offset int64) ([]*Notification, error)
	MarkAsRead(ctx context.Context, partnerID, id, userID int64) error
	MarkAllAsRead(ctx context.Context, partnerID, userID int64) error
	Delete(ctx context.Context, partnerID, id int64) error
	CountUnread(ctx context.Context, partnerID int64) (int64, error)
}

func (n *Notification) Validate() error {
	if n.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if n.Type == "" {
		return fmt.Errorf("type is required")
	}

	if n.Title == "" {
		return fmt.Errorf("title is required")
	}

	if n.Message == "" {
		return fmt.Errorf("message is required")
	}

	if n.Severity == "" {
		n.Severity = NotificationSeverityInfo
	}

	validSeverities := map[string]bool{
		NotificationSeverityInfo:     true,
		NotificationSeverityWarning:  true,
		NotificationSeverityCritical: true,
	}

	if !validSeverities[n.Severity] {
		return fmt.Errorf("invalid severity: %s", n.Severity)
	}

	return nil
}

package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	ActionPlanStatusPending    = "pending"
	ActionPlanStatusInProgress = "in_progress"
	ActionPlanStatusCompleted  = "completed"
	ActionPlanStatusCancelled  = "cancelled"
)

const (
	ActionPlanPriorityLow      = "low"
	ActionPlanPriorityMedium   = "medium"
	ActionPlanPriorityHigh     = "high"
	ActionPlanPriorityCritical = "critical"
)

type ActionPlan struct {
	ID                  int64      `json:"id"`
	PartnerID           int64      `json:"partner_id"`
	CompanyID           int64      `json:"company_id"`
	QuestionnaireID     int64      `json:"questionnaire_id"`
	DepartmentID        int64      `json:"department_id"`
	SnapshotID          *int64     `json:"snapshot_id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	RiskLevel           string     `json:"risk_level"`
	Priority            string     `json:"priority"`
	Category            string     `json:"category"`
	ResponsibleName     string     `json:"responsible_name"`
	ResponsibleID       *int64     `json:"responsible_id"`
	ResponsibleUsername *string    `json:"responsible_username,omitempty"`
	Status              string     `json:"status"`
	DueDate             *time.Time `json:"due_date"`
	CompletedAt         *time.Time `json:"completed_at"`
	EvidenceURLs        []string   `json:"evidence_urls"`
	Notes               string     `json:"notes"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ActionPlanRepository interface {
	Create(ctx context.Context, actionPlan *ActionPlan) error
	GetByID(ctx context.Context, partnerID, id int64) (*ActionPlan, error)
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*ActionPlan, error)
	ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*ActionPlan, error)
	ListBySnapshot(ctx context.Context, partnerID, snapshotID int64, limit, offset int64) ([]*ActionPlan, error)
	ListByQuestionnaire(ctx context.Context, partnerID, questionnaireID int64, limit, offset int64) ([]*ActionPlan, error)
	ListByStatus(ctx context.Context, partnerID int64, status string, limit, offset int64) ([]*ActionPlan, error)
	ListByResponsible(ctx context.Context, partnerID, responsibleID int64, limit, offset int64) ([]*ActionPlan, error)
	Update(ctx context.Context, actionPlan *ActionPlan) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (a *ActionPlan) Validate() error {
	if a.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if a.CompanyID <= 0 {
		return fmt.Errorf("company_id is required")
	}

	if a.QuestionnaireID <= 0 {
		return fmt.Errorf("questionnaire_id is required")
	}

	if a.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}

	if a.Title == "" {
		return fmt.Errorf("title is required")
	}

	if a.Description == "" {
		return fmt.Errorf("description is required")
	}

	if a.RiskLevel == "" {
		return fmt.Errorf("risk_level is required")
	}

	if a.Priority == "" {
		return fmt.Errorf("priority is required")
	}

	if a.ResponsibleName == "" {
		return fmt.Errorf("responsible_name is required")
	}

	if a.Status == "" {
		return fmt.Errorf("status is required")
	}

	return nil
}

func (a *ActionPlan) IsOverdue() bool {
	if a.DueDate == nil {
		return false
	}
	return time.Now().After(*a.DueDate) && a.Status != ActionPlanStatusCompleted && a.Status != ActionPlanStatusCancelled
}

func (a *ActionPlan) MarkAsCompleted() {
	now := time.Now()
	a.Status = ActionPlanStatusCompleted
	a.CompletedAt = &now
}

func (a *ActionPlan) MarkAsInProgress() {
	a.Status = ActionPlanStatusInProgress
}

func (a *ActionPlan) MarkAsCancelled() {
	a.Status = ActionPlanStatusCancelled
}

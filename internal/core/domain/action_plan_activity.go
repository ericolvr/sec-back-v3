package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	ActivityTypeNote     = "note"
	ActivityTypeAction   = "action"
	ActivityTypeEvidence = "evidence"
	ActivityTypeUpdate   = "update"
)

const (
	ActivityStatusPending    = "pending"
	ActivityStatusInProgress = "in_progress"
	ActivityStatusCompleted  = "completed"
	ActivityStatusCancelled  = "cancelled"
)

type ActionPlanActivity struct {
	ID            int64            `json:"id"`
	ActionPlanID  int64            `json:"action_plan_id"`
	Type          string           `json:"type"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Status        string           `json:"status"`
	DueDate       *time.Time       `json:"due_date"`
	CompletedAt   *time.Time       `json:"completed_at"`
	CreatedBy     *int64           `json:"created_by"`
	CreatedByName string           `json:"created_by_name"`
	Medias        []*ActivityMedia `json:"medias,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type ActionPlanActivityRepository interface {
	Create(ctx context.Context, activity *ActionPlanActivity) error
	GetByID(ctx context.Context, id int64) (*ActionPlanActivity, error)
	ListByActionPlan(ctx context.Context, actionPlanID int64) ([]*ActionPlanActivity, error)
	Update(ctx context.Context, activity *ActionPlanActivity) error
	Delete(ctx context.Context, id int64) error
}

func (a *ActionPlanActivity) Validate() error {
	if a.ActionPlanID <= 0 {
		return fmt.Errorf("action_plan_id is required")
	}

	if a.Type == "" {
		return fmt.Errorf("type is required")
	}

	validTypes := []string{ActivityTypeNote, ActivityTypeAction, ActivityTypeEvidence, ActivityTypeUpdate}
	isValidType := false
	for _, t := range validTypes {
		if a.Type == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid type: %s", a.Type)
	}

	if a.Title == "" {
		return fmt.Errorf("title is required")
	}

	if a.Status == "" {
		a.Status = ActivityStatusPending
	}

	validStatuses := []string{ActivityStatusPending, ActivityStatusInProgress, ActivityStatusCompleted, ActivityStatusCancelled}
	isValidStatus := false
	for _, s := range validStatuses {
		if a.Status == s {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return fmt.Errorf("invalid status: %s", a.Status)
	}

	return nil
}

func (a *ActionPlanActivity) MarkAsCompleted() {
	now := time.Now()
	a.Status = ActivityStatusCompleted
	a.CompletedAt = &now
}

package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	InvitationStatusPending = "pending"
	InvitationStatusSent    = "sent"
	InvitationStatusFailed  = "failed"
)

type Invitation struct {
	ID              int64      `json:"id"`
	PartnerID        int64      `json:"partner_id"`
	QuestionnaireID int64      `json:"questionnaire_id"`
	DepartmentID    int64      `json:"department_id"`
	ResponseID      int64      `json:"response_id"`
	EmployeeEmail   string     `json:"employee_email"`
	Status          string     `json:"status"`
	SentAt          *time.Time `json:"sent_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
	GetByID(ctx context.Context, tenantID, id int64) (*Invitation, error)
	List(ctx context.Context, tenantID, limit, offset int64) ([]*Invitation, error)
	ListByQuestionnaireAndDepartment(ctx context.Context, tenantID, questionnaireID, departmentID int64) ([]*Invitation, error)
	ListByStatus(ctx context.Context, tenantID int64, status string, limit, offset int64) ([]*Invitation, error)
	Update(ctx context.Context, invitation *Invitation) error
	Delete(ctx context.Context, tenantID, id int64) error
	CountByQuestionnaireAndDepartment(ctx context.Context, tenantID, questionnaireID, departmentID int64) (int64, error)
}

func (i *Invitation) ValidateInvitation() error {
	if i.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if i.QuestionnaireID <= 0 {
		return fmt.Errorf("questionnaire_id is required")
	}

	if i.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}

	if i.ResponseID <= 0 {
		return fmt.Errorf("response_id is required")
	}

	if i.EmployeeEmail == "" {
		return fmt.Errorf("employee_email is required")
	}

	if i.Status == "" {
		return fmt.Errorf("status is required")
	}

	return nil
}

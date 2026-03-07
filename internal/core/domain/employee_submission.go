package domain

import (
	"context"
	"fmt"
	"time"
)

// Status constants
const (
	SubmissionStatusPending    = "pending"
	SubmissionStatusInProgress = "in_progress"
	SubmissionStatusCompleted  = "completed"
)

// EmployeeSubmission representa a submissão de um questionário por um funcionário
// Anteriormente chamado de "Response"
type EmployeeSubmission struct {
	ID              int64      `json:"id"`
	PartnerID       int64      `json:"partner_id"`
	CompanyID       int64      `json:"company_id"`
	QuestionnaireID int64      `json:"questionnaire_id"`
	EmployeeID      int64      `json:"employee_id"`
	DepartmentID    int64      `json:"department_id"`
	InvitationToken string     `json:"invitation_token"`
	Status          string     `json:"status"` // pending, in_progress, completed
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type EmployeeSubmissionRepository interface {
	Create(ctx context.Context, submission *EmployeeSubmission) error
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*EmployeeSubmission, error)
	ListByCompany(ctx context.Context, partnerID, companyID int64, limit, offset int64) ([]*EmployeeSubmission, error)
	ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*EmployeeSubmission, error)
	GetByID(ctx context.Context, partnerID, id int64) (*EmployeeSubmission, error)
	GetByToken(ctx context.Context, token string) (*EmployeeSubmission, error)
	Update(ctx context.Context, submission *EmployeeSubmission) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (s *EmployeeSubmission) Validate() error {
	if s.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if s.CompanyID <= 0 {
		return fmt.Errorf("company_id is required")
	}
	if s.EmployeeID <= 0 {
		return fmt.Errorf("employee_id is required")
	}
	if s.QuestionnaireID <= 0 {
		return fmt.Errorf("questionnaire_id is required")
	}
	if s.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}
	if !IsValidSubmissionStatus(s.Status) {
		return fmt.Errorf("invalid status")
	}
	return nil
}

func ValidateEmployeeSubmission(e *EmployeeSubmission) error {
	return e.Validate()
}

func IsValidSubmissionStatus(status string) bool {
	validStatuses := []string{"pending", "in_progress", "completed"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func (s *EmployeeSubmission) MarkAsCompleted() {
	now := time.Now()
	s.Status = SubmissionStatusCompleted
	s.CompletedAt = &now
}

func (s *EmployeeSubmission) IsCompleted() bool {
	return s.Status == SubmissionStatusCompleted
}

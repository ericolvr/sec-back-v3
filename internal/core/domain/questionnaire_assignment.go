package domain

import (
	"context"
	"fmt"
	"time"
)

type QuestionnaireAssignment struct {
	ID                int64      `json:"id"`
	PartnerID          int64      `json:"partner_id"`
	QuestionnaireID   int64      `json:"questionnaire_id"`
	QuestionnaireName string     `json:"questionnaire_name,omitempty"`
	DepartmentID      int64      `json:"department_id"`
	DepartmentName    string     `json:"department_name,omitempty"`
	Active            bool       `json:"active"`
	StartedAt         time.Time  `json:"started_at"`
	ClosedAt          *time.Time `json:"closed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type QuestionnaireAssignmentRepository interface {
	Create(ctx context.Context, assignment *QuestionnaireAssignment) error
	GetByID(ctx context.Context, tenantID, id int64) (*QuestionnaireAssignment, error)
	GetByQuestionnaireAndDepartment(ctx context.Context, tenantID, questionnaireID, departmentID int64) (*QuestionnaireAssignment, error)
	List(ctx context.Context, tenantID, limit, offset int64) ([]*QuestionnaireAssignment, error)
	ListByQuestionnaire(ctx context.Context, tenantID, questionnaireID int64, limit, offset int64) ([]*QuestionnaireAssignment, error)
	ListByDepartment(ctx context.Context, tenantID, departmentID int64, limit, offset int64) ([]*QuestionnaireAssignment, error)
	ListActive(ctx context.Context, tenantID, departmentID int64) ([]*QuestionnaireAssignment, error)
	Update(ctx context.Context, assignment *QuestionnaireAssignment) error
	Delete(ctx context.Context, tenantID, id int64) error
}

func (qa *QuestionnaireAssignment) ValidateQuestionnaireAssignment() error {
	if qa.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if qa.QuestionnaireID <= 0 {
		return fmt.Errorf("questionnaire_id is required")
	}

	if qa.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}

	return nil
}

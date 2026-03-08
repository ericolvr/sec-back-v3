package domain

import (
	"context"
	"fmt"
	"time"
)

type DepartmentData struct {
	DepartmentName string `json:"department_name"`
	CompanyName    string `json:"company_name"`
}

type AssessmentAssignment struct {
	ID            int64            `json:"id"`
	PartnerID     int64            `json:"partner_id"`
	TemplateID    int64            `json:"template_id"`
	TemplateName  string           `json:"template_name,omitempty"`
	DepartmentIDs []int64          `json:"department_ids"`
	Data          []DepartmentData `json:"data,omitempty"`
	Count         int              `json:"count,omitempty"`
	Active        bool             `json:"active"`
	StartedAt     time.Time        `json:"-"` // Não serializar
	ClosedAt      *time.Time       `json:"-"` // Não serializar
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"-"` // Não serializar
}

type AssessmentAssignmentRepository interface {
	Create(ctx context.Context, assignment *AssessmentAssignment) error
	GetByID(ctx context.Context, partnerID, id int64) (*AssessmentAssignment, error)
	GetByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) (*AssessmentAssignment, error)
	List(ctx context.Context, partnerID, limit, offset int64) ([]*AssessmentAssignment, error)
	ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*AssessmentAssignment, error)
	ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*AssessmentAssignment, error)
	ListActive(ctx context.Context, partnerID, departmentID int64) ([]*AssessmentAssignment, error)
	Update(ctx context.Context, assignment *AssessmentAssignment) error
	CloseByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (aa *AssessmentAssignment) Validate() error {
	if aa.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if aa.TemplateID <= 0 {
		return fmt.Errorf("template_id is required")
	}

	if len(aa.DepartmentIDs) == 0 {
		return fmt.Errorf("at least one department_id is required")
	}

	return nil
}

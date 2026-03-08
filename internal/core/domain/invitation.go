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
	ID            int64      `json:"id"`
	PartnerID     int64      `json:"partner_id"`
	EmployeeID    int64      `json:"employee_id"` // Adicionado para compatibilidade com banco
	TemplateID    int64      `json:"template_id"`
	TemplateName  string     `json:"template_name,omitempty"`
	DepartmentID  int64      `json:"department_id"`
	ResponseID    int64      `json:"response_id"` // Mantido para compatibilidade
	Token         string     `json:"token"`       // Adicionado para compatibilidade com banco
	EmployeeEmail string     `json:"employee_email"`
	Status        string     `json:"status"`
	Sent          bool       `json:"sent"`
	SentAt        *time.Time `json:"sent_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
	GetByID(ctx context.Context, tenantID, id int64) (*Invitation, error)
	List(ctx context.Context, tenantID, limit, offset int64) ([]*Invitation, error)
	ListByTemplateAndDepartment(ctx context.Context, tenantID, templateID, departmentID int64) ([]*Invitation, error)
	ListByStatus(ctx context.Context, tenantID int64, status string, limit, offset int64) ([]*Invitation, error)
	Update(ctx context.Context, invitation *Invitation) error
	Delete(ctx context.Context, tenantID, id int64) error
	CountByTemplateAndDepartment(ctx context.Context, tenantID, templateID, departmentID int64) (int64, error)
}

func (i *Invitation) ValidateInvitation() error {
	if i.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if i.EmployeeID <= 0 && i.ResponseID <= 0 {
		return fmt.Errorf("employee_id or response_id is required")
	}

	if i.TemplateID <= 0 {
		return fmt.Errorf("template_id is required")
	}

	if i.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}

	if i.Token == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}

package domain

import (
	"context"
	"fmt"
	"time"
)

type Department struct {
	ID             int64     `json:"id"`
	PartnerID      int64     `json:"partner_id"` // Tenant raiz
	CompanyID      int64     `json:"company_id"` // Cliente do Partner
	Name           string    `json:"name"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	TotalEmployees int       `json:"total_employees,omitempty"`
}

type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*Department, error)
	ListByCompany(ctx context.Context, partnerID, companyID int64, limit, offset int64) ([]*Department, error)
	GetByID(ctx context.Context, partnerID, id int64) (*Department, error)
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (d *Department) Validate() error {
	if d.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if d.CompanyID <= 0 {
		return fmt.Errorf("company_id is required")
	}
	if d.Name == "" {
		return fmt.Errorf("department name is required")
	}
	return nil
}

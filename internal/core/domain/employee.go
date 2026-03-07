package domain

import (
	"context"
	"fmt"
	"time"
)

type Employee struct {
	ID           int64     `json:"id"`
	PartnerID    int64     `json:"partner_id"`
	CompanyID    int64     `json:"company_id"`
	DepartmentID int64     `json:"department_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Mobile       string    `json:"mobile"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Department *Department `json:"department,omitempty"`
}

type EmployeeRepository interface {
	Create(ctx context.Context, employee *Employee) error
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*Employee, error)
	ListByCompany(ctx context.Context, partnerID, companyID int64, limit, offset int64) ([]*Employee, error)
	ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*Employee, error)
	GetByID(ctx context.Context, partnerID, id int64) (*Employee, error)
	Update(ctx context.Context, employee *Employee) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (e *Employee) Validate() error {
	if e.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if e.CompanyID <= 0 {
		return fmt.Errorf("company_id is required")
	}
	if e.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}
	if e.Name == "" {
		return fmt.Errorf("employee name is required")
	}
	if e.Email == "" && e.Mobile == "" {
		return fmt.Errorf("employee must have email or mobile")
	}
	return nil
}

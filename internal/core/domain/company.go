package domain

import (
	"context"
	"fmt"
	"time"
)

// Company representa uma empresa CLIENTE do Partner
// Cada Company pertence a um Partner (multi-tenancy)
type Company struct {
	ID        int64     `json:"id"`
	PartnerID int64     `json:"partner_id"` // Referencia Partner (tenant)
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*Company, error)
	GetByID(ctx context.Context, partnerID, id int64) (*Company, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (c *Company) Validate() error {
	if c.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.CNPJ == "" {
		return fmt.Errorf("CNPJ is required")
	}
	return nil
}

func ValidateCompany(c *Company) error {
	return c.Validate()
}

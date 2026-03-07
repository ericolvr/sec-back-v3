package domain

import (
	"context"
	"fmt"
	"time"
)

// Partner representa a empresa parceira (Consultoria RH)
// Partner é o TENANT raiz do sistema multi-tenant
type Partner struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PartnerRepository interface {
	Create(ctx context.Context, partner *Partner) error
	List(ctx context.Context, limit, offset int64) ([]*Partner, error)
	GetByID(ctx context.Context, id int64) (*Partner, error)
	Update(ctx context.Context, partner *Partner) error
	Delete(ctx context.Context, id int64) error
}

func (p *Partner) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("partner name is required")
	}
	if p.CNPJ == "" {
		return fmt.Errorf("partner CNPJ is required")
	}
	return nil
}

func ValidatePartner(p *Partner) error {
	return p.Validate()
}

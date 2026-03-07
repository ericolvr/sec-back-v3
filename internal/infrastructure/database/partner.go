package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type PartnerRepository struct {
	db *sql.DB
}

func NewPartnerRepository(db *sql.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

func (r *PartnerRepository) Create(ctx context.Context, partner *domain.Partner) error {
	query := `
		INSERT INTO partners (name, cnpj, email, mobile, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, partner.Name, partner.CNPJ, partner.Email, partner.Mobile, partner.Active).
		Scan(&partner.ID, &partner.CreatedAt, &partner.UpdatedAt)
	return err
}

func (r *PartnerRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Partner, error) {
	query := `
		SELECT id, name, cnpj, email, mobile, active, created_at, updated_at
		FROM partners
		WHERE active = true
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partners []*domain.Partner
	for rows.Next() {
		var p domain.Partner
		err := rows.Scan(&p.ID, &p.Name, &p.CNPJ, &p.Email, &p.Mobile, &p.Active, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		partners = append(partners, &p)
	}
	return partners, nil
}

func (r *PartnerRepository) GetByID(ctx context.Context, id int64) (*domain.Partner, error) {
	query := `
		SELECT id, name, cnpj, email, mobile, active, created_at, updated_at
		FROM partners
		WHERE id = $1
	`
	var p domain.Partner
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.CNPJ, &p.Email, &p.Mobile, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PartnerRepository) Update(ctx context.Context, partner *domain.Partner) error {
	query := `
		UPDATE partners
		SET name = $1, cnpj = $2, email = $3, mobile = $4, active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, name, cnpj, email, mobile, active, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, partner.Name, partner.CNPJ, partner.Email, partner.Mobile, partner.Active, partner.ID).
		Scan(&partner.ID, &partner.Name, &partner.CNPJ, &partner.Email, &partner.Mobile, &partner.Active, &partner.CreatedAt, &partner.UpdatedAt)
	return err
}

func (r *PartnerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM partners WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

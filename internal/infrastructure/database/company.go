package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (partner_id, name, cnpj, email, mobile, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, company.PartnerID, company.Name, company.CNPJ, company.Email, company.Mobile, company.Active).
		Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
	return err
}

func (r *CompanyRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.Company, error) {
	query := `
		SELECT id, partner_id, name, cnpj, email, mobile, active, created_at, updated_at
		FROM companies
		WHERE partner_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*domain.Company
	for rows.Next() {
		var c domain.Company
		err := rows.Scan(&c.ID, &c.PartnerID, &c.Name, &c.CNPJ, &c.Email, &c.Mobile, &c.Active, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, &c)
	}
	return companies, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.Company, error) {
	query := `
		SELECT id, partner_id, name, cnpj, email, mobile, active, created_at, updated_at
		FROM companies
		WHERE partner_id = $1 AND id = $2
	`
	var c domain.Company
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(&c.ID, &c.PartnerID, &c.Name, &c.CNPJ, &c.Email, &c.Mobile, &c.Active, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `
		UPDATE companies
		SET name = $1, cnpj = $2, email = $3, mobile = $4, active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE partner_id = $6 AND id = $7
		RETURNING id, partner_id, name, cnpj, email, mobile, active, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, company.Name, company.CNPJ, company.Email, company.Mobile, company.Active, company.PartnerID, company.ID).
		Scan(&company.ID, &company.PartnerID, &company.Name, &company.CNPJ, &company.Email, &company.Mobile, &company.Active, &company.CreatedAt, &company.UpdatedAt)
	return err
}

func (r *CompanyRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM companies WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

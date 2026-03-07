package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(ctx context.Context, department *domain.Department) error {
	query := `
		INSERT INTO departments (partner_id, company_id, name, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, department.PartnerID, department.CompanyID, department.Name, department.Active).Scan(&department.ID, &department.CreatedAt, &department.UpdatedAt)
	return err
}

func (r *DepartmentRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.Department, error) {
	query := `
		SELECT 
			d.id, 
			d.partner_id, 
			d.company_id, 
			d.name, 
			d.active,
			d.created_at,
			d.updated_at,
			COUNT(e.id) as total_employees
		FROM departments d
		LEFT JOIN employees e ON d.id = e.department_id AND d.partner_id = e.partner_id
		WHERE d.partner_id = $1 AND d.active = true
		GROUP BY d.id, d.partner_id, d.company_id, d.name, d.active, d.created_at, d.updated_at
		ORDER BY d.name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*domain.Department

	for rows.Next() {
		var d domain.Department
		var totalEmployees int
		err := rows.Scan(&d.ID, &d.PartnerID, &d.CompanyID, &d.Name, &d.Active, &d.CreatedAt, &d.UpdatedAt, &totalEmployees)
		if err != nil {
			return nil, err
		}
		departments = append(departments, &d)
	}
	return departments, nil
}

func (r *DepartmentRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.Department, error) {
	query := `
		SELECT id, partner_id, company_id, name, active, created_at, updated_at 
		FROM departments 
		WHERE partner_id = $1 AND id = $2 AND active = true`

	var c domain.Department
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(&c.ID, &c.PartnerID, &c.CompanyID, &c.Name, &c.Active, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, department *domain.Department) error {
	query := `
		UPDATE departments
		SET name = $1, active = $2, updated_at = CURRENT_TIMESTAMP
		WHERE partner_id = $3 AND id = $4
		RETURNING id, partner_id, company_id, name, active, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, department.Name, department.Active, department.PartnerID, department.ID).
		Scan(&department.ID, &department.PartnerID, &department.CompanyID, &department.Name, &department.Active, &department.CreatedAt, &department.UpdatedAt)
	return err
}

func (r *DepartmentRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM departments WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *DepartmentRepository) ListByCompany(ctx context.Context, partnerID, companyID, limit, offset int64) ([]*domain.Department, error) {
	query := `SELECT id, partner_id, company_id, name, active, created_at, updated_at FROM departments WHERE partner_id = $1 AND company_id = $2 AND active = true ORDER BY name LIMIT $3 OFFSET $4`
	rows, err := r.db.QueryContext(ctx, query, partnerID, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*domain.Department
	for rows.Next() {
		var d domain.Department
		if err := rows.Scan(&d.ID, &d.PartnerID, &d.CompanyID, &d.Name, &d.Active, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		departments = append(departments, &d)
	}
	return departments, rows.Err()
}

func (r *DepartmentRepository) ToggleActive(ctx context.Context, partnerID, id int64, active bool) error {
	query := `
		UPDATE departments
		SET active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE partner_id = $2 AND id = $3
	`
	result, err := r.db.ExecContext(ctx, query, active, partnerID, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *DepartmentRepository) ListDeleted(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.Department, error) {
	query := `
		SELECT 
			d.id, 
			d.partner_id, 
			d.company_id, 
			d.name, 
			d.active,
			d.created_at,
			d.updated_at,
			COUNT(e.id) as total_employees
		FROM departments d
		LEFT JOIN employees e ON d.id = e.department_id AND d.partner_id = e.partner_id
		WHERE d.partner_id = $1 AND d.active = false
		GROUP BY d.id, d.partner_id, d.company_id, d.name, d.active, d.created_at, d.updated_at
		ORDER BY d.name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*domain.Department

	for rows.Next() {
		var d domain.Department
		var totalEmployees int
		err := rows.Scan(&d.ID, &d.PartnerID, &d.CompanyID, &d.Name, &d.Active, &d.CreatedAt, &d.UpdatedAt, &totalEmployees)
		if err != nil {
			return nil, err
		}
		departments = append(departments, &d)
	}
	return departments, nil
}

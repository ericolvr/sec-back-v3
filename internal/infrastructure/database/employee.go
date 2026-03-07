package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	query := `
		INSERT INTO employees (partner_id, department_id, name, email, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		employee.PartnerID,
		employee.DepartmentID,
		employee.Name,
		employee.Email,
		employee.Active,
		now,
	).Scan(&employee.ID)
	if err == nil {
		employee.CreatedAt = now
	}
	return err
}

func (r *EmployeeRepository) List(ctx context.Context, tenantID, limit, offset int64) ([]*domain.Employee, error) {
	query := `
		SELECT 
			e.id, 
			e.partner_id, 
			e.department_id, 
			e.name, 
			e.email, 
			e.active, 
			e.created_at,
			d.id,
			d.name,
			d.created_at
		FROM employees e
		INNER JOIN departments d ON e.department_id = d.id AND e.partner_id = d.partner_id
		WHERE e.partner_id = $1
		ORDER BY e.name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee

	for rows.Next() {
		var e domain.Employee
		var d domain.Department

		err := rows.Scan(
			&e.ID,
			&e.PartnerID,
			&e.DepartmentID,
			&e.Name,
			&e.Email,
			&e.Active,
			&e.CreatedAt,
			&d.ID,
			&d.Name,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		d.PartnerID = e.PartnerID
		e.Department = &d

		employees = append(employees, &e)
	}
	return employees, nil
}

func (r *EmployeeRepository) ListByDepartment(ctx context.Context, tenantID, departmentID, limit, offset int64) ([]*domain.Employee, error) {
	query := `
		SELECT 
			e.id, 
			e.partner_id, 
			e.department_id, 
			e.name, 
			e.email, 
			e.active, 
			e.created_at,
			d.id,
			d.name,
			d.created_at
		FROM employees e
		INNER JOIN departments d ON e.department_id = d.id AND e.partner_id = d.partner_id
		WHERE e.partner_id = $1 AND e.department_id = $2 AND e.active = true
		ORDER BY e.name ASC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, tenantID, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee

	for rows.Next() {
		var e domain.Employee
		var d domain.Department

		err := rows.Scan(
			&e.ID,
			&e.PartnerID,
			&e.DepartmentID,
			&e.Name,
			&e.Email,
			&e.Active,
			&e.CreatedAt,
			&d.ID,
			&d.Name,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		d.PartnerID = e.PartnerID
		e.Department = &d

		employees = append(employees, &e)
	}
	return employees, nil
}

func (r *EmployeeRepository) GetByID(ctx context.Context, tenantID, id int64) (*domain.Employee, error) {
	query := `
		SELECT 
			e.id, 
			e.partner_id, 
			e.department_id, 
			e.name, 
			e.email, 
			e.active, 
			e.created_at,
			d.id,
			d.name,
			d.created_at
		FROM employees e
		INNER JOIN departments d ON e.department_id = d.id AND e.partner_id = d.partner_id
		WHERE e.partner_id = $1 AND e.id = $2`

	var e domain.Employee
	var d domain.Department

	err := r.db.QueryRowContext(ctx, query, tenantID, id).Scan(
		&e.ID,
		&e.PartnerID,
		&e.DepartmentID,
		&e.Name,
		&e.Email,
		&e.Active,
		&e.CreatedAt,
		&d.ID,
		&d.Name,
		&d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	d.PartnerID = e.PartnerID
	e.Department = &d

	return &e, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	query := `
		UPDATE employees
		SET department_id = $1, name = $2, email = $3, active = $4
		WHERE partner_id = $5 AND id = $6
		RETURNING id, partner_id, department_id, name, email, active, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		employee.DepartmentID,
		employee.Name,
		employee.Email,
		employee.Active,
		employee.PartnerID,
		employee.ID,
	).Scan(
		&employee.ID,
		&employee.PartnerID,
		&employee.DepartmentID,
		&employee.Name,
		&employee.Email,
		&employee.Active,
		&employee.CreatedAt,
	)
	return err
}

func (r *EmployeeRepository) Delete(ctx context.Context, tenantID, id int64) error {
	query := `DELETE FROM employees WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, tenantID, id)
	return err
}

func (r *EmployeeRepository) ListByCompany(ctx context.Context, partnerID, companyID, limit, offset int64) ([]*domain.Employee, error) {
query := `SELECT e.id, e.partner_id, e.company_id, e.department_id, e.name, e.name, e.mobile, e.email, e.active, e.created_at, e.updated_at FROM employees e WHERE e.partner_id = $1 AND e.company_id = $2 ORDER BY e.name LIMIT $3 OFFSET $4`
rows, err := r.db.QueryContext(ctx, query, partnerID, companyID, limit, offset)
if err != nil {
return nil, err
}
defer rows.Close()

var employees []*domain.Employee
for rows.Next() {
var e domain.Employee
if err := rows.Scan(&e.ID, &e.PartnerID, &e.CompanyID, &e.DepartmentID, &e.Name, &e.Mobile, &e.Email, &e.Active, &e.CreatedAt, &e.UpdatedAt); err != nil {
return nil, err
}
employees = append(employees, &e)
}
return employees, rows.Err()
}

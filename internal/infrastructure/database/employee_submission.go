package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type EmployeeSubmissionRepository struct {
	db *sql.DB
}

func NewEmployeeSubmissionRepository(db *sql.DB) *EmployeeSubmissionRepository {
	return &EmployeeSubmissionRepository{db: db}
}

func (r *EmployeeSubmissionRepository) Create(ctx context.Context, submission *domain.EmployeeSubmission) error {
	query := `
		INSERT INTO employee_submissions (
			partner_id, company_id, template_id, employee_id, department_id,
			invitation_token, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		submission.PartnerID, submission.CompanyID, submission.TemplateID,
		submission.EmployeeID, submission.DepartmentID, submission.InvitationToken, submission.Status,
	).Scan(&submission.ID, &submission.CreatedAt, &submission.UpdatedAt)
	return err
}

func (r *EmployeeSubmissionRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.EmployeeSubmission, error) {
	query := `
		SELECT id, partner_id, company_id, template_id, employee_id, department_id,
			   invitation_token, status, completed_at, created_at, updated_at
		FROM employee_submissions
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*domain.EmployeeSubmission
	for rows.Next() {
		var s domain.EmployeeSubmission
		err := rows.Scan(&s.ID, &s.PartnerID, &s.CompanyID, &s.TemplateID, &s.EmployeeID, &s.DepartmentID,
			&s.InvitationToken, &s.Status, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, &s)
	}
	return submissions, nil
}

func (r *EmployeeSubmissionRepository) ListByCompany(ctx context.Context, partnerID, companyID int64, limit, offset int64) ([]*domain.EmployeeSubmission, error) {
	query := `
		SELECT id, partner_id, company_id, template_id, employee_id, department_id,
			   invitation_token, status, completed_at, created_at, updated_at
		FROM employee_submissions
		WHERE partner_id = $1 AND company_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*domain.EmployeeSubmission
	for rows.Next() {
		var s domain.EmployeeSubmission
		err := rows.Scan(&s.ID, &s.PartnerID, &s.CompanyID, &s.TemplateID, &s.EmployeeID, &s.DepartmentID,
			&s.InvitationToken, &s.Status, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, &s)
	}
	return submissions, nil
}

func (r *EmployeeSubmissionRepository) ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*domain.EmployeeSubmission, error) {
	query := `
		SELECT id, partner_id, company_id, template_id, employee_id, department_id,
			   invitation_token, status, completed_at, created_at, updated_at
		FROM employee_submissions
		WHERE partner_id = $1 AND department_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*domain.EmployeeSubmission
	for rows.Next() {
		var s domain.EmployeeSubmission
		err := rows.Scan(&s.ID, &s.PartnerID, &s.CompanyID, &s.TemplateID, &s.EmployeeID, &s.DepartmentID,
			&s.InvitationToken, &s.Status, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, &s)
	}
	return submissions, nil
}

func (r *EmployeeSubmissionRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.EmployeeSubmission, error) {
	query := `
		SELECT id, partner_id, company_id, template_id, employee_id, department_id,
			   invitation_token, status, completed_at, created_at, updated_at
		FROM employee_submissions
		WHERE partner_id = $1 AND id = $2
	`
	var s domain.EmployeeSubmission
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&s.ID, &s.PartnerID, &s.CompanyID, &s.TemplateID, &s.EmployeeID, &s.DepartmentID,
		&s.InvitationToken, &s.Status, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *EmployeeSubmissionRepository) GetByToken(ctx context.Context, token string) (*domain.EmployeeSubmission, error) {
	query := `
		SELECT id, partner_id, company_id, template_id, employee_id, department_id,
			   invitation_token, status, completed_at, created_at, updated_at
		FROM employee_submissions
		WHERE invitation_token = $1
	`
	var s domain.EmployeeSubmission
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&s.ID, &s.PartnerID, &s.CompanyID, &s.TemplateID, &s.EmployeeID, &s.DepartmentID,
		&s.InvitationToken, &s.Status, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *EmployeeSubmissionRepository) Update(ctx context.Context, submission *domain.EmployeeSubmission) error {
	query := `
		UPDATE employee_submissions
		SET status = $1, completed_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE partner_id = $3 AND id = $4
		RETURNING id, partner_id, company_id, template_id, employee_id, department_id,
				  invitation_token, status, completed_at, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, submission.Status, submission.CompletedAt, submission.PartnerID, submission.ID).Scan(
		&submission.ID, &submission.PartnerID, &submission.CompanyID, &submission.TemplateID, &submission.EmployeeID, &submission.DepartmentID,
		&submission.InvitationToken, &submission.Status, &submission.CompletedAt, &submission.CreatedAt, &submission.UpdatedAt)
	return err
}

func (r *EmployeeSubmissionRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM employee_submissions WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

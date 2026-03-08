package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentAssignmentRepository struct {
	db *sql.DB
}

func NewAssessmentAssignmentRepository(db *sql.DB) *AssessmentAssignmentRepository {
	return &AssessmentAssignmentRepository{db: db}
}

func (r *AssessmentAssignmentRepository) Create(ctx context.Context, assignment *domain.AssessmentAssignment) error {
	query := `
		INSERT INTO assessment_assignments (
			partner_id, template_id, department_id, active, started_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
		RETURNING id, started_at, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		assignment.PartnerID,
		assignment.TemplateID,
		assignment.DepartmentID,
		assignment.Active,
	).Scan(&assignment.ID, &assignment.StartedAt, &assignment.CreatedAt, &assignment.UpdatedAt)
}

func (r *AssessmentAssignmentRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentAssignment, error) {
	query := `
		SELECT id, partner_id, template_id, department_id, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND id = $2`

	var qa domain.AssessmentAssignment
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&qa.ID, &qa.PartnerID, &qa.TemplateID, &qa.DepartmentID, &qa.Active,
		&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &qa, nil
}

func (r *AssessmentAssignmentRepository) GetByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) (*domain.AssessmentAssignment, error) {
	query := `
		SELECT id, partner_id, template_id, department_id, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND template_id = $2 AND department_id = $3`

	var qa domain.AssessmentAssignment
	err := r.db.QueryRowContext(ctx, query, partnerID, templateID, departmentID).Scan(
		&qa.ID, &qa.PartnerID, &qa.TemplateID, &qa.DepartmentID, &qa.Active,
		&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &qa, nil
}

func (r *AssessmentAssignmentRepository) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	query := `
		SELECT 
			aa.id, aa.partner_id, aa.template_id, aa.department_id, aa.active, 
			aa.started_at, aa.closed_at, aa.created_at, aa.updated_at,
			at.name as template_name,
			d.company_id,
			c.name as company_name,
			d.name as department_name
		FROM assessment_assignments aa
		LEFT JOIN assessment_templates at ON aa.template_id = at.id
		LEFT JOIN departments d ON aa.department_id = d.id
		LEFT JOIN companies c ON d.company_id = c.id
		WHERE aa.partner_id = $1
		ORDER BY aa.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignmentsWithJoins(rows)
}

func (r *AssessmentAssignmentRepository) ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	query := `
		SELECT id, partner_id, template_id, department_id, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND template_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, templateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

func (r *AssessmentAssignmentRepository) ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	query := `
		SELECT id, partner_id, template_id, department_id, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND department_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

func (r *AssessmentAssignmentRepository) ListActive(ctx context.Context, partnerID, departmentID int64) ([]*domain.AssessmentAssignment, error) {
	query := `
		SELECT id, partner_id, template_id, department_id, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND department_id = $2 AND active = true
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

func (r *AssessmentAssignmentRepository) Update(ctx context.Context, assignment *domain.AssessmentAssignment) error {
	query := `
		UPDATE assessment_assignments SET
			active = $1, updated_at = NOW()
		WHERE partner_id = $2 AND id = $3`

	_, err := r.db.ExecContext(ctx, query, assignment.Active, assignment.PartnerID, assignment.ID)
	return err
}

func (r *AssessmentAssignmentRepository) CloseByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) error {
	query := `
		UPDATE assessment_assignments 
		SET active = false, closed_at = NOW(), updated_at = NOW()
		WHERE partner_id = $1 AND template_id = $2 AND department_id = $3 AND active = true`

	_, err := r.db.ExecContext(ctx, query, partnerID, templateID, departmentID)
	return err
}

func (r *AssessmentAssignmentRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM assessment_assignments WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *AssessmentAssignmentRepository) scanAssignments(rows *sql.Rows) ([]*domain.AssessmentAssignment, error) {
	var assignments []*domain.AssessmentAssignment

	for rows.Next() {
		var qa domain.AssessmentAssignment
		err := rows.Scan(
			&qa.ID, &qa.PartnerID, &qa.TemplateID, &qa.DepartmentID, &qa.Active,
			&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		assignments = append(assignments, &qa)
	}

	return assignments, rows.Err()
}

func (r *AssessmentAssignmentRepository) scanAssignmentsWithJoins(rows *sql.Rows) ([]*domain.AssessmentAssignment, error) {
	var assignments []*domain.AssessmentAssignment

	for rows.Next() {
		var qa domain.AssessmentAssignment
		var templateName, companyName, departmentName sql.NullString
		var companyID sql.NullInt64

		err := rows.Scan(
			&qa.ID, &qa.PartnerID, &qa.TemplateID, &qa.DepartmentID, &qa.Active,
			&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
			&templateName, &companyID, &companyName, &departmentName,
		)

		if err != nil {
			return nil, err
		}

		if templateName.Valid {
			qa.TemplateName = templateName.String
		}
		if companyID.Valid {
			qa.CompanyID = companyID.Int64
		}
		if companyName.Valid {
			qa.CompanyName = companyName.String
		}
		if departmentName.Valid {
			qa.DepartmentName = departmentName.String
		}

		assignments = append(assignments, &qa)
	}

	return assignments, rows.Err()
}

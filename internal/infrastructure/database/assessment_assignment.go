package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/lib/pq"
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
			partner_id, template_id, department_ids, active, started_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
		RETURNING id, started_at, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		assignment.PartnerID,
		assignment.TemplateID,
		pq.Array(assignment.DepartmentIDs),
		assignment.Active,
	).Scan(&assignment.ID, &assignment.StartedAt, &assignment.CreatedAt, &assignment.UpdatedAt)
}

func (r *AssessmentAssignmentRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentAssignment, error) {
	query := `
		SELECT 
			aa.id, aa.partner_id, aa.template_id, aa.department_ids, aa.active, 
			aa.started_at, aa.closed_at, aa.created_at, aa.updated_at,
			at.name as template_name,
			COALESCE(array_agg(d.name ORDER BY d.id) FILTER (WHERE d.id IS NOT NULL), '{}') as department_names,
			COALESCE(array_agg(c.name ORDER BY d.id) FILTER (WHERE c.id IS NOT NULL), '{}') as company_names
		FROM assessment_assignments aa
		LEFT JOIN assessment_templates at ON aa.template_id = at.id
		LEFT JOIN LATERAL unnest(aa.department_ids) dept_id ON true
		LEFT JOIN departments d ON d.id = dept_id AND d.partner_id = aa.partner_id
		LEFT JOIN companies c ON d.company_id = c.id
		WHERE aa.partner_id = $1 AND aa.id = $2
		GROUP BY aa.id, aa.partner_id, aa.template_id, aa.department_ids, aa.active, 
			aa.started_at, aa.closed_at, aa.created_at, aa.updated_at, at.name`

	var qa domain.AssessmentAssignment
	var templateName sql.NullString
	var departmentNames, companyNames []string

	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&qa.ID, &qa.PartnerID, &qa.TemplateID, pq.Array(&qa.DepartmentIDs), &qa.Active,
		&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
		&templateName,
		pq.Array(&departmentNames),
		pq.Array(&companyNames),
	)

	if err != nil {
		return nil, err
	}

	if templateName.Valid {
		qa.TemplateName = templateName.String
	}

	// Montar estrutura Data
	qa.Data = make([]domain.DepartmentData, 0, len(departmentNames))
	for i := 0; i < len(departmentNames) && i < len(companyNames); i++ {
		qa.Data = append(qa.Data, domain.DepartmentData{
			DepartmentName: departmentNames[i],
			CompanyName:    companyNames[i],
		})
	}
	qa.Count = len(qa.Data)

	return &qa, nil
}

func (r *AssessmentAssignmentRepository) GetByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) (*domain.AssessmentAssignment, error) {
	query := `
		SELECT 
			aa.id, aa.partner_id, aa.template_id, aa.department_ids, aa.active, 
			aa.started_at, aa.closed_at, aa.created_at, aa.updated_at,
			at.name as template_name
		FROM assessment_assignments aa
		LEFT JOIN assessment_templates at ON aa.template_id = at.id
		WHERE aa.partner_id = $1 AND aa.template_id = $2 AND $3 = ANY(aa.department_ids)`

	var qa domain.AssessmentAssignment
	err := r.db.QueryRowContext(ctx, query, partnerID, templateID, departmentID).Scan(
		&qa.ID, &qa.PartnerID, &qa.TemplateID, pq.Array(&qa.DepartmentIDs), &qa.Active,
		&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
		&qa.TemplateName,
	)

	if err != nil {
		return nil, err
	}

	return &qa, nil
}

func (r *AssessmentAssignmentRepository) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	query := `
		SELECT 
			aa.id, aa.partner_id, aa.template_id, aa.department_ids, aa.active, 
			aa.started_at, aa.closed_at, aa.created_at, aa.updated_at,
			at.name as template_name,
			COALESCE(array_agg(d.name ORDER BY d.id) FILTER (WHERE d.id IS NOT NULL), '{}') as department_names,
			COALESCE(array_agg(c.name ORDER BY d.id) FILTER (WHERE c.id IS NOT NULL), '{}') as company_names
		FROM assessment_assignments aa
		LEFT JOIN assessment_templates at ON aa.template_id = at.id
		LEFT JOIN LATERAL unnest(aa.department_ids) dept_id ON true
		LEFT JOIN departments d ON d.id = dept_id AND d.partner_id = aa.partner_id
		LEFT JOIN companies c ON d.company_id = c.id
		WHERE aa.partner_id = $1
		GROUP BY aa.id, aa.partner_id, aa.template_id, aa.department_ids, aa.active, 
			aa.started_at, aa.closed_at, aa.created_at, aa.updated_at, at.name
		ORDER BY aa.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignmentsWithDetails(rows)
}

func (r *AssessmentAssignmentRepository) ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	query := `
		SELECT id, partner_id, template_id, department_ids, active, started_at, closed_at, created_at, updated_at
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
		SELECT id, partner_id, template_id, department_ids, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND $2 = ANY(department_ids)
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
		SELECT id, partner_id, template_id, department_ids, active, started_at, closed_at, created_at, updated_at
		FROM assessment_assignments
		WHERE partner_id = $1 AND $2 = ANY(department_ids) AND active = true
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
			template_id = $1,
			department_ids = $2,
			active = $3,
			updated_at = NOW()
		WHERE partner_id = $4 AND id = $5`

	_, err := r.db.ExecContext(ctx, query,
		assignment.TemplateID,
		pq.Array(assignment.DepartmentIDs),
		assignment.Active,
		assignment.PartnerID,
		assignment.ID,
	)
	return err
}

func (r *AssessmentAssignmentRepository) CloseByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) error {
	query := `
		UPDATE assessment_assignments 
		SET active = false, closed_at = NOW(), updated_at = NOW()
		WHERE partner_id = $1 AND template_id = $2 AND $3 = ANY(department_ids) AND active = true`

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
			&qa.ID, &qa.PartnerID, &qa.TemplateID, pq.Array(&qa.DepartmentIDs), &qa.Active,
			&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		assignments = append(assignments, &qa)
	}

	return assignments, rows.Err()
}

func (r *AssessmentAssignmentRepository) scanAssignmentsWithDetails(rows *sql.Rows) ([]*domain.AssessmentAssignment, error) {
	var assignments []*domain.AssessmentAssignment

	for rows.Next() {
		var qa domain.AssessmentAssignment
		var templateName sql.NullString
		var departmentNames, companyNames []string

		err := rows.Scan(
			&qa.ID, &qa.PartnerID, &qa.TemplateID, pq.Array(&qa.DepartmentIDs), &qa.Active,
			&qa.StartedAt, &qa.ClosedAt, &qa.CreatedAt, &qa.UpdatedAt,
			&templateName,
			pq.Array(&departmentNames),
			pq.Array(&companyNames),
		)

		if err != nil {
			return nil, err
		}

		if templateName.Valid {
			qa.TemplateName = templateName.String
		}

		// Montar estrutura Data com department_name e company_name
		qa.Data = make([]domain.DepartmentData, 0, len(departmentNames))
		for i := 0; i < len(departmentNames) && i < len(companyNames); i++ {
			qa.Data = append(qa.Data, domain.DepartmentData{
				DepartmentName: departmentNames[i],
				CompanyName:    companyNames[i],
			})
		}
		qa.Count = len(qa.Data)

		assignments = append(assignments, &qa)
	}

	return assignments, rows.Err()
}

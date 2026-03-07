package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type QuestionnaireAssignmentRepository struct {
	db *sql.DB
}

func NewQuestionnaireAssignmentRepository(db *sql.DB) *QuestionnaireAssignmentRepository {
	return &QuestionnaireAssignmentRepository{db: db}
}

func (r *QuestionnaireAssignmentRepository) Create(ctx context.Context, assignment *domain.QuestionnaireAssignment) error {
	query := `
		INSERT INTO questionnaire_assignments (
			partner_id, questionnaire_id, department_id, active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		assignment.PartnerID,
		assignment.QuestionnaireID,
		assignment.DepartmentID,
		assignment.Active,
	).Scan(&assignment.ID, &assignment.CreatedAt, &assignment.UpdatedAt)
}

func (r *QuestionnaireAssignmentRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.QuestionnaireAssignment, error) {
	query := `
		SELECT id, partner_id, questionnaire_id, department_id, active, created_at, updated_at
		FROM questionnaire_assignments
		WHERE partner_id = $1 AND id = $2`

	var qa domain.QuestionnaireAssignment
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&qa.ID, &qa.PartnerID, &qa.QuestionnaireID, &qa.DepartmentID, &qa.Active, &qa.CreatedAt, &qa.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &qa, nil
}

func (r *QuestionnaireAssignmentRepository) GetByQuestionnaireAndDepartment(ctx context.Context, partnerID, questionnaireID, departmentID int64) (*domain.QuestionnaireAssignment, error) {
	query := `
		SELECT id, partner_id, questionnaire_id, department_id, active, created_at, updated_at
		FROM questionnaire_assignments
		WHERE partner_id = $1 AND questionnaire_id = $2 AND department_id = $3`

	var qa domain.QuestionnaireAssignment
	err := r.db.QueryRowContext(ctx, query, partnerID, questionnaireID, departmentID).Scan(
		&qa.ID, &qa.PartnerID, &qa.QuestionnaireID, &qa.DepartmentID, &qa.Active, &qa.CreatedAt, &qa.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &qa, nil
}

func (r *QuestionnaireAssignmentRepository) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.QuestionnaireAssignment, error) {
	query := `
		SELECT id, partner_id, questionnaire_id, department_id, active, created_at, updated_at
		FROM questionnaire_assignments
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

func (r *QuestionnaireAssignmentRepository) ListByQuestionnaire(ctx context.Context, partnerID, questionnaireID int64) ([]*domain.QuestionnaireAssignment, error) {
	query := `
		SELECT id, partner_id, questionnaire_id, department_id, active, created_at, updated_at
		FROM questionnaire_assignments
		WHERE partner_id = $1 AND questionnaire_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, partnerID, questionnaireID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

func (r *QuestionnaireAssignmentRepository) ListByDepartment(ctx context.Context, partnerID, departmentID int64) ([]*domain.QuestionnaireAssignment, error) {
	query := `
		SELECT id, partner_id, questionnaire_id, department_id, active, created_at, updated_at
		FROM questionnaire_assignments
		WHERE partner_id = $1 AND department_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

func (r *QuestionnaireAssignmentRepository) Update(ctx context.Context, assignment *domain.QuestionnaireAssignment) error {
	query := `
		UPDATE questionnaire_assignments SET
			active = $1, updated_at = NOW()
		WHERE partner_id = $2 AND id = $3`

	_, err := r.db.ExecContext(ctx, query, assignment.Active, assignment.PartnerID, assignment.ID)
	return err
}

func (r *QuestionnaireAssignmentRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM questionnaire_assignments WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *QuestionnaireAssignmentRepository) scanAssignments(rows *sql.Rows) ([]*domain.QuestionnaireAssignment, error) {
	var assignments []*domain.QuestionnaireAssignment

	for rows.Next() {
		var qa domain.QuestionnaireAssignment
		err := rows.Scan(
			&qa.ID, &qa.PartnerID, &qa.QuestionnaireID, &qa.DepartmentID, &qa.Active, &qa.CreatedAt, &qa.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		assignments = append(assignments, &qa)
	}

	return assignments, rows.Err()
}

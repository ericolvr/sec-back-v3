package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentTemplateRepository struct {
	db *sql.DB
}

func NewAssessmentTemplateRepository(db *sql.DB) *AssessmentTemplateRepository {
	return &AssessmentTemplateRepository{db: db}
}

func (r *AssessmentTemplateRepository) Create(ctx context.Context, template *domain.AssessmentTemplate) error {
	query := `
		INSERT INTO assessment_templates (partner_id, name, description, version, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, template.PartnerID, template.Name, template.Description, template.Version, template.Active).
		Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
	return err
}

func (r *AssessmentTemplateRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.AssessmentTemplate, error) {
	query := `
		SELECT id, partner_id, name, description, version, active, created_at, updated_at
		FROM assessment_templates
		WHERE partner_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.AssessmentTemplate
	for rows.Next() {
		var t domain.AssessmentTemplate
		err := rows.Scan(&t.ID, &t.PartnerID, &t.Name, &t.Description, &t.Version, &t.Active, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, nil
}

func (r *AssessmentTemplateRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentTemplate, error) {
	query := `
		SELECT id, partner_id, name, description, version, active, created_at, updated_at
		FROM assessment_templates
		WHERE partner_id = $1 AND id = $2
	`
	var t domain.AssessmentTemplate
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(&t.ID, &t.PartnerID, &t.Name, &t.Description, &t.Version, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *AssessmentTemplateRepository) Update(ctx context.Context, template *domain.AssessmentTemplate) error {
	query := `
		UPDATE assessment_templates
		SET name = $1, description = $2, active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE partner_id = $4 AND id = $5
		RETURNING id, partner_id, name, description, version, active, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, template.Name, template.Description, template.Active, template.PartnerID, template.ID).
		Scan(&template.ID, &template.PartnerID, &template.Name, &template.Description, &template.Version, &template.Active, &template.CreatedAt, &template.UpdatedAt)
	return err
}

func (r *AssessmentTemplateRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM assessment_templates WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *AssessmentTemplateRepository) IncrementVersion(ctx context.Context, partnerID, id int64) error {
	query := `
		UPDATE assessment_templates
		SET version = version + 1, updated_at = CURRENT_TIMESTAMP
		WHERE partner_id = $1 AND id = $2
	`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

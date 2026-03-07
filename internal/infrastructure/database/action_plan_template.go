package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActionPlanTemplateRepository struct {
	db *sql.DB
}

func NewActionPlanTemplateRepository(db *sql.DB) *ActionPlanTemplateRepository {
	return &ActionPlanTemplateRepository{db: db}
}

func (r *ActionPlanTemplateRepository) Create(ctx context.Context, template *domain.ActionPlanTemplate) error {
	query := `
		INSERT INTO action_plan_templates (
			partner_id, category, min_risk_level, title_template, description_template,
			priority, default_due_days, auto_create, active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		template.PartnerID,
		template.Category,
		template.MinRiskLevel,
		template.TitleTemplate,
		template.DescriptionTemplate,
		template.Priority,
		template.DefaultDueDays,
		template.AutoCreate,
		template.Active,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
}

func (r *ActionPlanTemplateRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.ActionPlanTemplate, error) {
	query := `
		SELECT id, partner_id, category, min_risk_level, title_template, description_template,
			   priority, default_due_days, auto_create, active, created_at, updated_at
		FROM action_plan_templates
		WHERE partner_id = $1 AND id = $2`

	var t domain.ActionPlanTemplate
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&t.ID, &t.PartnerID, &t.Category, &t.MinRiskLevel, &t.TitleTemplate, &t.DescriptionTemplate,
		&t.Priority, &t.DefaultDueDays, &t.AutoCreate, &t.Active, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *ActionPlanTemplateRepository) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.ActionPlanTemplate, error) {
	query := `
		SELECT id, partner_id, category, min_risk_level, title_template, description_template,
			   priority, default_due_days, auto_create, active, created_at, updated_at
		FROM action_plan_templates
		WHERE partner_id = $1
		ORDER BY category, min_risk_level
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *ActionPlanTemplateRepository) ListActive(ctx context.Context, partnerID int64) ([]*domain.ActionPlanTemplate, error) {
	query := `
		SELECT id, partner_id, category, min_risk_level, title_template, description_template,
			   priority, default_due_days, auto_create, active, created_at, updated_at
		FROM action_plan_templates
		WHERE partner_id = $1 AND active = true AND auto_create = true
		ORDER BY category, min_risk_level`

	rows, err := r.db.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *ActionPlanTemplateRepository) ListByCategory(ctx context.Context, partnerID int64, category string) ([]*domain.ActionPlanTemplate, error) {
	query := `
		SELECT id, partner_id, category, min_risk_level, title_template, description_template,
			   priority, default_due_days, auto_create, active, created_at, updated_at
		FROM action_plan_templates
		WHERE partner_id = $1 AND category = $2 AND active = true AND auto_create = true
		ORDER BY min_risk_level`

	rows, err := r.db.QueryContext(ctx, query, partnerID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *ActionPlanTemplateRepository) Update(ctx context.Context, template *domain.ActionPlanTemplate) error {
	query := `
		UPDATE action_plan_templates SET
			category = $1, min_risk_level = $2, title_template = $3, description_template = $4,
			priority = $5, default_due_days = $6, auto_create = $7, active = $8, updated_at = NOW()
		WHERE partner_id = $9 AND id = $10`

	_, err := r.db.ExecContext(
		ctx, query,
		template.Category, template.MinRiskLevel, template.TitleTemplate, template.DescriptionTemplate,
		template.Priority, template.DefaultDueDays, template.AutoCreate, template.Active,
		template.PartnerID, template.ID,
	)

	return err
}

func (r *ActionPlanTemplateRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM action_plan_templates WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *ActionPlanTemplateRepository) scanTemplates(rows *sql.Rows) ([]*domain.ActionPlanTemplate, error) {
	var templates []*domain.ActionPlanTemplate

	for rows.Next() {
		var t domain.ActionPlanTemplate
		err := rows.Scan(
			&t.ID, &t.PartnerID, &t.Category, &t.MinRiskLevel, &t.TitleTemplate, &t.DescriptionTemplate,
			&t.Priority, &t.DefaultDueDays, &t.AutoCreate, &t.Active, &t.CreatedAt, &t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		templates = append(templates, &t)
	}

	return templates, rows.Err()
}

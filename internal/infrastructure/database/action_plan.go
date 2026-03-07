package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActionPlanRepository struct {
	db *sql.DB
}

func NewActionPlanRepository(db *sql.DB) *ActionPlanRepository {
	return &ActionPlanRepository{db: db}
}

func (r *ActionPlanRepository) Create(ctx context.Context, actionPlan *domain.ActionPlan) error {
	query := `
		INSERT INTO action_plans (
			partner_id, company_id, questionnaire_id, department_id, snapshot_id,
			title, description, risk_level, priority, category,
			responsible_name, responsible_id, status, due_date,
			evidence_urls, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	evidenceJSON, _ := json.Marshal(actionPlan.EvidenceURLs)

	return r.db.QueryRowContext(
		ctx, query,
		actionPlan.PartnerID,
		actionPlan.CompanyID,
		actionPlan.QuestionnaireID,
		actionPlan.DepartmentID,
		actionPlan.SnapshotID,
		actionPlan.Title,
		actionPlan.Description,
		actionPlan.RiskLevel,
		actionPlan.Priority,
		actionPlan.Category,
		actionPlan.ResponsibleName,
		actionPlan.ResponsibleID,
		actionPlan.Status,
		actionPlan.DueDate,
		evidenceJSON,
		actionPlan.Notes,
	).Scan(&actionPlan.ID, &actionPlan.CreatedAt, &actionPlan.UpdatedAt)
}

func (r *ActionPlanRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.ActionPlan, error) {
	query := `
		SELECT id, partner_id, company_id, questionnaire_id, department_id, snapshot_id,
			   title, description, risk_level, priority, category,
			   responsible_name, responsible_id, status, due_date, completed_at,
			   evidence_urls, notes, created_at, updated_at
		FROM action_plans
		WHERE partner_id = $1 AND id = $2`

	var ap domain.ActionPlan
	var evidenceJSON []byte

	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&ap.ID, &ap.PartnerID, &ap.CompanyID, &ap.QuestionnaireID, &ap.DepartmentID, &ap.SnapshotID,
		&ap.Title, &ap.Description, &ap.RiskLevel, &ap.Priority, &ap.Category,
		&ap.ResponsibleName, &ap.ResponsibleID, &ap.Status, &ap.DueDate, &ap.CompletedAt,
		&evidenceJSON, &ap.Notes, &ap.CreatedAt, &ap.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if len(evidenceJSON) > 0 {
		json.Unmarshal(evidenceJSON, &ap.EvidenceURLs)
	}

	return &ap, nil
}

func (r *ActionPlanRepository) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.ActionPlan, error) {
	query := `
		SELECT id, partner_id, company_id, questionnaire_id, department_id, snapshot_id,
			   title, description, risk_level, priority, category,
			   responsible_name, responsible_id, status, due_date, completed_at,
			   evidence_urls, notes, created_at, updated_at
		FROM action_plans
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionPlans(rows)
}

func (r *ActionPlanRepository) ListByDepartment(ctx context.Context, partnerID, departmentID, limit, offset int64) ([]*domain.ActionPlan, error) {
	query := `
		SELECT id, partner_id, company_id, questionnaire_id, department_id, snapshot_id,
			   title, description, risk_level, priority, category,
			   responsible_name, responsible_id, status, due_date, completed_at,
			   evidence_urls, notes, created_at, updated_at
		FROM action_plans
		WHERE partner_id = $1 AND department_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionPlans(rows)
}

func (r *ActionPlanRepository) ListBySnapshot(ctx context.Context, partnerID, snapshotID, limit, offset int64) ([]*domain.ActionPlan, error) {
	query := `
		SELECT id, partner_id, company_id, questionnaire_id, department_id, snapshot_id,
			   title, description, risk_level, priority, category,
			   responsible_name, responsible_id, status, due_date, completed_at,
			   evidence_urls, notes, created_at, updated_at
		FROM action_plans
		WHERE partner_id = $1 AND snapshot_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, snapshotID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActionPlans(rows)
}

func (r *ActionPlanRepository) Update(ctx context.Context, actionPlan *domain.ActionPlan) error {
	query := `
		UPDATE action_plans SET
			title = $1, description = $2, risk_level = $3, priority = $4, category = $5,
			responsible_name = $6, responsible_id = $7, status = $8, due_date = $9,
			completed_at = $10, evidence_urls = $11, notes = $12, updated_at = NOW()
		WHERE partner_id = $13 AND id = $14`

	evidenceJSON, _ := json.Marshal(actionPlan.EvidenceURLs)

	_, err := r.db.ExecContext(
		ctx, query,
		actionPlan.Title, actionPlan.Description, actionPlan.RiskLevel, actionPlan.Priority, actionPlan.Category,
		actionPlan.ResponsibleName, actionPlan.ResponsibleID, actionPlan.Status, actionPlan.DueDate,
		actionPlan.CompletedAt, evidenceJSON, actionPlan.Notes,
		actionPlan.PartnerID, actionPlan.ID,
	)

	return err
}

func (r *ActionPlanRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM action_plans WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *ActionPlanRepository) scanActionPlans(rows *sql.Rows) ([]*domain.ActionPlan, error) {
	var actionPlans []*domain.ActionPlan

	for rows.Next() {
		var ap domain.ActionPlan
		var evidenceJSON []byte

		err := rows.Scan(
			&ap.ID, &ap.PartnerID, &ap.CompanyID, &ap.QuestionnaireID, &ap.DepartmentID, &ap.SnapshotID,
			&ap.Title, &ap.Description, &ap.RiskLevel, &ap.Priority, &ap.Category,
			&ap.ResponsibleName, &ap.ResponsibleID, &ap.Status, &ap.DueDate, &ap.CompletedAt,
			&evidenceJSON, &ap.Notes, &ap.CreatedAt, &ap.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if len(evidenceJSON) > 0 {
			json.Unmarshal(evidenceJSON, &ap.EvidenceURLs)
		}

		actionPlans = append(actionPlans, &ap)
	}

	return actionPlans, rows.Err()
}

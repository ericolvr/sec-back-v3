package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type RiskCategoryRepository struct {
	db *sql.DB
}

func NewRiskCategoryRepository(db *sql.DB) *RiskCategoryRepository {
	return &RiskCategoryRepository{db: db}
}

func (r *RiskCategoryRepository) Create(ctx context.Context, riskCategory *domain.RiskCategory) error {
	query := `
		INSERT INTO risk_categories (
			partner_id, snapshot_id, category, average_score, risk_level, question_count, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		riskCategory.PartnerID,
		riskCategory.SnapshotID,
		riskCategory.Category,
		riskCategory.AverageScore,
		riskCategory.RiskLevel,
		riskCategory.QuestionCount,
	).Scan(&riskCategory.ID, &riskCategory.CreatedAt)
}

func (r *RiskCategoryRepository) ListBySnapshot(ctx context.Context, partnerID, snapshotID int64) ([]*domain.RiskCategory, error) {
	query := `
		SELECT id, partner_id, snapshot_id, category, average_score, risk_level, question_count, created_at
		FROM risk_categories
		WHERE partner_id = $1 AND snapshot_id = $2
		ORDER BY category`

	rows, err := r.db.QueryContext(ctx, query, partnerID, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRiskCategories(rows)
}

func (r *RiskCategoryRepository) ListByCategory(ctx context.Context, partnerID int64, category string, limit, offset int64) ([]*domain.RiskCategory, error) {
	query := `
		SELECT id, partner_id, snapshot_id, category, average_score, risk_level, question_count, created_at
		FROM risk_categories
		WHERE partner_id = $1 AND category = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRiskCategories(rows)
}

func (r *RiskCategoryRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.RiskCategory, error) {
	query := `
		SELECT id, partner_id, snapshot_id, category, average_score, risk_level, question_count, created_at
		FROM risk_categories
		WHERE partner_id = $1 AND id = $2`

	var rc domain.RiskCategory
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&rc.ID, &rc.PartnerID, &rc.SnapshotID, &rc.Category, &rc.AverageScore, &rc.RiskLevel, &rc.QuestionCount, &rc.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &rc, nil
}

func (r *RiskCategoryRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM risk_categories WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *RiskCategoryRepository) scanRiskCategories(rows *sql.Rows) ([]*domain.RiskCategory, error) {
	var categories []*domain.RiskCategory

	for rows.Next() {
		var rc domain.RiskCategory
		err := rows.Scan(
			&rc.ID, &rc.PartnerID, &rc.SnapshotID, &rc.Category, &rc.AverageScore, &rc.RiskLevel, &rc.QuestionCount, &rc.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, &rc)
	}

	return categories, rows.Err()
}

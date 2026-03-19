package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type PartnerSettingsRepository struct {
	db *sql.DB
}

func NewPartnerSettingsRepository(db *sql.DB) *PartnerSettingsRepository {
	return &PartnerSettingsRepository{db: db}
}

func (r *PartnerSettingsRepository) Create(ctx context.Context, settings *domain.PartnerSettings) error {
	query := `
		INSERT INTO partner_settings (
			partner_id, min_response_rate_to_close, min_response_rate_reliable,
			notify_on_threshold, auto_generate_action_plans, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		settings.PartnerID,
		settings.MinResponseRateToClose,
		settings.MinResponseRateReliable,
		settings.NotifyOnThreshold,
		settings.AutoGenerateActionPlans,
	).Scan(&settings.CreatedAt, &settings.UpdatedAt)
}

func (r *PartnerSettingsRepository) GetByPartnerID(ctx context.Context, partnerID int64) (*domain.PartnerSettings, error) {
	query := `
		SELECT partner_id, min_response_rate_to_close, min_response_rate_reliable,
			   notify_on_threshold, auto_generate_action_plans, created_at, updated_at
		FROM partner_settings
		WHERE partner_id = $1`

	var settings domain.PartnerSettings
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&settings.PartnerID,
		&settings.MinResponseRateToClose,
		&settings.MinResponseRateReliable,
		&settings.NotifyOnThreshold,
		&settings.AutoGenerateActionPlans,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (r *PartnerSettingsRepository) Update(ctx context.Context, settings *domain.PartnerSettings) error {
	query := `
		UPDATE partner_settings SET
			min_response_rate_to_close = $1,
			min_response_rate_reliable = $2,
			notify_on_threshold = $3,
			auto_generate_action_plans = $4,
			updated_at = NOW()
		WHERE partner_id = $5`

	_, err := r.db.ExecContext(
		ctx, query,
		settings.MinResponseRateToClose,
		settings.MinResponseRateReliable,
		settings.NotifyOnThreshold,
		settings.AutoGenerateActionPlans,
		settings.PartnerID,
	)

	return err
}

func (r *PartnerSettingsRepository) Delete(ctx context.Context, partnerID int64) error {
	query := `DELETE FROM partner_settings WHERE partner_id = $1`
	_, err := r.db.ExecContext(ctx, query, partnerID)
	return err
}

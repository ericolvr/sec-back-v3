package domain

import (
	"context"
	"fmt"
	"time"
)

type PartnerSettings struct {
	PartnerID                 int64     `json:"partner_id"`
	MinResponseRateToClose    float64   `json:"min_response_rate_to_close"`
	MinResponseRateReliable   float64   `json:"min_response_rate_reliable"`
	NotifyOnThreshold         bool      `json:"notify_on_threshold"`
	AutoGenerateActionPlans   bool      `json:"auto_generate_action_plans"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type PartnerSettingsRepository interface {
	Create(ctx context.Context, settings *PartnerSettings) error
	GetByPartnerID(ctx context.Context, partnerID int64) (*PartnerSettings, error)
	Update(ctx context.Context, settings *PartnerSettings) error
	Delete(ctx context.Context, partnerID int64) error
}

func (ps *PartnerSettings) Validate() error {
	if ps.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if ps.MinResponseRateToClose < 0 || ps.MinResponseRateToClose > 100 {
		return fmt.Errorf("min_response_rate_to_close must be between 0 and 100")
	}

	if ps.MinResponseRateReliable < 0 || ps.MinResponseRateReliable > 100 {
		return fmt.Errorf("min_response_rate_reliable must be between 0 and 100")
	}

	if ps.MinResponseRateToClose < ps.MinResponseRateReliable {
		return fmt.Errorf("min_response_rate_to_close must be >= min_response_rate_reliable")
	}

	return nil
}

func DefaultPartnerSettings(partnerID int64) *PartnerSettings {
	return &PartnerSettings{
		PartnerID:               partnerID,
		MinResponseRateToClose:  70.0,
		MinResponseRateReliable: 30.0,
		NotifyOnThreshold:       true,
		AutoGenerateActionPlans: true,
	}
}

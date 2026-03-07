package domain

import (
	"context"
	"fmt"
	"time"
)

type RiskCategory struct {
	ID            int64     `json:"id"`
	PartnerID      int64     `json:"partner_id"`
	SnapshotID    int64     `json:"snapshot_id"`
	Category      string    `json:"category"`
	AverageScore  float64   `json:"average_score"`
	RiskLevel     string    `json:"risk_level"`
	QuestionCount int       `json:"question_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type RiskCategoryRepository interface {
	Create(ctx context.Context, riskCategory *RiskCategory) error
	ListBySnapshot(ctx context.Context, tenantID, snapshotID int64) ([]*RiskCategory, error)
	ListByCategory(ctx context.Context, tenantID int64, category string, limit, offset int64) ([]*RiskCategory, error)
	GetByID(ctx context.Context, tenantID, id int64) (*RiskCategory, error)
	Delete(ctx context.Context, tenantID, id int64) error
}

func (rc *RiskCategory) Validate() error {
	if rc.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if rc.SnapshotID <= 0 {
		return fmt.Errorf("snapshot_id is required")
	}

	if rc.Category == "" {
		return fmt.Errorf("category is required")
	}

	if rc.QuestionCount < 0 {
		return fmt.Errorf("question_count must be >= 0")
	}

	return nil
}

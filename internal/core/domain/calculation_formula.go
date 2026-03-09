package domain

import (
	"context"
	"fmt"
	"time"
)

type CalculationFormula struct {
	ID        int64  `json:"id"`
	PartnerID int64  `json:"partner_id"`
	Version   string `json:"version"`
	Active    bool   `json:"active"`

	// Thresholds de risco
	RiskLowMax    float64 `json:"risk_low_max"`
	RiskMediumMax float64 `json:"risk_medium_max"`

	// Thresholds de confiabilidade
	ReliabilityAcceptableMin float64 `json:"reliability_acceptable_min"`
	ReliabilityGoodMin       float64 `json:"reliability_good_min"`
	ReliabilityExcellentMin  float64 `json:"reliability_excellent_min"`

	// Metadados
	Description string     `json:"description,omitempty"`
	Changelog   string     `json:"changelog,omitempty"`
	CreatedBy   *int64     `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

type CalculationFormulaRepository interface {
	Create(ctx context.Context, formula *CalculationFormula) error
	GetActive(ctx context.Context, partnerID int64) (*CalculationFormula, error)
	GetByVersion(ctx context.Context, partnerID int64, version string) (*CalculationFormula, error)
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*CalculationFormula, error)
	Update(ctx context.Context, formula *CalculationFormula) error
	Activate(ctx context.Context, partnerID int64, version string) error
	Delete(ctx context.Context, partnerID int64, version string) error
}

func (cf *CalculationFormula) Validate() error {
	if cf.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if cf.Version == "" {
		return fmt.Errorf("version is required")
	}

	if cf.RiskLowMax <= 0 || cf.RiskLowMax >= 5 {
		return fmt.Errorf("risk_low_max must be between 0 and 5")
	}

	if cf.RiskMediumMax <= cf.RiskLowMax || cf.RiskMediumMax >= 5 {
		return fmt.Errorf("risk_medium_max must be between risk_low_max and 5")
	}

	if cf.ReliabilityAcceptableMin < 0 || cf.ReliabilityAcceptableMin > 100 {
		return fmt.Errorf("reliability_acceptable_min must be between 0 and 100")
	}

	if cf.ReliabilityGoodMin < cf.ReliabilityAcceptableMin || cf.ReliabilityGoodMin > 100 {
		return fmt.Errorf("reliability_good_min must be between reliability_acceptable_min and 100")
	}

	if cf.ReliabilityExcellentMin < cf.ReliabilityGoodMin || cf.ReliabilityExcellentMin > 100 {
		return fmt.Errorf("reliability_excellent_min must be between reliability_good_min and 100")
	}

	return nil
}

// DefaultCalculationFormula retorna fórmula padrão NR-1
func DefaultCalculationFormula(partnerID int64) *CalculationFormula {
	return &CalculationFormula{
		PartnerID:                partnerID,
		Version:                  "1.0",
		Active:                   true,
		RiskLowMax:               1.5,
		RiskMediumMax:            2.5,
		ReliabilityAcceptableMin: 30,
		ReliabilityGoodMin:       50,
		ReliabilityExcellentMin:  70,
		Description:              "Fórmula padrão baseada em estudos NR-1",
	}
}

// CalculateRiskLevel calcula nível de risco usando thresholds da fórmula
func (cf *CalculationFormula) CalculateRiskLevel(score float64) string {
	if score <= cf.RiskLowMax {
		return "low"
	}
	if score <= cf.RiskMediumMax {
		return "medium"
	}
	return "high"
}

// CalculateReliability calcula confiabilidade usando thresholds da fórmula
func (cf *CalculationFormula) CalculateReliability(responseRate float64) string {
	if responseRate >= cf.ReliabilityExcellentMin {
		return "excellent"
	}
	if responseRate >= cf.ReliabilityGoodMin {
		return "good"
	}
	if responseRate >= cf.ReliabilityAcceptableMin {
		return "acceptable"
	}
	return "insufficient"
}

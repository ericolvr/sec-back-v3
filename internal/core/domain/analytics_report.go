package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CalculationMetadata armazena metadados da fórmula usada no cálculo
type CalculationMetadata struct {
	FormulaVersion        string                    `json:"formula_version"`
	TemplateVersion       int                       `json:"template_version"`
	CalculatedAt          time.Time                 `json:"calculated_at"`
	RiskThresholds        map[string]ThresholdRange `json:"risk_thresholds"`
	ReliabilityThresholds map[string]float64        `json:"reliability_thresholds"`
}

type ThresholdRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// DepartmentAnalyticsWithMetadata extends DepartmentAnalytics with calculation metadata
type DepartmentAnalyticsWithMetadata struct {
	DepartmentAnalytics
	CalculationMetadata *CalculationMetadata `json:"calculation_metadata,omitempty"`
}

// AnalyticsReport representa um snapshot congelado de analytics
type AnalyticsReport struct {
	ID           int64     `json:"id"`
	PartnerID    int64     `json:"partner_id"`
	DepartmentID int64     `json:"department_id"`
	TemplateID   int64     `json:"template_id"`
	ReportData   []byte    `json:"report_data"` // JSON do DepartmentAnalyticsWithMetadata
	CreatedBy    *int64    `json:"created_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type AnalyticsReportRepository interface {
	Create(ctx context.Context, report *AnalyticsReport) error
	GetByID(ctx context.Context, partnerID, id int64) (*AnalyticsReport, error)
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*AnalyticsReport, error)
	ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*AnalyticsReport, error)
	ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*AnalyticsReport, error)
	ListByDepartmentAndTemplate(ctx context.Context, partnerID, departmentID, templateID int64, limit, offset int64) ([]*AnalyticsReport, error)
	Delete(ctx context.Context, partnerID, id int64) error
}

func (ar *AnalyticsReport) Validate() error {
	if ar.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if ar.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}

	if ar.TemplateID <= 0 {
		return fmt.Errorf("template_id is required")
	}

	if len(ar.ReportData) == 0 {
		return fmt.Errorf("report_data is required")
	}

	return nil
}

// GetAnalytics deserializa o JSON do snapshot
func (ar *AnalyticsReport) GetAnalytics() (*DepartmentAnalyticsWithMetadata, error) {
	var analytics DepartmentAnalyticsWithMetadata
	if err := json.Unmarshal(ar.ReportData, &analytics); err != nil {
		return nil, err
	}
	return &analytics, nil
}

// SetAnalytics serializa analytics para JSON
func (ar *AnalyticsReport) SetAnalytics(analytics *DepartmentAnalyticsWithMetadata) error {
	data, err := json.Marshal(analytics)
	if err != nil {
		return err
	}
	ar.ReportData = data
	return nil
}

// BuildCalculationMetadata cria metadados a partir da fórmula
func BuildCalculationMetadata(formula *CalculationFormula) *CalculationMetadata {
	return &CalculationMetadata{
		FormulaVersion: formula.Version,
		CalculatedAt:   time.Now(),
		RiskThresholds: map[string]ThresholdRange{
			"low":    {Min: 0, Max: formula.RiskLowMax},
			"medium": {Min: formula.RiskLowMax, Max: formula.RiskMediumMax},
			"high":   {Min: formula.RiskMediumMax, Max: 5.0},
		},
		ReliabilityThresholds: map[string]float64{
			"insufficient": 0,
			"acceptable":   formula.ReliabilityAcceptableMin,
			"good":         formula.ReliabilityGoodMin,
			"excellent":    formula.ReliabilityExcellentMin,
		},
	}
}

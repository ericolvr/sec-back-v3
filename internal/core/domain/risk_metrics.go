package domain

import (
	"context"
	"fmt"
	"time"
)

// RiskMetrics armazena métricas de risco PRÉ-CALCULADAS
// Substitui o cálculo em tempo real do Analytics
type RiskMetrics struct {
	ID              int64 `json:"id"`
	PartnerID       int64 `json:"partner_id"`
	CompanyID       int64 `json:"company_id"`
	DepartmentID    int64 `json:"department_id"`
	QuestionnaireID int64 `json:"questionnaire_id"`

	// Métricas pré-calculadas
	TotalEmployees       int     `json:"total_employees"`
	TotalSubmissions     int     `json:"total_submissions"`
	CompletedSubmissions int     `json:"completed_submissions"`
	ResponseRate         float64 `json:"response_rate"`
	AverageScore         float64 `json:"average_score"`
	RiskLevel            string  `json:"risk_level"`  // low, medium, high, critical
	Reliability          string  `json:"reliability"` // low, medium, high
	CanCalculateRisk     bool    `json:"can_calculate_risk"`

	// Scores por categoria (JSON)
	CategoryScores string `json:"category_scores"` // JSONB: {"Sobrecarga": 7.5, "Assédio": 3.2, ...}

	// Auditoria
	CalculatedAt time.Time `json:"calculated_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RiskMetricsRepository interface {
	Create(ctx context.Context, metrics *RiskMetrics) error
	Upsert(ctx context.Context, metrics *RiskMetrics) error // Insert or Update
	GetByDepartment(ctx context.Context, partnerID, departmentID, questionnaireID int64) (*RiskMetrics, error)
	GetByCompany(ctx context.Context, partnerID, companyID, questionnaireID int64) ([]*RiskMetrics, error)
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*RiskMetrics, error)
	Delete(ctx context.Context, partnerID, id int64) error
}

func (m *RiskMetrics) Validate() error {
	if m.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if m.CompanyID <= 0 {
		return fmt.Errorf("company_id is required")
	}
	if m.DepartmentID <= 0 {
		return fmt.Errorf("department_id is required")
	}
	return nil
}

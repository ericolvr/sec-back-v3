package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type RiskMetricsRepository struct {
	db *sql.DB
}

func NewRiskMetricsRepository(db *sql.DB) *RiskMetricsRepository {
	return &RiskMetricsRepository{db: db}
}

func (r *RiskMetricsRepository) Create(ctx context.Context, metrics *domain.RiskMetrics) error {
	query := `
		INSERT INTO risk_metrics (
			partner_id, company_id, department_id, questionnaire_id,
			total_employees, total_submissions, completed_submissions,
			response_rate, average_score, risk_level, reliability, can_calculate_risk,
			category_scores, calculated_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, calculated_at, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		metrics.PartnerID, metrics.CompanyID, metrics.DepartmentID, metrics.QuestionnaireID,
		metrics.TotalEmployees, metrics.TotalSubmissions, metrics.CompletedSubmissions,
		metrics.ResponseRate, metrics.AverageScore, metrics.RiskLevel, metrics.Reliability, metrics.CanCalculateRisk,
		metrics.CategoryScores,
	).Scan(&metrics.ID, &metrics.CalculatedAt, &metrics.CreatedAt, &metrics.UpdatedAt)
	return err
}

func (r *RiskMetricsRepository) Upsert(ctx context.Context, metrics *domain.RiskMetrics) error {
	query := `
		INSERT INTO risk_metrics (
			partner_id, company_id, department_id, questionnaire_id,
			total_employees, total_submissions, completed_submissions,
			response_rate, average_score, risk_level, reliability, can_calculate_risk,
			category_scores, calculated_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (partner_id, company_id, department_id, questionnaire_id)
		DO UPDATE SET
			total_employees = EXCLUDED.total_employees,
			total_submissions = EXCLUDED.total_submissions,
			completed_submissions = EXCLUDED.completed_submissions,
			response_rate = EXCLUDED.response_rate,
			average_score = EXCLUDED.average_score,
			risk_level = EXCLUDED.risk_level,
			reliability = EXCLUDED.reliability,
			can_calculate_risk = EXCLUDED.can_calculate_risk,
			category_scores = EXCLUDED.category_scores,
			calculated_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, calculated_at, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		metrics.PartnerID, metrics.CompanyID, metrics.DepartmentID, metrics.QuestionnaireID,
		metrics.TotalEmployees, metrics.TotalSubmissions, metrics.CompletedSubmissions,
		metrics.ResponseRate, metrics.AverageScore, metrics.RiskLevel, metrics.Reliability, metrics.CanCalculateRisk,
		metrics.CategoryScores,
	).Scan(&metrics.ID, &metrics.CalculatedAt, &metrics.CreatedAt, &metrics.UpdatedAt)
	return err
}

func (r *RiskMetricsRepository) GetByDepartment(ctx context.Context, partnerID, departmentID, questionnaireID int64) (*domain.RiskMetrics, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, questionnaire_id,
			   total_employees, total_submissions, completed_submissions,
			   response_rate, average_score, risk_level, reliability, can_calculate_risk,
			   category_scores, calculated_at, created_at, updated_at
		FROM risk_metrics
		WHERE partner_id = $1 AND department_id = $2 AND questionnaire_id = $3
	`
	var m domain.RiskMetrics
	err := r.db.QueryRowContext(ctx, query, partnerID, departmentID, questionnaireID).Scan(
		&m.ID, &m.PartnerID, &m.CompanyID, &m.DepartmentID, &m.QuestionnaireID,
		&m.TotalEmployees, &m.TotalSubmissions, &m.CompletedSubmissions,
		&m.ResponseRate, &m.AverageScore, &m.RiskLevel, &m.Reliability, &m.CanCalculateRisk,
		&m.CategoryScores, &m.CalculatedAt, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *RiskMetricsRepository) GetByCompany(ctx context.Context, partnerID, companyID, questionnaireID int64) ([]*domain.RiskMetrics, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, questionnaire_id,
			   total_employees, total_submissions, completed_submissions,
			   response_rate, average_score, risk_level, reliability, can_calculate_risk,
			   category_scores, calculated_at, created_at, updated_at
		FROM risk_metrics
		WHERE partner_id = $1 AND company_id = $2 AND questionnaire_id = $3
		ORDER BY department_id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, companyID, questionnaireID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*domain.RiskMetrics
	for rows.Next() {
		var m domain.RiskMetrics
		err := rows.Scan(
			&m.ID, &m.PartnerID, &m.CompanyID, &m.DepartmentID, &m.QuestionnaireID,
			&m.TotalEmployees, &m.TotalSubmissions, &m.CompletedSubmissions,
			&m.ResponseRate, &m.AverageScore, &m.RiskLevel, &m.Reliability, &m.CanCalculateRisk,
			&m.CategoryScores, &m.CalculatedAt, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &m)
	}
	return metrics, nil
}

func (r *RiskMetricsRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.RiskMetrics, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, questionnaire_id,
			   total_employees, total_submissions, completed_submissions,
			   response_rate, average_score, risk_level, reliability, can_calculate_risk,
			   category_scores, calculated_at, created_at, updated_at
		FROM risk_metrics
		WHERE partner_id = $1
		ORDER BY calculated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*domain.RiskMetrics
	for rows.Next() {
		var m domain.RiskMetrics
		err := rows.Scan(
			&m.ID, &m.PartnerID, &m.CompanyID, &m.DepartmentID, &m.QuestionnaireID,
			&m.TotalEmployees, &m.TotalSubmissions, &m.CompletedSubmissions,
			&m.ResponseRate, &m.AverageScore, &m.RiskLevel, &m.Reliability, &m.CanCalculateRisk,
			&m.CategoryScores, &m.CalculatedAt, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &m)
	}
	return metrics, nil
}

func (r *RiskMetricsRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM risk_metrics WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

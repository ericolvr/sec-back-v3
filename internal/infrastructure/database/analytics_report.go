package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AnalyticsReportRepository struct {
	db *sql.DB
}

func NewAnalyticsReportRepository(db *sql.DB) *AnalyticsReportRepository {
	return &AnalyticsReportRepository{db: db}
}

func (r *AnalyticsReportRepository) Create(ctx context.Context, report *domain.AnalyticsReport) error {
	query := `
		INSERT INTO analytics_reports (
			partner_id, department_id, questionnaire_id,
			report_data, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		report.PartnerID,
		report.DepartmentID,
		report.QuestionnaireID,
		report.ReportData,
		report.CreatedBy,
	).Scan(&report.ID, &report.CreatedAt)
}

func (r *AnalyticsReportRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.AnalyticsReport, error) {
	query := `
		SELECT id, partner_id, department_id, questionnaire_id,
			report_data, created_by, created_at
		FROM analytics_reports
		WHERE partner_id = $1 AND id = $2`

	var report domain.AnalyticsReport
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&report.ID,
		&report.PartnerID,
		&report.DepartmentID,
		&report.QuestionnaireID,
		&report.ReportData,
		&report.CreatedBy,
		&report.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *AnalyticsReportRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	query := `
		SELECT id, partner_id, department_id, questionnaire_id,
			report_data, created_by, created_at
		FROM analytics_reports
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *AnalyticsReportRepository) ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	query := `
		SELECT id, partner_id, department_id, questionnaire_id,
			report_data, created_by, created_at
		FROM analytics_reports
		WHERE partner_id = $1 AND department_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *AnalyticsReportRepository) ListByQuestionnaire(ctx context.Context, partnerID, questionnaireID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	query := `
		SELECT id, partner_id, department_id, questionnaire_id,
			report_data, created_by, created_at
		FROM analytics_reports
		WHERE partner_id = $1 AND questionnaire_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, questionnaireID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *AnalyticsReportRepository) ListByDepartmentAndQuestionnaire(ctx context.Context, partnerID, departmentID, questionnaireID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	query := `
		SELECT id, partner_id, department_id, questionnaire_id,
			report_data, created_by, created_at
		FROM analytics_reports
		WHERE partner_id = $1 AND department_id = $2 AND questionnaire_id = $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, questionnaireID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *AnalyticsReportRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM analytics_reports WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *AnalyticsReportRepository) scanReports(rows *sql.Rows) ([]*domain.AnalyticsReport, error) {
	var reports []*domain.AnalyticsReport

	for rows.Next() {
		var report domain.AnalyticsReport
		err := rows.Scan(
			&report.ID,
			&report.PartnerID,
			&report.DepartmentID,
			&report.QuestionnaireID,
			&report.ReportData,
			&report.CreatedBy,
			&report.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		reports = append(reports, &report)
	}

	return reports, rows.Err()
}

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
			partner_id, department_id, template_id,
			report_data, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		report.PartnerID,
		report.DepartmentID,
		report.TemplateID,
		report.ReportData,
		report.CreatedBy,
	).Scan(&report.ID, &report.CreatedAt)
}

func (r *AnalyticsReportRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.AnalyticsReport, error) {
	query := `
		SELECT id, partner_id, department_id, template_id,
			report_data, created_by, created_at
		FROM analytics_reports
		WHERE partner_id = $1 AND id = $2`

	var report domain.AnalyticsReport
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&report.ID,
		&report.PartnerID,
		&report.DepartmentID,
		&report.TemplateID,
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
		SELECT 
			ar.id, ar.partner_id, ar.department_id, ar.template_id,
			ar.report_data, ar.created_by, ar.created_at,
			d.name as department_name,
			at.name as template_name
		FROM analytics_reports ar
		LEFT JOIN departments d ON ar.department_id = d.id
		LEFT JOIN assessment_templates at ON ar.template_id = at.id
		WHERE ar.partner_id = $1
		ORDER BY ar.created_at DESC
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
		SELECT 
			ar.id, ar.partner_id, ar.department_id, ar.template_id,
			ar.report_data, ar.created_by, ar.created_at,
			d.name as department_name,
			at.name as template_name
		FROM analytics_reports ar
		LEFT JOIN departments d ON ar.department_id = d.id
		LEFT JOIN assessment_templates at ON ar.template_id = at.id
		WHERE ar.partner_id = $1 AND ar.department_id = $2
		ORDER BY ar.created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *AnalyticsReportRepository) ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	query := `
		SELECT 
			ar.id, ar.partner_id, ar.department_id, ar.template_id,
			ar.report_data, ar.created_by, ar.created_at,
			d.name as department_name,
			at.name as template_name
		FROM analytics_reports ar
		LEFT JOIN departments d ON ar.department_id = d.id
		LEFT JOIN assessment_templates at ON ar.template_id = at.id
		WHERE ar.partner_id = $1 AND ar.template_id = $2
		ORDER BY ar.created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, templateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *AnalyticsReportRepository) ListByDepartmentAndTemplate(ctx context.Context, partnerID, departmentID, templateID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	query := `
		SELECT 
			ar.id, ar.partner_id, ar.department_id, ar.template_id,
			ar.report_data, ar.created_by, ar.created_at,
			d.name as department_name,
			at.name as template_name
		FROM analytics_reports ar
		LEFT JOIN departments d ON ar.department_id = d.id
		LEFT JOIN assessment_templates at ON ar.template_id = at.id
		WHERE ar.partner_id = $1 AND ar.department_id = $2 AND ar.template_id = $3
		ORDER BY ar.created_at DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.db.QueryContext(ctx, query, partnerID, departmentID, templateID, limit, offset)
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
		var departmentName, templateName sql.NullString

		err := rows.Scan(
			&report.ID,
			&report.PartnerID,
			&report.DepartmentID,
			&report.TemplateID,
			&report.ReportData,
			&report.CreatedBy,
			&report.CreatedAt,
			&departmentName,
			&templateName,
		)

		if err != nil {
			return nil, err
		}

		if departmentName.Valid {
			report.DepartmentName = departmentName.String
		}
		if templateName.Valid {
			report.TemplateName = templateName.String
		}

		reports = append(reports, &report)
	}

	return reports, rows.Err()
}

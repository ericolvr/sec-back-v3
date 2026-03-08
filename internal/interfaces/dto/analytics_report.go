package dto

import "time"

type CreateSnapshotRequest struct {
	TemplateID int64 `json:"template_id" binding:"required"`
}

type AnalyticsReportResponse struct {
	ID              int64       `json:"id"`
	PartnerID    int64  `json:"partner_id"`
	DepartmentID    int64       `json:"department_id"`
	TemplateID int64       `json:"template_id"`
	ReportDate      string      `json:"report_date"`
	Data            interface{} `json:"data"`
	CreatedBy       *int64      `json:"created_by"`
	CreatedAt       string      `json:"created_at"`
}

func FormatAnalyticsReportResponse(report interface{}) *AnalyticsReportResponse {
	// This will be implemented in the handler
	return nil
}

func FormatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

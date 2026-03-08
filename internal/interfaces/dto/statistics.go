package dto

import "time"

// TemplateRiskOverviewResponse represents a questionnaire with aggregated risk statistics
type TemplateRiskOverviewResponse struct {
	// Basic Information
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Application Statistics
	TotalDepartments   int     `json:"total_departments"`
	TotalEmployees     int64   `json:"total_employees"`
	TotalResponses     int64   `json:"total_responses"`
	CompletedResponses int64   `json:"completed_responses"`
	ResponseRate       float64 `json:"response_rate"`

	// Risk Indicators (MOST IMPORTANT)
	DepartmentsWithHighRisk   int     `json:"departments_with_high_risk"`
	DepartmentsWithMediumRisk int     `json:"departments_with_medium_risk"`
	DepartmentsWithLowRisk    int     `json:"departments_with_low_risk"`
	DepartmentsPending        int     `json:"departments_pending"`
	OverallRiskLevel          string  `json:"overall_risk_level"` // "high", "medium", "low", "unknown"
	AverageRiskScore          float64 `json:"average_risk_score"`
	HighestRiskDepartment     string  `json:"highest_risk_department,omitempty"`

	// Alerts and Status
	RequiresAttention bool       `json:"requires_attention"`
	CanGenerateReport bool       `json:"can_generate_report"`
	LastUpdated       *time.Time `json:"last_updated,omitempty"`
}

// TemplateInProgressResponse represents a questionnaire in progress with completion metrics
type TemplateInProgressResponse struct {
	// Basic Information
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedAtFormatted string    `json:"created_at_formatted"`

	// Completion Metrics
	TotalDepartments      int   `json:"total_departments"`
	DepartmentsCompleted  int   `json:"departments_completed"`
	DepartmentsInProgress int   `json:"departments_in_progress"`
	DepartmentsNotStarted int   `json:"departments_not_started"`
	TotalEmployees        int64 `json:"total_employees"`
	TotalInvitations      int64 `json:"total_invitations"`
	CompletedResponses    int64 `json:"completed_responses"`

	// Partial Risk Analysis (from completed responses)
	DepartmentsWithHighRisk   int     `json:"departments_with_high_risk"`
	DepartmentsWithMediumRisk int     `json:"departments_with_medium_risk"`
	DepartmentsWithLowRisk    int     `json:"departments_with_low_risk"`
	OverallRiskLevel          string  `json:"overall_risk_level"`
	AverageScore              float64 `json:"average_score"`

	// Department Details
	Departments []DepartmentProgressDetail `json:"departments"`

	LastUpdated time.Time `json:"last_updated"`
}

// DepartmentProgressDetail represents progress details for a department
type DepartmentProgressDetail struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	TotalEmployees     int64   `json:"total_employees"`
	CompletedResponses int64   `json:"completed_responses"`
	PendingResponses   int64   `json:"pending_responses"`
	ResponseRate       float64 `json:"response_rate"`
	CanCalculateRisk   bool    `json:"can_calculate_risk"`
	Reliability        string  `json:"reliability"`
	AverageScore       float64 `json:"average_score"`
	RiskLevel          string  `json:"risk_level"`
	Status             string  `json:"status"` // "completed", "in_progress", "not_started"
	CanClose           bool    `json:"can_close"`
	CanCloseReason     string  `json:"can_close_reason,omitempty"`
}

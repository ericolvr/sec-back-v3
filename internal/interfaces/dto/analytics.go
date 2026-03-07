package dto

type DepartmentAnalyticsResponse struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	QuestionnaireID    int64   `json:"questionnaire_id"`
	QuestionnaireName  string  `json:"questionnaire_name"`
	TotalEmployees     int64   `json:"total_employees"`
	TotalResponses     int64   `json:"total_responses"`
	CompletedResponses int64   `json:"completed_responses"`
	ResponseRate       float64 `json:"response_rate"`
	Reliability        string  `json:"reliability"`
	CanCalculateRisk   bool    `json:"can_calculate_risk"`
	AverageScore       float64 `json:"average_score"`
	RiskLevel          string  `json:"risk_level"`
	Warning            string  `json:"warning,omitempty"`
}

type DepartmentSummaryResponse struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	TotalEmployees     int64   `json:"total_employees"`
	TotalResponses     int64   `json:"total_responses"`
	CompletedResponses int64   `json:"completed_responses"`
	ResponseRate       float64 `json:"response_rate"`
	Reliability        string  `json:"reliability"`
	CanCalculateRisk   bool    `json:"can_calculate_risk"`
	AverageScore       float64 `json:"average_score"`
	RiskLevel          string  `json:"risk_level"`
	Warning            string  `json:"warning,omitempty"`
}

type QuestionnaireAnalyticsResponse struct {
	QuestionnaireID   int64                        `json:"questionnaire_id"`
	QuestionnaireName string                       `json:"questionnaire_name"`
	Departments       []*DepartmentSummaryResponse `json:"departments"`
}

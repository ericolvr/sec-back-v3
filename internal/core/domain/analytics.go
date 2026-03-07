package domain

// DepartmentAnalytics representa o relatório detalhado de um departamento
type DepartmentAnalytics struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	QuestionnaireID    int64   `json:"questionnaire_id"`
	QuestionnaireName  string  `json:"questionnaire_name"`
	TotalEmployees     int64   `json:"total_employees"`
	TotalSubmissions   int64   `json:"total_submissions"`
	CompletedSubmissions int64 `json:"completed_submissions"`
	ResponseRate       float64 `json:"response_rate"`
	Reliability        string  `json:"reliability"`
	CanCalculateRisk   bool    `json:"can_calculate_risk"`
	AverageScore       float64 `json:"average_score"`
	RiskLevel          string  `json:"risk_level"`
	CanClose           bool    `json:"can_close"`
	CanCloseReason     string  `json:"can_close_reason,omitempty"`
	Warning            string  `json:"warning,omitempty"`
	
	// Riscos por categoria
	RiskCategories []*RiskCategorySummary `json:"risk_categories,omitempty"`
	
	// Action plans gerados
	ActionPlansCount int `json:"action_plans_count"`
}

// RiskCategorySummary representa o resumo de risco por categoria
type RiskCategorySummary struct {
	Category     string  `json:"category"`
	AverageScore float64 `json:"average_score"`
	RiskLevel    string  `json:"risk_level"`
	QuestionCount int    `json:"question_count"`
}

// CompanyAnalytics representa o relatório consolidado de uma empresa
type CompanyAnalytics struct {
	CompanyID          int64                  `json:"company_id"`
	CompanyName        string                 `json:"company_name"`
	QuestionnaireID    int64                  `json:"questionnaire_id"`
	QuestionnaireName  string                 `json:"questionnaire_name"`
	TotalDepartments   int                    `json:"total_departments"`
	OverallResponseRate float64               `json:"overall_response_rate"`
	OverallRiskLevel   string                 `json:"overall_risk_level"`
	DepartmentsAtRisk  int                    `json:"departments_at_risk"`
	Departments        []*DepartmentSummary   `json:"departments"`
}

// DepartmentSummary representa o resumo de um departamento (para visão de empresa)
type DepartmentSummary struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	TotalEmployees     int64   `json:"total_employees"`
	CompletedSubmissions int64 `json:"completed_submissions"`
	ResponseRate       float64 `json:"response_rate"`
	Reliability        string  `json:"reliability"`
	CanCalculateRisk   bool    `json:"can_calculate_risk"`
	AverageScore       float64 `json:"average_score"`
	RiskLevel          string  `json:"risk_level"`
	CanClose           bool    `json:"can_close"`
	CanCloseReason     string  `json:"can_close_reason,omitempty"`
	Warning            string  `json:"warning,omitempty"`
}

// PartnerAnalytics representa o relatório consolidado do Partner (visão geral)
type PartnerAnalytics struct {
	PartnerID              int64              `json:"partner_id"`
	PartnerName            string             `json:"partner_name"`
	TotalCompanies         int                `json:"total_companies"`
	TotalActiveQuestionnaires int             `json:"total_active_questionnaires"`
	OverallResponseRate    float64            `json:"overall_response_rate"`
	CompaniesAtRisk        int                `json:"companies_at_risk"`
	Companies              []*CompanySummary  `json:"companies"`
}

// CompanySummary representa o resumo de uma empresa (para visão de partner)
type CompanySummary struct {
	CompanyID          int64   `json:"company_id"`
	CompanyName        string  `json:"company_name"`
	TotalDepartments   int     `json:"total_departments"`
	ResponseRate       float64 `json:"response_rate"`
	RiskLevel          string  `json:"risk_level"`
	DepartmentsAtRisk  int     `json:"departments_at_risk"`
	ActiveQuestionnaires int   `json:"active_questionnaires"`
}

// CalculateReliability determina confiabilidade baseado na taxa de resposta
func CalculateReliability(responseRate float64) string {
	switch {
	case responseRate >= 70:
		return "excellent"
	case responseRate >= 50:
		return "good"
	case responseRate >= 30:
		return "acceptable"
	default:
		return "insufficient"
	}
}

// CalculateRiskLevel determina nível de risco baseado no score médio
func CalculateRiskLevel(averageScore float64) string {
	switch {
	case averageScore <= 1.5:
		return "low"
	case averageScore <= 2.5:
		return "medium"
	default:
		return "high"
	}
}

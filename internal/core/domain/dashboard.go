package domain

// CompanyDashboard representa o dashboard para gestor de uma empresa
type CompanyDashboard struct {
	CompanyID                 int64                           `json:"company_id"`
	CompanyName               string                          `json:"company_name"`
	InProgressQuestionnaires  []*TemplateInProgress      `json:"in_progress_questionnaires"`
	UnreadNotifications       int64                           `json:"unread_notifications"`
	NotificationsPreview      []*Notification                 `json:"notifications_preview"`
	PendingActionPlans        int                             `json:"pending_action_plans"`
	OverdueActionPlans        int                             `json:"overdue_action_plans"`
	PendingInvitations        int                             `json:"pending_invitations"`
	TotalEmployees            int                             `json:"total_employees"`
	TotalDepartments          int                             `json:"total_departments"`
	OverallRiskLevel          string                          `json:"overall_risk_level"`
	Alerts                    []string                        `json:"alerts"`
}

// TemplateInProgress representa um template em andamento
type TemplateInProgress struct {
	TemplateID     int64               `json:"template_id"`
	TemplateName   string              `json:"template_name"`
	TotalDepartments    int                 `json:"total_departments"`
	ResponseRate        float64             `json:"response_rate"`
	OverallRiskLevel    string              `json:"overall_risk_level"`
	DepartmentsAtRisk   int                 `json:"departments_at_risk"`
	Departments         []*DepartmentStatus `json:"departments"`
}

// DepartmentStatus representa o status de um departamento em um template
type DepartmentStatus struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	TotalEmployees     int64   `json:"total_employees"`
	ResponseRate       float64 `json:"response_rate"`
	RiskLevel          string  `json:"risk_level"`
	CanClose           bool    `json:"can_close"`
	CanCloseReason     string  `json:"can_close_reason,omitempty"`
}

// PartnerDashboard representa o dashboard para consultoria (Partner)
type PartnerDashboard struct {
	PartnerID                 int64            `json:"partner_id"`
	PartnerName               string           `json:"partner_name"`
	CompaniesSummary          []*CompanyStatus `json:"companies_summary"`
	TotalCompanies            int              `json:"total_companies"`
	TotalActiveTemplates int              `json:"total_active_questionnaires"`
	CompaniesAtRisk           int              `json:"companies_at_risk"`
	OverallResponseRate       float64          `json:"overall_response_rate"`
	Alerts                    []string         `json:"alerts"`
}

// CompanyStatus representa o status de uma empresa (visão do Partner)
type CompanyStatus struct {
	CompanyID            int64   `json:"company_id"`
	CompanyName          string  `json:"company_name"`
	ActiveTemplates int     `json:"active_questionnaires"`
	ResponseRate         float64 `json:"response_rate"`
	RiskLevel            string  `json:"risk_level"`
	DepartmentsAtRisk    int     `json:"departments_at_risk"`
	TotalDepartments     int     `json:"total_departments"`
}

// DepartmentDashboard representa o dashboard para supervisor de departamento
type DepartmentDashboard struct {
	DepartmentID         int64                        `json:"department_id"`
	DepartmentName       string                       `json:"department_name"`
	CompanyID            int64                        `json:"company_id"`
	CompanyName          string                       `json:"company_name"`
	ActiveTemplates []*DepartmentQuestionnaire   `json:"active_questionnaires"`
	EmployeesSummary     *EmployeesSummary            `json:"employees_summary"`
	ActionPlans          []*ActionPlanSummary         `json:"action_plans"`
	RiskCategories       []*RiskCategorySummary       `json:"risk_categories"`
	Alerts               []string                     `json:"alerts"`
}

// DepartmentQuestionnaire representa um template do departamento
type DepartmentQuestionnaire struct {
	TemplateID  int64   `json:"template_id"`
	TemplateName string `json:"template_name"`
	TotalEmployees   int64   `json:"total_employees"`
	Responses        int64   `json:"responses"`
	ResponseRate     float64 `json:"response_rate"`
	RiskLevel        string  `json:"risk_level"`
	CanClose         bool    `json:"can_close"`
}

// EmployeesSummary representa resumo de funcionários do departamento
type EmployeesSummary struct {
	Total     int      `json:"total"`
	Responded int      `json:"responded"`
	Pending   int      `json:"pending"`
	PendingNames []string `json:"pending_names,omitempty"`
}

// ActionPlanSummary representa resumo de um action plan
type ActionPlanSummary struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date,omitempty"`
	Status      string `json:"status"`
	IsOverdue   bool   `json:"is_overdue"`
}

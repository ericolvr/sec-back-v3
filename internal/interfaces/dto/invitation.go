package dto

type InvitationResponse struct {
	ID              int64  `json:"id"`
	PartnerID       int64  `json:"partner_id"`
	TemplateID      int64  `json:"template_id"`
	TemplateName    string `json:"template_name,omitempty"`
	DepartmentID    int64  `json:"department_id"`
	ResponseID      int64  `json:"response_id"`
	InvitationToken string `json:"invitation_token"`
	EmployeeEmail   string `json:"employee_email"`
	Status          string `json:"status"`
	SentAt          string `json:"sent_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type InvitationSummaryResponse struct {
	TemplateID       int64  `json:"template_id"`
	TemplateName     string `json:"template_name"`
	DepartmentID     int64  `json:"department_id"`
	DepartmentName   string `json:"department_name"`
	TotalInvitations int    `json:"total_invitations"`
	PendingCount     int    `json:"pending_count"`
	SentCount        int    `json:"sent_count"`
}

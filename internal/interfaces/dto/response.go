package dto

type ResponseRequest struct {
	QuestionnaireID int64  `json:"questionnaire_id" binding:"required"`
	InvitationToken string `json:"invitation_token" binding:"required"`
	CompanyID       int64  `json:"company_id" binding:"required"`
	DepartmentID    int64  `json:"department_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
}

type ResponseResponse struct {
	ID              int64  `json:"id"`
	PartnerID    int64  `json:"partner_id"`
	QuestionnaireID int64  `json:"questionnaire_id"`
	InvitationToken string `json:"invitation_token"`
	CompanyID       int64  `json:"company_id"`
	DepartmentID    int64  `json:"department_id"`
	Status          string `json:"status"`
	CompletedAt     string `json:"completed_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

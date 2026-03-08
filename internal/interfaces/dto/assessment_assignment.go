package dto

type CreateAssessmentAssignmentRequest struct {
	TemplateID   int64 `json:"template_id" binding:"required"`
	DepartmentID int64 `json:"department_id" binding:"required"`
}

type CreateAssessmentAssignmentBatchRequest struct {
	TemplateID    int64   `json:"template_id" binding:"required"`
	DepartmentIDs []int64 `json:"department_ids" binding:"required,min=1"`
}

type AssessmentAssignmentResponse struct {
	ID               int64  `json:"id"`
	PartnerID        int64  `json:"partner_id"`
	TemplateID       int64  `json:"template_id"`
	TemplateName     string `json:"template_name,omitempty"`
	DepartmentID     int64  `json:"department_id"`
	DepartmentName   string `json:"department_name,omitempty"`
	Active           bool   `json:"active"`
	ResponsesCreated int    `json:"responses_created,omitempty"`
	StartedAt        string `json:"started_at"`
	ClosedAt         string `json:"closed_at,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type BatchAssignmentResponse struct {
	Success []AssessmentAssignmentResponse `json:"success"`
	Errors  []BatchAssignmentError         `json:"errors"`
}

type BatchAssignmentError struct {
	DepartmentID int64  `json:"department_id"`
	Error        string `json:"error"`
}

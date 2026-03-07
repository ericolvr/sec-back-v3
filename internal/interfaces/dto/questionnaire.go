package dto

type QuestionnaireRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type QuestionnaireResponse struct {
	ID             int64  `json:"id"`
	PartnerID    int64  `json:"partner_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Active         bool   `json:"active"`
	TotalQuestions int    `json:"total_questions"`
	CreatedAt      string `json:"created_at"`
}

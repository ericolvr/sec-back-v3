package dto

type AnswerRequest struct {
	ResponseID int64  `json:"response_id" binding:"required"`
	QuestionID int64  `json:"question_id" binding:"required"`
	Value      string `json:"value" binding:"required"`
}

type AnswerResponse struct {
	ID         int64  `json:"id"`
	PartnerID    int64  `json:"partner_id"`
	ResponseID int64  `json:"response_id"`
	QuestionID int64  `json:"question_id"`
	Value      string `json:"value"`
	Score      *int   `json:"score,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

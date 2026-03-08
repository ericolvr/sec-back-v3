package dto

import "github.com/ericolvr/sec-back-v2/internal/core/domain"

type QuestionRequest struct {
	TemplateID  int64               `json:"template_id" binding:"required"`
	Question    string              `json:"question" binding:"required"`
	Type        domain.QuestionType `json:"type" binding:"required"`
	Category    string              `json:"category"`
	Options     []string            `json:"options"`
	ScoreValues []int               `json:"score_values"`
	Weight      float64             `json:"weight" binding:"required"`
	Required    bool                `json:"required"`
	OrderNum    int                 `json:"order_num" binding:"required"`
}

type QuestionResponse struct {
	ID           int64               `json:"id"`
	PartnerID    int64               `json:"partner_id"`
	TemplateID   int64               `json:"template_id"`
	TemplateName string              `json:"template_name,omitempty"`
	Question     string              `json:"question"`
	Type         domain.QuestionType `json:"type"`
	Category     string              `json:"category"`
	Options      []string            `json:"options"`
	ScoreValues  []int               `json:"score_values"`
	Weight       float64             `json:"weight"`
	Required     bool                `json:"required"`
	OrderNum     int                 `json:"order_num"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

type TemplateInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type QuestionListResponse struct {
	Template       TemplateInfo       `json:"template"`
	TotalQuestions int                `json:"total_questions"`
	Questions      []QuestionResponse `json:"questions"`
}

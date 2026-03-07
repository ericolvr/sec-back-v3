package domain

import (
	"context"
	"fmt"
	"time"
)

type QuestionType string

const (
	QuestionTypeText           QuestionType = "text"
	QuestionTypeNumber         QuestionType = "number"
	QuestionTypeScale          QuestionType = "scale"
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeYesNo          QuestionType = "yes_no"
)

func GetValidQuestionTypes() []QuestionType {
	return []QuestionType{
		QuestionTypeText,
		QuestionTypeNumber,
		QuestionTypeScale,
		QuestionTypeMultipleChoice,
		QuestionTypeYesNo,
	}
}

func (qt QuestionType) IsValid() bool {
	validTypes := GetValidQuestionTypes()
	for _, validType := range validTypes {
		if qt == validType {
			return true
		}
	}
	return false
}

func (qt QuestionType) String() string {
	return string(qt)
}

type Question struct {
	ID              int64        `json:"id"`
	PartnerID       int64        `json:"partner_id"`
	QuestionnaireID int64        `json:"questionnaire_id"`
	Question        string       `json:"question"`
	Type            QuestionType `json:"type"`         // QuestionType enum
	Category        string       `json:"category"`     // Categoria de risco: "Sobrecarga", "Autonomia", "Relacionamento", etc.
	Options         []string     `json:"options"`      // Para multiple_choice: ["Nunca", "Raramente", "Às vezes", "Frequentemente", "Sempre"]
	ScoreValues     []int        `json:"score_values"` // Scores correspondentes: [0, 1, 2, 3, 4]
	Weight          float64      `json:"weight"`       // Peso da pergunta (1.0 = normal, 2.0 = dobro de importância)
	Required        bool         `json:"required"`
	OrderNum        int          `json:"order_num"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`

	AssessmentTemplate *AssessmentTemplate `json:"assessment_template,omitempty"`
}

type QuestionRepository interface {
	Create(ctx context.Context, question *Question) error
	List(ctx context.Context, partnerID, questionnaireID, limit, offset int64) ([]*Question, error)
	GetByID(ctx context.Context, partnerID, id int64) (*Question, error)
	Update(ctx context.Context, question *Question) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (q *Question) Validate() error {
	if q.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if q.QuestionnaireID <= 0 {
		return fmt.Errorf("questionnaire_id is required")
	}
	if q.Question == "" {
		return fmt.Errorf("question text is required")
	}
	if q.Type == "" {
		return fmt.Errorf("question type is required")
	}
	if q.Type != QuestionTypeScale &&
		q.Type != QuestionTypeMultipleChoice &&
		q.Type != QuestionTypeText &&
		q.Type != QuestionTypeYesNo &&
		q.Type != QuestionTypeNumber {
		return fmt.Errorf("invalid question type")
	}
	if q.Weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}
	if q.Type == QuestionTypeMultipleChoice || q.Type == QuestionTypeScale {
		if len(q.Options) == 0 {
			return fmt.Errorf("options are required for %s type", q.Type)
		}
		if len(q.ScoreValues) != len(q.Options) {
			return fmt.Errorf("score_values must match the number of options")
		}
	}

	return nil
}

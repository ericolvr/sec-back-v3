package domain

import (
	"context"
	"fmt"
	"time"
)

type Answer struct {
	ID           int64     `json:"id"`
	PartnerID    int64     `json:"partner_id"`
	SubmissionID int64     `json:"submission_id"` // Referencia EmployeeSubmission
	QuestionID   int64     `json:"question_id"`
	Value        string    `json:"value"`
	Score        *int      `json:"score,omitempty"` // pontuação calculada
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Question *Question `json:"question,omitempty"`
}

type AnswerRepository interface {
	Create(ctx context.Context, answer *Answer) error
	List(ctx context.Context, partnerID, submissionID, limit, offset int64) ([]*Answer, error)
	GetByID(ctx context.Context, partnerID, id int64) (*Answer, error)
	CountBySubmission(ctx context.Context, partnerID, submissionID int64) (int64, error)
	Update(ctx context.Context, answer *Answer) error
	Delete(ctx context.Context, partnerID, id int64) error
}

func (a *Answer) Validate() error {
	if a.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if a.SubmissionID <= 0 {
		return fmt.Errorf("submission_id is required")
	}
	if a.QuestionID <= 0 {
		return fmt.Errorf("question_id is required")
	}
	if a.Value == "" {
		return fmt.Errorf("answer value is required")
	}
	return nil
}

func ValidateAnswer(a *Answer) error {
	return a.Validate()
}

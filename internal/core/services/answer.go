package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AnswerService struct {
	answerRepo     domain.AnswerRepository
	questionRepo   domain.QuestionRepository
	submissionRepo domain.EmployeeSubmissionRepository
}

func NewAnswerService(
	answerRepo domain.AnswerRepository,
	questionRepo domain.QuestionRepository,
	submissionRepo domain.EmployeeSubmissionRepository,
) *AnswerService {
	return &AnswerService{
		answerRepo:     answerRepo,
		questionRepo:   questionRepo,
		submissionRepo: submissionRepo,
	}
}

func (s *AnswerService) Create(ctx context.Context, answer *domain.Answer) error {
	if err := domain.ValidateAnswer(answer); err != nil {
		return err
	}

	// Auto-calculate score based on question
	if err := s.calculateScore(ctx, answer); err != nil {
		// Don't fail if score calculation fails, just proceed without score
		answer.Score = nil
	}

	if err := s.answerRepo.Create(ctx, answer); err != nil {
		return err
	}

	// Update Response status to "in_progress" after first answer
	if err := s.updateSubmissionStatus(ctx, answer.PartnerID, answer.SubmissionID); err != nil {
		// Log error but don't fail the answer creation
		// In production, you might want to log this
	}

	return nil
}

func (s *AnswerService) List(ctx context.Context, partnerID, responseID, limit, offset int64) ([]*domain.Answer, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.answerRepo.List(ctx, partnerID, responseID, limit, offset)
}

func (s *AnswerService) GetByID(ctx context.Context, partnerID, id int64) (*domain.Answer, error) {
	if id <= 0 {
		return nil, errors.New("ID is required")
	}

	return s.answerRepo.GetByID(ctx, partnerID, id)
}

func (s *AnswerService) Update(ctx context.Context, answer *domain.Answer) error {
	if answer.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if err := domain.ValidateAnswer(answer); err != nil {
		return err
	}

	existingAnswer, err := s.answerRepo.GetByID(ctx, answer.PartnerID, answer.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("answer not found")
		}
		return err
	}

	// Check if Response is already completed
	submission, err := s.submissionRepo.GetByID(ctx, answer.PartnerID, existingAnswer.SubmissionID)
	if err != nil {
		return err
	}

	if submission.Status == "completed" {
		return errors.New("cannot update answer: response already completed")
	}

	// Recalculate score if value changed
	if err := s.calculateScore(ctx, answer); err != nil {
		answer.Score = nil
	}

	return s.answerRepo.Update(ctx, answer)
}

// calculateScore calculates the score based on the question's score_values
func (s *AnswerService) calculateScore(ctx context.Context, answer *domain.Answer) error {
	question, err := s.questionRepo.GetByID(ctx, answer.PartnerID, answer.QuestionID)
	if err != nil {
		return err
	}

	// Only calculate score for questions with options and score_values
	if question.Type == domain.QuestionTypeScale || question.Type == domain.QuestionTypeMultipleChoice || question.Type == domain.QuestionTypeYesNo {
		// Find the index of the selected option
		for i, option := range question.Options {
			if option == answer.Value {
				if i < len(question.ScoreValues) {
					score := question.ScoreValues[i]
					answer.Score = &score
					return nil
				}
			}
		}
	}

	// For text/number questions, score is nil
	answer.Score = nil
	return nil
}

// updateSubmissionStatus updates the Response status based on answers count
// Automatically completes the Response when all required questions are answered
func (s *AnswerService) updateSubmissionStatus(ctx context.Context, partnerID, responseID int64) error {
	// Get the response
	submission, err := s.submissionRepo.GetByID(ctx, partnerID, responseID)
	if err != nil {
		return err
	}

	// Don't update if already completed
	if submission.Status == "completed" {
		return nil
	}

	// Count answers
	answersCount, err := s.answerRepo.CountBySubmission(ctx, partnerID, responseID)
	if err != nil {
		return err
	}

	// Update status based on answers count
	if answersCount == 0 {
		submission.Status = "pending"
	} else {
		// Check if all required questions are answered
		questions, err := s.questionRepo.List(ctx, partnerID, submission.TemplateID, MaxQuestionsPerTemplate, 0)
		if err != nil {
			// If can't get questions, just set to in_progress
			submission.Status = "in_progress"
		} else {
			// Count required questions
			requiredCount := int64(0)
			for _, q := range questions {
				if q.Required {
					requiredCount++
				}
			}

			// Auto-complete if all required questions are answered
			if answersCount >= requiredCount {
				submission.Status = "completed"
				now := time.Now()
				submission.CompletedAt = &now
			} else {
				submission.Status = "in_progress"
			}
		}
	}

	return s.submissionRepo.Update(ctx, submission)
}

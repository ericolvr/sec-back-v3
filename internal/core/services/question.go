package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type QuestionService struct {
	questionRepo domain.QuestionRepository
}

func NewQuestionService(repo domain.QuestionRepository) *QuestionService {
	return &QuestionService{questionRepo: repo}
}

func (s *QuestionService) Create(ctx context.Context, question *domain.Question) error {
	if err := question.Validate(); err != nil {
		return err
	}

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return err
	}

	return nil
}

func (s *QuestionService) List(ctx context.Context, partnerID, templateID, limit, offset int64) ([]*domain.Question, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.questionRepo.List(ctx, partnerID, templateID, limit, offset)
}

func (s *QuestionService) ListAll(ctx context.Context, partnerID, limit, offset int64) ([]*domain.Question, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.questionRepo.ListAllByPartner(ctx, partnerID, limit, offset)
}

func (s *QuestionService) GetByID(ctx context.Context, partnerID, id int64) (*domain.Question, error) {
	if id <= 0 {
		return nil, errors.New("ID is required")
	}

	return s.questionRepo.GetByID(ctx, partnerID, id)
}

func (s *QuestionService) Update(ctx context.Context, question *domain.Question) error {
	if question.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if err := question.Validate(); err != nil {
		return err
	}

	_, err := s.questionRepo.GetByID(ctx, question.PartnerID, question.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("question not found")
		}
		return err
	}

	return s.questionRepo.Update(ctx, question)
}

func (s *QuestionService) Delete(ctx context.Context, partnerID, id int64) (*domain.Question, error) {
	if id <= 0 {
		return nil, errors.New("id is required")
	}

	question, err := s.questionRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("question not found")
		}
		return nil, err
	}

	if err := s.questionRepo.Delete(ctx, partnerID, id); err != nil {
		return nil, err
	}
	return question, nil
}

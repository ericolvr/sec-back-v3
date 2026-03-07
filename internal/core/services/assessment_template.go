package services

import (
	"context"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentTemplateService struct {
	templateRepo domain.AssessmentTemplateRepository
	partnerRepo  domain.PartnerRepository
}

func NewAssessmentTemplateService(
	templateRepo domain.AssessmentTemplateRepository,
	partnerRepo domain.PartnerRepository,
) *AssessmentTemplateService {
	return &AssessmentTemplateService{
		templateRepo: templateRepo,
		partnerRepo:  partnerRepo,
	}
}

func (s *AssessmentTemplateService) Create(ctx context.Context, template *domain.AssessmentTemplate) error {
	if err := domain.ValidateAssessmentTemplate(template); err != nil {
		return err
	}

	partner, err := s.partnerRepo.GetByID(ctx, template.PartnerID)
	if err != nil {
		return errors.New("partner not found")
	}
	if !partner.Active {
		return errors.New("partner is not active")
	}

	if template.Version == 0 {
		template.Version = 1
	}

	return s.templateRepo.Create(ctx, template)
}

func (s *AssessmentTemplateService) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.AssessmentTemplate, error) {
	if partnerID <= 0 {
		return nil, errors.New("invalid partner ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.templateRepo.List(ctx, partnerID, limit, offset)
}

func (s *AssessmentTemplateService) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentTemplate, error) {
	if partnerID <= 0 || id <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.templateRepo.GetByID(ctx, partnerID, id)
}

func (s *AssessmentTemplateService) Update(ctx context.Context, template *domain.AssessmentTemplate) error {
	if template.ID <= 0 || template.PartnerID <= 0 {
		return errors.New("invalid IDs")
	}

	if err := domain.ValidateAssessmentTemplate(template); err != nil {
		return err
	}

	return s.templateRepo.Update(ctx, template)
}

func (s *AssessmentTemplateService) IncrementVersion(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.templateRepo.IncrementVersion(ctx, partnerID, id)
}

func (s *AssessmentTemplateService) Delete(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.templateRepo.Delete(ctx, partnerID, id)
}

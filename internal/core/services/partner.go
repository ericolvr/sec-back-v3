package services

import (
	"context"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type PartnerService struct {
	partnerRepo domain.PartnerRepository
}

func NewPartnerService(partnerRepo domain.PartnerRepository) *PartnerService {
	return &PartnerService{
		partnerRepo: partnerRepo,
	}
}

func (s *PartnerService) Create(ctx context.Context, partner *domain.Partner) error {
	if err := domain.ValidatePartner(partner); err != nil {
		return err
	}

	return s.partnerRepo.Create(ctx, partner)
}

func (s *PartnerService) List(ctx context.Context, limit, offset int64) ([]*domain.Partner, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.partnerRepo.List(ctx, limit, offset)
}

func (s *PartnerService) GetByID(ctx context.Context, id int64) (*domain.Partner, error) {
	if id <= 0 {
		return nil, errors.New("invalid partner ID")
	}

	return s.partnerRepo.GetByID(ctx, id)
}

func (s *PartnerService) Update(ctx context.Context, partner *domain.Partner) error {
	if partner.ID <= 0 {
		return errors.New("invalid partner ID")
	}

	if err := domain.ValidatePartner(partner); err != nil {
		return err
	}

	return s.partnerRepo.Update(ctx, partner)
}

func (s *PartnerService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid partner ID")
	}

	return s.partnerRepo.Delete(ctx, id)
}

package services

import (
	"context"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type CompanyService struct {
	companyRepo domain.CompanyRepository
	partnerRepo domain.PartnerRepository
}

func NewCompanyService(companyRepo domain.CompanyRepository, partnerRepo domain.PartnerRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
		partnerRepo: partnerRepo,
	}
}

func (s *CompanyService) Create(ctx context.Context, company *domain.Company) error {
	if err := domain.ValidateCompany(company); err != nil {
		return err
	}

	partner, err := s.partnerRepo.GetByID(ctx, company.PartnerID)
	if err != nil {
		return errors.New("partner not found")
	}
	if !partner.Active {
		return errors.New("partner is not active")
	}

	return s.companyRepo.Create(ctx, company)
}

func (s *CompanyService) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.Company, error) {
	if partnerID <= 0 {
		return nil, errors.New("invalid partner ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.companyRepo.List(ctx, partnerID, limit, offset)
}

func (s *CompanyService) GetByID(ctx context.Context, partnerID, id int64) (*domain.Company, error) {
	if partnerID <= 0 || id <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.companyRepo.GetByID(ctx, partnerID, id)
}

func (s *CompanyService) Update(ctx context.Context, company *domain.Company) error {
	if company.ID <= 0 || company.PartnerID <= 0 {
		return errors.New("invalid IDs")
	}

	if err := domain.ValidateCompany(company); err != nil {
		return err
	}

	return s.companyRepo.Update(ctx, company)
}

func (s *CompanyService) Delete(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.companyRepo.Delete(ctx, partnerID, id)
}

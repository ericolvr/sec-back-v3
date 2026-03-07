package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type DepartmentService struct {
	departmentRepo domain.DepartmentRepository
}

func NewDepartmentService(repo domain.DepartmentRepository) *DepartmentService {
	return &DepartmentService{departmentRepo: repo}
}

func (s *DepartmentService) Create(ctx context.Context, department *domain.Department) error {
	if err := department.Validate(); err != nil {
		return err
	}

	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return err
	}

	return nil
}

func (s *DepartmentService) List(ctx context.Context, partnerID int64, limit, offset int) ([]*domain.Department, error) {
	if limit <= 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.departmentRepo.List(ctx, partnerID, int64(limit), int64(offset))
}

func (s *DepartmentService) ListByCompany(ctx context.Context, partnerID, companyID int64, limit, offset int) ([]*domain.Department, error) {
	if limit <= 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.departmentRepo.ListByCompany(ctx, partnerID, companyID, int64(limit), int64(offset))
}

func (s *DepartmentService) GetByID(ctx context.Context, partnerID int64, id int) (*domain.Department, error) {
	if id <= 0 {
		return nil, errors.New("ID is required")
	}

	return s.departmentRepo.GetByID(ctx, partnerID, int64(id))
}

func (s *DepartmentService) Update(ctx context.Context, department *domain.Department) error {
	if department.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if department.Name == "" {
		return errors.New("name is required")
	}

	_, err := s.departmentRepo.GetByID(ctx, department.PartnerID, int64(department.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("department not found")
		}
		return err
	}

	return s.departmentRepo.Update(ctx, department)
}

func (s *DepartmentService) Delete(ctx context.Context, partnerID int64, id int) (*domain.Department, error) {
	if id <= 0 {
		return nil, errors.New("id is required")
	}

	department, err := s.departmentRepo.GetByID(ctx, partnerID, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("department not found")
		}
		return nil, err
	}

	if err := s.departmentRepo.Delete(ctx, partnerID, int64(id)); err != nil {
		return nil, err
	}
	return department, nil
}

func (s *DepartmentService) ToggleActive(ctx context.Context, partnerID, id int64, active bool) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.departmentRepo.ToggleActive(ctx, partnerID, id, active)
}

func (s *DepartmentService) ListDeleted(ctx context.Context, partnerID int64, limit, offset int) ([]*domain.Department, error) {
	if limit <= 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.departmentRepo.ListDeleted(ctx, partnerID, int64(limit), int64(offset))
}

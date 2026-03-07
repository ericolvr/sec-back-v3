package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type EmployeeService struct {
	employeeRepo domain.EmployeeRepository
}

func NewEmployeeService(repo domain.EmployeeRepository) *EmployeeService {
	return &EmployeeService{employeeRepo: repo}
}

func (s *EmployeeService) Create(ctx context.Context, employee *domain.Employee) error {
	if err := employee.Validate(); err != nil {
		return err
	}

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return err
	}

	return nil
}

func (s *EmployeeService) List(ctx context.Context, partnerID int64, limit, offset int) ([]*domain.Employee, error) {
	if limit <= 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.employeeRepo.List(ctx, partnerID, int64(limit), int64(offset))
}

func (s *EmployeeService) GetByID(ctx context.Context, partnerID, id int64) (*domain.Employee, error) {
	if id <= 0 {
		return nil, errors.New("ID is required")
	}

	return s.employeeRepo.GetByID(ctx, partnerID, id)
}

func (s *EmployeeService) Update(ctx context.Context, employee *domain.Employee) error {
	if employee.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if err := employee.Validate(); err != nil {
		return err
	}

	_, err := s.employeeRepo.GetByID(ctx, employee.PartnerID, employee.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("employee not found")
		}
		return err
	}

	return s.employeeRepo.Update(ctx, employee)
}

func (s *EmployeeService) Delete(ctx context.Context, partnerID, id int64) (*domain.Employee, error) {
	if id <= 0 {
		return nil, errors.New("id is required")
	}

	employee, err := s.employeeRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}

	if err := s.employeeRepo.Delete(ctx, partnerID, id); err != nil {
		return nil, err
	}
	return employee, nil
}

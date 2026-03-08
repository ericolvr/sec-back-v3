package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentAssignmentService struct {
	assignmentRepo domain.AssessmentAssignmentRepository
	departmentRepo domain.DepartmentRepository
}

func NewAssessmentAssignmentService(
	assignmentRepo domain.AssessmentAssignmentRepository,
	departmentRepo domain.DepartmentRepository,
) *AssessmentAssignmentService {
	return &AssessmentAssignmentService{
		assignmentRepo: assignmentRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *AssessmentAssignmentService) Create(ctx context.Context, assignment *domain.AssessmentAssignment) error {
	if err := assignment.Validate(); err != nil {
		return err
	}

	_, err := s.departmentRepo.GetByID(ctx, assignment.PartnerID, assignment.DepartmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("department not found")
		}
		return err
	}

	existing, err := s.assignmentRepo.GetByTemplateAndDepartment(ctx, assignment.PartnerID, assignment.TemplateID, assignment.DepartmentID)
	if err == nil && existing != nil {
		return errors.New("assessment already assigned to this department")
	}

	return s.assignmentRepo.Create(ctx, assignment)
}

func (s *AssessmentAssignmentService) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentAssignment, error) {
	return s.assignmentRepo.GetByID(ctx, partnerID, id)
}

func (s *AssessmentAssignmentService) GetByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) (*domain.AssessmentAssignment, error) {
	return s.assignmentRepo.GetByTemplateAndDepartment(ctx, partnerID, templateID, departmentID)
}

func (s *AssessmentAssignmentService) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.assignmentRepo.List(ctx, partnerID, limit, offset)
}

func (s *AssessmentAssignmentService) ListByTemplate(ctx context.Context, partnerID, templateID int64) ([]*domain.AssessmentAssignment, error) {
	return s.assignmentRepo.ListByTemplate(ctx, partnerID, templateID, MaxAssignmentsPerQuestionnaire, 0)
}

func (s *AssessmentAssignmentService) Deactivate(ctx context.Context, partnerID, id int64) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	assignment.Active = false
	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *AssessmentAssignmentService) Activate(ctx context.Context, partnerID, id int64) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	assignment.Active = true
	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *AssessmentAssignmentService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.assignmentRepo.Delete(ctx, partnerID, id)
}

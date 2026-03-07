package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type QuestionnaireAssignmentService struct {
	assignmentRepo domain.QuestionnaireAssignmentRepository
	departmentRepo domain.DepartmentRepository
}

func NewQuestionnaireAssignmentService(
	assignmentRepo domain.QuestionnaireAssignmentRepository,
	departmentRepo domain.DepartmentRepository,
) *QuestionnaireAssignmentService {
	return &QuestionnaireAssignmentService{
		assignmentRepo: assignmentRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *QuestionnaireAssignmentService) Create(ctx context.Context, assignment *domain.QuestionnaireAssignment) error {
	if err := assignment.ValidateQuestionnaireAssignment(); err != nil {
		return err
	}

	_, err := s.departmentRepo.GetByID(ctx, assignment.PartnerID, assignment.DepartmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("department not found")
		}
		return err
	}

	existing, err := s.assignmentRepo.GetByQuestionnaireAndDepartment(ctx, assignment.PartnerID, assignment.QuestionnaireID, assignment.DepartmentID)
	if err == nil && existing != nil {
		return errors.New("questionnaire already assigned to this department")
	}

	return s.assignmentRepo.Create(ctx, assignment)
}

func (s *QuestionnaireAssignmentService) GetByID(ctx context.Context, partnerID, id int64) (*domain.QuestionnaireAssignment, error) {
	return s.assignmentRepo.GetByID(ctx, partnerID, id)
}

func (s *QuestionnaireAssignmentService) GetByQuestionnaireAndDepartment(ctx context.Context, partnerID, questionnaireID, departmentID int64) (*domain.QuestionnaireAssignment, error) {
	return s.assignmentRepo.GetByQuestionnaireAndDepartment(ctx, partnerID, questionnaireID, departmentID)
}

func (s *QuestionnaireAssignmentService) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.QuestionnaireAssignment, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.assignmentRepo.List(ctx, partnerID, limit, offset)
}

func (s *QuestionnaireAssignmentService) ListByQuestionnaire(ctx context.Context, partnerID, questionnaireID int64) ([]*domain.QuestionnaireAssignment, error) {
	return s.assignmentRepo.ListByQuestionnaire(ctx, partnerID, questionnaireID, MaxAssignmentsPerQuestionnaire, 0)
}

func (s *QuestionnaireAssignmentService) Deactivate(ctx context.Context, partnerID, id int64) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	assignment.Active = false
	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *QuestionnaireAssignmentService) Activate(ctx context.Context, partnerID, id int64) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	assignment.Active = true
	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *QuestionnaireAssignmentService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.assignmentRepo.Delete(ctx, partnerID, id)
}

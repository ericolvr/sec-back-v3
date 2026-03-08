package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/google/uuid"
)

type AssessmentAssignmentService struct {
	assignmentRepo domain.AssessmentAssignmentRepository
	departmentRepo domain.DepartmentRepository
	employeeRepo   domain.EmployeeRepository
	submissionRepo domain.EmployeeSubmissionRepository
	invitationRepo domain.InvitationRepository
}

func NewAssessmentAssignmentService(
	assignmentRepo domain.AssessmentAssignmentRepository,
	departmentRepo domain.DepartmentRepository,
	employeeRepo domain.EmployeeRepository,
	submissionRepo domain.EmployeeSubmissionRepository,
	invitationRepo domain.InvitationRepository,
) *AssessmentAssignmentService {
	return &AssessmentAssignmentService{
		assignmentRepo: assignmentRepo,
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
		submissionRepo: submissionRepo,
		invitationRepo: invitationRepo,
	}
}

func (s *AssessmentAssignmentService) Create(ctx context.Context, assignment *domain.AssessmentAssignment) (int, error) {
	if err := assignment.Validate(); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	// Validar que todos os departments existem
	for _, deptID := range assignment.DepartmentIDs {
		_, err := s.departmentRepo.GetByID(ctx, assignment.PartnerID, deptID)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, fmt.Errorf("department %d not found", deptID)
			}
			return 0, fmt.Errorf("error fetching department %d: %w", deptID, err)
		}
	}

	// Criar o Assignment
	if err := s.assignmentRepo.Create(ctx, assignment); err != nil {
		return 0, fmt.Errorf("error creating assignment in database: %w", err)
	}

	// Criar EmployeeSubmissions e Invitations para cada department
	totalSubmissions := 0
	for _, deptID := range assignment.DepartmentIDs {
		count, err := s.createSubmissionsForDepartment(ctx, assignment, deptID)
		if err != nil {
			return totalSubmissions, fmt.Errorf("error creating submissions for department %d: %w", deptID, err)
		}
		totalSubmissions += count
	}

	return totalSubmissions, nil
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

func (s *AssessmentAssignmentService) Update(ctx context.Context, partnerID, id int64, templateID *int64, departmentIDs []int64, active *bool) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Atualizar campos se fornecidos
	if templateID != nil {
		assignment.TemplateID = *templateID
	}
	if len(departmentIDs) > 0 {
		// Validar que todos os departments existem
		for _, deptID := range departmentIDs {
			_, err := s.departmentRepo.GetByID(ctx, partnerID, deptID)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("department not found")
				}
				return err
			}
		}
		assignment.DepartmentIDs = departmentIDs
	}
	if active != nil {
		assignment.Active = *active
	}

	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *AssessmentAssignmentService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.assignmentRepo.Delete(ctx, partnerID, id)
}

func (s *AssessmentAssignmentService) createSubmissionsForDepartment(ctx context.Context, assignment *domain.AssessmentAssignment, departmentID int64) (int, error) {
	// Buscar todos os employees do department
	employees, err := s.employeeRepo.ListByDepartment(ctx, assignment.PartnerID, departmentID, 1000, 0)
	if err != nil {
		return 0, fmt.Errorf("error listing employees: %w", err)
	}

	count := 0
	for _, employee := range employees {
		// Criar EmployeeSubmission
		submission := &domain.EmployeeSubmission{
			PartnerID:       assignment.PartnerID,
			CompanyID:       employee.CompanyID,
			TemplateID:      assignment.TemplateID,
			EmployeeID:      employee.ID,
			DepartmentID:    departmentID,
			InvitationToken: generateInvitationToken(),
			Status:          domain.SubmissionStatusPending,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := s.submissionRepo.Create(ctx, submission); err != nil {
			return count, err
		}

		// Criar Invitation vinculada ao EmployeeSubmission
		invitation := &domain.Invitation{
			PartnerID:     assignment.PartnerID,
			EmployeeID:    employee.ID,
			TemplateID:    assignment.TemplateID,
			DepartmentID:  departmentID,
			ResponseID:    submission.ID,
			Token:         submission.InvitationToken,
			EmployeeEmail: employee.Email,
			Status:        domain.InvitationStatusPending,
			Sent:          false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := s.invitationRepo.Create(ctx, invitation); err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

func generateInvitationToken() string {
	return uuid.New().String()
}

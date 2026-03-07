package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type EmployeeSubmissionService struct {
	submissionRepo domain.EmployeeSubmissionRepository
	employeeRepo   domain.EmployeeRepository
	templateRepo   domain.AssessmentTemplateRepository
}

func NewEmployeeSubmissionService(
	submissionRepo domain.EmployeeSubmissionRepository,
	employeeRepo domain.EmployeeRepository,
	templateRepo domain.AssessmentTemplateRepository,
) *EmployeeSubmissionService {
	return &EmployeeSubmissionService{
		submissionRepo: submissionRepo,
		employeeRepo:   employeeRepo,
		templateRepo:   templateRepo,
	}
}

func (s *EmployeeSubmissionService) Create(ctx context.Context, submission *domain.EmployeeSubmission) error {
	if err := domain.ValidateEmployeeSubmission(submission); err != nil {
		return err
	}

	employee, err := s.employeeRepo.GetByID(ctx, submission.PartnerID, submission.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}
	if !employee.Active {
		return errors.New("employee is not active")
	}

	template, err := s.templateRepo.GetByID(ctx, submission.PartnerID, submission.QuestionnaireID)
	if err != nil {
		return errors.New("assessment template not found")
	}
	if !template.Active {
		return errors.New("assessment template is not active")
	}

	if submission.InvitationToken == "" {
		token, err := generateToken()
		if err != nil {
			return err
		}
		submission.InvitationToken = token
	}

	if submission.Status == "" {
		submission.Status = "pending"
	}

	return s.submissionRepo.Create(ctx, submission)
}

func (s *EmployeeSubmissionService) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.EmployeeSubmission, error) {
	if partnerID <= 0 {
		return nil, errors.New("invalid partner ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.submissionRepo.List(ctx, partnerID, limit, offset)
}

func (s *EmployeeSubmissionService) ListByCompany(ctx context.Context, partnerID, companyID int64, limit, offset int64) ([]*domain.EmployeeSubmission, error) {
	if partnerID <= 0 || companyID <= 0 {
		return nil, errors.New("invalid IDs")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.submissionRepo.ListByCompany(ctx, partnerID, companyID, limit, offset)
}

func (s *EmployeeSubmissionService) ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*domain.EmployeeSubmission, error) {
	if partnerID <= 0 || departmentID <= 0 {
		return nil, errors.New("invalid IDs")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.submissionRepo.ListByDepartment(ctx, partnerID, departmentID, limit, offset)
}

func (s *EmployeeSubmissionService) GetByID(ctx context.Context, partnerID, id int64) (*domain.EmployeeSubmission, error) {
	if partnerID <= 0 || id <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.submissionRepo.GetByID(ctx, partnerID, id)
}

func (s *EmployeeSubmissionService) GetByToken(ctx context.Context, token string) (*domain.EmployeeSubmission, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}

	return s.submissionRepo.GetByToken(ctx, token)
}

func (s *EmployeeSubmissionService) Complete(ctx context.Context, partnerID, id int64) error {
	submission, err := s.submissionRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	if submission.Status == "completed" {
		return errors.New("submission already completed")
	}

	submission.Status = "completed"
	now := time.Now()
	submission.CompletedAt = &now

	return s.submissionRepo.Update(ctx, submission)
}

func (s *EmployeeSubmissionService) UpdateStatus(ctx context.Context, partnerID, id int64, status string) error {
	if !domain.IsValidSubmissionStatus(status) {
		return errors.New("invalid status")
	}

	submission, err := s.submissionRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	submission.Status = status
	if status == "completed" && submission.CompletedAt == nil {
		now := time.Now()
		submission.CompletedAt = &now
	}

	return s.submissionRepo.Update(ctx, submission)
}

func (s *EmployeeSubmissionService) Delete(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.submissionRepo.Delete(ctx, partnerID, id)
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

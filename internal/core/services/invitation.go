package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type InvitationService struct {
	invitationRepo domain.InvitationRepository
	submissionRepo domain.EmployeeSubmissionRepository
	employeeRepo   domain.EmployeeRepository
}

func NewInvitationService(
	invitationRepo domain.InvitationRepository,
	submissionRepo domain.EmployeeSubmissionRepository,
	employeeRepo domain.EmployeeRepository,
) *InvitationService {
	return &InvitationService{
		invitationRepo: invitationRepo,
		submissionRepo: submissionRepo,
		employeeRepo:   employeeRepo,
	}
}

func (s *InvitationService) Create(ctx context.Context, invitation *domain.Invitation) error {
	if err := invitation.ValidateInvitation(); err != nil {
		return err
	}
	return s.invitationRepo.Create(ctx, invitation)
}

func (s *InvitationService) GetByID(ctx context.Context, partnerID, id int64) (*domain.Invitation, error) {
	return s.invitationRepo.GetByID(ctx, partnerID, id)
}

func (s *InvitationService) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.Invitation, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.invitationRepo.List(ctx, partnerID, limit, offset)
}

func (s *InvitationService) ListByQuestionnaireAndDepartment(ctx context.Context, partnerID, questionnaireID, departmentID int64) ([]*domain.Invitation, error) {
	return s.invitationRepo.ListByQuestionnaireAndDepartment(ctx, partnerID, questionnaireID, departmentID)
}

func (s *InvitationService) MarkAsSent(ctx context.Context, partnerID, id int64) error {
	invitation, err := s.invitationRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("invitation not found")
		}
		return err
	}

	now := time.Now()
	invitation.Status = domain.InvitationStatusSent
	invitation.SentAt = &now

	return s.invitationRepo.Update(ctx, invitation)
}

func (s *InvitationService) MarkAsFailed(ctx context.Context, partnerID, id int64) error {
	invitation, err := s.invitationRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("invitation not found")
		}
		return err
	}

	invitation.Status = domain.InvitationStatusFailed
	return s.invitationRepo.Update(ctx, invitation)
}

func (s *InvitationService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.invitationRepo.Delete(ctx, partnerID, id)
}

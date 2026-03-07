package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type InvitationRepository struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Create(ctx context.Context, invitation *domain.Invitation) error {
	query := `
		INSERT INTO invitations (
			partner_id, employee_id, questionnaire_id, department_id, 
			token, sent, sent_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		invitation.PartnerID,
		invitation.ResponseID, // employee_id
		invitation.QuestionnaireID,
		invitation.DepartmentID,
		invitation.EmployeeEmail, // token
		invitation.Status == domain.InvitationStatusSent,
		invitation.SentAt,
	).Scan(&invitation.ID, &invitation.CreatedAt)
}

func (r *InvitationRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.Invitation, error) {
	query := `
		SELECT id, partner_id, employee_id, questionnaire_id, department_id,
			   token, sent, sent_at, created_at
		FROM invitations
		WHERE partner_id = $1 AND id = $2`

	var inv domain.Invitation
	var sent bool

	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&inv.ID, &inv.PartnerID, &inv.ResponseID, &inv.QuestionnaireID, &inv.DepartmentID,
		&inv.EmployeeEmail, &sent, &inv.SentAt, &inv.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if sent {
		inv.Status = domain.InvitationStatusSent
	} else {
		inv.Status = domain.InvitationStatusPending
	}

	return &inv, nil
}

func (r *InvitationRepository) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.Invitation, error) {
	query := `
		SELECT id, partner_id, employee_id, questionnaire_id, department_id,
			   token, sent, sent_at, created_at
		FROM invitations
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanInvitations(rows)
}

func (r *InvitationRepository) ListByQuestionnaireAndDepartment(ctx context.Context, partnerID, questionnaireID, departmentID int64) ([]*domain.Invitation, error) {
	query := `
		SELECT id, partner_id, employee_id, questionnaire_id, department_id,
			   token, sent, sent_at, created_at
		FROM invitations
		WHERE partner_id = $1 AND questionnaire_id = $2 AND department_id = $3
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, partnerID, questionnaireID, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanInvitations(rows)
}

func (r *InvitationRepository) ListByStatus(ctx context.Context, partnerID int64, status string, limit, offset int64) ([]*domain.Invitation, error) {
	sent := status == domain.InvitationStatusSent

	query := `
		SELECT id, partner_id, employee_id, questionnaire_id, department_id,
			   token, sent, sent_at, created_at
		FROM invitations
		WHERE partner_id = $1 AND sent = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, sent, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanInvitations(rows)
}

func (r *InvitationRepository) Update(ctx context.Context, invitation *domain.Invitation) error {
	sent := invitation.Status == domain.InvitationStatusSent

	query := `
		UPDATE invitations SET
			sent = $1, sent_at = $2
		WHERE partner_id = $3 AND id = $4`

	_, err := r.db.ExecContext(ctx, query, sent, invitation.SentAt, invitation.PartnerID, invitation.ID)
	return err
}

func (r *InvitationRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM invitations WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *InvitationRepository) CountByQuestionnaireAndDepartment(ctx context.Context, partnerID, questionnaireID, departmentID int64) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM invitations 
		WHERE partner_id = $1 AND questionnaire_id = $2 AND department_id = $3`

	var count int64
	err := r.db.QueryRowContext(ctx, query, partnerID, questionnaireID, departmentID).Scan(&count)
	return count, err
}

func (r *InvitationRepository) scanInvitations(rows *sql.Rows) ([]*domain.Invitation, error) {
	var invitations []*domain.Invitation

	for rows.Next() {
		var inv domain.Invitation
		var sent bool

		err := rows.Scan(
			&inv.ID, &inv.PartnerID, &inv.ResponseID, &inv.QuestionnaireID, &inv.DepartmentID,
			&inv.EmployeeEmail, &sent, &inv.SentAt, &inv.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if sent {
			inv.Status = domain.InvitationStatusSent
		} else {
			inv.Status = domain.InvitationStatusPending
		}

		invitations = append(invitations, &inv)
	}

	return invitations, rows.Err()
}

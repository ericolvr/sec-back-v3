package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (
			partner_id, company_id, department_id, template_id,
			type, title, message, severity, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		notification.PartnerID,
		notification.CompanyID,
		notification.DepartmentID,
		notification.TemplateID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Severity,
		notification.Metadata,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *NotificationRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.Notification, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, template_id,
			   type, title, message, severity, metadata,
			   read, read_at, read_by, created_at
		FROM notifications
		WHERE partner_id = $1 AND id = $2`

	var n domain.Notification
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&n.ID, &n.PartnerID, &n.CompanyID, &n.DepartmentID, &n.TemplateID,
		&n.Type, &n.Title, &n.Message, &n.Severity, &n.Metadata,
		&n.Read, &n.ReadAt, &n.ReadBy, &n.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NotificationRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.Notification, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, template_id,
			   type, title, message, severity, metadata,
			   read, read_at, read_by, created_at
		FROM notifications
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *NotificationRepository) ListUnread(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.Notification, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, template_id,
			   type, title, message, severity, metadata,
			   read, read_at, read_by, created_at
		FROM notifications
		WHERE partner_id = $1 AND read = false
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *NotificationRepository) ListByType(ctx context.Context, partnerID int64, notifType string, limit, offset int64) ([]*domain.Notification, error) {
	query := `
		SELECT id, partner_id, company_id, department_id, template_id,
			   type, title, message, severity, metadata,
			   read, read_at, read_by, created_at
		FROM notifications
		WHERE partner_id = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, partnerID, notifType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, partnerID, id, userID int64) error {
	query := `
		UPDATE notifications SET
			read = true,
			read_at = $1,
			read_by = $2
		WHERE partner_id = $3 AND id = $4`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, userID, partnerID, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, partnerID, userID int64) error {
	query := `
		UPDATE notifications SET
			read = true,
			read_at = $1,
			read_by = $2
		WHERE partner_id = $3 AND read = false`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, userID, partnerID)
	return err
}

func (r *NotificationRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM notifications WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

func (r *NotificationRepository) CountUnread(ctx context.Context, partnerID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE partner_id = $1 AND read = false`

	var count int64
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) scanNotifications(rows *sql.Rows) ([]*domain.Notification, error) {
	var notifications []*domain.Notification

	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(
			&n.ID, &n.PartnerID, &n.CompanyID, &n.DepartmentID, &n.TemplateID,
			&n.Type, &n.Title, &n.Message, &n.Severity, &n.Metadata,
			&n.Read, &n.ReadAt, &n.ReadBy, &n.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		notifications = append(notifications, &n)
	}

	return notifications, rows.Err()
}

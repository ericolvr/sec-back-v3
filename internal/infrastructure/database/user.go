package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (partner_id, name, mobile, password, type, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.PartnerID,
		user.Name,
		user.Mobile,
		user.Password,
		user.Type,
		user.Active,
		now,
	).Scan(&user.ID)

	return err
}

func (r *UserRepository) List(ctx context.Context, tenantID int64, limit, offset int64) ([]*domain.User, error) {
	query := `
		SELECT id, partner_id, name, mobile, type, active, created_at 
		FROM users 
		WHERE partner_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID,
			&u.PartnerID,
			&u.Name,
			&u.Mobile,
			&u.Type,
			&u.Active,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, tenantID, id int64) (*domain.User, error) {
	query := `
		SELECT id, partner_id, name, mobile, type, active, created_at 
		FROM users 
		WHERE partner_id = $1 AND id = $2`

	var u domain.User
	err := r.db.QueryRowContext(
		ctx,
		query,
		tenantID,
		id,
	).Scan(
		&u.ID,
		&u.PartnerID,
		&u.Name,
		&u.Mobile,
		&u.Type,
		&u.Active,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByMobile(ctx context.Context, tenantID int64, mobile string) (*domain.User, error) {
	query := `
		SELECT id, partner_id, name, mobile, password, type, active, created_at 
		FROM users 
		WHERE partner_id = $1 AND mobile = $2`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, tenantID, mobile).Scan(
		&u.ID,
		&u.PartnerID,
		&u.Name,
		&u.Mobile,
		&u.Password,
		&u.Type,
		&u.Active,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByMobileGlobal(ctx context.Context, mobile string) (*domain.User, error) {
	query := `
		SELECT id, partner_id, name, mobile, password, type, active, created_at 
		FROM users 
		WHERE mobile = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, mobile).Scan(
		&u.ID,
		&u.PartnerID,
		&u.Name,
		&u.Mobile,
		&u.Password,
		&u.Type,
		&u.Active,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, tenantID int64, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, mobile = $2, type = $3, active = $4
		WHERE partner_id = $5 AND id = $6
		RETURNING id, partner_id, name, mobile, type, active, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Mobile,
		user.Type,
		user.Active,
		tenantID,
		user.ID,
	).Scan(
		&user.ID,
		&user.PartnerID,
		&user.Name,
		&user.Mobile,
		&user.Type,
		&user.Active,
		&user.CreatedAt,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, tenantID, id int64) error {
	query := `DELETE FROM users WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, tenantID, id)
	return err
}

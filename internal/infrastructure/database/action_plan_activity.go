package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActionPlanActivityRepository struct {
	db *sql.DB
}

func NewActionPlanActivityRepository(db *sql.DB) *ActionPlanActivityRepository {
	return &ActionPlanActivityRepository{db: db}
}

func (r *ActionPlanActivityRepository) Create(ctx context.Context, activity *domain.ActionPlanActivity) error {
	query := `
		INSERT INTO action_plan_activities (
			action_plan_id, type, title, description, status, 
			due_date, completed_at, created_by, created_by_name
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		activity.ActionPlanID,
		activity.Type,
		activity.Title,
		activity.Description,
		activity.Status,
		activity.DueDate,
		activity.CompletedAt,
		activity.CreatedBy,
		activity.CreatedByName,
	).Scan(&activity.ID, &activity.CreatedAt, &activity.UpdatedAt)
}

func (r *ActionPlanActivityRepository) GetByID(ctx context.Context, id int64) (*domain.ActionPlanActivity, error) {
	query := `
		SELECT id, action_plan_id, type, title, description, status,
			due_date, completed_at, created_by, created_by_name,
			created_at, updated_at
		FROM action_plan_activities
		WHERE id = $1
	`

	activity := &domain.ActionPlanActivity{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.ActionPlanID,
		&activity.Type,
		&activity.Title,
		&activity.Description,
		&activity.Status,
		&activity.DueDate,
		&activity.CompletedAt,
		&activity.CreatedBy,
		&activity.CreatedByName,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Carregar mídias associadas
	medias, err := r.loadMedias(ctx, activity.ID)
	if err != nil {
		return nil, err
	}
	activity.Medias = medias

	return activity, nil
}

func (r *ActionPlanActivityRepository) ListByActionPlan(ctx context.Context, actionPlanID int64) ([]*domain.ActionPlanActivity, error) {
	query := `
		SELECT id, action_plan_id, type, title, description, status,
			due_date, completed_at, created_by, created_by_name,
			created_at, updated_at
		FROM action_plan_activities
		WHERE action_plan_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, actionPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*domain.ActionPlanActivity

	for rows.Next() {
		activity := &domain.ActionPlanActivity{}

		err := rows.Scan(
			&activity.ID,
			&activity.ActionPlanID,
			&activity.Type,
			&activity.Title,
			&activity.Description,
			&activity.Status,
			&activity.DueDate,
			&activity.CompletedAt,
			&activity.CreatedBy,
			&activity.CreatedByName,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Carregar mídias associadas
		medias, err := r.loadMedias(ctx, activity.ID)
		if err != nil {
			return nil, err
		}
		activity.Medias = medias

		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *ActionPlanActivityRepository) Update(ctx context.Context, activity *domain.ActionPlanActivity) error {
	query := `
		UPDATE action_plan_activities
		SET title = $1, description = $2, status = $3, due_date = $4,
			completed_at = $5, updated_at = NOW()
		WHERE id = $6
	`

	_, err := r.db.ExecContext(
		ctx, query,
		activity.Title,
		activity.Description,
		activity.Status,
		activity.DueDate,
		activity.CompletedAt,
		activity.ID,
	)

	return err
}

func (r *ActionPlanActivityRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM action_plan_activities WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// loadMedias carrega todas as mídias associadas a uma atividade
func (r *ActionPlanActivityRepository) loadMedias(ctx context.Context, activityID int64) ([]*domain.ActivityMedia, error) {
	query := `
		SELECT id, activity_id, media_url, media_type, created_at
		FROM activity_media
		WHERE activity_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medias []*domain.ActivityMedia

	for rows.Next() {
		media := &domain.ActivityMedia{}
		err := rows.Scan(
			&media.ID,
			&media.ActivityID,
			&media.MediaURL,
			&media.MediaType,
			&media.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	return medias, nil
}

package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AnswerRepository struct {
	db *sql.DB
}

func NewAnswerRepository(db *sql.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Create(ctx context.Context, answer *domain.Answer) error {
	query := `
		INSERT INTO answers (
			partner_id, submission_id, question_id, value, score,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		answer.PartnerID,
		answer.SubmissionID,
		answer.QuestionID,
		answer.Value,
		answer.Score,
	).Scan(&answer.ID, &answer.CreatedAt, &answer.UpdatedAt)
	return err
}

func (r *AnswerRepository) List(ctx context.Context, tenantID, responseID, limit, offset int64) ([]*domain.Answer, error) {
	query := `
		SELECT 
			a.id, a.partner_id, a.submission_id, a.question_id, 
			a.value, a.score, a.created_at, a.updated_at,
			q.id, q.partner_id, q.questionnaire_id, q.question, 
			q.type, q.options, q.score_values, q.weight, 
			q.required, q.order_num, q.created_at, q.updated_at
		FROM answers a
		INNER JOIN questions q ON a.question_id = q.id AND a.partner_id = q.partner_id
		WHERE a.partner_id = $1 AND a.submission_id = $2
		ORDER BY q.order_num ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, responseID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*domain.Answer

	for rows.Next() {
		var a domain.Answer
		var q domain.Question
		var optionsJSON, scoreValuesJSON []byte

		err := rows.Scan(
			&a.ID, &a.PartnerID, &a.SubmissionID, &a.QuestionID,
			&a.Value, &a.Score, &a.CreatedAt, &a.UpdatedAt,
			&q.ID, &q.PartnerID, &q.TemplateID, &q.Question,
			&q.Type, &optionsJSON, &scoreValuesJSON, &q.Weight,
			&q.Required, &q.OrderNum, &q.CreatedAt, &q.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields from Question
		if len(optionsJSON) > 0 && string(optionsJSON) != "null" {
			if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
				return nil, err
			}
		}
		if len(scoreValuesJSON) > 0 && string(scoreValuesJSON) != "null" {
			if err := json.Unmarshal(scoreValuesJSON, &q.ScoreValues); err != nil {
				return nil, err
			}
		}

		a.Question = &q
		answers = append(answers, &a)
	}
	return answers, nil
}

func (r *AnswerRepository) GetByID(ctx context.Context, tenantID, id int64) (*domain.Answer, error) {
	query := `
		SELECT 
			a.id, a.partner_id, a.submission_id, a.question_id, 
			a.value, a.score, a.created_at, a.updated_at,
			q.id, q.partner_id, q.questionnaire_id, q.question, 
			q.type, q.options, q.score_values, q.weight, 
			q.required, q.order_num, q.created_at, q.updated_at
		FROM answers a
		INNER JOIN questions q ON a.question_id = q.id AND a.partner_id = q.partner_id
		WHERE a.partner_id = $1 AND a.id = $2
	`

	var a domain.Answer
	var q domain.Question
	var optionsJSON, scoreValuesJSON []byte

	err := r.db.QueryRowContext(ctx, query, tenantID, id).Scan(
		&a.ID, &a.PartnerID, &a.SubmissionID, &a.QuestionID,
		&a.Value, &a.Score, &a.CreatedAt, &a.UpdatedAt,
		&q.ID, &q.PartnerID, &q.TemplateID, &q.Question,
		&q.Type, &optionsJSON, &scoreValuesJSON, &q.Weight,
		&q.Required, &q.OrderNum, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields from Question
	if len(optionsJSON) > 0 && string(optionsJSON) != "null" {
		if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
			return nil, err
		}
	}
	if len(scoreValuesJSON) > 0 && string(scoreValuesJSON) != "null" {
		if err := json.Unmarshal(scoreValuesJSON, &q.ScoreValues); err != nil {
			return nil, err
		}
	}

	a.Question = &q
	return &a, nil
}

func (r *AnswerRepository) CountBySubmission(ctx context.Context, tenantID, responseID int64) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM answers 
		WHERE partner_id = $1 AND submission_id = $2
	`

	var count int64
	err := r.db.QueryRowContext(ctx, query, tenantID, responseID).Scan(&count)
	return count, err
}

func (r *AnswerRepository) Update(ctx context.Context, answer *domain.Answer) error {
	query := `
		UPDATE answers
		SET 
			value = $1, 
			score = $2, 
			updated_at = NOW()
		WHERE partner_id = $3 AND id = $4
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		answer.Value,
		answer.Score,
		answer.PartnerID,
		answer.ID,
	).Scan(&answer.UpdatedAt)
	return err
}

func (r *AnswerRepository) Delete(ctx context.Context, partnerID, id int64) error {
	query := `DELETE FROM answers WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, id)
	return err
}

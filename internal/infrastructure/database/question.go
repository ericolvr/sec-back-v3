package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(ctx context.Context, question *domain.Question) error {
	query := `
		INSERT INTO questions (partner_id, template_id, question, type, category, options, score_values, weight, required, order_num, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	now := time.Now()

	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return err
	}

	scoreValuesJSON, err := json.Marshal(question.ScoreValues)
	if err != nil {
		return err
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		question.PartnerID,
		question.TemplateID,
		question.Question,
		question.Type,
		question.Category,
		optionsJSON,
		scoreValuesJSON,
		question.Weight,
		question.Required,
		question.OrderNum,
		now,
		now,
	).Scan(&question.ID)

	if err == nil {
		question.CreatedAt = now
		question.UpdatedAt = now
	}

	return err
}

func (r *QuestionRepository) List(ctx context.Context, tenantID, questionnaireID, limit, offset int64) ([]*domain.Question, error) {
	query := `
		SELECT 
			q.id,
			q.partner_id,
			q.template_id,
			q.question,
			q.type,
			q.category,
			q.options,
			q.score_values,
			q.weight,
			q.required,
			q.order_num,
			q.created_at,
			q.updated_at,
			qn.id,
			qn.name,
			qn.description,
			qn.active,
			qn.created_at
		FROM questions q
		LEFT JOIN assessment_templates qn ON q.template_id = qn.id AND q.partner_id = qn.partner_id
		WHERE q.partner_id = $1 AND q.template_id = $2
		ORDER BY q.order_num ASC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, tenantID, questionnaireID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.Question

	for rows.Next() {
		var q domain.Question
		var qn domain.AssessmentTemplate
		var optionsJSON, scoreValuesJSON []byte

		err := rows.Scan(
			&q.ID,
			&q.PartnerID,
			&q.TemplateID,
			&q.Question,
			&q.Type,
			&q.Category,
			&optionsJSON,
			&scoreValuesJSON,
			&q.Weight,
			&q.Required,
			&q.OrderNum,
			&q.CreatedAt,
			&q.UpdatedAt,
			&qn.ID,
			&qn.Name,
			&qn.Description,
			&qn.Active,
			&qn.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

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

		qn.PartnerID = q.PartnerID
		q.AssessmentTemplate = &qn

		questions = append(questions, &q)
	}
	return questions, nil
}

func (r *QuestionRepository) ListAllByPartner(ctx context.Context, partnerID, limit, offset int64) ([]*domain.Question, error) {
	query := `
		SELECT 
			q.id,
			q.partner_id,
			q.template_id,
			q.question,
			q.type,
			q.category,
			q.options,
			q.score_values,
			q.weight,
			q.required,
			q.order_num,
			q.created_at,
			q.updated_at,
			qn.id,
			qn.name,
			qn.description,
			qn.active,
			qn.created_at
		FROM questions q
		LEFT JOIN assessment_templates qn ON q.template_id = qn.id AND q.partner_id = qn.partner_id
		WHERE q.partner_id = $1
		ORDER BY q.template_id, q.order_num ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.Question

	for rows.Next() {
		var q domain.Question
		var qn domain.AssessmentTemplate
		var optionsJSON, scoreValuesJSON []byte

		err := rows.Scan(
			&q.ID,
			&q.PartnerID,
			&q.TemplateID,
			&q.Question,
			&q.Type,
			&q.Category,
			&optionsJSON,
			&scoreValuesJSON,
			&q.Weight,
			&q.Required,
			&q.OrderNum,
			&q.CreatedAt,
			&q.UpdatedAt,
			&qn.ID,
			&qn.Name,
			&qn.Description,
			&qn.Active,
			&qn.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

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

		qn.PartnerID = q.PartnerID
		q.AssessmentTemplate = &qn

		questions = append(questions, &q)
	}
	return questions, nil
}

func (r *QuestionRepository) GetByID(ctx context.Context, tenantID, id int64) (*domain.Question, error) {
	query := `
		SELECT 
			q.id,
			q.partner_id,
			q.template_id,
			q.question,
			q.type,
			q.category,
			q.options,
			q.score_values,
			q.weight,
			q.required,
			q.order_num,
			q.created_at,
			q.updated_at,
			qn.id,
			qn.name,
			qn.description,
			qn.active,
			qn.created_at
		FROM questions q
		LEFT JOIN assessment_templates qn ON q.template_id = qn.id AND q.partner_id = qn.partner_id
		WHERE q.partner_id = $1 AND q.id = $2`

	var q domain.Question
	var qn domain.AssessmentTemplate
	var optionsJSON, scoreValuesJSON []byte

	err := r.db.QueryRowContext(
		ctx,
		query,
		tenantID,
		id,
	).Scan(
		&q.ID,
		&q.PartnerID,
		&q.TemplateID,
		&q.Question,
		&q.Type,
		&q.Category,
		&optionsJSON,
		&scoreValuesJSON,
		&q.Weight,
		&q.Required,
		&q.OrderNum,
		&q.CreatedAt,
		&q.UpdatedAt,
		&qn.ID,
		&qn.Name,
		&qn.Description,
		&qn.Active,
		&qn.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

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

	qn.PartnerID = q.PartnerID
	q.AssessmentTemplate = &qn

	return &q, nil
}

func (r *QuestionRepository) Update(ctx context.Context, question *domain.Question) error {
	query := `
		UPDATE questions
		SET 
			question = $1, 
			type = $2,
			category = $3,
			options = $4, 
			score_values = $5, 
			weight = $6, 
			required = $7, 
			order_num = $8, 
			updated_at = $9
		WHERE partner_id = $10 AND id = $11
	`

	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return err
	}

	scoreValuesJSON, err := json.Marshal(question.ScoreValues)
	if err != nil {
		return err
	}

	now := time.Now()

	_, err = r.db.ExecContext(
		ctx,
		query,
		question.Question,
		question.Type,
		question.Category,
		optionsJSON,
		scoreValuesJSON,
		question.Weight,
		question.Required,
		question.OrderNum,
		now,
		question.PartnerID,
		question.ID,
	)

	if err == nil {
		question.UpdatedAt = now
	}

	return err
}

func (r *QuestionRepository) Delete(ctx context.Context, tenantID, id int64) error {
	query := `DELETE FROM questions WHERE partner_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, tenantID, id)
	return err
}

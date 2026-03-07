package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type CalculationFormulaRepository struct {
	db *sql.DB
}

func NewCalculationFormulaRepository(db *sql.DB) *CalculationFormulaRepository {
	return &CalculationFormulaRepository{db: db}
}

func (r *CalculationFormulaRepository) Create(ctx context.Context, formula *domain.CalculationFormula) error {
	query := `
		INSERT INTO calculation_formulas (
			partner_id, version, active,
			risk_low_max, risk_medium_max,
			reliability_acceptable_min, reliability_good_min, reliability_excellent_min,
			description, changelog, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		formula.PartnerID,
		formula.Version,
		formula.Active,
		formula.RiskLowMax,
		formula.RiskMediumMax,
		formula.ReliabilityAcceptableMin,
		formula.ReliabilityGoodMin,
		formula.ReliabilityExcellentMin,
		formula.Description,
		formula.Changelog,
		formula.CreatedBy,
	).Scan(&formula.ID, &formula.CreatedAt)
}

func (r *CalculationFormulaRepository) GetActive(ctx context.Context, partnerID int64) (*domain.CalculationFormula, error) {
	query := `
		SELECT id, partner_id, version, active,
			risk_low_max, risk_medium_max,
			reliability_acceptable_min, reliability_good_min, reliability_excellent_min,
			description, changelog, created_by, created_at, activated_at
		FROM calculation_formulas
		WHERE partner_id = $1 AND active = true
		LIMIT 1`

	var formula domain.CalculationFormula
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&formula.ID,
		&formula.PartnerID,
		&formula.Version,
		&formula.Active,
		&formula.RiskLowMax,
		&formula.RiskMediumMax,
		&formula.ReliabilityAcceptableMin,
		&formula.ReliabilityGoodMin,
		&formula.ReliabilityExcellentMin,
		&formula.Description,
		&formula.Changelog,
		&formula.CreatedBy,
		&formula.CreatedAt,
		&formula.ActivatedAt,
	)

	if err == sql.ErrNoRows {
		// Se não tem fórmula ativa, retorna padrão
		return domain.DefaultCalculationFormula(partnerID), nil
	}

	if err != nil {
		return nil, err
	}

	return &formula, nil
}

func (r *CalculationFormulaRepository) GetByVersion(ctx context.Context, partnerID int64, version string) (*domain.CalculationFormula, error) {
	query := `
		SELECT id, partner_id, version, active,
			risk_low_max, risk_medium_max,
			reliability_acceptable_min, reliability_good_min, reliability_excellent_min,
			description, changelog, created_by, created_at, activated_at
		FROM calculation_formulas
		WHERE partner_id = $1 AND version = $2`

	var formula domain.CalculationFormula
	err := r.db.QueryRowContext(ctx, query, partnerID, version).Scan(
		&formula.ID,
		&formula.PartnerID,
		&formula.Version,
		&formula.Active,
		&formula.RiskLowMax,
		&formula.RiskMediumMax,
		&formula.ReliabilityAcceptableMin,
		&formula.ReliabilityGoodMin,
		&formula.ReliabilityExcellentMin,
		&formula.Description,
		&formula.Changelog,
		&formula.CreatedBy,
		&formula.CreatedAt,
		&formula.ActivatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &formula, nil
}

func (r *CalculationFormulaRepository) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.CalculationFormula, error) {
	query := `
		SELECT id, partner_id, version, active,
			risk_low_max, risk_medium_max,
			reliability_acceptable_min, reliability_good_min, reliability_excellent_min,
			description, changelog, created_by, created_at, activated_at
		FROM calculation_formulas
		WHERE partner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, partnerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanFormulas(rows)
}

func (r *CalculationFormulaRepository) Activate(ctx context.Context, partnerID int64, version string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Desativar todas as fórmulas do partner
	_, err = tx.ExecContext(ctx, `
		UPDATE calculation_formulas 
		SET active = false 
		WHERE partner_id = $1`, partnerID)
	if err != nil {
		return err
	}

	// Ativar a versão específica
	_, err = tx.ExecContext(ctx, `
		UPDATE calculation_formulas 
		SET active = true, activated_at = $1
		WHERE partner_id = $2 AND version = $3`,
		time.Now(), partnerID, version)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CalculationFormulaRepository) Delete(ctx context.Context, partnerID int64, version string) error {
	query := `DELETE FROM calculation_formulas WHERE partner_id = $1 AND version = $2`
	_, err := r.db.ExecContext(ctx, query, partnerID, version)
	return err
}

func (r *CalculationFormulaRepository) scanFormulas(rows *sql.Rows) ([]*domain.CalculationFormula, error) {
	var formulas []*domain.CalculationFormula

	for rows.Next() {
		var formula domain.CalculationFormula
		err := rows.Scan(
			&formula.ID,
			&formula.PartnerID,
			&formula.Version,
			&formula.Active,
			&formula.RiskLowMax,
			&formula.RiskMediumMax,
			&formula.ReliabilityAcceptableMin,
			&formula.ReliabilityGoodMin,
			&formula.ReliabilityExcellentMin,
			&formula.Description,
			&formula.Changelog,
			&formula.CreatedBy,
			&formula.CreatedAt,
			&formula.ActivatedAt,
		)

		if err != nil {
			return nil, err
		}

		formulas = append(formulas, &formula)
	}

	return formulas, rows.Err()
}

package services

import (
	"context"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type CalculationFormulaService struct {
	formulaRepo domain.CalculationFormulaRepository
}

func NewCalculationFormulaService(formulaRepo domain.CalculationFormulaRepository) *CalculationFormulaService {
	return &CalculationFormulaService{
		formulaRepo: formulaRepo,
	}
}

func (s *CalculationFormulaService) GetActive(ctx context.Context, partnerID int64) (*domain.CalculationFormula, error) {
	return s.formulaRepo.GetActive(ctx, partnerID)
}

func (s *CalculationFormulaService) UpdateActive(ctx context.Context, partnerID int64, formula *domain.CalculationFormula) (*domain.CalculationFormula, error) {
	// Validar
	if err := formula.Validate(); err != nil {
		return nil, err
	}

	// Buscar fórmula ativa atual
	current, err := s.formulaRepo.GetActive(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	current.RiskLowMax = formula.RiskLowMax
	current.RiskMediumMax = formula.RiskMediumMax
	current.ReliabilityAcceptableMin = formula.ReliabilityAcceptableMin
	current.ReliabilityGoodMin = formula.ReliabilityGoodMin
	current.ReliabilityExcellentMin = formula.ReliabilityExcellentMin
	if formula.Description != "" {
		current.Description = formula.Description
	}

	// Atualizar no banco
	if err := s.formulaRepo.Update(ctx, current); err != nil {
		return nil, err
	}

	return current, nil
}

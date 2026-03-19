package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type RiskMetricsService struct {
	metricsRepo       domain.RiskMetricsRepository
	submissionRepo    domain.EmployeeSubmissionRepository
	answerRepo        domain.AnswerRepository
	questionRepo      domain.QuestionRepository
	employeeRepo      domain.EmployeeRepository
	departmentRepo    domain.DepartmentRepository
	riskCategoryRepo  domain.RiskCategoryRepository
	formulaRepo       domain.CalculationFormulaRepository
	actionPlanService *ActionPlanService
}

func NewRiskMetricsService(
	metricsRepo domain.RiskMetricsRepository,
	submissionRepo domain.EmployeeSubmissionRepository,
	answerRepo domain.AnswerRepository,
	questionRepo domain.QuestionRepository,
	employeeRepo domain.EmployeeRepository,
	departmentRepo domain.DepartmentRepository,
	riskCategoryRepo domain.RiskCategoryRepository,
	formulaRepo domain.CalculationFormulaRepository,
	actionPlanService *ActionPlanService,
) *RiskMetricsService {
	return &RiskMetricsService{
		metricsRepo:       metricsRepo,
		submissionRepo:    submissionRepo,
		answerRepo:        answerRepo,
		questionRepo:      questionRepo,
		employeeRepo:      employeeRepo,
		departmentRepo:    departmentRepo,
		riskCategoryRepo:  riskCategoryRepo,
		formulaRepo:       formulaRepo,
		actionPlanService: actionPlanService,
	}
}

func (s *RiskMetricsService) CalculateAndStore(ctx context.Context, partnerID, companyID, departmentID, templateID int64) (*domain.RiskMetrics, error) {
	if partnerID <= 0 || companyID <= 0 || departmentID <= 0 || templateID <= 0 {
		return nil, errors.New("invalid IDs")
	}

	_, err := s.departmentRepo.GetByID(ctx, partnerID, departmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	employees, err := s.employeeRepo.ListByDepartment(ctx, partnerID, departmentID, MaxEmployeesForMetrics, 0)
	if err != nil {
		return nil, err
	}

	totalEmployees := int64(0)
	for _, emp := range employees {
		if emp.Active {
			totalEmployees++
		}
	}

	fmt.Printf("[DEBUG] Department %d - Total Active Employees: %d\n", departmentID, totalEmployees)

	submissions, err := s.submissionRepo.ListByDepartment(ctx, partnerID, departmentID, MaxSubmissionsForMetrics, 0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUG] Department %d - Total Submissions Found: %d\n", departmentID, len(submissions))

	totalSubmissions := int64(0)
	completedSubmissions := int64(0)
	var completedSubmissionIDs []int64

	for _, sub := range submissions {
		fmt.Printf("[DEBUG] Submission ID %d - TemplateID: %d (expected: %d), Status: %s\n", sub.ID, sub.TemplateID, templateID, sub.Status)
		if sub.TemplateID == templateID {
			totalSubmissions++
			if sub.Status == "completed" {
				completedSubmissions++
				completedSubmissionIDs = append(completedSubmissionIDs, sub.ID)
			}
		}
	}

	fmt.Printf("[DEBUG] Department %d - Completed Submissions: %d / %d employees\n", departmentID, completedSubmissions, totalEmployees)

	responseRate := float64(0)
	if totalEmployees > 0 {
		responseRate = (float64(completedSubmissions) / float64(totalEmployees)) * 100
	}

	fmt.Printf("[DEBUG] Department %d - Response Rate: %.2f%%\n", departmentID, responseRate)

	// Buscar todas as perguntas do template para obter os pesos
	questions, err := s.questionRepo.List(ctx, partnerID, templateID, MaxQuestionsPerTemplate, 0)
	if err != nil {
		return nil, err
	}

	// Criar mapa de pesos para acesso rápido
	questionWeights := make(map[int64]float64)
	for _, q := range questions {
		questionWeights[q.ID] = q.Weight
	}

	// Calcular score médio PONDERADO
	totalScore := 0.0
	totalWeight := 0.0

	for _, subID := range completedSubmissionIDs {
		answers, err := s.answerRepo.List(ctx, partnerID, subID, MaxAnswersForCalculation, 0)
		if err != nil {
			continue
		}

		for _, ans := range answers {
			if ans.Score == nil {
				continue
			}

			// Buscar peso da pergunta
			weight, exists := questionWeights[ans.QuestionID]
			if !exists {
				weight = 1.0 // Peso padrão se não encontrado
			}

			// Calcular score ponderado
			weightedScore := float64(*ans.Score) * weight
			totalScore += weightedScore
			totalWeight += weight
		}
	}

	// Calcular média ponderada
	averageScore := 0.0
	if totalWeight > 0 {
		averageScore = totalScore / totalWeight
	}

	// Buscar fórmula ativa do partner
	formula, err := s.formulaRepo.GetActive(ctx, partnerID)
	if err != nil {
		// Se falhar, usa fórmula padrão
		formula = domain.DefaultCalculationFormula(partnerID)
	}

	// Usar fórmula do banco para calcular risco e confiabilidade
	riskLevel := formula.CalculateRiskLevel(averageScore)
	reliability := formula.CalculateReliability(responseRate)
	canCalculateRisk := responseRate >= formula.ReliabilityAcceptableMin

	metrics := &domain.RiskMetrics{
		PartnerID:            partnerID,
		CompanyID:            companyID,
		DepartmentID:         departmentID,
		TemplateID:           templateID,
		TotalEmployees:       int(totalEmployees),
		TotalSubmissions:     int(totalSubmissions),
		CompletedSubmissions: int(completedSubmissions),
		ResponseRate:         responseRate,
		AverageScore:         averageScore,
		RiskLevel:            riskLevel,
		Reliability:          reliability,
		CanCalculateRisk:     canCalculateRisk,
		CategoryScores:       "{}",
	}

	if err := s.metricsRepo.Upsert(ctx, metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *RiskMetricsService) GetByDepartment(ctx context.Context, partnerID, departmentID, templateID int64) (*domain.RiskMetrics, error) {
	if partnerID <= 0 || departmentID <= 0 || templateID <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.metricsRepo.GetByDepartment(ctx, partnerID, departmentID, templateID)
}

func (s *RiskMetricsService) GetByCompany(ctx context.Context, partnerID, companyID, templateID int64) ([]*domain.RiskMetrics, error) {
	if partnerID <= 0 || companyID <= 0 || templateID <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.metricsRepo.GetByCompany(ctx, partnerID, companyID, templateID)
}

func (s *RiskMetricsService) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.RiskMetrics, error) {
	if partnerID <= 0 {
		return nil, errors.New("invalid partner ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.metricsRepo.List(ctx, partnerID, limit, offset)
}

func (s *RiskMetricsService) Delete(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.metricsRepo.Delete(ctx, partnerID, id)
}

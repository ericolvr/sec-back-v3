package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AnalyticsService struct {
	riskMetricsService     *RiskMetricsService
	departmentRepo         domain.DepartmentRepository
	companyRepo            domain.CompanyRepository
	partnerRepo            domain.PartnerRepository
	templateRepo           domain.AssessmentTemplateRepository
	settingsRepo           domain.PartnerSettingsRepository
	actionPlanRepo         domain.ActionPlanRepository
	actionPlanTemplateRepo domain.ActionPlanTemplateRepository
	riskCategoryRepo       domain.RiskCategoryRepository
	questionRepo           domain.QuestionRepository
	submissionRepo         domain.EmployeeSubmissionRepository
	answerRepo             domain.AnswerRepository
	formulaRepo            domain.CalculationFormulaRepository
	reportRepo             domain.AnalyticsReportRepository
	assignmentRepo         domain.AssessmentAssignmentRepository
}

func NewAnalyticsService(
	riskMetricsService *RiskMetricsService,
	departmentRepo domain.DepartmentRepository,
	companyRepo domain.CompanyRepository,
	partnerRepo domain.PartnerRepository,
	templateRepo domain.AssessmentTemplateRepository,
	settingsRepo domain.PartnerSettingsRepository,
	actionPlanRepo domain.ActionPlanRepository,
	actionPlanTemplateRepo domain.ActionPlanTemplateRepository,
	riskCategoryRepo domain.RiskCategoryRepository,
	questionRepo domain.QuestionRepository,
	submissionRepo domain.EmployeeSubmissionRepository,
	answerRepo domain.AnswerRepository,
	formulaRepo domain.CalculationFormulaRepository,
	reportRepo domain.AnalyticsReportRepository,
	assignmentRepo domain.AssessmentAssignmentRepository,
) *AnalyticsService {
	return &AnalyticsService{
		riskMetricsService:     riskMetricsService,
		departmentRepo:         departmentRepo,
		companyRepo:            companyRepo,
		partnerRepo:            partnerRepo,
		templateRepo:           templateRepo,
		settingsRepo:           settingsRepo,
		actionPlanRepo:         actionPlanRepo,
		actionPlanTemplateRepo: actionPlanTemplateRepo,
		riskCategoryRepo:       riskCategoryRepo,
		questionRepo:           questionRepo,
		submissionRepo:         submissionRepo,
		answerRepo:             answerRepo,
		formulaRepo:            formulaRepo,
		reportRepo:             reportRepo,
		assignmentRepo:         assignmentRepo,
	}
}

// GetDepartmentReport retorna relatório completo de um departamento (com nomes)
func (s *AnalyticsService) GetDepartmentReport(ctx context.Context, partnerID, companyID, departmentID, templateID int64) (*domain.DepartmentAnalytics, error) {
	// 1. Calcular ou buscar métricas
	metrics, err := s.riskMetricsService.CalculateAndStore(ctx, partnerID, companyID, departmentID, templateID)
	if err != nil {
		return nil, err
	}

	// 2. Buscar department (para nome)
	department, err := s.departmentRepo.GetByID(ctx, partnerID, departmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	// 3. Buscar assessment template (para nome)
	template, err := s.templateRepo.GetByID(ctx, partnerID, templateID)
	if err != nil {
		return nil, errors.New("questionnaire not found")
	}

	// 4. Buscar configurações do partner (para critério de fechamento)
	settings, err := s.settingsRepo.GetByPartnerID(ctx, partnerID)
	if err != nil {
		settings = domain.DefaultPartnerSettings(partnerID)
	}

	// 5. Determinar se pode fechar
	canClose := metrics.ResponseRate >= settings.MinResponseRateToClose && metrics.CanCalculateRisk
	canCloseReason := ""
	if !canClose {
		if metrics.ResponseRate < settings.MinResponseRateToClose {
			canCloseReason = "Taxa de resposta abaixo do mínimo configurado"
		} else if !metrics.CanCalculateRisk {
			canCloseReason = "Dados insuficientes para cálculo confiável"
		}
	} else {
		canCloseReason = "Critérios atingidos - pode fechar a coleta"
	}

	// 6. Buscar risk categories (se existirem)
	var riskCategories []*domain.RiskCategorySummary
	// TODO: Implementar busca de risk_categories quando snapshot for criado

	// 7. Contar action plans gerados
	actionPlansCount := 0
	// TODO: Implementar contagem de action plans

	// 8. Montar DTO completo
	analytics := &domain.DepartmentAnalytics{
		DepartmentID:         departmentID,
		DepartmentName:       department.Name,
		TemplateID:           templateID,
		TemplateName:         template.Name,
		TotalEmployees:       int64(metrics.TotalEmployees),
		TotalSubmissions:     int64(metrics.TotalSubmissions),
		CompletedSubmissions: int64(metrics.CompletedSubmissions),
		ResponseRate:         metrics.ResponseRate,
		Reliability:          metrics.Reliability,
		CanCalculateRisk:     metrics.CanCalculateRisk,
		AverageScore:         metrics.AverageScore,
		RiskLevel:            metrics.RiskLevel,
		CanClose:             canClose,
		CanCloseReason:       canCloseReason,
		RiskCategories:       riskCategories,
		ActionPlansCount:     actionPlansCount,
	}

	// 9. Adicionar warning se necessário
	if !metrics.CanCalculateRisk {
		analytics.Warning = "Taxa de resposta muito baixa. Dados podem não ser confiáveis."
	}

	return analytics, nil
}

// GetCompanyReport retorna relatório consolidado de uma empresa
func (s *AnalyticsService) GetCompanyReport(ctx context.Context, partnerID, companyID, templateID int64) (*domain.CompanyAnalytics, error) {
	// 1. Buscar company (para nome)
	company, err := s.companyRepo.GetByID(ctx, partnerID, companyID)
	if err != nil {
		return nil, errors.New("company not found")
	}

	// 2. Buscar template (para nome)
	template, err := s.templateRepo.GetByID(ctx, partnerID, templateID)
	if err != nil {
		return nil, errors.New("questionnaire not found")
	}

	// 3. Buscar todos os departamentos da empresa
	departments, err := s.departmentRepo.ListByCompany(ctx, partnerID, companyID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// 4. Buscar métricas de cada departamento
	var departmentSummaries []*domain.DepartmentSummary
	totalResponseRate := 0.0
	departmentsAtRisk := 0
	settings, _ := s.settingsRepo.GetByPartnerID(ctx, partnerID)
	if settings == nil {
		settings = domain.DefaultPartnerSettings(partnerID)
	}

	for _, dept := range departments {
		metrics, err := s.riskMetricsService.GetByDepartment(ctx, partnerID, dept.ID, templateID)
		if err != nil {
			continue // Departamento sem métricas ainda
		}

		canClose := metrics.ResponseRate >= settings.MinResponseRateToClose && metrics.CanCalculateRisk
		canCloseReason := ""
		if !canClose {
			if metrics.ResponseRate < settings.MinResponseRateToClose {
				canCloseReason = "Taxa de resposta insuficiente"
			} else {
				canCloseReason = "Dados insuficientes"
			}
		}

		summary := &domain.DepartmentSummary{
			DepartmentID:         dept.ID,
			DepartmentName:       dept.Name,
			TotalEmployees:       int64(metrics.TotalEmployees),
			CompletedSubmissions: int64(metrics.CompletedSubmissions),
			ResponseRate:         metrics.ResponseRate,
			Reliability:          metrics.Reliability,
			CanCalculateRisk:     metrics.CanCalculateRisk,
			AverageScore:         metrics.AverageScore,
			RiskLevel:            metrics.RiskLevel,
			CanClose:             canClose,
			CanCloseReason:       canCloseReason,
		}

		departmentSummaries = append(departmentSummaries, summary)
		totalResponseRate += metrics.ResponseRate

		if metrics.RiskLevel == "high" {
			departmentsAtRisk++
		}
	}

	// 5. Calcular médias gerais
	overallResponseRate := 0.0
	if len(departmentSummaries) > 0 {
		overallResponseRate = totalResponseRate / float64(len(departmentSummaries))
	}

	// 6. Determinar risco geral da empresa
	overallRiskLevel := "low"
	if departmentsAtRisk > len(departmentSummaries)/2 {
		overallRiskLevel = "high"
	} else if departmentsAtRisk > 0 {
		overallRiskLevel = "medium"
	}

	return &domain.CompanyAnalytics{
		CompanyID:           companyID,
		CompanyName:         company.Name,
		TemplateID:          templateID,
		TemplateName:        template.Name,
		TotalDepartments:    len(departmentSummaries),
		OverallResponseRate: overallResponseRate,
		OverallRiskLevel:    overallRiskLevel,
		DepartmentsAtRisk:   departmentsAtRisk,
		Departments:         departmentSummaries,
	}, nil
}

// GetPartnerReport retorna relatório consolidado do Partner (todas as empresas)
func (s *AnalyticsService) GetPartnerReport(ctx context.Context, partnerID int64) (*domain.PartnerAnalytics, error) {
	// 1. Buscar partner (para nome)
	partner, err := s.partnerRepo.GetByID(ctx, partnerID)
	if err != nil {
		return nil, errors.New("partner not found")
	}

	// 2. Buscar todas as empresas do partner
	companies, err := s.companyRepo.List(ctx, partnerID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// 3. Buscar métricas de cada empresa
	var companySummaries []*domain.CompanySummary
	totalResponseRate := 0.0
	companiesAtRisk := 0
	totalActiveTemplates := 0

	for _, company := range companies {
		// Buscar departamentos da empresa
		departments, err := s.departmentRepo.ListByCompany(ctx, partnerID, company.ID, 1000, 0)
		if err != nil {
			continue
		}

		// Calcular métricas da empresa
		companyResponseRate := 0.0
		companyRiskLevel := "low"
		departmentsAtRisk := 0
		activeQuestionnaires := 0

		for range departments {
			// Buscar métricas do departamento (assumindo template mais recente)
			// TODO: Melhorar para listar todos os templates ativos
			metrics, err := s.riskMetricsService.List(ctx, partnerID, 1, 0)
			if err == nil && len(metrics) > 0 {
				companyResponseRate += metrics[0].ResponseRate
				if metrics[0].RiskLevel == "high" {
					departmentsAtRisk++
				}
				activeQuestionnaires++
			}
		}

		if len(departments) > 0 {
			companyResponseRate = companyResponseRate / float64(len(departments))
		}

		if departmentsAtRisk > len(departments)/2 {
			companyRiskLevel = "high"
			companiesAtRisk++
		} else if departmentsAtRisk > 0 {
			companyRiskLevel = "medium"
		}

		summary := &domain.CompanySummary{
			CompanyID:         company.ID,
			CompanyName:       company.Name,
			TotalDepartments:  len(departments),
			ResponseRate:      companyResponseRate,
			RiskLevel:         companyRiskLevel,
			DepartmentsAtRisk: departmentsAtRisk,
			ActiveTemplates:   activeQuestionnaires,
		}

		companySummaries = append(companySummaries, summary)
		totalResponseRate += companyResponseRate
		totalActiveTemplates += activeQuestionnaires
	}

	// 4. Calcular médias gerais
	overallResponseRate := 0.0
	if len(companySummaries) > 0 {
		overallResponseRate = totalResponseRate / float64(len(companySummaries))
	}

	return &domain.PartnerAnalytics{
		PartnerID:            partnerID,
		PartnerName:          partner.Name,
		TotalCompanies:       len(companySummaries),
		TotalActiveTemplates: totalActiveTemplates,
		OverallResponseRate:  overallResponseRate,
		CompaniesAtRisk:      companiesAtRisk,
		Companies:            companySummaries,
	}, nil
}

// GetInProgressTemplates retorna todos os templates em andamento de uma empresa
func (s *AnalyticsService) GetInProgressTemplates(ctx context.Context, partnerID, companyID int64) ([]*domain.TemplateInProgress, error) {
	// 1. Buscar assignments ativos
	assignments, err := s.assignmentRepo.List(ctx, partnerID, 1000, 0)
	if err != nil {
		return nil, err
	}

	settings, _ := s.settingsRepo.GetByPartnerID(ctx, partnerID)
	if settings == nil {
		settings = domain.DefaultPartnerSettings(partnerID)
	}

	// 2. Mapear templates ativos (agrupados por template_id)
	templateMap := make(map[int64]*domain.TemplateInProgress)

	for _, assignment := range assignments {
		if !assignment.Active {
			continue // Pular assignments encerrados
		}

		// Buscar template
		template, err := s.templateRepo.GetByID(ctx, partnerID, assignment.TemplateID)
		if err != nil {
			continue
		}

		// Verificar se já temos esse template no mapa
		qip, exists := templateMap[assignment.TemplateID]
		if !exists {
			qip = &domain.TemplateInProgress{
				ID:                        assignment.TemplateID,
				Name:                      template.Name,
				Description:               template.Description,
				Status:                    "active",
				CreatedAt:                 assignment.CreatedAt.Format(time.RFC3339),
				Departments:               []*domain.DepartmentStatus{},
				DepartmentsCompleted:      0,
				DepartmentsInProgress:     0,
				DepartmentsNotStarted:     0,
				TotalEmployees:            0,
				TotalInvitations:          0,
				CompletedResponses:        0,
				DepartmentsWithHighRisk:   0,
				DepartmentsWithMediumRisk: 0,
				DepartmentsWithLowRisk:    0,
				AverageScore:              0,
				LastUpdated:               time.Now().Format(time.RFC3339),
			}
			templateMap[assignment.TemplateID] = qip
		}

		// Para cada departamento do assignment
		for _, deptID := range assignment.DepartmentIDs {
			// Buscar department
			dept, err := s.departmentRepo.GetByID(ctx, partnerID, deptID)
			if err != nil {
				continue
			}

			// Buscar métricas do departamento
			metrics, err := s.riskMetricsService.GetByDepartment(ctx, partnerID, deptID, assignment.TemplateID)
			if err != nil {
				// Métricas não existem - tentar calcular agora
				// Buscar company_id do departamento
				companyID := dept.CompanyID

				// Tentar calcular e armazenar métricas
				metrics, err = s.riskMetricsService.CalculateAndStore(ctx, partnerID, companyID, deptID, assignment.TemplateID)
				if err != nil {
					// Realmente não há dados - departamento não iniciado
					deptStatus := &domain.DepartmentStatus{
						DepartmentID:       deptID,
						DepartmentName:     dept.Name,
						TotalEmployees:     0,
						CompletedResponses: 0,
						PendingResponses:   0,
						ResponseRate:       0,
						CanCalculateRisk:   false,
						Reliability:        "insufficient",
						AverageScore:       0,
						RiskLevel:          "unknown",
						Status:             "not_started",
						IsActive:           true,
						CanClose:           false,
						CanCloseReason:     "Nenhuma resposta ainda",
					}
					qip.Departments = append(qip.Departments, deptStatus)
					qip.DepartmentsNotStarted++
					continue
				}
			}

			// Determinar status do departamento
			status := "not_started"
			if metrics.CompletedSubmissions > 0 {
				if metrics.CompletedSubmissions >= metrics.TotalSubmissions {
					status = "completed"
					qip.DepartmentsCompleted++
				} else {
					status = "in_progress"
					qip.DepartmentsInProgress++
				}
			} else {
				qip.DepartmentsNotStarted++
			}

			// Determinar se pode fechar
			canClose := metrics.ResponseRate >= settings.MinResponseRateToClose && metrics.CanCalculateRisk
			canCloseReason := ""
			if !canClose {
				if metrics.ResponseRate < settings.MinResponseRateToClose {
					canCloseReason = fmt.Sprintf("Taxa de resposta abaixo de %.0f%% (atual: %.0f%%)", settings.MinResponseRateToClose, metrics.ResponseRate)
				} else {
					canCloseReason = "Dados insuficientes para cálculo confiável"
				}
			} else {
				canCloseReason = "Critérios atingidos - pode fechar a coleta"
			}

			// Criar status do departamento
			deptStatus := &domain.DepartmentStatus{
				DepartmentID:       deptID,
				DepartmentName:     dept.Name,
				TotalEmployees:     int64(metrics.TotalEmployees),
				CompletedResponses: int64(metrics.CompletedSubmissions),
				PendingResponses:   int64(metrics.TotalEmployees - metrics.CompletedSubmissions),
				ResponseRate:       metrics.ResponseRate,
				CanCalculateRisk:   metrics.CanCalculateRisk,
				Reliability:        metrics.Reliability,
				AverageScore:       metrics.AverageScore,
				RiskLevel:          metrics.RiskLevel,
				Status:             status,
				IsActive:           true,
				CanClose:           canClose,
				CanCloseReason:     canCloseReason,
			}

			qip.Departments = append(qip.Departments, deptStatus)

			// Acumular métricas gerais
			qip.TotalEmployees += metrics.TotalEmployees
			qip.TotalInvitations += metrics.TotalSubmissions
			qip.CompletedResponses += metrics.CompletedSubmissions
			qip.AverageScore += metrics.AverageScore

			// Contar departamentos por nível de risco
			switch metrics.RiskLevel {
			case "high":
				qip.DepartmentsWithHighRisk++
			case "medium":
				qip.DepartmentsWithMediumRisk++
			case "low":
				qip.DepartmentsWithLowRisk++
			}
		}
	}

	// 3. Calcular médias e determinar risco geral de cada template
	var result []*domain.TemplateInProgress
	for _, qip := range templateMap {
		qip.TotalDepartments = len(qip.Departments)

		if qip.TotalDepartments > 0 {
			qip.AverageScore = qip.AverageScore / float64(qip.TotalDepartments)

			// Determinar risco geral
			if qip.DepartmentsWithHighRisk > qip.TotalDepartments/2 {
				qip.OverallRiskLevel = "high"
			} else if qip.DepartmentsWithHighRisk > 0 || qip.DepartmentsWithMediumRisk > qip.TotalDepartments/2 {
				qip.OverallRiskLevel = "medium"
			} else {
				qip.OverallRiskLevel = "low"
			}
		}
		result = append(result, qip)
	}

	return result, nil
}

// GetAllInProgressTemplates retorna todos os templates em andamento de todas as companies do partner
func (s *AnalyticsService) GetAllInProgressTemplates(ctx context.Context, partnerID int64) ([]*domain.TemplateInProgress, error) {
	// 1. Buscar todas as companies do partner
	companies, err := s.companyRepo.List(ctx, partnerID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// 2. Mapear templates ativos de todas as companies (agrupados por template_id)
	templateMap := make(map[int64]*domain.TemplateInProgress)

	for _, company := range companies {
		// Buscar templates em andamento da company
		templates, err := s.GetInProgressTemplates(ctx, partnerID, company.ID)
		if err != nil {
			continue
		}

		// Agregar templates
		for _, template := range templates {
			existing, exists := templateMap[template.ID]
			if !exists {
				templateMap[template.ID] = template
			} else {
				// Agregar métricas de múltiplas companies
				existing.TotalDepartments += template.TotalDepartments
				existing.DepartmentsCompleted += template.DepartmentsCompleted
				existing.DepartmentsInProgress += template.DepartmentsInProgress
				existing.DepartmentsNotStarted += template.DepartmentsNotStarted
				existing.TotalEmployees += template.TotalEmployees
				existing.TotalInvitations += template.TotalInvitations
				existing.CompletedResponses += template.CompletedResponses
				existing.DepartmentsWithHighRisk += template.DepartmentsWithHighRisk
				existing.DepartmentsWithMediumRisk += template.DepartmentsWithMediumRisk
				existing.DepartmentsWithLowRisk += template.DepartmentsWithLowRisk
				existing.Departments = append(existing.Departments, template.Departments...)
			}
		}
	}

	// 3. Recalcular médias e risco geral
	var result []*domain.TemplateInProgress
	for _, qip := range templateMap {
		if qip.TotalDepartments > 0 {
			// Recalcular average score
			totalScore := 0.0
			for _, dept := range qip.Departments {
				totalScore += dept.AverageScore
			}
			qip.AverageScore = totalScore / float64(len(qip.Departments))

			// Determinar risco geral
			if qip.DepartmentsWithHighRisk > qip.TotalDepartments/2 {
				qip.OverallRiskLevel = "high"
			} else if qip.DepartmentsWithHighRisk > 0 || qip.DepartmentsWithMediumRisk > qip.TotalDepartments/2 {
				qip.OverallRiskLevel = "medium"
			} else {
				qip.OverallRiskLevel = "low"
			}
		}
		result = append(result, qip)
	}

	return result, nil
}

// calculateRiskByCategory calcula risco por categoria (obrigatório NR-1)
func (s *AnalyticsService) calculateRiskByCategory(ctx context.Context, partnerID, companyID, departmentID, templateID int64) ([]*domain.RiskCategory, error) {
	// Buscar fórmula ativa do partner
	formula, err := s.formulaRepo.GetActive(ctx, partnerID)
	if err != nil {
		formula = domain.DefaultCalculationFormula(partnerID)
	}

	// 1. Buscar todas as perguntas do template
	allQuestions, err := s.questionRepo.List(ctx, partnerID, templateID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// 2. Buscar todas as submissions completadas do departamento
	allSubmissions, err := s.submissionRepo.List(ctx, partnerID, 10000, 0)
	if err != nil {
		return nil, err
	}

	var completedSubmissions []*domain.EmployeeSubmission
	for _, sub := range allSubmissions {
		if sub.DepartmentID == departmentID && sub.TemplateID == templateID && sub.Status == "completed" {
			completedSubmissions = append(completedSubmissions, sub)
		}
	}

	if len(completedSubmissions) == 0 {
		return nil, nil
	}

	// 3. Agrupar perguntas por categoria
	categoryQuestions := make(map[string][]*domain.Question)
	for _, q := range allQuestions {
		if q.Category != "" {
			categoryQuestions[q.Category] = append(categoryQuestions[q.Category], q)
		}
	}

	// 4. Calcular risco para cada categoria
	var riskCategories []*domain.RiskCategory

	for category, catQuestions := range categoryQuestions {
		totalScore := 0.0
		totalWeight := 0.0

		// Para cada submission completada
		for _, submission := range completedSubmissions {
			// Buscar respostas
			answers, err := s.answerRepo.List(ctx, partnerID, submission.ID, 10000, 0)
			if err != nil {
				continue
			}

			// Calcular score ponderado para esta categoria
			for _, answer := range answers {
				if answer.Score == nil {
					continue
				}

				// Verificar se a pergunta pertence a esta categoria
				for _, q := range catQuestions {
					if q.ID == answer.QuestionID {
						weightedScore := float64(*answer.Score) * q.Weight
						totalScore += weightedScore
						totalWeight += q.Weight
						break
					}
				}
			}
		}

		// Calcular média ponderada da categoria usando fórmula do banco
		if totalWeight > 0 {
			averageScore := totalScore / totalWeight
			riskLevel := formula.CalculateRiskLevel(averageScore)

			riskCategory := &domain.RiskCategory{
				PartnerID:     partnerID,
				Category:      category,
				AverageScore:  averageScore,
				RiskLevel:     riskLevel,
				QuestionCount: len(catQuestions),
			}

			riskCategories = append(riskCategories, riskCategory)
		}
	}

	return riskCategories, nil
}

// CreateSnapshot congela dados de analytics em um snapshot imutável
func (s *AnalyticsService) CreateSnapshot(ctx context.Context, partnerID, companyID, departmentID, templateID int64, createdBy *int64) (*domain.AnalyticsReport, error) {
	// 1. Buscar fórmula ativa do partner
	formula, err := s.formulaRepo.GetActive(ctx, partnerID)
	if err != nil {
		formula = domain.DefaultCalculationFormula(partnerID)
	}

	// 2. Obter analytics atuais do departamento
	analytics, err := s.GetDepartmentReport(ctx, partnerID, companyID, departmentID, templateID)
	if err != nil {
		return nil, err
	}

	// 3. Criar metadados de cálculo com fórmula usada
	metadata := domain.BuildCalculationMetadata(formula)

	// 4. Criar analytics com metadados
	analyticsWithMetadata := &domain.DepartmentAnalyticsWithMetadata{
		DepartmentAnalytics: *analytics,
		CalculationMetadata: metadata,
	}

	// 5. Criar snapshot
	report := &domain.AnalyticsReport{
		PartnerID:    partnerID,
		DepartmentID: departmentID,
		TemplateID:   templateID,
		CreatedBy:    createdBy,
	}

	// 6. Serializar analytics + metadados para JSON
	if err := report.SetAnalytics(analyticsWithMetadata); err != nil {
		return nil, err
	}

	// 7. Validar e salvar snapshot
	if err := report.Validate(); err != nil {
		return nil, err
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, err
	}

	// 8. Calcular e salvar risk categories vinculadas ao snapshot
	riskCategories, err := s.calculateRiskByCategory(ctx, partnerID, companyID, departmentID, templateID)
	if err == nil && len(riskCategories) > 0 {
		for _, rc := range riskCategories {
			rc.SnapshotID = report.ID
			if err := s.riskCategoryRepo.Create(ctx, rc); err != nil {
				// Log error but don't fail snapshot creation
			}
		}

		// 9. Auto-gerar Action Plans baseados em risk categories
		if err := s.autoGenerateActionPlans(ctx, partnerID, companyID, departmentID, templateID, report.ID, riskCategories); err != nil {
			// Log error but don't fail snapshot creation
		}
	}

	// 10. Fechar questionnaire assignment deste departamento
	if err := s.assignmentRepo.CloseByTemplateAndDepartment(ctx, partnerID, templateID, departmentID); err != nil {
		// Log error but don't fail snapshot creation
		// Assignment might not exist if created before this feature
	}

	return report, nil
}

// ListReports lista todos os snapshots de um partner
func (s *AnalyticsService) ListReports(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	return s.reportRepo.List(ctx, partnerID, limit, offset)
}

// GetReport busca um snapshot específico
func (s *AnalyticsService) GetReport(ctx context.Context, partnerID, reportID int64) (*domain.AnalyticsReport, error) {
	return s.reportRepo.GetByID(ctx, partnerID, reportID)
}

// ListReportsByDepartment lista snapshots de um departamento
func (s *AnalyticsService) ListReportsByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	return s.reportRepo.ListByDepartment(ctx, partnerID, departmentID, limit, offset)
}

// ListReportsByTemplate lista snapshots de um template
func (s *AnalyticsService) ListReportsByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*domain.AnalyticsReport, error) {
	return s.reportRepo.ListByTemplate(ctx, partnerID, templateID, limit, offset)
}

// GetRiskCategoriesBySnapshot busca risk categories de um snapshot
func (s *AnalyticsService) GetRiskCategoriesBySnapshot(ctx context.Context, partnerID, snapshotID int64) ([]*domain.RiskCategory, error) {
	return s.riskCategoryRepo.ListBySnapshot(ctx, partnerID, snapshotID)
}

// autoGenerateActionPlans gera automaticamente action plans baseado em templates e risk categories
func (s *AnalyticsService) autoGenerateActionPlans(ctx context.Context, partnerID, companyID, departmentID, templateID, snapshotID int64, riskCategories []*domain.RiskCategory) error {
	// Buscar department para obter nome
	department, err := s.departmentRepo.GetByID(ctx, partnerID, departmentID)
	if err != nil {
		return err
	}

	// Buscar templates ativos do partner
	templates, err := s.actionPlanTemplateRepo.ListActive(ctx, partnerID)
	if err != nil {
		return err
	}

	// Para cada risk category, verificar se há templates que devem ser acionados
	for _, riskCat := range riskCategories {
		for _, template := range templates {
			// Verificar se template corresponde à categoria e nível de risco
			if template.Category == riskCat.Category && template.ShouldTrigger(riskCat.RiskLevel) {
				// Criar action plan a partir do template
				actionPlan := s.buildActionPlanFromTemplate(template, riskCat, partnerID, companyID, departmentID, templateID, snapshotID, department.Name)

				// Salvar action plan
				if err := s.actionPlanRepo.Create(ctx, actionPlan); err != nil {
					// Log error but continue with other templates
					continue
				}
			}
		}
	}

	return nil
}

// buildActionPlanFromTemplate constrói um ActionPlan a partir de um template
func (s *AnalyticsService) buildActionPlanFromTemplate(
	template *domain.ActionPlanTemplate,
	riskCat *domain.RiskCategory,
	partnerID, companyID, departmentID, templateID, snapshotID int64,
	departmentName string,
) *domain.ActionPlan {
	// Substituir variáveis no título e descrição
	title := s.replaceTemplateVariables(template.TitleTemplate, riskCat, departmentName)
	description := s.replaceTemplateVariables(template.DescriptionTemplate, riskCat, departmentName)

	// Calcular due date
	dueDate := time.Now().AddDate(0, 0, template.DefaultDueDays)

	return &domain.ActionPlan{
		PartnerID:       partnerID,
		CompanyID:       companyID,
		TemplateID:      templateID,
		DepartmentID:    departmentID,
		SnapshotID:      &snapshotID,
		Title:           title,
		Description:     description,
		RiskLevel:       riskCat.RiskLevel,
		Priority:        template.Priority,
		Category:        riskCat.Category,
		ResponsibleName: "A definir",
		Status:          domain.ActionPlanStatusPending,
		DueDate:         &dueDate,
		EvidenceURLs:    []string{},
	}
}

// replaceTemplateVariables substitui variáveis no template
func (s *AnalyticsService) replaceTemplateVariables(template string, riskCat *domain.RiskCategory, departmentName string) string {
	result := template

	// Substituir {category}
	result = strings.ReplaceAll(result, "{category}", riskCat.Category)

	// Substituir {department_name}
	result = strings.ReplaceAll(result, "{department_name}", departmentName)

	// Substituir {average_score}
	result = strings.ReplaceAll(result, "{average_score}", fmt.Sprintf("%.2f", riskCat.AverageScore))

	// Substituir {risk_level}
	result = strings.ReplaceAll(result, "{risk_level}", riskCat.RiskLevel)

	return result
}

// GetClosedDepartmentSnapshots retorna todos os snapshots de departamentos encerrados
func (s *AnalyticsService) GetClosedDepartmentSnapshots(ctx context.Context, partnerID int64) ([]*domain.DepartmentSnapshot, error) {
	// 1. Buscar todos os assignments encerrados (active=false)
	assignments, err := s.assignmentRepo.List(ctx, partnerID, 1000, 0)
	if err != nil {
		return nil, err
	}

	var snapshots []*domain.DepartmentSnapshot

	// 2. Para cada assignment encerrado, buscar métricas
	for _, assignment := range assignments {
		if assignment.Active {
			continue // Pular assignments ainda ativos
		}

		// Para cada departamento do assignment
		for _, deptID := range assignment.DepartmentIDs {
			// Buscar department para nome
			department, err := s.departmentRepo.GetByID(ctx, partnerID, deptID)
			if err != nil {
				continue
			}

			// Buscar template para nome
			template, err := s.templateRepo.GetByID(ctx, partnerID, assignment.TemplateID)
			if err != nil {
				continue
			}

			// Buscar métricas do departamento
			metrics, err := s.riskMetricsService.GetByDepartment(ctx, partnerID, deptID, assignment.TemplateID)
			if err != nil {
				continue // Departamento sem métricas
			}

			// Buscar total de funcionários do departamento
			employees, err := s.submissionRepo.ListByDepartment(ctx, partnerID, deptID, 10000, 0)
			if err != nil {
				continue
			}

			// Filtrar submissions deste template
			var templateSubmissions []*domain.EmployeeSubmission
			completedCount := 0
			for _, sub := range employees {
				if sub.TemplateID == assignment.TemplateID {
					templateSubmissions = append(templateSubmissions, sub)
					if sub.Status == "completed" {
						completedCount++
					}
				}
			}

			// Calcular taxa de resposta
			responseRate := 0.0
			if len(templateSubmissions) > 0 {
				responseRate = (float64(completedCount) / float64(len(templateSubmissions))) * 100
			}

			// Criar snapshot
			snapshot := &domain.DepartmentSnapshot{
				SnapshotID:         assignment.ID, // Usando assignment ID como snapshot ID
				TemplateID:         assignment.TemplateID,
				TemplateName:       template.Name,
				DepartmentID:       deptID,
				DepartmentName:     department.Name,
				ClosedAt:           assignment.ClosedAt,
				TotalEmployees:     int64(len(templateSubmissions)),
				CompletedResponses: int64(completedCount),
				ResponseRate:       responseRate,
				Reliability:        metrics.Reliability,
				AverageScore:       metrics.AverageScore,
				RiskLevel:          metrics.RiskLevel,
				CanCalculateRisk:   metrics.CanCalculateRisk,
			}

			snapshots = append(snapshots, snapshot)
		}
	}

	return snapshots, nil
}

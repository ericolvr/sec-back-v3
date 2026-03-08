package services

import (
	"context"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type DashboardService struct {
	analyticsService    *AnalyticsService
	notificationService *NotificationService
	employeeRepo        domain.EmployeeRepository
	departmentRepo      domain.DepartmentRepository
	companyRepo         domain.CompanyRepository
	partnerRepo         domain.PartnerRepository
}

func NewDashboardService(
	analyticsService *AnalyticsService,
	notificationService *NotificationService,
	employeeRepo domain.EmployeeRepository,
	departmentRepo domain.DepartmentRepository,
	companyRepo domain.CompanyRepository,
	partnerRepo domain.PartnerRepository,
) *DashboardService {
	return &DashboardService{
		analyticsService:    analyticsService,
		notificationService: notificationService,
		employeeRepo:        employeeRepo,
		departmentRepo:      departmentRepo,
		companyRepo:         companyRepo,
		partnerRepo:         partnerRepo,
	}
}

// GetCompanyDashboard retorna dashboard completo para gestor de empresa
func (s *DashboardService) GetCompanyDashboard(ctx context.Context, partnerID, companyID int64) (*domain.CompanyDashboard, error) {
	// 1. Buscar company
	company, err := s.companyRepo.GetByID(ctx, partnerID, companyID)
	if err != nil {
		return nil, err
	}

	// 2. Buscar templates em andamento
	inProgressQuestionnaires, err := s.analyticsService.GetInProgressTemplates(ctx, partnerID, companyID)
	if err != nil {
		inProgressQuestionnaires = []*domain.TemplateInProgress{}
	}

	// 3. Buscar notificações não lidas
	unreadCount, _ := s.notificationService.CountUnread(ctx, partnerID)
	notificationsPreview, _ := s.notificationService.ListUnread(ctx, partnerID, 5, 0)

	// 4. Contar employees e departments
	departments, _ := s.departmentRepo.ListByCompany(ctx, partnerID, companyID, 1000, 0)
	employees, _ := s.employeeRepo.List(ctx, partnerID, 10000, 0)

	totalEmployees := 0
	for _, emp := range employees {
		if emp.CompanyID == companyID && emp.Active {
			totalEmployees++
		}
	}

	// 5. Determinar risco geral e gerar alertas
	overallRiskLevel := "low"
	var alerts []string
	departmentsAtRisk := 0

	for _, qip := range inProgressQuestionnaires {
		for _, dept := range qip.Departments {
			if dept.RiskLevel == "high" {
				departmentsAtRisk++
				alerts = append(alerts, "Departamento "+dept.DepartmentName+" com risco alto")
			}
			if dept.CanClose {
				alerts = append(alerts, "Departamento "+dept.DepartmentName+" pode ser fechado")
			}
		}
	}

	if departmentsAtRisk > len(departments)/2 {
		overallRiskLevel = "high"
	} else if departmentsAtRisk > 0 {
		overallRiskLevel = "medium"
	}

	// 6. Montar dashboard
	return &domain.CompanyDashboard{
		CompanyID:                companyID,
		CompanyName:              company.Name,
		InProgressQuestionnaires: inProgressQuestionnaires,
		UnreadNotifications:      unreadCount,
		NotificationsPreview:     notificationsPreview,
		PendingActionPlans:       0, // TODO: Implementar quando ActionPlanService estiver pronto
		OverdueActionPlans:       0, // TODO: Implementar
		PendingInvitations:       0, // TODO: Implementar quando InvitationService estiver pronto
		TotalEmployees:           totalEmployees,
		TotalDepartments:         len(departments),
		OverallRiskLevel:         overallRiskLevel,
		Alerts:                   alerts,
	}, nil
}

// GetPartnerDashboard retorna dashboard completo para consultoria (Partner)
func (s *DashboardService) GetPartnerDashboard(ctx context.Context, partnerID int64) (*domain.PartnerDashboard, error) {
	// 1. Buscar partner
	partner, err := s.partnerRepo.GetByID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	// 2. Buscar todas as empresas
	companies, err := s.companyRepo.List(ctx, partnerID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// 3. Buscar status de cada empresa
	var companiesSummary []*domain.CompanyStatus
	totalActiveTemplates := 0
	totalResponseRate := 0.0
	companiesAtRisk := 0
	var alerts []string

	for _, company := range companies {
		// Buscar templates em andamento da empresa
		questionnaires, _ := s.analyticsService.GetInProgressTemplates(ctx, partnerID, company.ID)

		// Buscar departamentos
		departments, _ := s.departmentRepo.ListByCompany(ctx, partnerID, company.ID, 1000, 0)

		companyResponseRate := 0.0
		companyRiskLevel := "low"
		departmentsAtRisk := 0

		for _, q := range questionnaires {
			// Calcular response rate a partir dos departamentos
			totalDepts := len(q.Departments)
			if totalDepts > 0 {
				deptResponseSum := 0.0
				for _, dept := range q.Departments {
					deptResponseSum += dept.ResponseRate
					if dept.RiskLevel == "high" {
						departmentsAtRisk++
					}
				}
				companyResponseRate += deptResponseSum / float64(totalDepts)
			}
			totalActiveTemplates++
		}

		if len(questionnaires) > 0 {
			companyResponseRate = companyResponseRate / float64(len(questionnaires))
		}

		if departmentsAtRisk > len(departments)/2 {
			companyRiskLevel = "high"
			companiesAtRisk++
			alerts = append(alerts, company.Name+": múltiplos departamentos em risco")
		} else if departmentsAtRisk > 0 {
			companyRiskLevel = "medium"
		}

		companiesSummary = append(companiesSummary, &domain.CompanyStatus{
			CompanyID:         company.ID,
			CompanyName:       company.Name,
			ActiveTemplates:   len(questionnaires),
			ResponseRate:      companyResponseRate,
			RiskLevel:         companyRiskLevel,
			DepartmentsAtRisk: departmentsAtRisk,
			TotalDepartments:  len(departments),
		})

		totalResponseRate += companyResponseRate
	}

	overallResponseRate := 0.0
	if len(companies) > 0 {
		overallResponseRate = totalResponseRate / float64(len(companies))
	}

	return &domain.PartnerDashboard{
		PartnerID:            partnerID,
		PartnerName:          partner.Name,
		CompaniesSummary:     companiesSummary,
		TotalCompanies:       len(companies),
		TotalActiveTemplates: totalActiveTemplates,
		CompaniesAtRisk:      companiesAtRisk,
		OverallResponseRate:  overallResponseRate,
		Alerts:               alerts,
	}, nil
}

// GetDepartmentDashboard retorna dashboard completo para supervisor de departamento
func (s *DashboardService) GetDepartmentDashboard(ctx context.Context, partnerID, companyID, departmentID int64) (*domain.DepartmentDashboard, error) {
	// 1. Buscar department e company
	department, err := s.departmentRepo.GetByID(ctx, partnerID, departmentID)
	if err != nil {
		return nil, err
	}

	company, err := s.companyRepo.GetByID(ctx, partnerID, companyID)
	if err != nil {
		return nil, err
	}

	// 2. Buscar employees do departamento
	allEmployees, _ := s.employeeRepo.List(ctx, partnerID, 10000, 0)
	var deptEmployees []*domain.Employee
	for _, emp := range allEmployees {
		if emp.DepartmentID == departmentID && emp.Active {
			deptEmployees = append(deptEmployees, emp)
		}
	}

	// 3. Buscar templates ativos do departamento
	// TODO: Melhorar para buscar apenas templates com assignments ativos deste departamento
	var activeQuestionnaires []*domain.DepartmentQuestionnaire
	var alerts []string

	// 4. Montar summary de employees
	employeesSummary := &domain.EmployeesSummary{
		Total:     len(deptEmployees),
		Responded: 0, // TODO: Calcular baseado em submissions
		Pending:   len(deptEmployees),
	}

	return &domain.DepartmentDashboard{
		DepartmentID:     departmentID,
		DepartmentName:   department.Name,
		CompanyID:        companyID,
		CompanyName:      company.Name,
		ActiveTemplates:  activeQuestionnaires,
		EmployeesSummary: employeesSummary,
		ActionPlans:      []*domain.ActionPlanSummary{},   // TODO: Implementar
		RiskCategories:   []*domain.RiskCategorySummary{}, // TODO: Implementar
		Alerts:           alerts,
	}, nil
}

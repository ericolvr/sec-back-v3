package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActionPlanService struct {
	actionPlanRepo   domain.ActionPlanRepository
	templateRepo     domain.ActionPlanTemplateRepository
	riskCategoryRepo domain.RiskCategoryRepository
	departmentRepo   domain.DepartmentRepository
}

func NewActionPlanService(
	actionPlanRepo domain.ActionPlanRepository,
	templateRepo domain.ActionPlanTemplateRepository,
	riskCategoryRepo domain.RiskCategoryRepository,
	departmentRepo domain.DepartmentRepository,
) *ActionPlanService {
	return &ActionPlanService{
		actionPlanRepo:   actionPlanRepo,
		templateRepo:     templateRepo,
		riskCategoryRepo: riskCategoryRepo,
		departmentRepo:   departmentRepo,
	}
}

// AutoGenerateFromRiskCategories auto-gera Action Plans baseado em categorias de risco
// Chamado após criar um snapshot com risk_categories
func (s *ActionPlanService) AutoGenerateFromRiskCategories(
	ctx context.Context,
	partnerID, companyID, departmentID, templateID, snapshotID int64,
) ([]*domain.ActionPlan, error) {

	// Buscar categorias de risco do snapshot
	riskCategories, err := s.riskCategoryRepo.ListBySnapshot(ctx, partnerID, snapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to list risk categories: %w", err)
	}

	// Buscar department info para usar nos templates
	department, err := s.departmentRepo.GetByID(ctx, partnerID, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	var createdPlans []*domain.ActionPlan

	// Para cada categoria de risco
	for _, riskCat := range riskCategories {
		// Buscar templates ativos para esta categoria
		templates, err := s.templateRepo.ListByCategory(ctx, partnerID, riskCat.Category)
		if err != nil {
			continue // Ignora erro e continua
		}

		// Para cada template, verificar se deve disparar
		for _, template := range templates {
			if template.ShouldTrigger(riskCat.RiskLevel) {
				// Criar Action Plan a partir do template
				actionPlan := s.createFromTemplate(
					template,
					riskCat,
					partnerID,
					companyID,
					departmentID,
					templateID,
					snapshotID,
					department.Name,
				)

				// Salvar Action Plan
				if err := s.actionPlanRepo.Create(ctx, actionPlan); err != nil {
					continue // Ignora erro e continua
				}

				createdPlans = append(createdPlans, actionPlan)
			}
		}
	}

	return createdPlans, nil
}

// CRUD methods
func (s *ActionPlanService) Create(ctx context.Context, actionPlan *domain.ActionPlan) error {
	if err := actionPlan.Validate(); err != nil {
		return err
	}
	return s.actionPlanRepo.Create(ctx, actionPlan)
}

func (s *ActionPlanService) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.ActionPlan, error) {
	return s.actionPlanRepo.List(ctx, partnerID, limit, offset)
}

func (s *ActionPlanService) GetByID(ctx context.Context, partnerID, id int64) (*domain.ActionPlan, error) {
	return s.actionPlanRepo.GetByID(ctx, partnerID, id)
}

func (s *ActionPlanService) Update(ctx context.Context, actionPlan *domain.ActionPlan) error {
	if err := actionPlan.Validate(); err != nil {
		return err
	}
	return s.actionPlanRepo.Update(ctx, actionPlan)
}

func (s *ActionPlanService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.actionPlanRepo.Delete(ctx, partnerID, id)
}

func (s *ActionPlanService) ListByDepartment(ctx context.Context, partnerID, departmentID int64, limit, offset int64) ([]*domain.ActionPlan, error) {
	return s.actionPlanRepo.ListByDepartment(ctx, partnerID, departmentID, limit, offset)
}

func (s *ActionPlanService) ListByStatus(ctx context.Context, partnerID int64, status string, limit, offset int64) ([]*domain.ActionPlan, error) {
	return s.actionPlanRepo.ListByStatus(ctx, partnerID, status, limit, offset)
}

// createFromTemplate cria um Action Plan a partir de um template
func (s *ActionPlanService) createFromTemplate(
	template *domain.ActionPlanTemplate,
	riskCat *domain.RiskCategory,
	partnerID, companyID, departmentID, templateID, snapshotID int64,
	departmentName string,
) *domain.ActionPlan {

	// Substituir variáveis no template
	title := s.replaceTemplateVars(template.TitleTemplate, riskCat, departmentName)
	description := s.replaceTemplateVars(template.DescriptionTemplate, riskCat, departmentName)

	// Calcular due date
	dueDate := time.Now().AddDate(0, 0, template.DefaultDueDays)

	return &domain.ActionPlan{
		PartnerID:       partnerID,
		CompanyID:       companyID,
		TemplateID: templateID,
		DepartmentID:    departmentID,
		SnapshotID:      &snapshotID,
		Title:           title,
		Description:     description,
		RiskLevel:       riskCat.RiskLevel,
		Priority:        template.Priority,
		Category:        riskCat.Category,
		ResponsibleName: "Sistema", // Pode ser configurado depois
		Status:          domain.ActionPlanStatusPending,
		DueDate:         &dueDate,
		Notes:           fmt.Sprintf("Auto-gerado a partir do template: %s", template.TitleTemplate),
	}
}

// replaceTemplateVars substitui variáveis no template
func (s *ActionPlanService) replaceTemplateVars(
	template string,
	riskCat *domain.RiskCategory,
	departmentName string,
) string {
	result := template

	// Substituir variáveis
	result = strings.ReplaceAll(result, "{category}", riskCat.Category)
	result = strings.ReplaceAll(result, "{department_name}", departmentName)
	result = strings.ReplaceAll(result, "{average_score}", fmt.Sprintf("%.2f", riskCat.AverageScore))
	result = strings.ReplaceAll(result, "{risk_level}", riskCat.RiskLevel)
	result = strings.ReplaceAll(result, "{question_count}", fmt.Sprintf("%d", riskCat.QuestionCount))

	return result
}

// MarkAsCompleted marca um Action Plan como completado
func (s *ActionPlanService) MarkAsCompleted(ctx context.Context, partnerID, id int64) error {
	actionPlan, err := s.actionPlanRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	actionPlan.MarkAsCompleted()
	return s.actionPlanRepo.Update(ctx, actionPlan)
}

// MarkAsInProgress marca um Action Plan como em progresso
func (s *ActionPlanService) MarkAsInProgress(ctx context.Context, partnerID, id int64) error {
	actionPlan, err := s.actionPlanRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	actionPlan.MarkAsInProgress()
	return s.actionPlanRepo.Update(ctx, actionPlan)
}

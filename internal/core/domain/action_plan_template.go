package domain

import (
	"context"
	"fmt"
	"time"
)

type ActionPlanTemplate struct {
	ID                  int64     `json:"id"`
	PartnerID            int64     `json:"partner_id"`
	Category            string    `json:"category"`
	MinRiskLevel        string    `json:"min_risk_level"`
	TitleTemplate       string    `json:"title_template"`
	DescriptionTemplate string    `json:"description_template"`
	Priority            string    `json:"priority"`
	DefaultDueDays      int       `json:"default_due_days"`
	AutoCreate          bool      `json:"auto_create"`
	Active              bool      `json:"active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type ActionPlanTemplateRepository interface {
	Create(ctx context.Context, template *ActionPlanTemplate) error
	GetByID(ctx context.Context, tenantID, id int64) (*ActionPlanTemplate, error)
	List(ctx context.Context, tenantID int64, limit, offset int64) ([]*ActionPlanTemplate, error)
	ListActive(ctx context.Context, tenantID int64) ([]*ActionPlanTemplate, error)
	ListByCategory(ctx context.Context, tenantID int64, category string) ([]*ActionPlanTemplate, error)
	Update(ctx context.Context, template *ActionPlanTemplate) error
	Delete(ctx context.Context, tenantID, id int64) error
}

func (t *ActionPlanTemplate) Validate() error {
	if t.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}

	if t.Category == "" {
		return fmt.Errorf("category is required")
	}

	if t.MinRiskLevel == "" {
		return fmt.Errorf("min_risk_level is required")
	}

	// Validate min_risk_level values
	validRiskLevels := []string{"low", "medium", "high"}
	isValid := false
	for _, level := range validRiskLevels {
		if t.MinRiskLevel == level {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("min_risk_level must be one of: low, medium, high")
	}

	if t.TitleTemplate == "" {
		return fmt.Errorf("title_template is required")
	}

	if t.DescriptionTemplate == "" {
		return fmt.Errorf("description_template is required")
	}

	if t.Priority == "" {
		return fmt.Errorf("priority is required")
	}

	// Validate priority values
	validPriorities := []string{
		"low",
		"medium",
		"high",
		"critical",
	}
	isValidPriority := false
	for _, p := range validPriorities {
		if t.Priority == p {
			isValidPriority = true
			break
		}
	}
	if !isValidPriority {
		return fmt.Errorf("invalid priority: %s", t.Priority)
	}

	if t.DefaultDueDays <= 0 {
		return fmt.Errorf("default_due_days must be greater than 0")
	}

	return nil
}

// ShouldTrigger checks if this template should trigger for a given risk level
func (t *ActionPlanTemplate) ShouldTrigger(riskLevel string) bool {
	if !t.Active || !t.AutoCreate {
		return false
	}

	// Risk level hierarchy: low < medium < high
	riskLevels := map[string]int{
		"low":    1,
		"medium": 2,
		"high":   3,
	}

	minLevel, minExists := riskLevels[t.MinRiskLevel]
	currentLevel, currentExists := riskLevels[riskLevel]

	if !minExists || !currentExists {
		return false
	}

	return currentLevel >= minLevel
}

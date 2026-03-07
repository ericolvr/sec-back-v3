package domain

import (
	"context"
	"fmt"
	"time"
)

// AssessmentTemplate representa um questionário de avaliação de riscos NR-1
// Anteriormente chamado de "Questionnaire"
type AssessmentTemplate struct {
	ID          int64     `json:"id"`
	PartnerID   int64     `json:"partner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     int       `json:"version"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AssessmentTemplateRepository interface {
	Create(ctx context.Context, template *AssessmentTemplate) error
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*AssessmentTemplate, error)
	GetByID(ctx context.Context, partnerID, id int64) (*AssessmentTemplate, error)
	Update(ctx context.Context, template *AssessmentTemplate) error
	Delete(ctx context.Context, partnerID, id int64) error
	IncrementVersion(ctx context.Context, partnerID, id int64) error
}

func (t *AssessmentTemplate) Validate() error {
	if t.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if t.Name == "" {
		return fmt.Errorf("template name is required")
	}
	return nil
}

func ValidateAssessmentTemplate(a *AssessmentTemplate) error {
	return a.Validate()
}

// AssessmentVersion registra mudanças em um template (auditoria/versionamento)
type AssessmentVersion struct {
	ID         int64     `json:"id"`
	TemplateID int64     `json:"template_id"`
	PartnerID  int64     `json:"partner_id"`
	Version    int       `json:"version"`
	Changes    string    `json:"changes"` // JSON descrevendo o que mudou
	CreatedBy  int64     `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type AssessmentVersionRepository interface {
	Create(ctx context.Context, version *AssessmentVersion) error
	ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*AssessmentVersion, error)
	GetByID(ctx context.Context, partnerID, id int64) (*AssessmentVersion, error)
}

func (v *AssessmentVersion) Validate() error {
	if v.TemplateID <= 0 {
		return fmt.Errorf("template_id is required")
	}
	if v.PartnerID <= 0 {
		return fmt.Errorf("partner_id is required")
	}
	if v.Version <= 0 {
		return fmt.Errorf("version must be > 0")
	}
	return nil
}

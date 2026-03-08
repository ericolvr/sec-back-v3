package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type NotificationService struct {
	notificationRepo domain.NotificationRepository
	settingsRepo     domain.PartnerSettingsRepository
}

func NewNotificationService(
	notificationRepo domain.NotificationRepository,
	settingsRepo domain.PartnerSettingsRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		settingsRepo:     settingsRepo,
	}
}

// CreateThresholdReached cria notificação quando threshold é atingido
func (s *NotificationService) CreateThresholdReached(ctx context.Context, partnerID, companyID, departmentID, templateID int64, responseRate float64, departmentName string) error {
	settings, err := s.settingsRepo.GetByPartnerID(ctx, partnerID)
	if err != nil || !settings.NotifyOnThreshold {
		return nil // Não notifica se configuração desabilitada
	}

	metadata := map[string]interface{}{
		"response_rate": responseRate,
		"threshold":     settings.MinResponseRateToClose,
	}
	metadataJSON, _ := json.Marshal(metadata)

	notification := &domain.Notification{
		PartnerID:       partnerID,
		CompanyID:       &companyID,
		DepartmentID:    &departmentID,
		TemplateID: &templateID,
		Type:            domain.NotificationTypeThresholdReached,
		Title:           fmt.Sprintf("Departamento %s atingiu %.0f%% de resposta", departmentName, responseRate),
		Message:         fmt.Sprintf("O departamento %s atingiu %.1f%% de taxa de resposta. Você pode fechar a coleta agora.", departmentName, responseRate),
		Severity:        domain.NotificationSeverityInfo,
		Metadata:        string(metadataJSON),
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	return s.notificationRepo.Create(ctx, notification)
}

// CreateCanClose cria notificação quando departamento pode ser fechado
func (s *NotificationService) CreateCanClose(ctx context.Context, partnerID, companyID, departmentID, templateID int64, departmentName string, responseRate float64) error {
	metadata := map[string]interface{}{
		"response_rate": responseRate,
		"can_close":     true,
	}
	metadataJSON, _ := json.Marshal(metadata)

	notification := &domain.Notification{
		PartnerID:       partnerID,
		CompanyID:       &companyID,
		DepartmentID:    &departmentID,
		TemplateID: &templateID,
		Type:            domain.NotificationTypeCanClose,
		Title:           fmt.Sprintf("Departamento %s pronto para encerramento", departmentName),
		Message:         fmt.Sprintf("O departamento %s atingiu os critérios necessários (%.1f%% de resposta). Você pode fechar a coleta e gerar o relatório.", departmentName, responseRate),
		Severity:        domain.NotificationSeverityInfo,
		Metadata:        string(metadataJSON),
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	return s.notificationRepo.Create(ctx, notification)
}

// CreateRiskDetected cria notificação quando risco alto é detectado
func (s *NotificationService) CreateRiskDetected(ctx context.Context, partnerID, companyID, departmentID, templateID int64, departmentName, category, riskLevel string, averageScore float64) error {
	severity := domain.NotificationSeverityWarning
	if riskLevel == "high" {
		severity = domain.NotificationSeverityCritical
	}

	metadata := map[string]interface{}{
		"category":      category,
		"risk_level":    riskLevel,
		"average_score": averageScore,
	}
	metadataJSON, _ := json.Marshal(metadata)

	notification := &domain.Notification{
		PartnerID:       partnerID,
		CompanyID:       &companyID,
		DepartmentID:    &departmentID,
		TemplateID: &templateID,
		Type:            domain.NotificationTypeRiskDetected,
		Title:           fmt.Sprintf("⚠️ Risco %s detectado em %s - %s", riskLevel, category, departmentName),
		Message:         fmt.Sprintf("Foi detectado risco %s na categoria %s do departamento %s (score médio: %.2f). Ações corretivas são recomendadas.", riskLevel, category, departmentName, averageScore),
		Severity:        severity,
		Metadata:        string(metadataJSON),
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	return s.notificationRepo.Create(ctx, notification)
}

// CreateSnapshotCreated cria notificação quando snapshot é gerado
func (s *NotificationService) CreateSnapshotCreated(ctx context.Context, partnerID, companyID, departmentID, templateID, snapshotID int64, departmentName string) error {
	metadata := map[string]interface{}{
		"snapshot_id": snapshotID,
	}
	metadataJSON, _ := json.Marshal(metadata)

	notification := &domain.Notification{
		PartnerID:       partnerID,
		CompanyID:       &companyID,
		DepartmentID:    &departmentID,
		TemplateID: &templateID,
		Type:            domain.NotificationTypeSnapshotCreated,
		Title:           fmt.Sprintf("✅ Snapshot gerado para %s", departmentName),
		Message:         fmt.Sprintf("O snapshot do departamento %s foi criado com sucesso. Os dados foram congelados e o relatório está disponível.", departmentName),
		Severity:        domain.NotificationSeverityInfo,
		Metadata:        string(metadataJSON),
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	return s.notificationRepo.Create(ctx, notification)
}

// CreateActionPlanGenerated cria notificação quando action plan é auto-gerado
func (s *NotificationService) CreateActionPlanGenerated(ctx context.Context, partnerID, companyID, departmentID, templateID, actionPlanID int64, title string) error {
	metadata := map[string]interface{}{
		"action_plan_id": actionPlanID,
	}
	metadataJSON, _ := json.Marshal(metadata)

	notification := &domain.Notification{
		PartnerID:       partnerID,
		CompanyID:       &companyID,
		DepartmentID:    &departmentID,
		TemplateID: &templateID,
		Type:            domain.NotificationTypeActionPlanCreated,
		Title:           "Plano de ação criado automaticamente",
		Message:         fmt.Sprintf("Um novo plano de ação foi gerado automaticamente: %s", title),
		Severity:        domain.NotificationSeverityInfo,
		Metadata:        string(metadataJSON),
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	return s.notificationRepo.Create(ctx, notification)
}

// List retorna notificações paginadas
func (s *NotificationService) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.List(ctx, partnerID, limit, offset)
}

// ListUnread retorna apenas notificações não lidas
func (s *NotificationService) ListUnread(ctx context.Context, partnerID, limit, offset int64) ([]*domain.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.ListUnread(ctx, partnerID, limit, offset)
}

// MarkAsRead marca notificação como lida
func (s *NotificationService) MarkAsRead(ctx context.Context, partnerID, id, userID int64) error {
	return s.notificationRepo.MarkAsRead(ctx, partnerID, id, userID)
}

// MarkAllAsRead marca todas como lidas
func (s *NotificationService) MarkAllAsRead(ctx context.Context, partnerID, userID int64) error {
	return s.notificationRepo.MarkAllAsRead(ctx, partnerID, userID)
}

// CountUnread retorna quantidade de notificações não lidas
func (s *NotificationService) CountUnread(ctx context.Context, partnerID int64) (int64, error) {
	return s.notificationRepo.CountUnread(ctx, partnerID)
}

// Delete remove notificação
func (s *NotificationService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.notificationRepo.Delete(ctx, partnerID, id)
}

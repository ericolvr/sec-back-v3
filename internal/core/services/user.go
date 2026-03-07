package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type UserService struct {
	userRepo   domain.UserRepository
	smsService domain.SMSService
}

func NewUserService(repo domain.UserRepository, smsService domain.SMSService) *UserService {
	return &UserService{
		userRepo:   repo,
		smsService: smsService,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) CreateWithPassword(ctx context.Context, user *domain.User, plainPassword string) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Print da senha gerada para visualização
	fmt.Printf("\n🔑 ========================================\n")
	fmt.Printf("   SENHA GERADA PARA USUÁRIO\n")
	fmt.Printf("========================================\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Nome: %s\n", user.Name)
	fmt.Printf("   Mobile: %s\n", user.Mobile)
	fmt.Printf("   Senha: %s\n", plainPassword)
	fmt.Printf("========================================\n\n")

	// Temporariamente comentado - envio de SMS
	/*
		if s.smsService != nil && user.Mobile != "" {
			fmt.Printf("[SMS_SECURITY] Enviando senha via SMS para usuário ID:%d, Mobile:%s\n", user.ID, user.Mobile)
			msg := domain.SMSMessage{
				To:      user.Mobile,
				Message: fmt.Sprintf("Olá %s, sua senha é: %s", user.Name, plainPassword),
			}
			if err := s.smsService.SendSMS(msg); err != nil {
				fmt.Printf("⚠️  [SMS_ERROR] Falha ao enviar SMS para %s: %v\n", user.Mobile, err)
			} else {
				fmt.Printf("✅ [SMS_SUCCESS] SMS enviado com sucesso para %s\n", user.Mobile)
			}
		} else {
			fmt.Printf("[SMS_SKIP] Usuário ID:%d sem número de celular ou SMS service não configurado\n", user.ID)
		}
	*/

	return nil
}

func (s *UserService) List(ctx context.Context, partnerID int64, limit, offset int) ([]*domain.User, error) {
	if limit <= 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepo.List(ctx, partnerID, int64(limit), int64(offset))
}

func (s *UserService) GetByID(ctx context.Context, partnerID, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("ID is required")
	}

	return s.userRepo.GetByID(ctx, partnerID, id)
}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	if user.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if err := user.Validate(); err != nil {
		return err
	}

	_, err := s.userRepo.GetByID(ctx, user.PartnerID, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Update(ctx, user.PartnerID, user)
}

func (s *UserService) UpdatePasswordWithSMS(ctx context.Context, user *domain.User, plainPassword string) error {
	if user.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if err := s.userRepo.Update(ctx, user.PartnerID, user); err != nil {
		return err
	}

	// Print da nova senha gerada
	fmt.Printf("\n🔄 ========================================\n")
	fmt.Printf("   SENHA RESETADA PARA USUÁRIO\n")
	fmt.Printf("========================================\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Nome: %s\n", user.Name)
	fmt.Printf("   Mobile: %s\n", user.Mobile)
	fmt.Printf("   Nova Senha: %s\n", plainPassword)
	fmt.Printf("========================================\n\n")

	// Temporariamente comentado - envio de SMS
	/*
		if s.smsService != nil && user.Mobile != "" {
			fmt.Printf("[SMS_SECURITY] Reset de senha - Enviando SMS para usuário ID:%d, Mobile:%s\n", user.ID, user.Mobile)
			msg := domain.SMSMessage{
				To:      user.Mobile,
				Message: fmt.Sprintf("Olá %s, sua nova senha é: %s", user.Name, plainPassword),
			}
			if err := s.smsService.SendSMS(msg); err != nil {
				fmt.Printf("⚠️  [SMS_ERROR] Falha ao enviar SMS para %s: %v\n", user.Mobile, err)
			} else {
				fmt.Printf("✅ [SMS_SUCCESS] SMS de reset enviado com sucesso para %s\n", user.Mobile)
			}
		} else {
			fmt.Printf("[SMS_SKIP] Usuário ID:%d sem número de celular ou SMS service não configurado\n", user.ID)
		}
	*/

	return nil
}

func (s *UserService) Delete(ctx context.Context, partnerID, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("id is required")
	}

	user, err := s.userRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.userRepo.Delete(ctx, partnerID, id); err != nil {
		return nil, err
	}
	return user, nil
}

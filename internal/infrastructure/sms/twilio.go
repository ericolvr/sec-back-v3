package sms

import (
	"fmt"
	"os"
	"strings"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioProvider struct {
	client *twilio.RestClient
	from   string
}

func NewTwilioProvider() *TwilioProvider {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM")

	if from == "" {
		from = "+15017122661"
	}

	if accountSID == "" || authToken == "" {
		fmt.Printf("[TWILIO] ⚠️  Credenciais não configuradas - SMS será desabilitado\n")
	}

	client := twilio.NewRestClient()

	return &TwilioProvider{
		client: client,
		from:   from,
	}
}

func (t *TwilioProvider) SendSMS(message domain.SMSMessage) error {
	if t.client == nil {
		return fmt.Errorf("Twilio client not initialized")
	}

	params := &api.CreateMessageParams{}
	params.SetBody(message.Message)
	params.SetFrom(t.from)
	params.SetTo(message.To)

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("erro ao enviar SMS via Twilio: %w", err)
	}

	fmt.Printf("✅ SMS enviado para %s (SID: %s)\n", message.To, *resp.Sid)
	return nil
}

func (t *TwilioProvider) FormatPhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	if strings.HasPrefix(phone, "+55") {
		return phone
	}

	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}

	if strings.HasPrefix(phone, "55") && len(phone) >= 12 {
		return "+" + phone
	}

	if len(phone) == 10 || len(phone) == 11 {
		return "+55" + phone
	}

	return "+55" + phone
}

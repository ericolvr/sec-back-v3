package domain

type SMSMessage struct {
	To      string
	Message string
}

type SMSService interface {
	SendSMS(message SMSMessage) error
}

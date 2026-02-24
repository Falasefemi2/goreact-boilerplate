package service

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client    *resend.Client
	fromEmail string
}

func NewEmailService(apiKey, fromEmail string) *EmailService {
	return &EmailService{
		client:    resend.NewClient(apiKey),
		fromEmail: fromEmail,
	}
}

func (s *EmailService) SendWelcome(toEmail string) error {
	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{toEmail},
		Subject: "Welcome aboard!",
		Html: fmt.Sprintf(`
			<h1>Welcome!</h1>
			<p>Thanks for signing up with email: %s</p>
			<p>You're all set to get started.</p>
		`, toEmail),
	}

	_, err := s.client.Emails.Send(params)
	return err
}

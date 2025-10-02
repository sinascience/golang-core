package service

import (
	"log"
	"strings"
	"venturo-core/internal/adapter/email"
	"venturo-core/internal/adapter/rabbitmq"
)

type EmailService struct {
	rabbit *rabbitmq.RabbitMQAdapter
	sender email.EmailSenderAdapter
}

func NewEmailService(rabbit *rabbitmq.RabbitMQAdapter, sender email.EmailSenderAdapter) *EmailService {
	return &EmailService{rabbit: rabbit, sender: sender}
}

func (s *EmailService) SendAsyncEmail(to, subject, body string) error {
	message := to + "|" + subject + "|" + body
	return s.rabbit.PublishMessage(message)
}

func (s *EmailService) StartConsumer() {
	s.rabbit.ConsumeMessages(func(msg string) {
		parts := strings.SplitN(msg, "|", 3)
		if len(parts) < 3 {
			log.Printf("Failed to parse message: not enough parts, msg=%s", msg)
			return
		}

		to := parts[0]
		subject := parts[1]
		body := parts[2]

		// Send email
		err := s.sender.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	})
}

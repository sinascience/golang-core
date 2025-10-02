package email

import (
	"fmt"
)

type EmailSenderAdapter interface {
	SendEmail(to, subject, body string) error
}

type ConsoleEmailAdapter struct{}

func NewConsoleEmailAdapter() *ConsoleEmailAdapter {
	return &ConsoleEmailAdapter{}
}

func (c *ConsoleEmailAdapter) SendEmail(to, subject, body string) error {
	fmt.Printf("Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}

package service

import (
	"errors"

	"notification-service/internal/domain/models"
)

type EmailSender interface {
	Send(to, subject, body string) error
}

type NotificationService struct {
	emailSender EmailSender
}

func NewNotificationService(emailSender EmailSender) *NotificationService {
	return &NotificationService{
		emailSender: emailSender,
	}
}

func (s *NotificationService) SendEmailNotification(n models.EmailNotification) error {
	if n.Email == "" {
		return errors.New("email is empty")
	}
	if n.Subject == "" {
		return errors.New("subject is empty")
	}
	if n.Message == "" {
		return errors.New("message is empty")
	}

	return s.emailSender.Send(n.Email, n.Subject, n.Message)
}

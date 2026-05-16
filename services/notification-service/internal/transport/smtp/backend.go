package smtp

import (
	gosmtp "github.com/emersion/go-smtp"
)

type Backend struct {
	notificationService NotificationService
}

func NewBackend(notificationService NotificationService) *Backend {
	return &Backend{
		notificationService: notificationService,
	}
}

func (b *Backend) NewSession(_ *gosmtp.Conn) (gosmtp.Session, error) {
	return NewSession(b.notificationService), nil
}

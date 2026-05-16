package smtp

import (
	"io"

	gosmtp "github.com/emersion/go-smtp"
)

type NotificationService interface {
	SendEmail(from string, to []string, data []byte) error
}

type Session struct {
	notificationService NotificationService

	from string
	to   []string
}

func NewSession(notificationService NotificationService) *Session {
	return &Session{
		notificationService: notificationService,
	}
}

func (s *Session) Mail(from string, _ *gosmtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, _ *gosmtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return s.notificationService.SendEmail(s.from, s.to, data)
}

func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

func (s *Session) Logout() error {
	return nil
}

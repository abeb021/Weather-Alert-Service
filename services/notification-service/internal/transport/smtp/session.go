package smtp

import (
	"errors"
	"io"
	"log/slog"

	gosmtp "github.com/emersion/go-smtp"
)

var (
	ErrReadData = errors.New("failed to read data of message")
	ErrAuth     = errors.New("invalid username or password")
)

type NotificationService interface {
	SendEmail(from string, to []string, data []byte) error
}

type Session struct {
	logger              *slog.Logger
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
	s.logger.Info("Mail from:", from)

	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, _ *gosmtp.RcptOptions) error {
	s.logger.Info("Rcpt to:", to)

	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		s.logger.Error(ErrReadData.Error(), err)
		return err
	}

	s.logger.Info("Received message:", string(data))

	return s.notificationService.SendEmail(s.from, s.to, data)
}

func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "testuser" || password != "testpass" {
		s.logger.Error(ErrAuth.Error(), username, password)
		return ErrAuth
	}

	return nil
}

func (s *Session) Logout() error {
	return nil
}

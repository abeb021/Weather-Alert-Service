package smtp

import (
	"notification-service/config"

	gosmtp "github.com/emersion/go-smtp"
)

type Server struct {
	server *gosmtp.Server
}

func NewServer(cfg *config.Config, backend gosmtp.Backend) *Server {
	s := gosmtp.NewServer(backend)

	s.Addr = cfg.SMTP.Address
	s.Domain = cfg.SMTP.Domain
	s.WriteTimeout = cfg.SMTP.WriteTimeout
	s.ReadTimeout = cfg.SMTP.ReadTimeout
	s.MaxMessageBytes = cfg.SMTP.MaxMessageBytes
	s.MaxRecipients = cfg.SMTP.MaxRecipients
	s.AllowInsecureAuth = cfg.SMTP.AllowInsecureAuth

	return &Server{
		server: s,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

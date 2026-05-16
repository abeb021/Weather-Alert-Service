package smtp

import (
	"net/smtp"
	"notification-service/config"
)

type Client struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewClient(config *config.Config) *Client {
	return &Client{
		host:     config.Client.Host,
		port:     config.Client.Port,
		username: config.Client.Username,
		password: config.Client.Password,
		from:     config.Client.From,
	}
}

func (c *Client) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	msg := []byte(
		"From: " + c.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	addr := c.host + ":" + c.port
	return smtp.SendMail(addr, auth, c.from, []string{to}, msg)
}

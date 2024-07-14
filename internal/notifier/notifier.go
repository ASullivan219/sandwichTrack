package notifier

import (
	"log/slog"
	"net/smtp"
)

type I_Notifier interface {
	NotifyAll(string)
	NotifyOne(string, string) error
}

type emailNotifier struct {
	FromEmail string
	password  string
	smtpHost  string
	smtpPort  string
	I_Notifier
}

func NewEmailNotifier(from string, password string, host string, port string) emailNotifier {
	return emailNotifier{
		FromEmail: from,
		password:  password,
		smtpHost:  host,
		smtpPort:  port,
	}
}

func (e *emailNotifier) NotifyAll(message string) {
	e.NotifyOne(message, "alexander.sullivan219@gmail.com")
}

func (e *emailNotifier) NotifyOne(message string, email string) error {
	toEmails := []string{
		email,
	}
	auth := smtp.PlainAuth("", e.FromEmail, e.password, e.smtpHost)
	err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.FromEmail, toEmails, []byte(message))
	if err != nil {
		slog.Error(
			"error sending email",
			slog.String("error", err.Error()))
		return err
	}
	slog.Info(
		"Email successfully sent",
		slog.String("email", email),
		slog.String("message", message))
	return nil
}

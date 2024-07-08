package notifier

import (
	"fmt"
	"net/smtp"
)

type I_Notifier interface {
	NotifyAll()
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

func (e *emailNotifier) NotifyAll() {
	//TODO:
	return
}

func (e *emailNotifier) NotifyOne(message string, email string) error {
	fmt.Println(email)
	fmt.Println(message)
	fmt.Println(e.smtpHost)
	fmt.Println(e.smtpPort)

	toEmails := []string{
		email,
	}
	auth := smtp.PlainAuth("", e.FromEmail, e.password, e.smtpHost)
	err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.FromEmail, toEmails, []byte(message))
	if err != nil {
		fmt.Println("ERROR: sending email", err)
		return err
	}
	fmt.Printf("Email successfully sent to %s\n", email)
	return nil
}

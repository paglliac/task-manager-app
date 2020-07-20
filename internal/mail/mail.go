package mail

import (
	"log"
	"net/smtp"
)

type DefaultMail struct {
	Subject string
	Body    string
	To      string
}

func (d DefaultMail) getSubject() string {
	return d.Subject
}

func (d DefaultMail) getBody() string {
	return d.Body
}

func (d DefaultMail) getTo() string {
	return d.To
}

type Mail interface {
	getSubject() string
	getBody() string
	getTo() string
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func Send(mail Mail) {
	auth := unencryptedAuth{smtp.PlainAuth(
		"",
		"user@example.com",
		"password",
		"mailhog",
	),
	} // Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{mail.getTo()}
	msg := []byte(
		"From: admin@tasks17.com\r\n" +
			"To: " + mail.getTo() + "\r\n" +
			"Subject: " + mail.getSubject() + "\r\n" +
			"\r\n" +
			mail.getBody() + "\r\n")

	err := smtp.SendMail(
		"mailhog:1025",
		auth,
		"admin@tasks17.com",
		to,
		msg,
	)
	if err != nil {
		log.Fatal(err)
	}
}

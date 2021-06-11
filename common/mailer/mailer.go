package mailer

import (
	"crypto/tls"
	"strings"

	gomail "gopkg.in/mail.v2"
)

func Mailer(receiver string, name string, pass string) {

	host := "smtp.gmail.com"

	port := 587

	user := "example@gmail.com" // Sender Email ID

	password := "12345" // Password

	mail := gomail.NewMessage()

	mail.SetHeader("From", "noreplymailed2020@gmail.com") // sender email

	mail.SetHeader("To", receiver) //receiver email

	// mail.SetAddressHeader("Cc", "dan@example.com", "Dan")	// cc emails

	mail.SetHeader("Subject", "Password reset") // Subject mail

	// message
	m := "Hi <b>name</b>, <br> <p>Here is the new password- <b>Password</b></p><br><p>If password reset wasn’t intended: If you didn't make the request, just ignore this email.</p><br><br>Thanks"
	message := strings.ReplaceAll(m, "name", name)
	message1 := strings.ReplaceAll(message, "Password", pass)

	mail.SetBody("text/html", message1) // email body

	// mail.Attach("/home/Alex/lolcat.jpg")		// Attachments

	d := gomail.NewDialer(host, port, user, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(mail); err != nil {
		panic(err)
	}
}

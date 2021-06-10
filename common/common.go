package common

import (
	"crypto/tls"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomPasswordGenerator() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)
	return str
}

func Mailer(receiver string, name string, pass string) {
	host := "smtp.gmail.com"
	port := 587
	user := "example@gmail.com" // Sender Email ID
	password := "12345"         // Password
	mail := gomail.NewMessage()
	mail.SetHeader("From", "noreplymailed2020@gmail.com")
	mail.SetHeader("To", receiver)
	// mail.SetAddressHeader("Cc", "dan@example.com", "Dan")
	mail.SetHeader("Subject", "Hello!")
	m := "Hi <b>name</b>, <br> <p>Here is the new password- <b>Password</b></p><br><p>If password reset wasn’t intended: If you didn't make the request, just ignore this email.</p><br><br>Thanks"
	message := strings.ReplaceAll(m, "name", name)
	message1 := strings.ReplaceAll(message, "Password", pass)
	mail.SetBody("text/html", message1)
	// mail.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(host, port, user, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(mail); err != nil {
		panic(err)
	}
}

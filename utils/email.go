package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"warunk-bem/domain"
	"warunk-bem/helpers"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	Code      int
	FirstName string
	Subject   string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *domain.User, data *EmailData) {
	_, err := helpers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	SMTPPortINT, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		fmt.Printf("Failed to convert SMTP_PORT to integer: %s\n", err.Error())
		return
	}

	// Sender data.
	from := os.Getenv("EMAIL_FROM")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := SMTPPortINT
	smtpPass := os.Getenv("SMTP_PASS")
	smtpUser := os.Getenv("SMTP_USER")
	fromName := os.Getenv("FROM_NAME")
	to := user.Email

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	mailer := gomail.NewMessage()
	mailer.SetAddressHeader("From", from, fromName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", data.Subject)
	mailer.SetBody("text/html", body.String())
	mailer.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpUser,
		smtpPass,
	)

	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// // Send Email
	// if err := d.DialAndSend(m); err != nil {
	// 	log.Fatal("Could not send email: ", err)
	// }

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal("Could not send email: ", err.Error())
	}

	log.Println("Mail sent!")
}

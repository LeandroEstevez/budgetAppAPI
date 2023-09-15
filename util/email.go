package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
	ToEmail   string
}

// ? Email template parser
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

	fmt.Println("Am parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

// ? Email template parser
func SendEmail(data *EmailData, templateName string) error {

	// Sender data.
	from := "leandroest111298@gmail.com"
	smtpPass := "iptidjnbeilmnrgn"
	smtpUser := "leandroest111298@gmail.com"
	to := data.ToEmail
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	fmt.Println(to)
	fmt.Println("Parsed the template")

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())

	fmt.Println("Got the template")

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("Sent the email")
	return nil
}

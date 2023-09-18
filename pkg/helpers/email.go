package helpers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"myapp/pkg/configs"
	"myapp/src/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL          string
	FirstName    string
	Subject      string
	Message      string
	TypeOfAction string
	SiteData     configs.SiteData
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

func SendEmail(user models.User, emailData *EmailData, emailTemplatename string) {
	config, err := configs.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// check for an empty struct
	if emailData.SiteData == (configs.SiteData{}) {
		siteData, _ := configs.GetSiteData(".")
		// update siteData
		emailData.SiteData = siteData
	}

	// Sender data.
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	emailTemplate := emailTemplatename
	isHtml := strings.HasSuffix(emailTemplate, ".html")
	if !isHtml {
		emailTemplate = fmt.Sprintf("%s.html", emailTemplate)
	}
	template.ExecuteTemplate(&body, emailTemplate, &emailData)

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(config.EmailFrom, config.AppName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", emailData.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	// Settings for SMTP server
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	tlsConfig := config.IsDebug
	log.Println(tlsConfig)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.TLSConfig = &tls.Config{InsecureSkipVerify: tlsConfig}

	// TODO: send email with celery task
	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}

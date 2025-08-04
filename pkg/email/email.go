package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

type EmailRepository interface {
	SendEmail(to, subject, templateName string, data interface{}) error
}

type EmailConfig struct {
	Host                   string `mapstructure:"EMAIL_HOST" validate:"required"`
	Port                   string `mapstructure:"EMAIL_PORT" validate:"required"`
	Username               string `mapstructure:"EMAIL_USERNAME"`
	Password               string `mapstructure:"EMAIL_PASSWORD"`
	From                   string `mapstructure:"EMAIL_FROM" validate:"required,email"`
	FromName               string `mapstructure:"EMAIL_FROM_NAME" validate:"required"`
	ResetPasswordURL       string `mapstructure:"RESET_PASSWORD_URL" validate:"required"`
	ResetPasswordExpiresIn int64  `mapstructure:"RESET_PASSWORD_EXPIRES_IN" validate:"required"`
}

type emailRepo struct {
	cfg *EmailConfig
}

func NewEmailRepo(cfg *EmailConfig) EmailRepository {
	return &emailRepo{
		cfg: cfg,
	}
}

func (m *emailRepo) SendEmail(to, subject, templateName string, data interface{}) error {
	// Read the template
	templatePath := filepath.Join("pkg", "email", "templates", templateName+".html")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing email template: %w", err)
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}

	// Prepare MIME message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	fromHeader := fmt.Sprintf("From: %s <%s>\n", m.cfg.FromName, m.cfg.From)
	subjectHeader := "Subject: " + subject + "\n"
	msg := []byte(fromHeader + subjectHeader + mime + body.String())

	// Authenticate if credentials are provided
	var auth smtp.Auth
	if m.cfg.Username != "" && m.cfg.Password != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}

	// Send the email
	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)
	if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

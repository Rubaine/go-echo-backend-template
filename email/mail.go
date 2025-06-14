package email

import "html/template"

// Vars contains template variables for email content
type Vars struct {
	Title      string
	Text       string
	ButtonText string
	ButtonURL  string
}

// Structure represents an email message structure
type Structure struct {
	To      string
	Subject string
	Vars    Vars
}

// Config contains email server configuration
// Deprecated: Use GmailConfig for Gmail-specific configuration
type Config struct {
	From            string
	User            string
	Password        string
	Host            string
	Port            string
	RecoverTemplate *template.Template
}

// EmailService represents a service for sending emails
type EmailService struct {
	SMTPConfig      *SMTPConfig
	RecoverTemplate *template.Template
	FromName        string
}

// NewEmailService creates a new email service with SMTP configuration
func NewEmailService(smtpConfig *SMTPConfig, fromName string, recoverTemplate *template.Template) *EmailService {
	return &EmailService{
		SMTPConfig:      smtpConfig,
		FromName:        fromName,
		RecoverTemplate: recoverTemplate,
	}
}

// NewGmailService creates a new email service with Gmail configuration
func NewGmailService(username, password, fromName string, recoverTemplate *template.Template) *EmailService {
	return NewEmailService(NewGmailConfig(username, password), fromName, recoverTemplate)
}

// NewOVHService creates a new email service with OVH configuration
func NewOVHService(username, password, fromName string, recoverTemplate *template.Template) *EmailService {
	return NewEmailService(NewOVHConfig(username, password), fromName, recoverTemplate)
}

// New creates a new email structure
func New(to, subject, title, text, btnText, btnUrl string) Structure {
	return Structure{
		To:      to,
		Subject: subject,
		Vars: Vars{
			Title:      title,
			Text:       text,
			ButtonText: btnText,
			ButtonURL:  btnUrl,
		},
	}
}

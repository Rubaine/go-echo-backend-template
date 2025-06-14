package email

import (
	"bytes"
	"fmt"
	"html"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

// Send sends an email using the provided configuration
func (email Structure) Send(config Config) error {
	// Legacy method - convert to SMTP config for backward compatibility
	smtpConfig := &SMTPConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.User,
		Password: config.Password,
		AuthType: AuthTypePlain,
		Provider: ProviderCustom,
	}

	service := &EmailService{
		SMTPConfig:      smtpConfig,
		FromName:        config.From,
		RecoverTemplate: config.RecoverTemplate,
	}

	return service.SendEmail(email)
}

// SendEmail sends an email using the EmailService with proper SMTP configuration
func (service *EmailService) SendEmail(email Structure) error {
	if err := service.SMTPConfig.Validate(); err != nil {
		return fmt.Errorf("invalid SMTP configuration: %w", err)
	}

	// Prepare sender and recipient addresses
	from := mail.Address{
		Name:    service.FromName,
		Address: service.SMTPConfig.Username,
	}

	to := mail.Address{
		Name:    "",
		Address: email.To,
	}

	// Validate recipient email
	if _, err := mail.ParseAddress(email.To); err != nil {
		return fmt.Errorf("invalid recipient email address: %w", err)
	}

	// Build email message
	message, err := service.buildMessage(from, to, email)
	if err != nil {
		return fmt.Errorf("failed to build email message: %w", err)
	}

	// Send the email
	return service.sendSMTP(from.Address, to.Address, message)
}

// buildMessage constructs the email message with headers and body
func (service *EmailService) buildMessage(from, to mail.Address, email Structure) (string, error) {
	headers := map[string]string{
		"From":         from.String(),
		"To":           to.String(),
		"Subject":      html.EscapeString(email.Subject),
		"MIME-Version": "1.0",
		"Content-Type": `text/html; charset="UTF-8"`,
		"Date":         time.Now().Format(time.RFC1123Z),
		"Message-ID":   fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), service.SMTPConfig.Host),
	}

	var message strings.Builder

	// Add headers
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	message.WriteString("\r\n")

	// Add body from template
	if service.RecoverTemplate != nil {
		var buff bytes.Buffer
		if err := service.RecoverTemplate.ExecuteTemplate(&buff, "email", email.Vars); err != nil {
			return "", fmt.Errorf("failed to execute email template: %w", err)
		}
		message.WriteString(buff.String())
	} else {
		// Fallback to simple text if no template
		message.WriteString(fmt.Sprintf("<html><body><h1>%s</h1><p>%s</p></body></html>",
			html.EscapeString(email.Vars.Title),
			html.EscapeString(email.Vars.Text)))
	}

	return message.String(), nil
}

// sendSMTP handles the SMTP connection and message transmission
func (service *EmailService) sendSMTP(from, to, message string) error {
	// Establish connection
	conn, err := net.DialTimeout("tcp", service.SMTPConfig.GetServerAddress(), 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, service.SMTPConfig.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Start TLS
	if err := client.StartTLS(service.SMTPConfig.GetTLSConfig()); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	auth := service.SMTPConfig.GetAuth()
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to initialize data transfer: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write message data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data transfer: %w", err)
	}

	return nil
}

// SendSimpleEmail is a convenience method for sending simple emails
func (service *EmailService) SendSimpleEmail(to, subject, htmlBody string) error {
	email := Structure{
		To:      to,
		Subject: subject,
		Vars: Vars{
			Title: subject,
			Text:  htmlBody,
		},
	}
	return service.SendEmail(email)
}

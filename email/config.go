package email

import (
	"fmt"
	"html/template"
	"os"
)

// ConfigFromEnv creates an SMTP configuration from environment variables
func ConfigFromEnv() (*SMTPConfig, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	if host == "" || port == "" || user == "" || password == "" {
		return nil, fmt.Errorf("missing required SMTP environment variables (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASSWORD)")
	}

	config := NewCustomConfig(host, port, user, password)

	// Auto-detect provider based on host
	switch host {
	case "smtp.gmail.com":
		config.Provider = ProviderGmail
	case "ssl0.ovh.net":
		config.Provider = ProviderOVH
	default:
		config.Provider = ProviderCustom
	}

	return config, nil
}

// ServiceFromEnv creates an EmailService from environment variables
func ServiceFromEnv(recoverTemplate *template.Template) (*EmailService, error) {
	config, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}

	displayName := os.Getenv("SMTP_DISPLAYNAME")
	if displayName == "" {
		displayName = config.Username // fallback to username
	}

	return NewEmailService(config, displayName, recoverTemplate), nil
}

// OVHConfigFromEnv creates an OVH configuration from the specified environment variables
func OVHConfigFromEnv() (*SMTPConfig, error) {
	host := os.Getenv("SMTP_HOST")         // ssl0.ovh.net
	port := os.Getenv("SMTP_PORT")         // 587
	user := os.Getenv("SMTP_USER")         // noreply@example.com
	password := os.Getenv("SMTP_PASSWORD") // UP2Q2w4uDM967wXdRueT

	if host == "" {
		host = "ssl0.ovh.net"
	}
	if port == "" {
		port = "587"
	}
	if user == "" || password == "" {
		return nil, fmt.Errorf("missing required OVH SMTP credentials (SMTP_USER, SMTP_PASSWORD)")
	}

	return &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
		AuthType: AuthTypePlain,
		Provider: ProviderOVH,
	}, nil
}

// OVHServiceFromEnv creates an EmailService configured for OVH from environment variables
func OVHServiceFromEnv(recoverTemplate *template.Template) (*EmailService, error) {
	config, err := OVHConfigFromEnv()
	if err != nil {
		return nil, err
	}

	displayName := os.Getenv("SMTP_DISPLAYNAME")
	if displayName == "" {
		displayName = "ExampleApp" // default display name
	}

	return NewEmailService(config, displayName, recoverTemplate), nil
}

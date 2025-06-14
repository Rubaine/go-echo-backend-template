package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// AuthType represents the type of authentication to use
type AuthType string

const (
	AuthTypePlain AuthType = "PLAIN"
	AuthTypeLogin AuthType = "LOGIN"
	AuthTypeOAuth AuthType = "OAUTH2"
)

// Provider represents different email providers
type Provider string

const (
	ProviderGmail  Provider = "gmail"
	ProviderOVH    Provider = "ovh"
	ProviderCustom Provider = "custom"
)

// SMTPConfig contains generic SMTP configuration for any provider
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	AuthType AuthType
	Provider Provider
}

// NewGmailConfig creates a new Gmail configuration with secure defaults
func NewGmailConfig(username, password string) *SMTPConfig {
	return &SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: username,
		Password: password,
		AuthType: AuthTypePlain,
		Provider: ProviderGmail,
	}
}

// NewOVHConfig creates a new OVH configuration with secure defaults
func NewOVHConfig(username, password string) *SMTPConfig {
	return &SMTPConfig{
		Host:     "ssl0.ovh.net",
		Port:     "587",
		Username: username,
		Password: password,
		AuthType: AuthTypePlain,
		Provider: ProviderOVH,
	}
}

// NewCustomConfig creates a custom SMTP configuration
func NewCustomConfig(host, port, username, password string) *SMTPConfig {
	return &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		AuthType: AuthTypePlain,
		Provider: ProviderCustom,
	}
}

// GetAuth returns the appropriate SMTP auth based on the configuration
func (sc *SMTPConfig) GetAuth() smtp.Auth {
	switch sc.AuthType {
	case AuthTypePlain:
		return smtp.PlainAuth("", sc.Username, sc.Password, sc.Host)
	case AuthTypeLogin:
		return &loginAuth{username: sc.Username, password: sc.Password}
	default:
		// Default to PLAIN auth
		return smtp.PlainAuth("", sc.Username, sc.Password, sc.Host)
	}
}

// GetTLSConfig returns a secure TLS configuration
func (sc *SMTPConfig) GetTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         sc.Host,
		MinVersion:         tls.VersionTLS12,
	}
}

// Validate checks if the SMTP configuration is valid
func (sc *SMTPConfig) Validate() error {
	if sc.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if sc.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if !strings.Contains(sc.Username, "@") {
		return fmt.Errorf("username must be a valid email address")
	}
	if sc.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if sc.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

// GetServerAddress returns the full server address for connections
func (sc *SMTPConfig) GetServerAddress() string {
	return net.JoinHostPort(sc.Host, sc.Port)
}

// loginAuth implements the LOGIN SMTP authentication mechanism
// This is kept for compatibility but PLAIN is preferred for Gmail
type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		prompt := strings.ToLower(strings.TrimSpace(string(fromServer)))
		switch {
		case strings.Contains(prompt, "username"):
			return []byte(a.username), nil
		case strings.Contains(prompt, "password"):
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown authentication prompt: %s", string(fromServer))
		}
	}
	return nil, nil
}

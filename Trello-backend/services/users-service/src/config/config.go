package config

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUser:     os.Getenv("SMTP_EMAIL"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}

	if config.SMTPHost == "" || config.SMTPPort == "" || config.SMTPUser == "" || config.SMTPPassword == "" {
		return nil, fmt.Errorf("missing required SMTP configuration environment variables")
	}

	return config, nil
}
var PasswordBlacklist map[string]bool

func LoadPasswordBlacklist(filePath string) error {
    PasswordBlacklist = make(map[string]bool)

    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        PasswordBlacklist[strings.TrimSpace(scanner.Text())] = true
    }

    return scanner.Err()
}

func IsBlacklistedPassword(password string) bool {
    _, exists := PasswordBlacklist[password]
    return exists
}
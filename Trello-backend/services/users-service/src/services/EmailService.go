package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendVerificationEmail sends a verification email with the given code to the specified email address.
func SendVerificationEmail(email, code string) error {
	SMTPHost := os.Getenv("SMTP_HOST")
	SMTPPort := os.Getenv("SMTP_PORT")
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", SMTPEmail, SMTPPassword, SMTPHost)

	msg := []byte(fmt.Sprintf("Subject: Verification Code\n\nYour verification code is: %s", code))

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTPHost,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", SMTPHost, SMTPPort), tlsConfig)
	if err != nil {
		log.Println("Error while connecting to SMTP server:", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, SMTPHost)
	if err != nil {
		log.Println("Error while creating SMTP client:", err)
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		log.Println("Error while authenticating:", err)
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	if err = client.Mail(SMTPEmail); err != nil {
		log.Println("Error while setting sender email:", err)
		return fmt.Errorf("failed to set sender email: %w", err)
	}

	if err = client.Rcpt(email); err != nil {
		log.Println("Error while setting recipient email:", err)
		return fmt.Errorf("failed to set recipient email: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		log.Println("Error while sending email data:", err)
		return fmt.Errorf("failed to send email data: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Println("Error while writing message:", err)
		return fmt.Errorf("failed to write message content: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Println("Error while closing message writer:", err)
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	client.Quit()

	log.Println("Verification email sent successfully to:", email)
	return nil
}

func SendMagicLinkEmail(email, magicLink string) error {
	SMTPHost := os.Getenv("SMTP_HOST")
	SMTPPort := os.Getenv("SMTP_PORT")
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", SMTPEmail, SMTPPassword, SMTPHost)

	msg := []byte(fmt.Sprintf("Subject: Magic Login Link\n\nClick the link below to log in:\n%s", magicLink))

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTPHost,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", SMTPHost, SMTPPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	if err = client.Mail(SMTPEmail); err != nil {
		return fmt.Errorf("failed to set sender email: %w", err)
	}

	if err = client.Rcpt(email); err != nil {
		return fmt.Errorf("failed to set recipient email: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send email data: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message content: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	client.Quit()
	return nil
}

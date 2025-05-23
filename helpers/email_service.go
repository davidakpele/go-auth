package helpers

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendVerificationEmail(toEmail, fullname, otp string) error {
    from :=os.Getenv("FROM_EMAIL")
    password := os.Getenv("FROM_EMAIL_PASSWORD")

    subject := "Verify Your Account – OTP Inside"
    body := fmt.Sprintf("Hello %s,\n\nYour OTP code is: %s\n\nThanks!", fullname, otp)

    msg := []byte("From: " + from + "\r\n" +
        "To: " + toEmail + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
        body)

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
    if err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }
    log.Println("Email sent successfully to", toEmail)
    return nil
}

// ResendSendOTPEmail sends an OTP email using SendGrid
func ResendSendOTPEmail(email, fullname, otp string) error {
	sendgridAPIKey := "your-sendgrid-api-key"
	fromEmail := "your-email@example.com"
	fromName := "Your Company"

	client := sendgrid.NewSendClient(sendgridAPIKey)
	from := mail.NewEmail(fromName, fromEmail)
	subject := "Your OTP Verification Code"
	to := mail.NewEmail(fullname, email)
	plainTextContent := fmt.Sprintf("Hello %s,\n\nYour OTP code is: %s\n\nBest Regards,\nYour Company", fullname, otp)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, plainTextContent)

	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully with status:", response.StatusCode)
	return nil
}
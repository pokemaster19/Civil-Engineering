package utils

import (
	"ControlSystem/models"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/resend/resend-go/v2"
)

func GenerateCorporateEmail(firstName, lastName string, db *gorm.DB) (string, error) {
	firstName = Transliterate(firstName)
	lastName = Transliterate(lastName)

	base := fmt.Sprintf("%s.%s@controlsystem.ru", strings.ToLower(firstName), strings.ToLower(lastName))
	email := base
	suffix := 1

	for {
		var user models.User
		err := db.Where("email = ?", email).First(&user).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return email, nil
		}
		if err != nil {
			return "", fmt.Errorf("ошибка проверки email: %w", err)
		}
		email = fmt.Sprintf("%s.%s%d@controlsystem.ru", strings.ToLower(firstName), strings.ToLower(lastName), suffix)
		suffix++
	}
}

func GeneratePassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate password: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func SendRegistrationEmail(personalEmail, corporateEmail, password string, resendClient *resend.Client) error {
	htmlBody := fmt.Sprintf("<p>Ваш email: %s</p><p>Пароль: %s</p>", corporateEmail, password)

	params := &resend.SendEmailRequest{
		From:    "Система контроля <onboarding@resend.dev>",
		To:      []string{personalEmail},
		Subject: "Ваши учетные данные",
		Html:    htmlBody,
	}

	sent, err := resendClient.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Письмо отправлено, ID:", sent.Id)
	return nil
}

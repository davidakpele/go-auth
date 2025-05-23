package utils

import (
	"crypto/rand"
	"regexp"
	"fmt"
	"math/big"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// IsValidEmail validates an email format
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// CheckPassword compares a hashed password with a plain text password
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil 
}

func GenerateSecureOTP(length int) (string, error) {
	otp := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10)) 
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num.Int64())
	}
	return otp, nil
}
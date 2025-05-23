package repositories

import (
	"api-service/models"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// IsEmailExist checks if an email already exists in the database
func (ar *AuthRepository) IsEmailExist(email string) bool {
	var count int64
	err := ar.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		log.Println("Error checking email existence:", err)
		return false
	}
	return count > 0
}

// RegisterUser inserts a new user into the database
func (ar *AuthRepository) RegisterUser(user *models.User) error {
	err := ar.db.Create(user).Error
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func (ar *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ar.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) StoreOTP(userID uint, otp string) error {
    otpRecord := models.AccountVerification{
        UserId:    userID,
        OTP:       otp,
        ExpiredAt: time.Now().Add(45 * time.Minute),
    }
    return r.db.Create(&otpRecord).Error
}

// Get OTP record by user ID
func (r *AuthRepository) GetOTPByUserID(userID uint) (*models.AccountVerification, error) {
	var otp models.AccountVerification
	if err := r.db.Where("user_id = ?", userID).First(&otp).Error; err != nil {
		return nil, errors.New("OTP not found")
	}
	return &otp, nil
}

// Update user status to ACTIVATED
func (r *AuthRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete OTP after verification
func (r *AuthRepository) DeleteOTP(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.AccountVerification{}).Error
}

func (r *AuthRepository) GenerateAndStoreOTP(userID uint) (string, error) {
	// Generate new OTP
	rand.Seed(time.Now().UnixNano())
	newOTP := fmt.Sprintf("%06d", rand.Intn(1000000)) 

	// Set expiration time (10 minutes)
	expirationTime := time.Now().Add(10 * time.Minute)

	// Create OTP record
	newOTPRecord := models.AccountVerification{
		UserId:    userID,
		OTP:       newOTP,
		ExpiredAt: expirationTime,
	}

	// Store in DB
	if err := r.db.Create(&newOTPRecord).Error; err != nil {
		return "", err
	}

	return newOTP, nil
}


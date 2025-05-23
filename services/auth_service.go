package services

import (
	"api-service/models"
	"api-service/repositories"
	"api-service/requests"
	"api-service/utils"
	"api-service/helpers"
	"net/http"
	"time"
)

type AuthService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (as *AuthService) Register(req requests.RegisterRequest) (map[string]interface{}, int) {
	// Validation
	if req.Fullname == "" {
		return map[string]interface{}{"status": "error", "message": "Fullname is required"}, http.StatusBadRequest
	}
	if req.Email == "" {
		return map[string]interface{}{"status": "error", "message": "Email address is required"}, http.StatusBadRequest
	}
	if !utils.IsValidEmail(req.Email) {
		return map[string]interface{}{"status": "error", "message": "Invalid email address"}, http.StatusBadRequest
	}
	if req.Password == "" {
		return map[string]interface{}{"status": "error", "message": "Password is required"}, http.StatusBadRequest
	}

	// Check if email already exists
	if as.authRepo.IsEmailExist(req.Email) {
		return map[string]interface{}{"status": "error", "message": "Email is already registered"}, http.StatusBadRequest
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "Error hashing password"}, http.StatusInternalServerError
	}

	// Save user
	user := models.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := as.authRepo.RegisterUser(&user); err != nil {
		return map[string]interface{}{"status": "error", "message": "Registration failed"}, http.StatusInternalServerError
	}

	// Generate Secure OTP
	otp, err := utils.GenerateSecureOTP(6) 
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "Error generating OTP"}, http.StatusInternalServerError
	}

	// Send OTP via Email
	if err := helpers.SendVerificationEmail(req.Email, req.Fullname, otp); err != nil {
		return map[string]interface{}{"status": "error", "message": "Failed to send verification email"}, http.StatusInternalServerError
	}

	if err := as.authRepo.StoreOTP(user.ID, otp); err != nil {
		return map[string]interface{}{
			"status":  "error",
			"message": "Failed to store OTP",
		}, http.StatusInternalServerError
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "A message has been sent to your email to complete the sign-up process.",
	}, http.StatusCreated
}

// Login function
func (as *AuthService) Login(req requests.LoginRequest) (map[string]interface{}, int) {
	// Fetch user from the database by email
	user, err := as.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "Invalid email or password"}, http.StatusUnauthorized
	}

	// Compare the entered password with the stored hashed password
	if !utils.CheckPassword(user.Password, req.Password) {
		return map[string]interface{}{"status": "error", "message": "Invalid email or password"}, http.StatusUnauthorized
	}

	if user.Status != "VERIFIED" || !user.Enabled {
		return map[string]interface{}{
			"status": "error",
			"message": "This account has not been verified or is disabled. Please verify your account.",
		}, http.StatusUnauthorized
	}
	
	// Generate JWT token
	token, _ := utils.GenerateJWT(user.ID, user.Email, user.Role)

	// Return user details along with the token
	return map[string]interface{}{
		"status":   "success",
		"message":  "Login successful",
		"token":    token,
		"id":       user.ID,
		"email":    user.Email,
		"fullname": user.Fullname,
	}, http.StatusOK
}

func (as *AuthService) VerifyAccount(req requests.VerifyAccountRequest) (map[string]interface{}, int) {
	// Fetch user by email
	user, err := as.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "User not found"}, http.StatusNotFound
	}

	// Check if OTP exists in account_verification table
	otpRecord, err := as.authRepo.GetOTPByUserID(user.ID)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "Invalid OTP provided"}, http.StatusBadRequest
	}

	// Check if OTP matches
	if otpRecord.OTP != req.OTP {
		return map[string]interface{}{"status": "error", "message": "Invalid OTP provided"}, http.StatusBadRequest
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiredAt) {
		return map[string]interface{}{"status": "error", "message": "OTP has expired"}, http.StatusBadRequest
	}

	// Update user status to "ACTIVATED"
	user.Status = "ACTIVATED"
	user.Enabled=true
	if err := as.authRepo.UpdateUser(user); err != nil {
		return map[string]interface{}{"status": "error", "message": "Failed to activate account"}, http.StatusInternalServerError
	}

	// Delete OTP from database
	if err := as.authRepo.DeleteOTP(user.ID); err != nil {
		return map[string]interface{}{"status": "error", "message": "Failed to clean up OTP"}, http.StatusInternalServerError
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Account verified successfully. You can now log in.",
	}, http.StatusOK
}

// Resend OTP function
func (as *AuthService) ResendOTP(email string) (map[string]interface{}, int) {
	// Fetch user by email
	user, err := as.authRepo.GetUserByEmail(email)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "User not found"}, http.StatusNotFound
	}

	// Delete old OTP & Generate new OTP in repository
	newOTP, err := as.authRepo.GenerateAndStoreOTP(user.ID)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "Failed to generate OTP"}, http.StatusInternalServerError
	}

	// Send new OTP via email
	helpers.ResendSendOTPEmail(user.Email, user.Fullname, newOTP)

	return map[string]interface{}{
		"status":  "success",
		"message": "A new OTP has been sent to your email.",
	}, http.StatusOK
}
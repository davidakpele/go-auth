package controllers

import (
	"api-service/requests"
	"api-service/services"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req requests.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request payload"})
		return
	}
	// Call service layer to handle registration
	response, statusCode := ac.authService.Register(req)
	c.JSON(statusCode, response)
}

// Login handles user login using Gin context
func (ac *AuthController) Login(c *gin.Context) {
	var req requests.LoginRequest

	// Bind JSON payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request payload"})
		return
	}

	// Call service layer to handle login
	response, statusCode := ac.authService.Login(req)
	c.JSON(statusCode, response)
}

func (ac *AuthController) Logout(c *gin.Context) {
	// Destroy the session
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout User Successful",
	})
}

func (ac *AuthController) VerifyAccount(c *gin.Context) {
	var req requests.VerifyAccountRequest
	// Bind JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	// Call the service layer
	response, status := ac.authService.VerifyAccount(req)
	c.JSON(status, response)
}

// Resend OTP handler
func (ac *AuthController) ResendOTP(c *gin.Context) {
	// Parse request body
	var requestData struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid email format"})
		return
	}

	// Call service method
	response, statusCode := ac.authService.ResendOTP(requestData.Email)

	// Send JSON response
	c.JSON(statusCode, response)
}
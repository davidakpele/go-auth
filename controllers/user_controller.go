package controllers

import (
	"api-service/services"
	"net/http"
	"strconv"
	"api-service/requests"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ac *UserController) GetUserByID(c *gin.Context) {
	// Get user ID from the URL param
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	// Call service
	response, statusCode := ac.userService.GetUserByID(uint(userID))
	c.JSON(statusCode, response)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var req requests.UpdateUserRequest
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request payload"})
		return
	}

	response, statusCode := uc.userService.UpdateUser(uint(userID), req)
	c.JSON(statusCode, response)
}

// Delete User
func (ctrl *UserController) Delete(c *gin.Context) {
	// Get user ID from the URL param
	userID := c.Param("id")

	// Validate User ID
	uid, err := strconv.Atoi(userID)
	if err != nil || uid <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID must be an integer"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User ID required"})
		return
	}

	// Call Service to Delete User
	err = ctrl.userService.DeleteUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}

// GetAllUsers fetches all users and sends the response
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	// Get pagination parameters from query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Ensure valid pagination values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Fetch users with pagination
	users, totalCount, err := ctrl.userService.GetAllUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch users"})
		return
	}

	// Prepare paginated response
	response := gin.H{
		"status":       "success",
		"totalUsers":   totalCount,
		"currentPage":  page,
		"pageSize":     pageSize,
		"totalPages":   (totalCount + int64(pageSize) - 1) / int64(pageSize),
		"users":        users,
	}

	c.JSON(http.StatusOK, response)
}

// PasswordUpdate handles the request to update a user's password
func (ctrl *UserController) PasswordUpdate(c *gin.Context) {
	 var req requests.UpdatePasswordRequest

	// Bind the incoming JSON request to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}

	// Validate required fields
	if req.UserID == 0 || req.OldPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "All fields are required"})
		return
	}

	// Call the service layer to update the password
	err := ctrl.userService.UpdatePassword(req.UserID, req.OldPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Send success response
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password updated successfully"})
}

// func (ctrl *UserController) Profile(c *gin.Context) {
// 	// Validate token and get user
// 	security := security.SecurityFilterChain{}
// 	user, err := security.IsValidToken(c)
// 	if err != nil || user == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access. Please login."})
// 		return
// 	}

// 	// Validate HTTP method
// 	if c.Request.Method != http.MethodGet {
// 		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": "error", "message": "Method Not Allowed"})
// 		return
// 	}

// 	// Get pagination parameters
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

// 	// Fetch user profile
// 	userData, err := ctrl.userService.GetUserProfile(int(user.ID), page, pageSize)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found."})
// 		return
// 	}

// 	// Return user profile data
// 	c.JSON(http.StatusOK, gin.H{"user": userData})
// }

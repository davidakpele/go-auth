package services

import (
	"api-service/models"
	"api-service/repositories"
	"api-service/requests"
	"api-service/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (as *UserService) GetUserByID(userID uint) (map[string]interface{}, int) {
	user, err := as.userRepo.GetUserByID(userID)
	if err != nil {
		return map[string]interface{}{"status": "error", "message": "User not found"}, http.StatusNotFound
	}

	// Update user views if needed
	as.userRepo.UpdateUserView(userID)

	// Return user details
	return map[string]interface{}{
		"status": "success",
		"user": gin.H{
			"id":       user.ID,
			"fullname": user.Fullname,
			"email":    user.Email,
			"views": user.Views,
			"createdAt":user.CreatedAt,
		},
	}, http.StatusOK
}

func (us *UserService) UpdateUser(userID uint, req requests.UpdateUserRequest) (map[string]interface{}, int) {
	// Fetch user
	user, err := us.userRepo.GetUserByID(userID)
	if err != nil {
		return gin.H{"status": "error", "message": "User not found"}, http.StatusNotFound
	}

	// Validate required fields
	if req.Fullname == "" || req.Email == "" || req.ContactTitle == "" ||
		req.FacebookLink == "" || req.InstagramLink == "" || req.TwitterLink == "" || req.LinkedinLink == "" {
		return gin.H{"status": "error", "message": "All fields are required"}, http.StatusBadRequest
	}

	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return gin.H{"status": "error", "message": "Invalid email format"}, http.StatusBadRequest
	}

	// Prepare update data
	user.Fullname = req.Fullname
	user.Email = req.Email
	user.ContactTitle = req.ContactTitle
	user.FacebookLink = req.FacebookLink
	user.InstagramLink = req.InstagramLink
	user.TwitterLink = req.TwitterLink
	user.LinkedInLink = req.LinkedinLink

	// Update user
	err = us.userRepo.UpdateUser(userID, *user)
	if err != nil {
		return gin.H{"status": "error", "message": "Update failed"}, http.StatusInternalServerError
	}

	return gin.H{"status": "success", "message": "User updated successfully"}, http.StatusOK
}

// DeleteUser deletes a user by ID
func (service *UserService) DeleteUser(userID int) error {
	// Call the repository to delete the user
	err := service.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserView updates user-specific data like view count
func (service *UserService) UpdateUserView(userID uint) error {
	// Call the repository to update user view
	err := service.userRepo.UpdateUserView(userID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllUsers fetches all users from the repository
func (service *UserService) GetAllUsers(page, pageSize int) ([]models.User, int64, error) {
	return service.userRepo.GetAllUsers(page, pageSize)
}

// UpdatePassword updates the user's password using the repository method
func (service *UserService) UpdatePassword(userID int, oldPassword, newPassword, confirmPassword string) error {
	// Validate password match
	if newPassword != confirmPassword {
		return fmt.Errorf("new password does not match confirm password")
	}

	// Call repository to update the password
	err := service.userRepo.UpdateUserPassword(userID, oldPassword, newPassword)
	if err != nil {
		return err
	}

	return nil
}

// func (service *UserService) GetUserProfile(userID, page, pageSize int) (*models.User, error) {
// 	offset := (page - 1) * pageSize
// 	return service.userRepo.GetUserViews(userID, offset, pageSize)
// }

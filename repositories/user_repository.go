package repositories

import (
	"api-service/models"
	"api-service/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID fetches a user by their ID
func (ar *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := ar.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserView increments the user's view count by 1
func (ar *UserRepository) UpdateUserView(userID uint) error {
	// Use GORM's Update feature to increment views
	err := ar.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("views", gorm.Expr("views + 1")).Error

	if err != nil {
		return err
	}
	return nil
}


func (ur *UserRepository) UpdateUser(userID uint, user models.User) error {
	// Update user record
	err := ur.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"fullname":      user.Fullname,
			"email":         user.Email,
			"contact_title": user.ContactTitle,
			"facebook_link": user.FacebookLink,
			"instagram_link": user.InstagramLink,
			"twitter_link":  user.TwitterLink,
			"linkedin_link": user.LinkedInLink,
		}).Error

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID
func (repo *UserRepository) DeleteUser(userID int) error {
	var user models.User
	if err := repo.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if err := repo.db.Delete(&user).Error; err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

// GetAllUsers fetches all users from the database
func (repo *UserRepository) GetAllUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64

	// Count total users
	if err := repo.db.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Fetch users with pagination
	offset := (page - 1) * pageSize
	if err := repo.db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// UpdateUserPassword updates the password for a given user
func (repo *UserRepository) UpdateUserPassword(userID int, oldPassword, newPassword string) error {
	var user models.User
	if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err // User not found
	}

	// Check if the old password matches
	if !utils.CheckPassword(oldPassword, user.Password) {
		return fmt.Errorf("old password does not match")
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password in the database
	if err := repo.db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) VerifyUserAccount(email string, userID int) (*models.User, error) {
	var user models.User
	err := repo.db.Where("id = ? AND email = ?", userID, email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// func (repo *UserRepository) GetUserViews(userID, offset, itemsPerPage int) (map[string]interface{}, error) {
// 	var user models.User
// 	err := repo.db.Select("id, contact_title, mobile, enabled, status, facebook_link, instagram_link, twitter_link, linkedin_link, image, fullname, email, views").
// 		Where("id = ?", userID).
// 		First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Set default image if empty
// 	if user.Image == "" {
// 		user.Image = "ASSETS/images/profile/avatar-place-holder.png"
// 	}

// 	// Get total number of resources created by the user
// 	var totalResources int64
// 	err = repo.db.Model(&models.Resource{}).
// 		Where("user_id = ?", userID).
// 		Count(&totalResources).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get user resource views
// 	var userResourceViews struct {
// 		ID    int
// 		Views int
// 	}
// 	err = repo.db.Table("resources").
// 		Select("id, views").
// 		Where("user_id = ?", userID).
// 		First(&userResourceViews).Error
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, err
// 	}

// 	// Get total number of bookmarks on user's documents
// 	var totalBookmarks int64
// 	err = repo.db.Table("resources r").
// 		Select("COUNT(b.resource_id) AS bookmark_count").
// 		Joins("INNER JOIN bookmark b ON r.id = b.resource_id").
// 		Where("r.user_id = ?", userID).
// 		Count(&totalBookmarks).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Fetch all user-created resources with pagination
// 	var resources []models.Resource
// 	err = repo.db.Table("resources r").
// 		Select("r.*, u.id, u.image").
// 		Joins("INNER JOIN users u ON r.user_id = u.id").
// 		Where("r.user_id = ?", userID).
// 		Offset(offset).
// 		Limit(itemsPerPage).
// 		Find(&resources).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Process resources
// 	for i := range resources {
// 		if resources[i].ResourceTypeOptionId != "" {
// 			resources[i].ResourceTypeArray = splitResourceType(resources[i].ResourceType)
// 		} else {
// 			resources[i].ResourceTypeArray = []string{}
// 		}

// 		if resources[i].Image == "" {
// 			resources[i].Image = "ASSETS/images/profile/avatar-place-holder.png"
// 		}
// 	}

// 	// Prepare response
// 	response := map[string]interface{}{
// 		"user": user,
// 		"info": map[string]interface{}{
// 			"total_no_profile_views":          user.Views,
// 			"total_number_of_resources":       totalResources,
// 			"total_number_of_user_bookmarked": totalBookmarks,
// 			"total_no_user_view_resource":     userResourceViews.Views,
// 		},
// 		"resources": resources,
// 	}

// 	return response, nil
// }

// Helper function to split resource_type field
// func splitResourceType(resourceType string) []string {
// 	return strings.Split(resourceType, ",")
// }
package requests

type UpdatePasswordRequest  struct {
	UserID        int    `json:"userId"`
	OldPassword   string `json:"oldPassword"`
	NewPassword   string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
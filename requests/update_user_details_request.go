package requests

type UpdateUserRequest struct {
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	ContactTitle string `json:"contact_title"`
	FacebookLink string `json:"facebook_link"`
	InstagramLink string `json:"instagram_link"`
	TwitterLink  string `json:"twitter_link"`
	LinkedinLink string `json:"linkedin_link"`
}

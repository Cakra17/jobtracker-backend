package user

type UserResponse struct {
	Email string `json:"email"`
	Username string `json:"username"`
	AccessToken string `json:"access_token"`
}
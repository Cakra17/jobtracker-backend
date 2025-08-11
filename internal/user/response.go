package user

import "github.com/Cakra17/JobTracker-Api/internal/job"

type UserResponse struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type LoggedUserData struct {
	Id            string        `json:"id"`
	Email         string        `json:"email"`
	Username      string        `json:"username"`
	DisplayName   string        `json:"display_name"`
  AvatarUrl     job.NullString`json:"avatar_url"`
	EmailVerified bool          `json:"email_verified"`
}

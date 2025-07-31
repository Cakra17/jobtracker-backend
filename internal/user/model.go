package user

import "database/sql"

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=5,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=20"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type ChangeAvatarRequest struct {
	AvatarUrl string `json:"avatar_url" validate:"required"`
}

type VerificationEmailRequest struct {
	Email    string `json:"email" validate:"required"`
	Verified bool   `json:"verified" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangeNameRequest struct {
	Username string `json:"username"`
	DisplayName string `json:"display_name"`

	ID string
}

type User struct {
	ID            string `db:"id"`
	Email         string `db:"email"`
	Username      string `db:"username"`
	PasswordHash  string `db:"password_hash"`
	DisplayName   string `db:"display_name"`
	AvatarUrl     sql.NullString `db:"avatar_url"`
	EmailVerified bool   `db:"email_verified"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
}

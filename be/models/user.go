package models

import (
	"mime/multipart"
	"time"
)

const (
	UserTypeGuest      = "guest"
	UserTypeRegistered = "registered"

	AuthProviderPassword = "password"
	AuthProviderGoogle   = "google"
	AuthProviderFacebook = "facebook"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DisplayName string    `gorm:"column:display_name;size:100;not null" json:"display_name"`
	AvatarURL   *string   `gorm:"column:avatar_url;type:text" json:"avatar_url,omitempty"`
	UserType    string    `gorm:"column:user_type;size:20;not null;default:guest" json:"user_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserAuthIdentity struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"column:user_id;not null" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider       string    `gorm:"column:provider;size:30;not null" json:"provider"`
	ProviderUserID *string   `gorm:"column:provider_user_id;size:255" json:"provider_user_id,omitempty"`
	Email          *string   `gorm:"column:email;size:255" json:"email,omitempty"`
	EmailVerified  bool      `gorm:"column:email_verified;not null;default:false" json:"email_verified"`
	PasswordHash   *string   `gorm:"column:password_hash;type:text" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (UserAuthIdentity) TableName() string {
	return "user_auth_identities"
}

type CreateGuestUserRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
}

type RegisterRequest struct {
	DisplayName     string                `json:"display_name" form:"display_name" binding:"required,max=100"`
	Email           string                `json:"email" form:"email" binding:"required,email,max=255"`
	Password        string                `json:"password" form:"password" binding:"required,min=8"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
	ProfilePicture  *multipart.FileHeader `json:"-" form:"profile_picture"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID             uint    `json:"id"`
	DisplayName    string  `json:"display_name"`
	Email          *string `json:"email"`
	ProfilePicture *string `json:"profile_picture"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
	UserType       string  `json:"user_type"`
}

type UserGameStats struct {
	CorrectAnswers int64 `json:"correct_answers"`
	Score          int64 `json:"score"`
}

type AuthResponse struct {
	UserGameStats
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

func NewUserResponse(user User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		DisplayName:    user.DisplayName,
		ProfilePicture: user.AvatarURL,
		AvatarURL:      user.AvatarURL,
		UserType:       user.UserType,
	}
}

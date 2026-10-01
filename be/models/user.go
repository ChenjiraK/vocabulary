package models

import "time"

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
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID          uint    `json:"id"`
	DisplayName string  `json:"display_name"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	UserType    string  `json:"user_type"`
}

type AuthResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

func NewUserResponse(user User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		UserType:    user.UserType,
	}
}

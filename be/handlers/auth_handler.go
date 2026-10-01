package handlers

import (
	"errors"
	"net/http"
	"os"

	"vocabulary/models"
	"vocabulary/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) CreateGuestUser(c *gin.Context) {
	var request models.CreateGuestUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "display_name is required"})
		return
	}

	user, err := h.authService.CreateGuestUser(c.Request.Context(), request.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot create guest user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "guest user created successfully",
		"user":    models.NewUserResponse(*user),
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request models.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "display_name, email, and password are required"})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, gin.H{"message": "email already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot register user"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email and password are required"})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	loginURL, err := h.authService.GoogleLoginURL()
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": loginURL})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	if !hasValidOAuthState(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code is required"})
		return
	}

	response, err := h.authService.LoginWithGoogleCode(c.Request.Context(), code)
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) FacebookLogin(c *gin.Context) {
	loginURL, err := h.authService.FacebookLoginURL()
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": loginURL})
}

func (h *AuthHandler) FacebookCallback(c *gin.Context) {
	if !hasValidOAuthState(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code is required"})
		return
	}

	response, err := h.authService.LoginWithFacebookCode(c.Request.Context(), code)
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func hasValidOAuthState(c *gin.Context) bool {
	expectedState := os.Getenv("OAUTH_STATE")
	if expectedState == "" {
		return true
	}

	return c.Query("state") == expectedState
}

func writeOAuthError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrOAuthNotConfigured) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "oauth provider is not configured"})
		return
	}
	if errors.Is(err, services.ErrOAuthProfileNotFound) {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cannot read oauth profile"})
		return
	}

	c.JSON(http.StatusBadGateway, gin.H{"message": "oauth login failed"})
}

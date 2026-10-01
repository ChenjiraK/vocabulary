package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"vocabulary/middlewares"
	"vocabulary/models"
	"vocabulary/services"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService services.GameService
}

func NewGameHandler(gameService services.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) StartSession(c *gin.Context) {
	var request models.StartGameSessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id is required"})
		return
	}

	session, err := h.gameService.StartSession(c.Request.Context(), request.UserID)
	if err != nil {
		writeGameError(c, err)
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *GameHandler) SubmitAnswer(c *gin.Context) {
	sessionID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var request models.SubmitGameAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id and word are required"})
		return
	}

	response, err := h.gameService.SubmitAnswer(c.Request.Context(), sessionID, request)
	if err != nil {
		writeGameError(c, err)
		return
	}

	if !response.Answer.IsCorrect {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *GameHandler) FinishSession(c *gin.Context) {
	sessionID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	session, err := h.gameService.FinishSession(c.Request.Context(), sessionID)
	if err != nil {
		writeGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *GameHandler) Leaderboard(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	entries, err := h.gameService.Leaderboard(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot load leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": entries})
}

func (h *GameHandler) UserSessions(c *gin.Context) {
	userID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	sessions, err := h.gameService.UserSessions(c.Request.Context(), userID)
	if err != nil {
		writeGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *GameHandler) GetVocabularyForGame(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	vocabulary, err := h.gameService.VocabularyForGame(c.Request.Context(), userID)
	if err != nil {
		writeGameError(c, err)
		return
	}
	if vocabulary == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":    "no vocabulary available",
			"vocabulary": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vocabulary": vocabulary})
}

func (h *GameHandler) SubmitGuessWordAnswer(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	var request models.SubmitGameAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id and word are required"})
		return
	}

	response, err := h.gameService.SubmitGuessWordAnswer(c.Request.Context(), userID, request)
	if err != nil {
		writeGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": name + " must be a positive number"})
		return 0, false
	}

	return uint(value), true
}

func authenticatedUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get(middlewares.UserIDContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user is not authenticated"})
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user token"})
		return 0, false
	}

	return userID, true
}

func writeGameError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	case errors.Is(err, services.ErrGameSessionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": "game session not found"})
	case errors.Is(err, services.ErrVocabularyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": "vocabulary not found"})
	case errors.Is(err, services.ErrGameSessionClosed):
		c.JSON(http.StatusConflict, gin.H{"message": "game session is not in progress"})
	case errors.Is(err, services.ErrInvalidAnswerWord):
		c.JSON(http.StatusBadRequest, gin.H{"message": "word must contain english letters only"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot process game request"})
	}
}

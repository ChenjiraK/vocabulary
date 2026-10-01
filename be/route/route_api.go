package route

import (
	"net/http"

	"vocabulary/handlers"
	"vocabulary/middlewares"
	"vocabulary/repositories"
	"vocabulary/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAPIRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Initialize repositories, services, and handlers ของ Module vocabulary
	vocabularyRepository := repositories.NewVocabularyRepository(db)
	vocabularyService := services.NewVocabularyService(vocabularyRepository)
	vocabularyHandler := handlers.NewVocabularyHandler(vocabularyService)

	// Initialize repositories, services, and handlers ของ Module user, auth
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize repositories, services, and handlers ของ Module game
	gameRepository := repositories.NewGameRepository(db)
	gameService := services.NewGameService(gameRepository, userRepository)
	gameHandler := handlers.NewGameHandler(gameService)

	router.POST("/users/guest", authHandler.CreateGuestUser)
	router.GET("/users/:id/game-sessions", gameHandler.UserSessions)

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)
	router.GET("/auth/facebook/login", authHandler.FacebookLogin)
	router.GET("/auth/facebook/callback", authHandler.FacebookCallback)

	router.POST("/game-sessions", gameHandler.StartSession)             //รอบที่ 1-10
	router.POST("/game-sessions/:id/answers", gameHandler.SubmitAnswer) //คำตอบของรอบที่ 1-10
	router.POST("/game-sessions/:id/finish", gameHandler.FinishSession) //จบเกม รอบที่ 1-10
	router.GET("/leaderboard", gameHandler.Leaderboard)

	router.POST("/vocabularies/upload", vocabularyHandler.UploadExcel)
	router.GET("/vocabulary/guess-word", middlewares.AuthMiddleware(), gameHandler.GetVocabularyForGame)
	router.POST("/vocabulary/guess-word/answer", middlewares.AuthMiddleware(), gameHandler.SubmitGuessWordAnswer)
}

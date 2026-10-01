package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"vocabulary/models"
	"vocabulary/repositories"
)

const pointsPerCorrectAnswer = 10

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrGameSessionNotFound    = errors.New("game session not found")
	ErrGameSessionClosed      = errors.New("game session is not in progress")
	ErrVocabularyNotFound     = errors.New("vocabulary not found")
	ErrInvalidAnswerWord      = errors.New("word must contain english letters only")
	ErrInvalidVocabularyLevel = errors.New("invalid vocabulary level")
)

type GameService interface {
	StartSession(ctx context.Context, userID uint) (*models.GameSession, error)
	SubmitAnswer(ctx context.Context, sessionID uint, request models.SubmitGameAnswerRequest) (*models.GameAnswerResponse, error)
	FinishSession(ctx context.Context, sessionID uint) (*models.GameSession, error)
	Leaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error)
	UserSessions(ctx context.Context, userID uint) ([]models.GameSession, error)
	VocabularyForGame(ctx context.Context, userID uint) (*models.GameVocabularyResponse, error)
	SubmitGuessWordAnswer(ctx context.Context, userID uint, request models.SubmitGameAnswerRequest) (*models.GuessWordAnswerResponse, error)
}

type gameService struct {
	gameRepository       repositories.GameRepository
	userRepository       repositories.UserRepository
	vocabularyRepository repositories.VocabularyRepository
}

func NewGameService(
	gameRepository repositories.GameRepository,
	userRepository repositories.UserRepository,
	vocabularyRepositories ...repositories.VocabularyRepository,
) GameService {
	var vocabularyRepository repositories.VocabularyRepository
	if len(vocabularyRepositories) > 0 {
		vocabularyRepository = vocabularyRepositories[0]
	}

	return &gameService{
		gameRepository:       gameRepository,
		userRepository:       userRepository,
		vocabularyRepository: vocabularyRepository,
	}
}

func (s *gameService) StartSession(ctx context.Context, userID uint) (*models.GameSession, error) {
	user, err := s.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	session := &models.GameSession{
		UserID:    userID,
		Status:    models.GameSessionStatusInProgress,
		StartedAt: time.Now(),
	}
	if err := s.gameRepository.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("create game session: %w", err)
	}

	return session, nil
}

func (s *gameService) VocabularyForGame(ctx context.Context, userID uint) (*models.GameVocabularyResponse, error) {
	user, err := s.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.gameRepository.VocabGame(ctx, userID)
}

func (s *gameService) SubmitAnswer(ctx context.Context, sessionID uint, request models.SubmitGameAnswerRequest) (*models.GameAnswerResponse, error) {
	if !isEnglishLettersOnly(request.Word) {
		return nil, ErrInvalidAnswerWord
	}

	session, err := s.gameRepository.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("find game session: %w", err)
	}
	if session == nil {
		return nil, ErrGameSessionNotFound
	}
	if session.Status != models.GameSessionStatusInProgress {
		return nil, ErrGameSessionClosed
	}

	vocabulary, err := s.gameRepository.FindVocabularyByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("find vocabulary: %w", err)
	}
	if vocabulary == nil {
		return nil, ErrVocabularyNotFound
	}

	isCorrect := normalizeAnswer(request.Word) == normalizeAnswer(vocabulary.Eng)
	answer := &models.GameAnswer{
		GameSessionID: session.ID,
		VocabularyID:  vocabulary.ID,
		QuestionText:  vocabulary.Thai,
		CorrectAnswer: vocabulary.Eng,
		UserAnswer:    strings.TrimSpace(request.Word),
		IsCorrect:     isCorrect,
		AnsweredAt:    time.Now(),
	}
	if err := s.gameRepository.CreateAnswer(ctx, answer); err != nil {
		return nil, fmt.Errorf("create game answer: %w", err)
	}

	session.TotalQuestions++
	if isCorrect {
		session.CorrectAnswers++
		session.Score += pointsPerCorrectAnswer
	}
	if err := s.gameRepository.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("update game session: %w", err)
	}

	return &models.GameAnswerResponse{
		Answer:  *answer,
		Session: *session,
	}, nil
}

func (s *gameService) FinishSession(ctx context.Context, sessionID uint) (*models.GameSession, error) {
	session, err := s.gameRepository.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("find game session: %w", err)
	}
	if session == nil {
		return nil, ErrGameSessionNotFound
	}
	if session.Status != models.GameSessionStatusInProgress {
		return session, nil
	}

	now := time.Now()
	session.Status = models.GameSessionStatusCompleted
	session.FinishedAt = &now
	if err := s.gameRepository.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("finish game session: %w", err)
	}

	return session, nil
}

func (s *gameService) Leaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.gameRepository.ListLeaderboard(ctx, limit)
}

func (s *gameService) UserSessions(ctx context.Context, userID uint) ([]models.GameSession, error) {
	user, err := s.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.gameRepository.ListSessionsByUserID(ctx, userID)
}

func normalizeAnswer(answer string) string {
	return strings.ToLower(strings.TrimSpace(answer))
}

func isEnglishLettersOnly(word string) bool {
	word = strings.TrimSpace(word)
	if word == "" {
		return false
	}

	for _, character := range word {
		if (character < 'A' || character > 'Z') && (character < 'a' || character > 'z') {
			return false
		}
	}

	return true
}

func scoreByVocabularyLevel(level string) (int, error) {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "A1":
		return 1, nil
	case "A2":
		return 2, nil
	case "B1":
		return 3, nil
	case "B2":
		return 4, nil
	case "C1", "C2":
		return 5, nil
	default:
		return 0, ErrInvalidVocabularyLevel
	}
}

func (s *gameService) SubmitGuessWordAnswer(
	ctx context.Context,
	userID uint,
	request models.SubmitGameAnswerRequest,
) (*models.GuessWordAnswerResponse, error) {
	if !isEnglishLettersOnly(request.Word) {
		return nil, ErrInvalidAnswerWord
	}

	user, err := s.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	vocabulary, err := s.gameRepository.FindVocabularyByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("find vocabulary: %w", err)
	}
	if vocabulary == nil {
		return nil, ErrVocabularyNotFound
	}

	session, err := s.gameRepository.FindOrCreateInProgressSessionByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find or create session: %w", err)
	}

	answer, err := s.gameRepository.FindUserAnswerByVocabularyID(ctx, userID, vocabulary.ID)
	if err != nil {
		return nil, fmt.Errorf("find existing answer: %w", err)
	}

	if answer != nil && answer.IsCorrect {
		return &models.GuessWordAnswerResponse{
			IsCorrect:  true,
			ScoreAdded: 0,
			Message:    "already answered correctly",
		}, nil
	}

	isCorrect := normalizeAnswer(request.Word) == normalizeAnswer(vocabulary.Eng)
	scoreAdded := 0
	if isCorrect {
		scoreAdded, err = scoreByVocabularyLevel(vocabulary.Level)
		if err != nil {
			return nil, err
		}
	}

	now := time.Now()
	if answer == nil {
		answer = &models.GameAnswer{
			UserID:        userID,
			GameSessionID: session.ID,
			VocabularyID:  vocabulary.ID,
			QuestionText:  vocabulary.Thai,
			CorrectAnswer: vocabulary.Eng,
		}
	}

	wasCorrect := answer.IsCorrect
	answer.UserAnswer = strings.TrimSpace(request.Word)
	answer.IsCorrect = isCorrect
	answer.AnsweredAt = now

	if answer.ID == 0 {
		if err := s.gameRepository.CreateAnswer(ctx, answer); err != nil {
			return nil, fmt.Errorf("create game answer: %w", err)
		}
	} else {
		if err := s.gameRepository.UpdateAnswer(ctx, answer); err != nil {
			return nil, fmt.Errorf("update game answer: %w", err)
		}
	}

	if isCorrect && !wasCorrect {
		session.TotalQuestions++
		session.CorrectAnswers++
		session.Score += scoreAdded
		if err := s.gameRepository.UpdateSession(ctx, session); err != nil {
			return nil, fmt.Errorf("update session: %w", err)
		}
	}

	if !isCorrect {
		return &models.GuessWordAnswerResponse{
			IsCorrect:  false,
			ScoreAdded: 0,
			Message:    "incorrect",
		}, nil
	}

	return &models.GuessWordAnswerResponse{
		IsCorrect:  true,
		ScoreAdded: scoreAdded,
		Message:    "correct",
	}, nil
}

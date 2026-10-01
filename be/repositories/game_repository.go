package repositories

import (
	"context"
	"fmt"
	"time"

	"vocabulary/models"

	"gorm.io/gorm"
)

type GameRepository interface {
	CreateSession(ctx context.Context, session *models.GameSession) error
	FindSessionByID(ctx context.Context, id uint) (*models.GameSession, error)
	UpdateSession(ctx context.Context, session *models.GameSession) error
	CreateAnswer(ctx context.Context, answer *models.GameAnswer) error
	FindVocabularyByID(ctx context.Context, id uint) (*models.Vocabulary, error)
	VocabGame(ctx context.Context, userID uint) (*models.GameVocabularyResponse, error)
	ListLeaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error)
	ListSessionsByUserID(ctx context.Context, userID uint) ([]models.GameSession, error)
	UpdateAnswer(ctx context.Context, answer *models.GameAnswer) error
	FindUserAnswerByVocabularyID(ctx context.Context, userID uint, vocabularyID uint) (*models.GameAnswer, error)
	FindOrCreateInProgressSessionByUserID(ctx context.Context, userID uint) (*models.GameSession, error)
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) CreateSession(ctx context.Context, session *models.GameSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *gameRepository) FindSessionByID(ctx context.Context, id uint) (*models.GameSession, error) {
	var session models.GameSession
	if err := r.db.WithContext(ctx).First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}

func (r *gameRepository) UpdateSession(ctx context.Context, session *models.GameSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *gameRepository) CreateAnswer(ctx context.Context, answer *models.GameAnswer) error {
	if answer.UserID == 0 {
		userID, err := r.findSessionUserID(ctx, answer.GameSessionID)
		if err != nil {
			return fmt.Errorf("find session user: %w", err)
		}
		answer.UserID = userID
	}

	return r.db.WithContext(ctx).Create(answer).Error
}

func (r *gameRepository) UpdateAnswer(ctx context.Context, answer *models.GameAnswer) error {
	return r.db.WithContext(ctx).Save(answer).Error
}

func (r *gameRepository) FindUserAnswerByVocabularyID(ctx context.Context, userID uint, vocabularyID uint) (*models.GameAnswer, error) {
	var answer models.GameAnswer
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND vocabulary_id = ?", userID, vocabularyID).
		Order("answered_at DESC, id DESC").
		First(&answer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &answer, nil
}

func (r *gameRepository) FindOrCreateInProgressSessionByUserID(ctx context.Context, userID uint) (*models.GameSession, error) {
	var session models.GameSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, models.GameSessionStatusInProgress).
		Order("started_at DESC, id DESC").
		First(&session).Error
	if err == nil {
		return &session, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	session = models.GameSession{
		UserID:    userID,
		Status:    models.GameSessionStatusInProgress,
		StartedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *gameRepository) findSessionUserID(ctx context.Context, sessionID uint) (uint, error) {
	var session models.GameSession
	err := r.db.WithContext(ctx).
		Select("user_id").
		First(&session, sessionID).Error
	if err != nil {
		return 0, err
	}

	return session.UserID, nil
}

func (r *gameRepository) FindVocabularyByID(ctx context.Context, id uint) (*models.Vocabulary, error) {
	var vocabulary models.Vocabulary
	if err := r.db.WithContext(ctx).First(&vocabulary, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &vocabulary, nil
}

func (r *gameRepository) VocabGame(ctx context.Context, userID uint) (*models.GameVocabularyResponse, error) {
	var vocabulary models.GameVocabularyResponse

	playedCorrectSubQuery := r.db.
		Table("game_answers AS ga").
		Select("1").
		Where("ga.vocabulary_id = v.id").
		Where("ga.user_id = ?", userID).
		Where("ga.is_correct = ?", true)

	//NOTE: CHAR_LENGTH นับจำนวนตัวอักษรของ field eng
	err := r.db.WithContext(ctx).
		Table("vocabulary AS v").
		Select(`
			v.id,
			CHAR_LENGTH(v.eng) AS length_of_word,
			v.parts_of_speech,
			v.thai,
			v.meaning,
			v.synonyms
		`).
		Where("NOT EXISTS (?)", playedCorrectSubQuery).
		Order("v.level ASC, v.id ASC").
		Limit(1).
		Scan(&vocabulary).Error
	if err != nil {
		return nil, err
	}
	if vocabulary.ID == 0 {
		return nil, nil
	}

	return &vocabulary, err
}

func (r *gameRepository) ListLeaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry
	err := r.db.WithContext(ctx).
		Table("game_sessions AS gs").
		Select(`
			gs.id AS game_session_id,
			u.id AS user_id,
			u.display_name,
			u.avatar_url,
			gs.score,
			gs.correct_answers,
			gs.total_questions,
			gs.finished_at,
			gs.started_at
		`).
		Joins("JOIN users AS u ON u.id = gs.user_id").
		Where("gs.status = ?", models.GameSessionStatusCompleted).
		Order("gs.score DESC, gs.finished_at ASC NULLS LAST, gs.id ASC").
		Limit(limit).
		Scan(&entries).Error

	return entries, err
}

func (r *gameRepository) ListSessionsByUserID(ctx context.Context, userID uint) ([]models.GameSession, error) {
	var sessions []models.GameSession
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("started_at DESC, id DESC").
		Find(&sessions).Error

	return sessions, err
}

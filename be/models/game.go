package models

import "time"

const (
	GameSessionStatusInProgress = "in_progress"
	GameSessionStatusCompleted  = "completed"
	GameSessionStatusAbandoned  = "abandoned"
)

type GameSession struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"column:user_id;not null" json:"user_id"`
	User           User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalQuestions int        `gorm:"column:total_questions;not null;default:0" json:"total_questions"`
	CorrectAnswers int        `gorm:"column:correct_answers;not null;default:0" json:"correct_answers"`
	Score          int        `gorm:"column:score;not null;default:0" json:"score"`
	Status         string     `gorm:"column:status;size:20;not null;default:in_progress" json:"status"`
	StartedAt      time.Time  `gorm:"column:started_at;not null" json:"started_at"`
	FinishedAt     *time.Time `gorm:"column:finished_at" json:"finished_at,omitempty"`
}

func (GameSession) TableName() string {
	return "game_sessions"
}

type GameAnswer struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"column:user_id;not null" json:"user_id"`
	GameSessionID uint      `gorm:"column:game_session_id;not null" json:"game_session_id"`
	VocabularyID  uint      `gorm:"column:vocabulary_id;not null" json:"vocabulary_id"`
	QuestionText  string    `gorm:"column:question_text;type:text;not null" json:"question_text"`
	CorrectAnswer string    `gorm:"column:correct_answer;type:text;not null" json:"correct_answer"`
	UserAnswer    string    `gorm:"column:user_answer;type:text" json:"user_answer"`
	IsCorrect     bool      `gorm:"column:is_correct;not null;default:false" json:"is_correct"`
	AnsweredAt    time.Time `gorm:"column:answered_at;not null" json:"answered_at"`
}

type GuessWordAnswerResponse struct {
	IsCorrect  bool   `json:"is_correct"`
	ScoreAdded int    `json:"score_added"`
	Message    string `json:"message"`
}

func (GameAnswer) TableName() string {
	return "game_answers"
}

type StartGameSessionRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

type SubmitGameAnswerRequest struct {
	ID   uint   `json:"id" binding:"required"`
	Word string `json:"word" binding:"required"`
}

type GameAnswerResponse struct {
	Answer  GameAnswer  `json:"answer"`
	Session GameSession `json:"session"`
}

type LeaderboardEntry struct {
	GameSessionID  uint       `json:"game_session_id"`
	UserID         uint       `json:"user_id"`
	DisplayName    string     `json:"display_name"`
	AvatarURL      *string    `json:"avatar_url,omitempty"`
	Score          int        `json:"score"`
	CorrectAnswers int        `json:"correct_answers"`
	TotalQuestions int        `json:"total_questions"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
	StartedAt      time.Time  `json:"started_at"`
}

// vocabulary response สำหรับ api /game
type GameVocabularyResponse struct {
	ID            uint   `json:"id"`
	LengthOfWord  int    `json:"length_of_word"` //length of field eng
	PartsOfSpeech string `json:"parts_of_speech"`
	Thai          string `json:"thai"`
	Meaning       string `json:"meaning"`
	Synonyms      string `json:"synonyms"`
}

package services

import (
	"context"
	"testing"
	"time"

	"vocabulary/models"
)

func TestGuestUserCanStartSessionAndSubmitCorrectAnswer(t *testing.T) {
	userRepository := &fakeUserRepository{
		users: []models.User{
			{
				ID:          1,
				DisplayName: "Guest",
				UserType:    models.UserTypeGuest,
			},
		},
	}
	gameRepository := &fakeGameRepository{
		vocabularies: []models.Vocabulary{
			{
				ID:   10,
				Eng:  "apple",
				Thai: "แอปเปิล",
			},
		},
	}
	service := NewGameService(gameRepository, userRepository)

	session, err := service.StartSession(context.Background(), 1)
	if err != nil {
		t.Fatalf("start session: %v", err)
	}

	response, err := service.SubmitAnswer(context.Background(), session.ID, models.SubmitGameAnswerRequest{
		ID:   10,
		Word: "apple",
	})
	if err != nil {
		t.Fatalf("submit answer: %v", err)
	}

	if !response.Answer.IsCorrect {
		t.Fatal("expected answer to be correct")
	}
	if response.Session.Score != pointsPerCorrectAnswer {
		t.Fatalf("expected score %d, got %d", pointsPerCorrectAnswer, response.Session.Score)
	}
	if response.Session.TotalQuestions != 1 {
		t.Fatalf("expected total questions 1, got %d", response.Session.TotalQuestions)
	}
}

func TestSubmitAnswerRecordsIncorrectAnswer(t *testing.T) {
	gameRepository := &fakeGameRepository{
		sessions: []models.GameSession{
			{
				ID:     1,
				UserID: 1,
				Status: models.GameSessionStatusInProgress,
			},
		},
		vocabularies: []models.Vocabulary{
			{
				ID:   10,
				Eng:  "apple",
				Thai: "แอปเปิล",
			},
		},
	}
	service := NewGameService(gameRepository, &fakeUserRepository{})

	response, err := service.SubmitAnswer(context.Background(), 1, models.SubmitGameAnswerRequest{
		ID:   10,
		Word: "banana",
	})
	if err != nil {
		t.Fatalf("submit answer: %v", err)
	}

	if response.Answer.IsCorrect {
		t.Fatal("expected answer to be incorrect")
	}
	if len(gameRepository.answers) != 1 {
		t.Fatalf("expected one recorded answer, got %d", len(gameRepository.answers))
	}
	if gameRepository.answers[0].IsCorrect {
		t.Fatal("expected recorded answer is_correct false")
	}
	if response.Session.Score != 0 {
		t.Fatalf("expected score 0, got %d", response.Session.Score)
	}
	if response.Session.TotalQuestions != 1 {
		t.Fatalf("expected total questions 1, got %d", response.Session.TotalQuestions)
	}
}

func TestSubmitAnswerRejectsNonEnglishWord(t *testing.T) {
	service := NewGameService(&fakeGameRepository{}, &fakeUserRepository{})

	_, err := service.SubmitAnswer(context.Background(), 1, models.SubmitGameAnswerRequest{
		ID:   10,
		Word: "แอปเปิล",
	})

	if err != ErrInvalidAnswerWord {
		t.Fatalf("expected invalid answer word error, got %v", err)
	}
}

func TestSubmitAnswerRejectsCompletedSession(t *testing.T) {
	userRepository := &fakeUserRepository{}
	gameRepository := &fakeGameRepository{
		sessions: []models.GameSession{
			{
				ID:     1,
				UserID: 1,
				Status: models.GameSessionStatusCompleted,
			},
		},
	}
	service := NewGameService(gameRepository, userRepository)

	_, err := service.SubmitAnswer(context.Background(), 1, models.SubmitGameAnswerRequest{
		ID:   10,
		Word: "answer",
	})

	if err != ErrGameSessionClosed {
		t.Fatalf("expected closed session error, got %v", err)
	}
}

func TestFinishSessionMarksSessionCompleted(t *testing.T) {
	service := NewGameService(&fakeGameRepository{
		sessions: []models.GameSession{
			{
				ID:     1,
				UserID: 1,
				Status: models.GameSessionStatusInProgress,
			},
		},
	}, &fakeUserRepository{})

	session, err := service.FinishSession(context.Background(), 1)
	if err != nil {
		t.Fatalf("finish session: %v", err)
	}

	if session.Status != models.GameSessionStatusCompleted {
		t.Fatalf("expected completed status, got %q", session.Status)
	}
	if session.FinishedAt == nil {
		t.Fatal("expected finished_at")
	}
}

func TestLeaderboardDefaultsInvalidLimitToTen(t *testing.T) {
	gameRepository := &fakeGameRepository{}
	service := NewGameService(gameRepository, &fakeUserRepository{})

	_, err := service.Leaderboard(context.Background(), -1)
	if err != nil {
		t.Fatalf("leaderboard: %v", err)
	}

	if gameRepository.lastLeaderboardLimit != 10 {
		t.Fatalf("expected default limit 10, got %d", gameRepository.lastLeaderboardLimit)
	}
}

type fakeGameRepository struct {
	sessions             []models.GameSession
	answers              []models.GameAnswer
	vocabularies         []models.Vocabulary
	leaderboard          []models.LeaderboardEntry
	lastLeaderboardLimit int
	nextSessionID        uint
}

func (r *fakeGameRepository) CreateSession(ctx context.Context, session *models.GameSession) error {
	session.ID = r.nextID()
	r.sessions = append(r.sessions, *session)
	return nil
}

func (r *fakeGameRepository) FindSessionByID(ctx context.Context, id uint) (*models.GameSession, error) {
	for index := range r.sessions {
		if r.sessions[index].ID == id {
			return &r.sessions[index], nil
		}
	}

	return nil, nil
}

func (r *fakeGameRepository) UpdateSession(ctx context.Context, session *models.GameSession) error {
	for index := range r.sessions {
		if r.sessions[index].ID == session.ID {
			r.sessions[index] = *session
			return nil
		}
	}

	r.sessions = append(r.sessions, *session)
	return nil
}

func (r *fakeGameRepository) CreateAnswer(ctx context.Context, answer *models.GameAnswer) error {
	answer.ID = uint(len(r.answers) + 1)
	r.answers = append(r.answers, *answer)
	return nil
}

func (r *fakeGameRepository) UpdateAnswer(ctx context.Context, answer *models.GameAnswer) error {
	for index := range r.answers {
		if r.answers[index].ID == answer.ID {
			r.answers[index] = *answer
			return nil
		}
	}

	r.answers = append(r.answers, *answer)
	return nil
}

func (r *fakeGameRepository) FindUserAnswerByVocabularyID(ctx context.Context, userID uint, vocabularyID uint) (*models.GameAnswer, error) {
	for index := range r.answers {
		if r.answers[index].UserID == userID && r.answers[index].VocabularyID == vocabularyID {
			return &r.answers[index], nil
		}
	}

	return nil, nil
}

func (r *fakeGameRepository) FindOrCreateInProgressSessionByUserID(ctx context.Context, userID uint) (*models.GameSession, error) {
	for index := range r.sessions {
		if r.sessions[index].UserID == userID && r.sessions[index].Status == models.GameSessionStatusInProgress {
			return &r.sessions[index], nil
		}
	}

	session := models.GameSession{
		ID:        r.nextID(),
		UserID:    userID,
		Status:    models.GameSessionStatusInProgress,
		StartedAt: time.Now(),
	}
	r.sessions = append(r.sessions, session)

	return &r.sessions[len(r.sessions)-1], nil
}

func (r *fakeGameRepository) FindVocabularyByID(ctx context.Context, id uint) (*models.Vocabulary, error) {
	for index := range r.vocabularies {
		if r.vocabularies[index].ID == id {
			return &r.vocabularies[index], nil
		}
	}

	return nil, nil
}

func (r *fakeGameRepository) VocabGame(ctx context.Context, userID uint) (*models.GameVocabularyResponse, error) {
	return nil, nil
}

func (r *fakeGameRepository) ListLeaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error) {
	r.lastLeaderboardLimit = limit
	if len(r.leaderboard) > limit {
		return r.leaderboard[:limit], nil
	}

	return r.leaderboard, nil
}

func (r *fakeGameRepository) ListSessionsByUserID(ctx context.Context, userID uint) ([]models.GameSession, error) {
	var sessions []models.GameSession
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (r *fakeGameRepository) nextID() uint {
	if r.nextSessionID == 0 {
		r.nextSessionID = 1
	}

	id := r.nextSessionID
	r.nextSessionID++
	return id
}

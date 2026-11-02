package repositories

import (
	"context"
	"errors"

	"vocabulary/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrDuplicateIdentity = errors.New("identity already exists")

type UserRepository interface {
	FindEmailByUserID(ctx context.Context, userID uint) (*string, error)
	GetUserGameStats(ctx context.Context, userID uint) (models.UserGameStats, error)
	CreateUser(ctx context.Context, user *models.User) error
	CreateUserWithIdentity(ctx context.Context, user *models.User, identity *models.UserAuthIdentity) error
	CreateIdentity(ctx context.Context, identity *models.UserAuthIdentity) error
	FindUserByID(ctx context.Context, id uint) (*models.User, error)
	FindIdentityByProviderAndEmail(ctx context.Context, provider string, email string) (*models.UserAuthIdentity, error)
	FindIdentityByProviderAndProviderUserID(ctx context.Context, provider string, providerUserID string) (*models.UserAuthIdentity, error)
	FindIdentityByEmail(ctx context.Context, email string) (*models.UserAuthIdentity, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) CreateUserWithIdentity(ctx context.Context, user *models.User, identity *models.UserAuthIdentity) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		identity.UserID = user.ID
		return tx.Create(identity).Error
	})
	var postgresError *pgconn.PgError
	if errors.Is(err, gorm.ErrDuplicatedKey) || (errors.As(err, &postgresError) && postgresError.Code == "23505" && postgresError.ConstraintName == "user_auth_identities_email_provider_unique") {
		return ErrDuplicateIdentity
	}
	return err
}

func (r *userRepository) FindEmailByUserID(ctx context.Context, userID uint) (*string, error) {
	var identity models.UserAuthIdentity
	err := r.db.WithContext(ctx).
		Select("email").
		Where("user_id = ? AND email IS NOT NULL", userID).
		Order("CASE WHEN provider = 'password' THEN 0 ELSE 1 END, email_verified DESC, id ASC").
		Take(&identity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return identity.Email, nil
}

func (r *userRepository) CreateIdentity(ctx context.Context, identity *models.UserAuthIdentity) error {
	return r.db.WithContext(ctx).Create(identity).Error
}

func (r *userRepository) FindUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindIdentityByProviderAndEmail(ctx context.Context, provider string, email string) (*models.UserAuthIdentity, error) {
	var identity models.UserAuthIdentity
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("provider = ? AND email = ?", provider, email).
		First(&identity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &identity, nil
}

func (r *userRepository) FindIdentityByProviderAndProviderUserID(ctx context.Context, provider string, providerUserID string) (*models.UserAuthIdentity, error) {
	var identity models.UserAuthIdentity
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&identity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &identity, nil
}

func (r *userRepository) FindIdentityByEmail(ctx context.Context, email string) (*models.UserAuthIdentity, error) {
	var identity models.UserAuthIdentity
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("email = ?", email).
		Order("id ASC").
		First(&identity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &identity, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) GetUserGameStats(ctx context.Context, userID uint) (models.UserGameStats, error) {
	var stats models.UserGameStats
	err := r.db.WithContext(ctx).Model(&models.GameSession{}).
		Select("COALESCE(SUM(correct_answers), 0) AS correct_answers, COALESCE(SUM(score), 0) AS score").
		Where("user_id = ?", userID).
		Scan(&stats).Error
	return stats, err
}

package repositories

import (
	"context"

	"vocabulary/models"

	"gorm.io/gorm"
)

type UserRepository interface {
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		identity.UserID = user.ID
		return tx.Create(identity).Error
	})
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

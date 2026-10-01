package services

import (
	"context"
	"testing"

	"vocabulary/models"

	"golang.org/x/crypto/bcrypt"
)

func TestRegisterReturnsDuplicateEmailForExistingPasswordIdentity(t *testing.T) {
	email := "player@example.com"
	repository := &fakeUserRepository{
		identities: []models.UserAuthIdentity{
			{
				UserID:   1,
				Provider: models.AuthProviderPassword,
				Email:    &email,
				User: models.User{
					ID:          1,
					DisplayName: "Player",
					UserType:    models.UserTypeRegistered,
				},
			},
		},
	}
	service := NewAuthService(repository)

	_, err := service.Register(context.Background(), models.RegisterRequest{
		DisplayName: "Player",
		Email:       email,
		Password:    "password123",
	})

	if err != ErrDuplicateEmail {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}

func TestRegisterLinksPasswordIdentityToExistingSocialUser(t *testing.T) {
	email := "player@example.com"
	repository := &fakeUserRepository{
		users: []models.User{
			{
				ID:          1,
				DisplayName: "Google Player",
				UserType:    models.UserTypeRegistered,
			},
		},
		identities: []models.UserAuthIdentity{
			{
				UserID:   1,
				Provider: models.AuthProviderGoogle,
				Email:    &email,
				User: models.User{
					ID:          1,
					DisplayName: "Google Player",
					UserType:    models.UserTypeRegistered,
				},
			},
		},
	}
	service := NewAuthService(repository)

	response, err := service.Register(context.Background(), models.RegisterRequest{
		DisplayName: "Password Player",
		Email:       email,
		Password:    "password123",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if response.User.ID != 1 {
		t.Fatalf("expected linked user ID 1, got %d", response.User.ID)
	}
	if len(repository.createdIdentities) != 1 {
		t.Fatalf("expected one created identity, got %d", len(repository.createdIdentities))
	}
	if repository.createdIdentities[0].Provider != models.AuthProviderPassword {
		t.Fatalf("expected password provider, got %q", repository.createdIdentities[0].Provider)
	}
}

func TestLoginWithSocialProfileReusesExistingProviderIdentity(t *testing.T) {
	providerUserID := "google-1"
	email := "player@example.com"
	repository := &fakeUserRepository{
		identities: []models.UserAuthIdentity{
			{
				UserID:         3,
				Provider:       models.AuthProviderGoogle,
				ProviderUserID: &providerUserID,
				Email:          &email,
				User: models.User{
					ID:          3,
					DisplayName: "Existing Player",
					UserType:    models.UserTypeRegistered,
				},
			},
		},
	}
	service := NewAuthService(repository).(*authService)

	response, err := service.loginWithSocialProfile(context.Background(), models.AuthProviderGoogle, socialProfile{
		ID:            providerUserID,
		Email:         email,
		EmailVerified: true,
		Name:          "New Name",
	})
	if err != nil {
		t.Fatalf("social login: %v", err)
	}

	if response.User.ID != 3 {
		t.Fatalf("expected existing user ID 3, got %d", response.User.ID)
	}
	if len(repository.createdUsers) != 0 {
		t.Fatalf("expected no new users, got %d", len(repository.createdUsers))
	}
}

func TestLoginWithFacebookProfileWithoutEmailCreatesUser(t *testing.T) {
	repository := &fakeUserRepository{}
	service := NewAuthService(repository).(*authService)

	response, err := service.loginWithSocialProfile(context.Background(), models.AuthProviderFacebook, socialProfile{
		ID:   "facebook-1",
		Name: "Facebook Player",
	})
	if err != nil {
		t.Fatalf("facebook login: %v", err)
	}

	if response.User.ID == 0 {
		t.Fatal("expected created user ID")
	}
	if len(repository.createdUsersWithIdentity) != 1 {
		t.Fatalf("expected one created user with identity, got %d", len(repository.createdUsersWithIdentity))
	}
	identity := repository.createdUsersWithIdentity[0].identity
	if identity.Provider != models.AuthProviderFacebook {
		t.Fatalf("expected facebook provider, got %q", identity.Provider)
	}
	if identity.Email != nil {
		t.Fatalf("expected nil email for facebook profile without email, got %q", *identity.Email)
	}
}

func TestPasswordLoginValidatesBcryptHash(t *testing.T) {
	email := "player@example.com"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	hash := string(hashBytes)

	repository := &fakeUserRepository{
		identities: []models.UserAuthIdentity{
			{
				UserID:       1,
				Provider:     models.AuthProviderPassword,
				Email:        &email,
				PasswordHash: &hash,
				User: models.User{
					ID:          1,
					DisplayName: "Player",
					UserType:    models.UserTypeRegistered,
				},
			},
		},
	}
	service := NewAuthService(repository)

	response, err := service.Login(context.Background(), models.LoginRequest{
		Email:    email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if response.Token == "" {
		t.Fatal("expected token")
	}

	userID, err := ValidateJWT(response.Token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if userID != 1 {
		t.Fatalf("expected token user ID 1, got %d", userID)
	}
}

type fakeUserRepository struct {
	users                    []models.User
	identities               []models.UserAuthIdentity
	createdUsers             []models.User
	createdIdentities        []models.UserAuthIdentity
	createdUsersWithIdentity []struct {
		user     models.User
		identity models.UserAuthIdentity
	}
	updatedUsers []models.User
	nextUserID   uint
}

func (r *fakeUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.ID = r.nextID()
	r.createdUsers = append(r.createdUsers, *user)
	r.users = append(r.users, *user)
	return nil
}

func (r *fakeUserRepository) CreateUserWithIdentity(ctx context.Context, user *models.User, identity *models.UserAuthIdentity) error {
	user.ID = r.nextID()
	identity.UserID = user.ID
	identity.User = *user

	r.createdUsersWithIdentity = append(r.createdUsersWithIdentity, struct {
		user     models.User
		identity models.UserAuthIdentity
	}{user: *user, identity: *identity})
	r.users = append(r.users, *user)
	r.identities = append(r.identities, *identity)
	return nil
}

func (r *fakeUserRepository) CreateIdentity(ctx context.Context, identity *models.UserAuthIdentity) error {
	identity.ID = uint(len(r.identities) + 1)
	r.createdIdentities = append(r.createdIdentities, *identity)
	r.identities = append(r.identities, *identity)
	return nil
}

func (r *fakeUserRepository) FindUserByID(ctx context.Context, id uint) (*models.User, error) {
	for index := range r.users {
		if r.users[index].ID == id {
			return &r.users[index], nil
		}
	}

	return nil, nil
}

func (r *fakeUserRepository) FindIdentityByProviderAndEmail(ctx context.Context, provider string, email string) (*models.UserAuthIdentity, error) {
	for index := range r.identities {
		identity := &r.identities[index]
		if identity.Provider == provider && identity.Email != nil && *identity.Email == email {
			return identity, nil
		}
	}

	return nil, nil
}

func (r *fakeUserRepository) FindIdentityByProviderAndProviderUserID(ctx context.Context, provider string, providerUserID string) (*models.UserAuthIdentity, error) {
	for index := range r.identities {
		identity := &r.identities[index]
		if identity.Provider == provider && identity.ProviderUserID != nil && *identity.ProviderUserID == providerUserID {
			return identity, nil
		}
	}

	return nil, nil
}

func (r *fakeUserRepository) FindIdentityByEmail(ctx context.Context, email string) (*models.UserAuthIdentity, error) {
	for index := range r.identities {
		identity := &r.identities[index]
		if identity.Email != nil && *identity.Email == email {
			return identity, nil
		}
	}

	return nil, nil
}

func (r *fakeUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	r.updatedUsers = append(r.updatedUsers, *user)
	for index := range r.users {
		if r.users[index].ID == user.ID {
			r.users[index] = *user
			return nil
		}
	}

	r.users = append(r.users, *user)
	return nil
}

func (r *fakeUserRepository) nextID() uint {
	if r.nextUserID == 0 {
		r.nextUserID = 1
	}

	id := r.nextUserID
	r.nextUserID++
	return id
}

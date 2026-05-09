package services

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/wardhana-org/ekstra/backend/internal/auth"
	"github.com/wardhana-org/ekstra/backend/internal/models"
	"github.com/wardhana-org/ekstra/backend/internal/repository"
)

const (
	defaultClientType = "web"
	tokenTypeAccess   = "access"
	tokenTypeRefresh  = "refresh"

	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour

	minPasswordLength = 12
	maxPasswordBytes  = 1024
)

type AuthService struct {
	users    *repository.UserRepository
	sessions *repository.AuthRepository
}

func NewAuthService(users *repository.UserRepository, sessions *repository.AuthRepository) *AuthService {
	return &AuthService{
		users:    users,
		sessions: sessions,
	}
}

type PasswordLoginInput struct {
	Email      string
	Password   string
	ClientType string
	DeviceName *string
	UserAgent  *string
	IPAddress  *string
}

type LoginResult struct {
	User                  *models.User
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type userSessionInput struct {
	User       *models.User
	ClientType string
	DeviceName *string
	UserAgent  *string
	IPAddress  *string
}

func (s *AuthService) LoginWithPassword(ctx context.Context, input PasswordLoginInput) (*LoginResult, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" {
		return nil, ErrEmailRequired
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

	user, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	provider, err := s.users.FindPasswordAuthProviderByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if provider.PasswordHash == nil || *provider.PasswordHash == "" {
		return nil, ErrInvalidCredentials
	}

	if err := auth.VerifyPassword(*provider.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.createUserSession(ctx, userSessionInput{
		User:       user,
		ClientType: input.ClientType,
		DeviceName: input.DeviceName,
		UserAgent:  input.UserAgent,
		IPAddress:  input.IPAddress,
	})
}

func (s *AuthService) createUserSession(ctx context.Context, input userSessionInput) (*LoginResult, error) {
	accessToken, err := auth.GenerateRawToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRawToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	accessTokenExpiresAt := now.Add(accessTokenTTL)
	refreshTokenExpiresAt := now.Add(refreshTokenTTL)
	clientType := strings.TrimSpace(input.ClientType)
	if clientType == "" {
		clientType = defaultClientType
	}

	_, _, err = s.sessions.CreateSessionWithTokens(ctx, repository.CreateSessionWithTokensInput{
		Session: repository.CreateSessionInput{
			UserID:     input.User.ID,
			ClientType: clientType,
			DeviceName: input.DeviceName,
			UserAgent:  input.UserAgent,
			IPAddress:  input.IPAddress,
			ExpiresAt:  refreshTokenExpiresAt,
		},
		Tokens: []repository.CreateTokenInput{
			{
				TokenHash: auth.HashToken(accessToken),
				TokenType: tokenTypeAccess,
				ExpiresAt: accessTokenExpiresAt,
			},
			{
				TokenHash: auth.HashToken(refreshToken),
				TokenType: tokenTypeRefresh,
				ExpiresAt: refreshTokenExpiresAt,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:                  input.User,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

type RegisterInput struct {
	Email      string
	Username   string
	Password   string
	ClientType string
	DeviceName *string
	UserAgent  *string
	IPAddress  *string
}

type RegisterResult struct {
	User                  *models.User
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func validatePasswordStrength(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if len(password) > maxPasswordBytes {
		return ErrPasswordTooLong
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*RegisterResult, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" {
		return nil, ErrEmailRequired
	}

	if input.Username == "" {
		return nil, ErrUsernameRequired
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

	// check credential availability
	credentialCheck, err := s.users.CheckCredentialAvailability(ctx, input.Email, input.Username)
	if err != nil {
		return nil, err
	}

	if !credentialCheck.EmailAvailable {
		return nil, ErrRegisterExistingEmail
	}

	if !credentialCheck.UsernameAvailable {
		return nil, ErrRegisterExistingUsername
	}

	if err := validatePasswordStrength(input.Password); err != nil {
		return nil, err
	}
}

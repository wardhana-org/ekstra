package services

import (
	"context"
	"errors"
	"strings"
	"time"

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

type LoginInput struct {
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

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
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
			UserID:     user.ID,
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
		User:                  user,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

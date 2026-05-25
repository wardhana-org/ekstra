package services

import (
	"context"
	"strings"
	"time"

	"github.com/wardhana-org/ekstra/backend/internal/auth"
	"github.com/wardhana-org/ekstra/backend/internal/repository"
)

func (s *AuthService) createUserSession(ctx context.Context, input userSessionInput) (*AuthSessionResult, error) {
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
	absoluteSessionExpiresAt := now.Add(absoluteSessionTTL)
	clientType := strings.TrimSpace(input.ClientType)
	if clientType == "" {
		clientType = defaultClientType
	}

	_, _, err = s.sessions.CreateSessionWithTokens(ctx, repository.CreateSessionWithTokensInput{
		Session: repository.CreateSessionInput{
			UserID:            input.User.ID,
			ClientType:        clientType,
			DeviceName:        input.DeviceName,
			UserAgent:         input.UserAgent,
			IPAddress:         input.IPAddress,
			ExpiresAt:         refreshTokenExpiresAt,
			AbsoluteExpiresAt: absoluteSessionExpiresAt,
		},
		Tokens: []repository.CreateTokenInput{
			{
				TokenHash: auth.HashToken(accessToken),
				TokenType: tokenTypeAccess,
				ExpiresAt: accessTokenExpiresAt,
			},
			{
				TokenHash:  auth.HashToken(refreshToken),
				TokenType:  tokenTypeRefresh,
				ExpiresAt:  refreshTokenExpiresAt,
				Generation: intPointer(1),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &AuthSessionResult{
		User:                  input.User,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func intPointer(value int) *int {
	return &value
}

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

func (s *AuthService) AuthenticateAccessToken(ctx context.Context, rawToken string) (*models.User, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return nil, ErrUnauthenticated
	}

	user, err := s.sessions.FindUserByAccessTokenHash(ctx, auth.HashToken(rawToken), time.Now().UTC())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUnauthenticated
		}

		return nil, err
	}

	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	rawRefreshToken = strings.TrimSpace(rawRefreshToken)
	if rawRefreshToken == "" {
		return nil
	}

	err := s.sessions.RevokeSessionByRefreshTokenHash(ctx, auth.HashToken(rawRefreshToken), time.Now().UTC())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}

		return err
	}

	return nil
}

func (s *AuthService) Refresh(ctx context.Context, input RefreshInput) (*RefreshResult, error) {
	rawRefreshToken := strings.TrimSpace(input.RefreshToken)
	if rawRefreshToken == "" {
		return nil, ErrUnauthenticated
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

	rotation, err := s.sessions.RotateRefreshToken(ctx, repository.RotateRefreshTokenInput{
		RefreshTokenHash:      auth.HashToken(rawRefreshToken),
		NewAccessTokenHash:    auth.HashToken(accessToken),
		NewRefreshTokenHash:   auth.HashToken(refreshToken),
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		RotatedAt:             now,
		ReuseGracePeriod:      refreshReuseGracePeriod,
		UserAgent:             input.UserAgent,
		IPAddress:             input.IPAddress,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenRace) {
			return nil, ErrRefreshTokenRace
		}

		if errors.Is(err, repository.ErrNotFound) ||
			errors.Is(err, repository.ErrRefreshTokenReused) {
			return nil, ErrUnauthenticated
		}

		return nil, err
	}

	return &RefreshResult{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  rotation.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: rotation.RefreshTokenExpiresAt,
	}, nil
}

package services

import (
	"context"
	"errors"
	"strings"

	"github.com/wardhana-org/ekstra/backend/internal/auth"
	"github.com/wardhana-org/ekstra/backend/internal/repository"
)

func (s *AuthService) LoginWithPassword(ctx context.Context, input PasswordLoginInput) (*AuthSessionResult, error) {
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

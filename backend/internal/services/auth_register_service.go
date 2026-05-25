package services

import (
	"context"
	"errors"
	"strings"

	"github.com/wardhana-org/ekstra/backend/internal/auth"
	"github.com/wardhana-org/ekstra/backend/internal/repository"
)

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthSessionResult, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Username = strings.TrimSpace(input.Username)
	if input.Email == "" {
		return nil, ErrEmailRequired
	}

	if input.Username == "" {
		return nil, ErrUsernameRequired
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

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

	hashPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.users.CreatePasswordUser(ctx, repository.CreatePasswordUserInput{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: hashPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrRegisterExistingEmail
		}

		if errors.Is(err, repository.ErrDuplicateUsername) {
			return nil, ErrRegisterExistingUsername
		}

		return nil, err
	}

	return s.createUserSession(ctx, userSessionInput{
		User:       user,
		ClientType: input.ClientType,
		DeviceName: input.DeviceName,
		UserAgent:  input.UserAgent,
		IPAddress:  input.IPAddress,
	})
}

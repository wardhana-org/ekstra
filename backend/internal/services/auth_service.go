package services

import "github.com/wardhana-org/ekstra/backend/internal/repository"

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

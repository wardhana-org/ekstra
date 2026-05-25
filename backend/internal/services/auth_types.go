package services

import (
	"time"

	"github.com/wardhana-org/ekstra/backend/internal/models"
)

type PasswordLoginInput struct {
	Email      string
	Password   string
	ClientType string
	DeviceName *string
	UserAgent  *string
	IPAddress  *string
}

type AuthSessionResult struct {
	User                  *models.User
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
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

type RefreshInput struct {
	RefreshToken string
	UserAgent    *string
	IPAddress    *string
}

type RefreshResult struct {
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

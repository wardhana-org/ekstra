package services

import "time"

const (
	defaultClientType = "web"
	tokenTypeAccess   = "access"
	tokenTypeRefresh  = "refresh"

	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour

	absoluteSessionTTL      = 90 * 24 * time.Hour
	refreshReuseGracePeriod = 5 * time.Second

	minPasswordLength = 12
	maxPasswordBytes  = 1024
)

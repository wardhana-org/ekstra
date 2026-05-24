package repository

const (
	authTokenTypeAccess  = "access"
	authTokenTypeRefresh = "refresh"

	authSessionRevokedReasonLogout            = "logout"
	authSessionRevokedReasonRefreshTokenReuse = "refresh_token_reuse"

	authTokenRevokedReasonLogout            = "logout"
	authTokenRevokedReasonRotated           = "rotated"
	authTokenRevokedReasonRefreshTokenReuse = "refresh_token_reuse"

	authSecurityEventRefreshTokenRace   = "refresh_token_race"
	authSecurityEventRefreshTokenReused = "refresh_token_reused"
)

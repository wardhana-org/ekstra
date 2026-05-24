package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wardhana-org/ekstra/backend/internal/models"
	"github.com/wardhana-org/ekstra/backend/internal/services"
)

const (
	defaultAccessTokenCookieName  = "access_token"
	defaultRefreshTokenCookieName = "refresh_token"
)

type AuthCookieConfig struct {
	AccessTokenName  string
	RefreshTokenName string
	Secure           bool
	SameSite         http.SameSite
}

func DefaultAuthCookieConfig() AuthCookieConfig {
	return AuthCookieConfig{
		AccessTokenName:  defaultAccessTokenCookieName,
		RefreshTokenName: defaultRefreshTokenCookieName,
		Secure:           false,
		SameSite:         http.SameSiteLaxMode,
	}
}

type WebAuthHandler struct {
	authService *services.AuthService
	cookies     AuthCookieConfig
}

func NewWebAuthHandler(authService *services.AuthService, cookies AuthCookieConfig) *WebAuthHandler {
	if cookies.AccessTokenName == "" {
		cookies.AccessTokenName = defaultAccessTokenCookieName
	}
	if cookies.RefreshTokenName == "" {
		cookies.RefreshTokenName = defaultRefreshTokenCookieName
	}
	if cookies.SameSite == 0 {
		cookies.SameSite = http.SameSiteLaxMode
	}

	return &WebAuthHandler{
		authService: authService,
		cookies:     cookies,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User UserResponse `json:"user"`
}

type MeResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

func (h *WebAuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeInvalidRequest(c)
		return
	}

	userAgent, clientIP := requestMetadata(c)

	result, err := h.authService.LoginWithPassword(
		c.Request.Context(),
		services.PasswordLoginInput{
			Email:      req.Email,
			Password:   req.Password,
			ClientType: "web",
			UserAgent:  userAgent,
			IPAddress:  clientIP,
		},
	)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			writeError(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		if errors.Is(err, services.ErrEmailRequired) ||
			errors.Is(err, services.ErrPasswordRequired) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}

		writeInternalServerError(c)
		return
	}

	h.setAuthCookies(c, result.AccessToken, result.AccessTokenExpiresAt, result.RefreshToken, result.RefreshTokenExpiresAt)

	c.JSON(http.StatusOK, LoginResponse{
		User: toUserResponse(result.User),
	})
}

func (h *WebAuthHandler) setAuthCookies(c *gin.Context, accessToken string, accessExpiresAt time.Time, refreshToken string, refreshExpiresAt time.Time) {
	h.setAuthCookie(c, h.cookies.AccessTokenName, accessToken, accessExpiresAt)
	h.setAuthCookie(c, h.cookies.RefreshTokenName, refreshToken, refreshExpiresAt)
}

func (h *WebAuthHandler) setAuthCookie(c *gin.Context, name string, value string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   h.cookies.Secure,
		SameSite: h.cookies.SameSite,
	})
}

func (h *WebAuthHandler) clearAuthCookies(c *gin.Context) {
	h.clearAuthCookie(c, h.cookies.AccessTokenName)
	h.clearAuthCookie(c, h.cookies.RefreshTokenName)
}

func (h *WebAuthHandler) clearAuthCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookies.Secure,
		SameSite: h.cookies.SameSite,
	})
}

func (h *WebAuthHandler) Me(c *gin.Context) {
	accessToken, err := c.Cookie(h.cookies.AccessTokenName)
	if err != nil || accessToken == "" {
		writeUnauthenticated(c)
		return
	}

	user, err := h.authService.AuthenticateAccessToken(c.Request.Context(), accessToken)
	if err != nil {
		if errors.Is(err, services.ErrUnauthenticated) {
			writeUnauthenticated(c)
			return
		}

		writeInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, MeResponse{
		User: toUserResponse(user),
	})
}

func (h *WebAuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie(h.cookies.RefreshTokenName)
	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		writeInternalServerError(c)
		return
	}

	h.clearAuthCookies(c)

	c.Status(http.StatusNoContent)
}

func (h *WebAuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(h.cookies.RefreshTokenName)
	if err != nil || refreshToken == "" {
		h.clearAuthCookies(c)
		writeUnauthenticated(c)
		return
	}

	userAgent, clientIP := requestMetadata(c)

	result, err := h.authService.Refresh(c.Request.Context(), services.RefreshInput{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    clientIP,
	})
	if err != nil {
		if errors.Is(err, services.ErrRefreshTokenRace) {
			// Another request already rotated this refresh token. The raw replacement
			// token is not recoverable, so the client must verify auth state again.
			writeError(c, http.StatusConflict, "refresh conflict")
			return
		}

		if errors.Is(err, services.ErrUnauthenticated) {
			h.clearAuthCookies(c)
			writeUnauthenticated(c)
			return
		}

		writeInternalServerError(c)
		return
	}

	h.setAuthCookies(c, result.AccessToken, result.AccessTokenExpiresAt, result.RefreshToken, result.RefreshTokenExpiresAt)

	c.Status(http.StatusNoContent)
}

func toUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Status:   user.Status,
	}
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func requestMetadata(c *gin.Context) (*string, *string) {
	return optionalString(c.GetHeader("User-Agent")), optionalString(c.ClientIP())
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

func (h *WebAuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeInvalidRequest(c)
		return
	}

	userAgent, clientIP := requestMetadata(c)

	result, err := h.authService.Register(
		c.Request.Context(),
		services.RegisterInput{
			Email:      req.Email,
			Username:   req.Username,
			Password:   req.Password,
			ClientType: "web",
			UserAgent:  userAgent,
			IPAddress:  clientIP,
		},
	)
	if err != nil {
		if errors.Is(err, services.ErrRegisterExistingEmail) ||
			errors.Is(err, services.ErrRegisterExistingUsername) {
			writeError(c, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, services.ErrEmailRequired) ||
			errors.Is(err, services.ErrUsernameRequired) ||
			errors.Is(err, services.ErrPasswordRequired) ||
			errors.Is(err, services.ErrPasswordTooShort) ||
			errors.Is(err, services.ErrPasswordTooLong) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}

		writeInternalServerError(c)
		return
	}

	h.setAuthCookies(c, result.AccessToken, result.AccessTokenExpiresAt, result.RefreshToken, result.RefreshTokenExpiresAt)

	c.JSON(http.StatusCreated, RegisterResponse{
		User: toUserResponse(result.User),
	})
}

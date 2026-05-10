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

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

func (h *WebAuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userAgent := optionalString(c.GetHeader("User-Agent"))
	clientIP := optionalString(c.ClientIP())

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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		if errors.Is(err, services.ErrEmailRequired) ||
			errors.Is(err, services.ErrPasswordRequired) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	h.setAuthCookie(c, h.cookies.AccessTokenName, result.AccessToken, result.AccessTokenExpiresAt)
	h.setAuthCookie(c, h.cookies.RefreshTokenName, result.RefreshToken, result.RefreshTokenExpiresAt)

	c.JSON(http.StatusOK, LoginResponse{
		User: toUserResponse(result.User),
	})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userAgent := optionalString(c.GetHeader("User-Agent"))
	clientIP := optionalString(c.ClientIP())

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
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrEmailRequired) ||
			errors.Is(err, services.ErrUsernameRequired) ||
			errors.Is(err, services.ErrPasswordRequired) ||
			errors.Is(err, services.ErrPasswordTooShort) ||
			errors.Is(err, services.ErrPasswordTooLong) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	h.setAuthCookie(c, h.cookies.AccessTokenName, result.AccessToken, result.AccessTokenExpiresAt)
	h.setAuthCookie(c, h.cookies.RefreshTokenName, result.RefreshToken, result.RefreshTokenExpiresAt)

	c.JSON(http.StatusCreated, RegisterResponse{
		User: toUserResponse(result.User),
	})
}

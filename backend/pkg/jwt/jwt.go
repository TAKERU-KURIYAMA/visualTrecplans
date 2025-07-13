package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrTokenClaims  = errors.New("invalid token claims")
)

// CustomClaims represents the custom claims for JWT tokens
type CustomClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey        []byte
	accessExpiry     time.Duration
	refreshExpiry    time.Duration
	issuer           string
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secret string, accessExpiry, refreshExpiry time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:     []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
	}
}

// GenerateAccessToken generates an access token for a user
func (j *JWTService) GenerateAccessToken(userID uuid.UUID, email string, isActive bool) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessExpiry)

	claims := CustomClaims{
		UserID:    userID,
		Email:     email,
		IsActive:  isActive,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			Audience:  []string{"trecplans"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a refresh token for a user
func (j *JWTService) GenerateRefreshToken(userID uuid.UUID, email string, isActive bool) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshExpiry)

	claims := CustomClaims{
		UserID:    userID,
		Email:     email,
		IsActive:  isActive,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			Audience:  []string{"trecplans"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTService) GenerateTokenPair(userID uuid.UUID, email string, isActive bool) (accessToken, refreshToken string, accessExpiry, refreshExpiry time.Time, err error) {
	accessToken, accessExpiry, err = j.GenerateAccessToken(userID, email, isActive)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	refreshToken, refreshExpiry, err = j.GenerateRefreshToken(userID, email, isActive)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessExpiry, refreshExpiry, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenClaims
	}

	return claims, nil
}

// ValidateAccessToken validates specifically an access token
func (j *JWTService) ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, fmt.Errorf("%w: expected access token", ErrTokenClaims)
	}

	return claims, nil
}

// ValidateRefreshToken validates specifically a refresh token
func (j *JWTService) ValidateRefreshToken(tokenString string) (*CustomClaims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("%w: expected refresh token", ErrTokenClaims)
	}

	return claims, nil
}

// ExtractUserID extracts user ID from a token without full validation
func (j *JWTService) ExtractUserID(tokenString string) (uuid.UUID, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return uuid.Nil, ErrTokenClaims
	}

	return claims.UserID, nil
}

// GetTokenExpiry returns the expiry time of a token
func (j *JWTService) GetTokenExpiry(tokenString string) (time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

// IsTokenExpired checks if a token is expired
func (j *JWTService) IsTokenExpired(tokenString string) bool {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Time.Before(time.Now())
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (j *JWTService) RefreshAccessToken(refreshTokenString string) (string, time.Time, error) {
	// Validate refresh token
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	// Generate new access token
	return j.GenerateAccessToken(claims.UserID, claims.Email, claims.IsActive)
}

// GetAccessTokenExpiry returns the access token expiry duration
func (j *JWTService) GetAccessTokenExpiry() time.Duration {
	return j.accessExpiry
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (j *JWTService) GetRefreshTokenExpiry() time.Duration {
	return j.refreshExpiry
}
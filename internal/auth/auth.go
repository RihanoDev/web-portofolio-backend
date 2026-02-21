package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTService interface defines JWT operations
type JWTService interface {
	GenerateToken(userID int, email, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (a *AuthService) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (a *AuthService) GenerateToken(userID int, username, role string) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   username,
			ID:        fmt.Sprintf("%d-%d", userID, now.Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Add the token type as a header
	token.Header["typ"] = "JWT"

	// Generate the signed token
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// Better validation with detailed error messages
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	// Handle parsing errors
	if err != nil {
		// Simple error handling with more descriptive messages
		errMessage := err.Error()
		switch {
		case strings.Contains(errMessage, "expired"):
			return nil, errors.New("token expired")
		case strings.Contains(errMessage, "not valid yet"):
			return nil, errors.New("token not valid yet")
		case strings.Contains(errMessage, "malformed"):
			return nil, errors.New("malformed token")
		case strings.Contains(errMessage, "signature"):
			return nil, errors.New("invalid token signature")
		default:
			return nil, fmt.Errorf("token validation error: %v", err)
		}
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation could be done here if needed
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

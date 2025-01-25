package jwt_manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mestvl-shop-app/auth/internal/domain"
)

// NewToken creates new JWT token for given user and app.
func NewToken(user *domain.Client, app *domain.App) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	duration, err := time.ParseDuration(fmt.Sprintf("%dm", app.JwtAccessTokenTtlMinutes))
	if err != nil {
		return "", fmt.Errorf("time.ParseDuration failed: %w", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.JwtSigningKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetAppID(accessToken string) (int, error) {
	token, _ := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, fmt.Errorf("get token claims failed")
	}

	return int(claims["app_id"].(float64)), nil
}

func Parse(accessToken string, jwtSigningKey string) (uuid.UUID, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSigningKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.UUID{}, jwt.ErrTokenExpired
		}
		return uuid.UUID{}, fmt.Errorf("parse token failed: %w", err)
	}

	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("get token claims failed")
	}

	idInterface := claims["uid"]
	if idInterface == nil {
		return uuid.UUID{}, fmt.Errorf("uid is empty")
	}

	idStr, ok := idInterface.(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("uid is not a string")
	}

	id, innerErr := uuid.Parse(idStr)
	if innerErr != nil {
		return uuid.UUID{}, fmt.Errorf("parse uuid failed: %w", innerErr)
	}

	return id, nil
}

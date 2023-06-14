package validators

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func UserIsAdmin(tokenString string) bool {

	return true
}

func IsValidRefreshToken(refreshToken string) bool {
	// Check if refresh token is blacklisted or expired
	if isBlacklistedRefreshToken(refreshToken) || isExpiredRefreshToken(refreshToken) {
		return false
	}

	// Check if refresh token has been tampered with
	_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return secret key for signing
		return []byte("secret_key"), nil
	})
	if err != nil {
		return false
	}

	// Refresh token is valid
	return true
}

func isBlacklistedRefreshToken(refreshToken string) bool {
	// Check if refresh token is in blacklist
	// You can implement the logic for maintaining a blacklist of refresh tokens
	// Here's an example implementation that stores blacklisted tokens in a map
	// and checks if the given refresh token is present in the map
	blacklist := make(map[string]bool)
	// Assume blacklist is populated with some refresh tokens
	if _, ok := blacklist[refreshToken]; ok {
		return true
	}

	// Refresh token is not blacklisted
	return false
}

func isExpiredRefreshToken(refreshToken string) bool {
	// Parse the refresh token to get the expiration time
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return secret key for signing
		return []byte("secret_key"), nil
	})
	if err != nil {
		return true // Refresh token is considered expired if it can't be parsed
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true // Refresh token is considered expired if claims can't be parsed
	}
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	// Check if the expiration time has passed
	return time.Now().After(expirationTime)
}

func ValidatePassword(password string) (string, bool) {
	if len(password) < 8 {
		return "Password must be at least 8 characters", false
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return "Password must contain at least one uppercase letter", false
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return "Password must contain at least one lowercase letter", false
	}
	if !strings.ContainsAny(password, "0123456789") {
		return "Password must contain at least one number", false
	}
	return "", true
}

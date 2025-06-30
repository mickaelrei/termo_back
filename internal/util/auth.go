package util

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"time"
)

const keyUserID = "id"

// GenerateAuthToken generates a new auth token storing the provided user ID
func GenerateAuthToken(userID int64, privateKey paseto.V4AsymmetricSecretKey) (string, error) {
	token := paseto.NewToken()
	now := time.Now()

	token.SetIssuer("termo")
	token.SetSubject("auth")
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(1 * time.Hour))

	err := token.Set(keyUserID, userID)
	if err != nil {
		return "", fmt.Errorf("[token.Set] | %v", err)
	}

	signed := token.V4Sign(privateKey, nil)
	return signed, nil
}

// GetIDFromToken attempts to extract the stored user ID from an auth token string
func GetIDFromToken(tokenString string, publicKey paseto.V4AsymmetricPublicKey) (int64, error) {
	parser := paseto.NewParser()

	parsedToken, err := parser.ParseV4Public(publicKey, tokenString, nil)
	if err != nil {
		return 0, fmt.Errorf("[parser.ParseV4Public] | %v", err)
	}

	var userID int64
	err = parsedToken.Get(keyUserID, &userID)
	if err != nil {
		return 0, fmt.Errorf("[parsedToken.Get] | %v", err)
	}

	return userID, nil
}

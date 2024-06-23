package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"memtravel/configs"
)

// CreateToken creates a jwt token for a specific user
func CreateToken(userID string) (string, error) {
	//TODO: add expiration time
	claims := jwt.MapClaims{
		"user": userID,
		"iss":  configs.Envs.JWTIssuer,
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(configs.Envs.JWTSecret)
}

// VerifyToken verifies a jwt token for a specific user
func VerifyToken(signedToken string) (bool, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return configs.Envs.JWTSecret, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, err
}

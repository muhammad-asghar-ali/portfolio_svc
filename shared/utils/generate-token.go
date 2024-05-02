package utils

import (
	"time"

	"github.com/0xbase-Corp/portfolio_svc/shared/configs"
	"github.com/golang-jwt/jwt"
)

type (
	// TODO: need to finalize
	Claims struct {
		UserID    int    `json:"user_id"`
		Email     string `json:"email"`
		PublicKey string `json:"public_key"`
		ExpiresAt int64  `json:"exp"`
		Issuer    string `json:"iss"`
		Audience  string `json:"aud"`
		IssuedAt  int64  `json:"iat"`
	}
)

func GenerateAccessToken(userID int, email, publicKey string) (string, error) {
	return generateToken(userID, email, publicKey, 438000*time.Hour)
}

func GenerateRefreshToken(userID int, email, publicKey string) (string, error) {
	return generateToken(userID, email, publicKey, 24*time.Hour)
}

func (c *Claims) Valid() error {
	if time.Unix(c.ExpiresAt, 0).Before(time.Now()) {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}
	return nil
}

func generateToken(userID int, email, publicKey string, expiresAt time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		PublicKey: publicKey,
		ExpiresAt: time.Now().Add(expiresAt).Unix(),
		Issuer:    "web3-",
		Audience:  "client",
		IssuedAt:  time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	signedToken, err := token.SignedString([]byte(configs.EnvConfigVars.GetSecret()))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

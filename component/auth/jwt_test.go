package auth

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {

	PublicKey, PrivateKey, _ = ed25519.GenerateKey(nil)

	tokenString, err := GenerateToken(1, "testuser", "testIssuer")

	assert.NoError(t, err)

	assert.NotEmpty(t, tokenString)

}

func TestParseToken_Expired(t *testing.T) {
	PublicKey, PrivateKey, _ = ed25519.GenerateKey(nil)

	now := time.Now()
	claims := Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)), // 设置为1小时前过期
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)), // 设置为2小时前签发
			Issuer:    "testIssuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenString, err := token.SignedString(PrivateKey)
	assert.NoError(t, err)

	parsedClaims, err := ParseToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestParseToken_Valid(t *testing.T) {
	PublicKey, PrivateKey, _ = ed25519.GenerateKey(nil)

	tokenString, err := GenerateToken(1, "testuser", "testIssuer")
	assert.NoError(t, err)

	parsedClaims, err := ParseToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)
	assert.Equal(t, int64(1), parsedClaims.UserID)
	assert.Equal(t, "testuser", parsedClaims.Username)
}

func TestErrEqual(t *testing.T) {
	assert.Equal(t, ErrTokenExpired, jwt.ErrTokenExpired)
}

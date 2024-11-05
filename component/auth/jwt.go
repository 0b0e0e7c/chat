package auth

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	TokenExpireDuration = time.Hour * 48
)

var (
	ErrTokenExpired = jwt.ErrTokenExpired
	PublicKey       ed25519.PublicKey
	PrivateKey      ed25519.PrivateKey
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`

	jwt.RegisteredClaims
}

func GenerateToken(userID int64, username, issuer string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(PrivateKey)
}

// ParseToken parses a JWT token and returns the claims
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return PublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func RefreshToken(userID int64, username, issuer string) (string, error) {
	return GenerateToken(userID, username, issuer)
}

func ValidateToken(tokenString string) (valid bool, uid int64, err error) {
	c, err := ParseToken(tokenString)

	if c != nil {
		logx.Info("token.uid: ", c.UserID)
		return true, c.UserID, nil
	}
	return false, 0, err

}

func GetTokenExpiry(tokenString string) (time.Time, error) {
	c, err := ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return c.ExpiresAt.Time, nil
}

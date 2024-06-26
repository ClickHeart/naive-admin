package jwt

import (
	"naive-admin/pkg/config"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SigningKey []byte
}

type MyCustomClaims struct {
	UserId int
	jwt.RegisteredClaims
}

func NewJwt(cfg *config.Jwt) *JWT {
	return &JWT{SigningKey: []byte(cfg.SigningKey)}
}

func (j *JWT) GenToken(userId int, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "",
			Subject:   "",
			ID:        "",
			Audience:  []string{},
		},
	})

	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.SigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	re := regexp.MustCompile(`(?i)Bearer `)
	tokenString = re.ReplaceAllString(tokenString, "")

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

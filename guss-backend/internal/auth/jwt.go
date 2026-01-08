// internal/auth/jwt.go
package auth

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5" // JWT 라이브러리 (표준에 가장 가까움)
)

var secretKey = []byte("GUSS_SECRET_KEY_2026") // 실제 운영 시 환경변수 처리

type Claims struct {
	UserNumber int64  `json:"user_number"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken: 로그인 성공 시 토큰 생성
func GenerateToken(userNumber int64, userID, role string) (string, error) {
	claims := &Claims{
		UserNumber: userNumber,
		UserID:     userID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken: 미들웨어에서 토큰 검증 시 사용
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
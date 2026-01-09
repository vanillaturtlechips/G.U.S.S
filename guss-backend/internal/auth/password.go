package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword: 평문 비밀번호를 bcrypt로 해싱합니다.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash: 입력된 비밀번호와 저장된 해시를 비교합니다.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
package api

import (
	"context"
	"guss-backend/internal/auth"
	"net/http"
	"strings"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Authorization 헤더 확인
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			s.errorJSON(w, "인증 토큰이 없습니다.", http.StatusUnauthorized)
			return
		}

		// 2. 토큰 문자열 추출 (Bearer 접두사 제거)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. [중요] 토큰 검증 및 Claims 추출
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			s.errorJSON(w, "유효하지 않거나 만료된 토큰입니다.", http.StatusUnauthorized)
			return
		}

		// 4. 검증된 진짜 Claims(유저 번호 포함)를 Context에 저장
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		// 5. 다음 핸들러(HandleReserve 등)로 전달
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

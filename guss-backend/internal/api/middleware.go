package api

import (
	"context"
	"guss-backend/internal/auth"
	"net/http"
	"strings"
)

// AuthMiddleware: 토큰의 유효성을 검증하고 유저 정보를 Context에 주입
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			s.errorJSON(w, "인증 토큰이 없습니다.", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 토큰 검증 및 Claims 추출
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			s.errorJSON(w, "유효하지 않거나 만료된 토큰입니다.", http.StatusUnauthorized)
			return
		}

		// 검증된 진짜 Claims(유저 번호, Role 포함)를 Context에 저장
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware: ADMIN 권한이 있는 유저만 허용 (AuthMiddleware 뒤에 배치해야 함)
func (s *Server) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Context에서 주입된 Claims 확인
		claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
		if !ok || claims.Role != "ADMIN" {
			s.errorJSON(w, "관리자 권한이 필요합니다.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

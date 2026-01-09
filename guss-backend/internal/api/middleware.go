package api

import (
	"context"
	"net/http"
	"strings"

	"guss-backend/internal/auth"
)

// handlers.go와 공유하는 컨텍스트 키
const UserContextKey = "user_number"

// AuthMiddleware: JWT 토큰 유효성 검사 및 사용자 정보 추출
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "인증이 필요합니다.", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// JWT 토큰 검증 및 Claims 추출
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "유효하지 않은 토큰입니다.", http.StatusUnauthorized)
			return
		}

		// 컨텍스트에 사용자 번호와 전체 Claims 저장
		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserNumber)
		ctx = context.WithValue(ctx, "full_claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware: ADMIN 역할(Role) 소유 여부 검사
func (s *Server) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 컨텍스트에서 Claims를 꺼내 Role 확인
		userInfo, ok := r.Context().Value("full_claims").(*auth.Claims)
		if !ok || userInfo.Role != "ADMIN" {
			http.Error(w, "관리자 권한이 없습니다.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

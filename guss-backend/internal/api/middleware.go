package api

import (
	"context"
	"net/http"
	"strings"

	"guss-backend/internal/auth"
)

// handlers.go와 약속한 키 이름
const UserContextKey = "user_number"

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "인증이 필요합니다.", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 1. JWT 토큰 검증 (반환 타입: *auth.Claims)
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "유효하지 않은 토큰입니다.", http.StatusUnauthorized)
			return
		}

		// 2. [수정 포인트] 구조체이므로 claims.UserNumber로 바로 접근합니다.
		// 이미 int64 타입일 것이므로 복잡한 변환도 필요 없습니다.
		userNum := claims.UserNumber

		// 3. 컨텍스트 주입
		ctx := context.WithValue(r.Context(), UserContextKey, userNum)
		ctx = context.WithValue(ctx, "full_claims", claims) // 구조체 포인터 그대로 저장
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value("full_claims").(*auth.Claims)
		if !ok || userInfo.Role != "ADMIN" {
			http.Error(w, "관리자 권한이 없습니다.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
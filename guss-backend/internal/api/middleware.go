package api

import (
	"context"
	"guss-backend/internal/auth"
	"net/http"
	"strings"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// 권한 없으면 여기서 에러 메세지 딱 보내고 끝! (이래야 EMPTY RESPONSE 안 남)
			s.errorJSON(w, "인증 토큰이 없습니다.", http.StatusUnauthorized)
			return
		}

		// 임시 인증 성공 처리
		ctx := context.WithValue(r.Context(), UserContextKey, &auth.Claims{UserID: "admin"})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package api

import (
	"context"
	"guss-backend/internal/auth"
	"net/http"
	"strings"
)

// [중요] 여기서 UserContextKey 정의를 삭제하세요!!

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}
		// ... 인증 로직 ...
		ctx := context.WithValue(r.Context(), UserContextKey, &auth.Claims{UserID: "admin"})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

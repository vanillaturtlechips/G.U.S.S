package api

import (
	"encoding/json"
	"guss-backend/internal/repository"
	"net/http"
)

// 여기서 한 번만 정의
type contextKey string

const UserContextKey contextKey = "user"

type Server struct {
	Repo    repository.Repository
	LogRepo repository.LogRepository
	Algo    any
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success", "token": "admin-token", "userRole": "ADMIN",
	})
}

func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gyms, _ := s.Repo.GetGyms()
	json.NewEncoder(w).Encode(gyms)
}

// 중복 에러 방지를 위해 다른 파일에 정의된 메서드는 절대 적지 마세요.
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request)     {}
func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {}
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request)      {}

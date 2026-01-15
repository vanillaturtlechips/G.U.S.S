package api

import (
	"encoding/json"
	"guss-backend/internal/algo"
	"guss-backend/internal/auth"
	"guss-backend/internal/domain"
	"guss-backend/internal/repository"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type contextKey string

const UserContextKey contextKey = "user"

type Server struct {
	Repo    repository.Repository
	LogRepo repository.LogRepository
	Algo    any
}

func (s *Server) errorJSON(w http.ResponseWriter, message string, code int) {
	log.Printf("[ERROR] %d: %s", code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input struct {
		UserID string `json:"user_id"`
		UserPW string `json:"user_pw"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.errorJSON(w, "잘못된 요청", 400)
		return
	}

	user, err := s.Repo.GetUserByID(input.UserID)
	var userNumber int64
	var userName string
	var hashedPassword string
	var role string = "USER"
	var gymID int64 = 0

	if err == nil {
		userNumber = user.UserNumber
		userName = user.UserName
		hashedPassword = user.UserPW
		if user.UserID == "admin" {
			role = "ADMIN"
		}
	} else {
		admin, err := s.Repo.GetAdminByID(input.UserID)
		if err != nil {
			s.errorJSON(w, "아이디/비밀번호 불일치", 401)
			return
		}
		userNumber = admin.AdminNumber
		userName = "관리자(" + admin.AdminID + ")"
		hashedPassword = admin.AdminPW
		role = "ADMIN"
		if admin.FKGussID.Valid {
			gymID = admin.FKGussID.Int64
		}
	}

	if !auth.CheckPasswordHash(input.UserPW, hashedPassword) {
		s.errorJSON(w, "비밀번호 불일치", 401)
		return
	}

	token, _ := auth.GenerateToken(userNumber, input.UserID, role)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success", "token": token, "user_name": userName, "user_role": role, "gym_id": gymID,
	})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	json.NewDecoder(r.Body).Decode(&u)
	hashed, _ := auth.HashPassword(u.UserPW)
	u.UserPW = hashed
	s.Repo.CreateUser(&u)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GymID     int64  `json:"gym_id"`
		StartTime string `json:"start_time"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	st, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
	claims := r.Context().Value(UserContextKey).(*auth.Claims)
	s.Repo.CreateReservationWithTime(claims.UserNumber, req.GymID, st, st.Add(30*time.Minute))
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleCancelReservation(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	resID, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	claims := r.Context().Value(UserContextKey).(*auth.Claims)
	s.Repo.UpdateReservationStatus(resID, claims.UserNumber, "cancelled")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleGetReservationStats(w http.ResponseWriter, r *http.Request) {
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gymId"), 10, 64)
	stats, _ := s.Repo.GetHourlyReservationStats(gymID)
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// main.go 에러 해결을 위한 매출 핸들러
	json.NewEncoder(w).Encode([]map[string]interface{}{{"type": "일일권", "amount": 10000, "date": "2025-07-24"}})
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "Running", "time": time.Now()})
}

func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	gyms, _ := s.Repo.GetGyms()
	if gyms == nil {
		gyms = []domain.Gym{}
	}
	json.NewEncoder(w).Encode(gyms)
}

func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	gym, _ := s.Repo.GetGymDetail(id)
	calculator := s.Algo.(*algo.RealTimeCalculator)
	json.NewEncoder(w).Encode(map[string]interface{}{"gym": gym, "congestion": calculator.Calculate(gym.GussUserCount, gym.GussSize)})
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("gymId"), 10, 64)
	list, _ := s.Repo.GetEquipmentsByGymID(id)
	json.NewEncoder(w).Encode(list)
}

func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	var eq domain.Equipment
	json.NewDecoder(r.Body).Decode(&eq)
	s.Repo.AddEquipment(&eq)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) {
	var eq domain.Equipment
	json.NewDecoder(r.Body).Decode(&eq)
	s.Repo.UpdateEquipment(&eq)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	s.Repo.DeleteEquipment(id)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("gymId"), 10, 64)
	list, _ := s.Repo.GetReservationsByGym(id)
	json.NewEncoder(w).Encode(list)
}

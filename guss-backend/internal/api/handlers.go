package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"guss-backend/internal/algo"
	"guss-backend/internal/auth"
	"guss-backend/internal/domain"
	"guss-backend/internal/repository"
)

type contextKey string

const UserContextKey contextKey = "user"

type Server struct {
	Repo    repository.Repository
	LogRepo repository.LogRepository
	Algo    any
	SQSURL  string // [추가] 환경별 SQS FIFO 큐 주소
}

func (s *Server) errorJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID string `json:"user_id"`
		UserPW string `json:"user_pw"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.errorJSON(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var userNumber int64
	var userName string
	var hashedPassword string
	var role string = "USER"
	var gymID int64 = 0

	user, err := s.Repo.GetUserByID(input.UserID)
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
			s.errorJSON(w, "인증 실패", http.StatusUnauthorized)
			return
		}
		userNumber = admin.AdminNumber
		userName = "관리자"
		hashedPassword = admin.AdminPW
		role = "ADMIN"
		if admin.FKGussID.Valid {
			gymID = admin.FKGussID.Int64
		}
	}

	if !auth.CheckPasswordHash(input.UserPW, hashedPassword) {
		s.errorJSON(w, "비밀번호 불일치", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(userNumber, input.UserID, role)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":     token,
		"user_role": role,
		"gym_id":    gymID,
		"user_name": userName,
		"status":    "success",
	})
}

func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	gyms, _ := s.Repo.GetGyms(search)
	json.NewEncoder(w).Encode(gyms)
}

func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GymID        int64  `json:"gym_id"`
		FkGussNumber int64  `json:"fk_guss_number"`
		VisitTime    string `json:"visit_time"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	targetID := req.GymID
	if targetID == 0 {
		targetID = req.FkGussNumber
	}

	t, err := time.Parse("2006-01-02 15:04:05", req.VisitTime)
	if err != nil || (t.Minute() != 0 && t.Minute() != 30) {
		s.errorJSON(w, "30분 단위로만 예약 가능합니다", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(UserContextKey).(*auth.Claims)
	_, err = s.Repo.CreateReservation(claims.UserNumber, targetID, t)
	if err != nil {
		s.errorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleCancelReservation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RevsNumber int64 `json:"revs_number"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	claims := r.Context().Value(UserContextKey).(*auth.Claims)
	err := s.Repo.CancelReservation(req.RevsNumber, claims.UserNumber, claims.Role)
	if err != nil {
		s.errorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	gym, err := s.Repo.GetGymDetail(id)
	if err != nil {
		s.errorJSON(w, "지점 정보 없음", http.StatusNotFound)
		return
	}
	calc := s.Algo.(*algo.RealTimeCalculator)
	congestion := calc.Calculate(gym.GussUserCount, gym.GussSize)
	json.NewEncoder(w).Encode(map[string]interface{}{"gym": gym, "congestion": congestion})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	json.NewDecoder(r.Body).Decode(&u)
	hashed, _ := auth.HashPassword(u.UserPW)
	u.UserPW = hashed
	if err := s.Repo.CreateUser(&u); err != nil {
		s.errorJSON(w, "가입 실패", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
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
	id, _ := strconv.ParseInt(r.URL.Query().Get("gym_id"), 10, 64)
	list, _ := s.Repo.GetReservationsByGym(id)
	json.NewEncoder(w).Encode(list)
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "Running", "server_time": time.Now().Format("2006-01-02 15:04:05")})
}

func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]map[string]interface{}{{"type": "일일권", "amount": 10000, "date": time.Now().Format("2006-01-02")}})
}

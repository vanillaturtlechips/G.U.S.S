package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"guss-backend/internal/algo"
	"guss-backend/internal/auth"
	"guss-backend/internal/domain"
	"guss-backend/internal/repository"
)

type Server struct {
	Repo    repository.Repository
	LogRepo repository.LogRepository
	Algo    algo.CongestionCalculator
}

// HandleRegister: 회원가입
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserName  string `json:"user_name"`
		UserPhone string `json:"user_phone"`
		UserID    string `json:"user_id"`
		UserPW    string `json:"user_pw"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "잘못된 요청 양식입니다.", http.StatusBadRequest)
		return
	}

	hashedPW, err := auth.HashPassword(input.UserPW)
	if err != nil {
		http.Error(w, "암호화 오류", http.StatusInternalServerError)
		return
	}

	user := domain.User{
		UserName:  input.UserName,
		UserPhone: input.UserPhone,
		UserID:    input.UserID,
		UserPW:    hashedPW,
	}

	if err := s.Repo.CreateUser(&user); err != nil {
		http.Error(w, "가입 실패 (아이디 중복 확인 요망)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration success"})
}

// HandleLogin: 로그인
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID  string `json:"user_id"`
		PWD string `json:"user_pw"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := s.Repo.GetUserByID(req.ID)
	if err != nil {
		http.Error(w, "사용자를 찾을 수 없습니다.", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(req.PWD, user.UserPW) {
		http.Error(w, "비밀번호가 일치하지 않습니다.", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.UserNumber, user.UserID, "USER")
	if err != nil {
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	_ = s.LogRepo.SaveUserLog(user.UserID, "LOGIN_SUCCESS")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// HandleGetGyms: 전체 목록
func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	gyms, err := s.Repo.GetAllGyms()
	if err != nil {
		http.Error(w, "목록 조회 실패", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gyms)
}

// HandleGetGymDetail: 상세 정보 (패닉 해결 버전)
func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/gyms/")
	idStr = strings.TrimSuffix(idStr, "/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "잘못된 ID 형식입니다.", http.StatusBadRequest)
		return
	}

	// [수정 포인트 1] DB에서 gym 데이터를 먼저 가져와야 합니다!
	gym, err := s.Repo.GetGymDetail(id)
	if err != nil || gym == nil {
		http.Error(w, "지점 정보를 찾을 수 없습니다.", http.StatusNotFound)
		return
	}

	// [수정 포인트 2] Algo가 nil일 경우를 대비한 방어 코드
	if s.Algo == nil {
		http.Error(w, "서버 설정 오류 (Algo is nil)", http.StatusInternalServerError)
		return
	}

	congestion := s.Algo.Calculate(gym.GussUserCount, gym.GussSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"gym":        gym,
		"congestion": congestion,
	})
}

// HandleReserve: 예약 (JWT 인증 연동)
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	userNum, ok := r.Context().Value(UserContextKey).(int64)
	if !ok {
		http.Error(w, "인증이 필요합니다.", http.StatusUnauthorized)
		return
	}

	var req struct {
		GymID int64 `json:"fk_guss_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청입니다.", http.StatusBadRequest)
		return
	}

	if err := s.Repo.CreateReservation(userNum, req.GymID); err != nil {
		http.Error(w, "예약 실패", http.StatusInternalServerError)
		return
	}

	_ = s.LogRepo.SaveUserLog(strconv.FormatInt(userNum, 10), "RESERVE_CREATED")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation completed"})
}
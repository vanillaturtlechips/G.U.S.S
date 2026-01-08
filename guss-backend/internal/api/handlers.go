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
	Repo    repository.Repository     // MySQL 처리
	LogRepo repository.LogRepository  // DynamoDB 로그 처리
	Algo    algo.CongestionCalculator // 혼잡도 알고리즘
}

// --- 1. 일반 사용자 핸들러 (Public API) ---

// HandleRegister: 회원가입 (입력 전용 구조체 사용으로 보안과 기능 동시 해결)
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// [보안 포인트] 입력받을 때만 사용할 구조체를 별도로 선언합니다.
	// domain.User의 json:"-" 태그 영향을 받지 않고 비밀번호를 읽어올 수 있습니다.
	var input struct {
		UserName  string `json:"user_name"`
		UserPhone string `json:"user_phone"`
		UserID    string `json:"user_id"`
		UserPW    string `json:"user_pw"` // 입력값 매핑
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "잘못된 요청 양식입니다.", http.StatusBadRequest)
		return
	}

	// 비밀번호 bcrypt 해싱
	hashedPW, err := auth.HashPassword(input.UserPW)
	if err != nil {
		http.Error(w, "암호화 오류", http.StatusInternalServerError)
		return
	}

	// 진짜 DB 저장용 객체(domain.User)에 데이터를 옮겨 담습니다.
	user := domain.User{
		UserName:  input.UserName,
		UserPhone: input.UserPhone,
		UserID:    input.UserID,
		UserPW:    hashedPW, // 해싱된 비밀번호 주입
	}

	if err := s.Repo.CreateUser(&user); err != nil {
		http.Error(w, "가입 실패 (아이디 중복 확인 요망)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration success"})
}

// HandleLogin: 로그인 및 JWT 발급
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

	// DB에서 가져온 해싱된 비번과 사용자가 입력한 평문 비번 비교
	if !auth.CheckPasswordHash(req.PWD, user.UserPW) {
		http.Error(w, "비밀번호가 일치하지 않습니다.", http.StatusUnauthorized)
		return
	}

	// JWT 토큰 생성
	token, err := auth.GenerateToken(user.UserNumber, user.UserID, "USER")
	if err != nil {
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	_ = s.LogRepo.SaveUserLog(user.UserID, "LOGIN_SUCCESS")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// HandleGetGyms: 전체 체육관 목록 조회
func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	gyms, err := s.Repo.GetAllGyms()
	if err != nil {
		http.Error(w, "체육관 목록 조회 실패", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gyms)
}

// HandleGetGymDetail: 상세 정보 조회 및 혼잡도 계산
func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	
	id, _ := strconv.ParseInt(idStr, 10, 64)

	gym, err := s.Repo.GetGymDetail(id)
	if err != nil {
		http.Error(w, "체육관 정보를 찾을 수 없습니다.", http.StatusNotFound)
		return
	}

	congestion := s.Algo.Calculate(gym.GussUserCount, gym.GussSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"gym":        gym,
		"congestion": congestion,
	})
}

// HandleReserve: 이용 예약 신청 (인증 필요)
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	// middleware.go에 정의된 UserContextKey를 사용
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
		http.Error(w, "예약 실패 (이미 예약 중이거나 데이터 오류)", http.StatusInternalServerError)
		return
	}

	_ = s.LogRepo.SaveUserLog(strconv.FormatInt(userNum, 10), "RESERVE_CREATED")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation completed"})
}
package api

import (
	"encoding/json"
	"guss-backend/internal/algo"
	"guss-backend/internal/auth" // JWT 및 Bcrypt 인증 패키지
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

// errorJSON: 공통 에러 응답 처리용 헬퍼 함수
func (s *Server) errorJSON(w http.ResponseWriter, message string, code int) {
	log.Printf("[ERROR] 코드: %d, 메시지: %s", code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// HandleLogin: 실제 DB 조회, 비밀번호 검증 및 JWT 발급
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input struct {
		UserID string `json:"user_id"`
		UserPW string `json:"user_pw"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.errorJSON(w, "잘못된 요청 형식입니다.", http.StatusBadRequest)
		return
	}

	// 1. DB에서 해당 유저 아이디로 정보 조회
	user, err := s.Repo.GetUserByID(input.UserID)
	if err != nil {
		s.errorJSON(w, "존재하지 않는 사용자 아이디입니다.", http.StatusUnauthorized)
		return
	}

	// 2. 비밀번호 검증 (Bcrypt 비교)
	if !auth.CheckPasswordHash(input.UserPW, user.UserPW) {
		s.errorJSON(w, "비밀번호가 일치하지 않습니다.", http.StatusUnauthorized)
		return
	}

	// 3. 진짜 JWT 토큰 생성 (관리자 여부 판단 로직 필요 시 추가)
	// 예: ID가 admin이면 ADMIN 권한 부여
	role := "USER"
	if user.UserID == "admin" {
		role = "ADMIN"
	}
	
	token, err := auth.GenerateToken(user.UserNumber, user.UserID, role)
	if err != nil {
		s.errorJSON(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"token":     token,
		"user_name": user.UserName,
		"user_role": role,
	})
}

// HandleRegister: 비밀번호 해싱 후 DB 저장
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u domain.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.errorJSON(w, "데이터 형식 오류", 400)
		return
	}

	// 비밀번호 해싱 처리
	hashedPW, err := auth.HashPassword(u.UserPW)
	if err != nil {
		s.errorJSON(w, "비밀번호 처리 중 오류 발생", 500)
		return
	}
	u.UserPW = hashedPW

	if err := s.Repo.CreateUser(&u); err != nil {
		s.errorJSON(w, "회원가입 실패 (아이디 중복 등)", 500)
		return
	}

	log.Printf("[SUCCESS] 신규 유저 가입: %s (No: %d)", u.UserName, u.UserNumber)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"user_number": u.UserNumber,
	})
}

// HandleReserve: JWT 토큰에서 유저 정보를 추출하여 예약 처리
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		GymID        int64 `json:"gym_id"`
		FkGussNumber int64 `json:"fk_guss_number"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	if req.GymID == 0 && req.FkGussNumber > 0 {
		req.GymID = req.FkGussNumber
	}

	// [중요] 미들웨어에서 넣어준 Claims에서 유저 번호 추출
	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		s.errorJSON(w, "인증 정보가 없습니다.", http.StatusUnauthorized)
		return
	}

	// DB 예약 로직 호출 (중복 예약 방지 로직이 포함된 Repo 메서드)
	_, err := s.Repo.CreateReservation(claims.UserNumber, req.GymID)
	if err != nil {
		// [수정] 500 에러가 아닌 400 에러를 반환하여 프론트에서 경고 모달을 띄우게 함
		s.errorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[SUCCESS] 유저 %d번 -> 체육관 %d번 예약 완료", claims.UserNumber, req.GymID)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// --- 관리자 핸들러 (Mock 데이터 포함) ---

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := map[string]interface{}{
		"status":     "Running",
		"active_now": 12,
	}
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs := []map[string]interface{}{
		{"type": "일일권", "amount": 10000, "date": time.Now().Format("2006-01-02")},
	}
	json.NewEncoder(w).Encode(logs)
}

// --- 공통 조회 핸들러들 ---

func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gyms, err := s.Repo.GetGyms()
	if err != nil {
		s.errorJSON(w, "조회 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(gyms)
}

func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]
	id, _ := strconv.ParseInt(idStr, 10, 64)

	gym, err := s.Repo.GetGymDetail(id)
	if err != nil {
		s.errorJSON(w, "체육관 정보를 찾을 수 없음", 404)
		return
	}

	calculator := s.Algo.(*algo.RealTimeCalculator)
	utilization := calculator.Calculate(gym.GussUserCount, gym.GussSize)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"gym":        gym,
		"congestion": utilization,
	})
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("gymId")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	list, err := s.Repo.GetEquipmentsByGymID(id)
	if err != nil {
		s.errorJSON(w, "조회 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	var eq domain.Equipment
	json.NewDecoder(r.Body).Decode(&eq)
	if err := s.Repo.AddEquipment(&eq); err != nil {
		s.errorJSON(w, "등록 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var eq domain.Equipment

	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		s.errorJSON(w, "데이터 형식 오류", 400)
		return
	}

	// [디버그 로그] 프론트에서 넘어온 데이터 확인
	log.Printf("[DEBUG] 수정 요청 데이터: %+v", eq)

	// JSON 바디에 ID가 없으면 URL 경로에서 추출
	if eq.ID <= 0 {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			// 경로의 가장 마지막 요소를 ID로 간주
			id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
			eq.ID = id
		}
	}

	if eq.ID <= 0 {
		s.errorJSON(w, "수정할 기구 ID를 찾을 수 없습니다.", 400)
		return
	}

	err := s.Repo.UpdateEquipment(&eq)
	if err != nil {
		s.errorJSON(w, "DB 수정 실패: "+err.Error(), 500)
		return
	}

	log.Printf("[SUCCESS] 기구 %d번(%s) 수정 완료", eq.ID, eq.Name)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err := s.Repo.DeleteEquipment(id); err != nil {
		s.errorJSON(w, "삭제 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. 쿼리 스트링에서 gym_id 추출
	idStr := r.URL.Query().Get("gym_id")
	if idStr == "" {
		idStr = r.URL.Query().Get("gymId")
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if id <= 0 {
		s.errorJSON(w, "체육관 ID가 유효하지 않습니다.", http.StatusBadRequest)
		return
	}

	// 2. DB에서 예약 목록 조회
	list, err := s.Repo.GetReservationsByGym(id)
	if err != nil {
		s.errorJSON(w, "예약 목록 조회 실패", http.StatusInternalServerError)
		return
	}

	// 3. JSON 응답
	json.NewEncoder(w).Encode(list)
}

// test
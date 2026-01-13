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

	// [디버그] 입력값 확인 (비밀번호 불일치 원인 파악용)
	log.Printf("[DEBUG] 로그인 시도 - ID: %s, PW길이: %d", input.UserID, len(input.UserPW))

	if input.UserPW == "" {
		s.errorJSON(w, "비밀번호가 입력되지 않았습니다.", http.StatusBadRequest)
		return
	}

	// 1. DB에서 해당 유저 아이디로 정보 조회
	user, err := s.Repo.GetUserByID(input.UserID)
	if err != nil {
		s.errorJSON(w, "존재하지 않는 사용자 아이디입니다.", http.StatusUnauthorized)
		return
	}

	// 2. 비밀번호 검증 (DB의 해시값 vs 입력된 평문)
	if !auth.CheckPasswordHash(input.UserPW, user.UserPW) {
		s.errorJSON(w, "비밀번호가 일치하지 않습니다.", http.StatusUnauthorized)
		return
	}

	// 3. 진짜 JWT 토큰 생성
	token, err := auth.GenerateToken(user.UserNumber, user.UserID, "USER")
	if err != nil {
		s.errorJSON(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	// 4. 성공 응답
	log.Printf("[SUCCESS] 로그인 완료: %s (No: %d)", user.UserName, user.UserNumber)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"token":     token,
		"user_name": user.UserName,
		"user_role": "USER",
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

	// [주의] domain.User 구조체에 UserPW가 json:"-"로 되어 있으면 여기서 비어있게 됨
	if u.UserPW == "" {
		s.errorJSON(w, "비밀번호가 전달되지 않았습니다 (JSON 태그 확인 필요)", 400)
		return
	}

	// 비밀번호 저장 전 해싱 처리
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
		UserNumber   int64 `json:"user_number"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	// GymID 필드명 보정
	if req.GymID == 0 && req.FkGussNumber > 0 {
		req.GymID = req.FkGussNumber
	}

	// [중요] 미들웨어에서 넣어준 Claims에서 진짜 유저 번호를 꺼냄
	if claims, ok := r.Context().Value(UserContextKey).(*auth.Claims); ok {
		req.UserNumber = claims.UserNumber
	}

	if req.UserNumber <= 0 {
		s.errorJSON(w, "로그인이 필요한 서비스입니다.", http.StatusUnauthorized)
		return
	}

	if req.GymID <= 0 {
		s.errorJSON(w, "체육관 번호가 누락되었습니다.", 400)
		return
	}

	_, err := s.Repo.CreateReservation(req.UserNumber, req.GymID)
	if err != nil {
		s.errorJSON(w, "예약 실패: "+err.Error(), 500)
		return
	}

	log.Printf("[SUCCESS] 유저 %d번 -> 체육관 %d번 예약 완료", req.UserNumber, req.GymID)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleGetGyms: 체육관 목록 조회
func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gyms, err := s.Repo.GetGyms()
	if err != nil {
		s.errorJSON(w, "조회 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(gyms)
}

// HandleGetGymDetail: 특정 체육관 상세 정보 및 실시간 혼잡도 계산
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

	// 알고리즘을 사용한 혼잡도 계산
	current := gym.GussUserCount
	max := gym.GussSize
	if max <= 0 {
		max = 20
	}

	calculator := s.Algo.(*algo.RealTimeCalculator)
	utilization := calculator.Calculate(current, max)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"gym":        gym,
		"congestion": utilization,
	})
}

// HandleGetEquipments: 체육관별 기구 목록 조회
func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("gymId")
	if idStr == "" {
		idStr = r.URL.Query().Get("gym_id")
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)

	list, err := s.Repo.GetEquipmentsByGymID(id)
	if err != nil {
		s.errorJSON(w, "기구 목록 조회 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// HandleAddEquipment: 신규 기구 등록
func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var eq domain.Equipment

	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		s.errorJSON(w, "데이터 형식 오류", 400)
		return
	}

	if eq.GymID <= 0 {
		idStr := r.URL.Query().Get("gym_id")
		eq.GymID, _ = strconv.ParseInt(idStr, 10, 64)
	}

	if eq.PurchaseDate == "" {
		eq.PurchaseDate = time.Now().Format("2006-01-02")
	}

	err := s.Repo.AddEquipment(&eq)
	if err != nil {
		s.errorJSON(w, "DB 등록 실패", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleUpdateEquipment: 기구 정보 수정
func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var eq domain.Equipment

	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		s.errorJSON(w, "데이터 형식 오류", 400)
		return
	}

	if eq.ID <= 0 {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 3 {
			id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
			eq.ID = id
		}
	}

	err := s.Repo.UpdateEquipment(&eq)
	if err != nil {
		s.errorJSON(w, "DB 수정 실패", 500)
		return
	}

	log.Printf("[SUCCESS] 기구 %d번 수정 완료", eq.ID)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleDeleteEquipment: 기구 삭제
func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	if id <= 0 {
		s.errorJSON(w, "ID 오류", 400)
		return
	}

	err := s.Repo.DeleteEquipment(id)
	if err != nil {
		s.errorJSON(w, "삭제 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

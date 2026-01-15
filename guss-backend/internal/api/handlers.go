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

// HandleLogin: 유저/관리자 통합 로그인 및 지점별 권한 부여
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

	var userNumber int64
	var userName string
	var hashedPassword string
	var role string = "USER"
	var gymID int64 = 0 // 관리자의 담당 지점 ID (0은 전체 권한 의미)

	// 1. 먼저 일반 유저 테이블에서 조회
	user, err := s.Repo.GetUserByID(input.UserID)
	if err == nil {
		userNumber = user.UserNumber
		userName = user.UserName
		hashedPassword = user.UserPW
		// 아이디가 admin인 유저는 USER 테이블에 있더라도 ADMIN으로 취급
		if user.UserID == "admin" {
			role = "ADMIN"
		}
	} else {
		// 2. 유저가 없으면 관리자 전용 테이블(admin_table) 조회
		admin, err := s.Repo.GetAdminByID(input.UserID)
		if err != nil {
			s.errorJSON(w, "아이디 또는 비밀번호가 일치하지 않습니다.", http.StatusUnauthorized)
			return
		}
		userNumber = admin.AdminNumber
		userName = "관리자(" + admin.AdminID + ")"
		hashedPassword = admin.AdminPW
		role = "ADMIN"

		// [중요] sql.NullInt64 안전하게 처리 (super_admin은 NULL이므로 Valid가 false)
		if admin.FKGussID.Valid {
			gymID = admin.FKGussID.Int64 // 지점 관리자
		} else {
			gymID = 0 // 최고 관리자 (NULL)
		}
	}

	// 3. 비밀번호 검증 (Bcrypt)
	if !auth.CheckPasswordHash(input.UserPW, hashedPassword) {
		s.errorJSON(w, "아이디 또는 비밀번호가 일치하지 않습니다.", http.StatusUnauthorized)
		return
	}

	// 4. 최고 관리자 ID 별도 판단 (로직 보강 가능)
	if input.UserID == "super_admin" {
		role = "SUPER_ADMIN"
	}

	// 5. 토큰 생성 및 응답
	token, err := auth.GenerateToken(userNumber, input.UserID, role)
	if err != nil {
		s.errorJSON(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	log.Printf("[LOGIN] %s 접속 (Role: %s, GymID: %d)", input.UserID, role, gymID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"token":     token,
		"user_name": userName,
		"user_role": role,
		"gym_id":    gymID, // 프론트엔드에서 지점 필터링에 사용
	})
}

// HandleRegister: 회원가입 (Bcrypt 적용)
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u domain.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.errorJSON(w, "데이터 형식 오류", 400)
		return
	}

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

// HandleReserve: 중복 예약 방지 로직 적용
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

	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		s.errorJSON(w, "인증 정보가 없습니다.", http.StatusUnauthorized)
		return
	}

	_, err := s.Repo.CreateReservation(claims.UserNumber, req.GymID)
	if err != nil {
		s.errorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[SUCCESS] 유저 %d번 -> 체육관 %d번 예약 완료", claims.UserNumber, req.GymID)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleDashboard: 지점별 실시간 통계
func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := map[string]interface{}{
		"status":      "Running",
		"active_now":  12,
		"server_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	json.NewEncoder(w).Encode(stats)
}

// HandleGetSales: 매출 데이터 조회 (향후 DynamoDB 마이그레이션 대상)
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

	if eq.ID <= 0 {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
			eq.ID = id
		}
	}

	if eq.ID <= 0 {
		s.errorJSON(w, "수정할 기구 ID를 찾을 수 없습니다.", 400)
		return
	}

	if err := s.Repo.UpdateEquipment(&eq); err != nil {
		s.errorJSON(w, "DB 수정 실패", 500)
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

	idStr := r.URL.Query().Get("gym_id")
	if idStr == "" {
		idStr = r.URL.Query().Get("gymId")
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if id <= 0 {
		s.errorJSON(w, "체육관 ID가 유효하지 않습니다.", http.StatusBadRequest)
		return
	}

	list, err := s.Repo.GetReservationsByGym(id)
	if err != nil {
		s.errorJSON(w, "예약 목록 조회 실패", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}

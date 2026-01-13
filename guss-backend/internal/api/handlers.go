package api

import (
	"encoding/json"
	"guss-backend/internal/algo"
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
	log.Printf("[ERROR] 코드: %d, 메시지: %s", code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "token": "admin-token", "userRole": "ADMIN"})
}

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
		s.errorJSON(w, "정보 없음", 404)
		return
	}
	current, max := gym.GussUserCount, gym.GussSize
	if max <= 0 {
		max = 20
	}
	calculator := s.Algo.(*algo.RealTimeCalculator)
	utilization := calculator.Calculate(current, max)
	json.NewEncoder(w).Encode(map[string]interface{}{"gym": gym, "congestion": utilization})
}

func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		GymID int64 `json:"gym_id"`
	} // 구조체와 일치시킴
	json.NewDecoder(r.Body).Decode(&req)
	if req.GymID <= 0 {
		idStr := r.URL.Query().Get("gymId")
		if idStr == "" {
			idStr = r.URL.Query().Get("gym_id")
		}
		req.GymID, _ = strconv.ParseInt(idStr, 10, 64)
	}
	s.Repo.CreateReservation(1, req.GymID)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("gymId")
	if idStr == "" {
		idStr = r.URL.Query().Get("gym_id")
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	list, _ := s.Repo.GetEquipmentsByGymID(id)
	json.NewEncoder(w).Encode(list)
}

// HandleAddEquipment: 이제 구조체 태그 수정으로 깔끔하게 작동합니다.
func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var eq domain.Equipment

	// 1. JSON 읽기 (이제 구조체의 json 태그가 맞으므로 자동으로 값이 들어옵니다)
	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		s.errorJSON(w, "JSON 파싱 실패", 400)
		return
	}

	// 2. 만약 JSON에 번호가 빠졌을 경우를 대비해 URL 쿼리도 확인
	if eq.GymID <= 0 {
		idStr := r.URL.Query().Get("gymId")
		if idStr == "" {
			idStr = r.URL.Query().Get("gym_id")
		}
		eq.GymID, _ = strconv.ParseInt(idStr, 10, 64)
	}

	// 3. 번호가 없으면 에러 (꼼수 없이 정확하게 체크)
	if eq.GymID <= 0 {
		s.errorJSON(w, "체육관 번호(gym_id)가 없습니다.", 400)
		return
	}

	// 4. 날짜 보정
	if eq.PurchaseDate == "" {
		eq.PurchaseDate = time.Now().Format("2006-01-02")
	}

	// 5. DB 저장
	err := s.Repo.AddEquipment(&eq)
	if err != nil {
		s.errorJSON(w, "DB 저장 실패: "+err.Error(), 500)
		return
	}

	log.Printf("[SUCCESS] 체육관 %d번에 기구 '%s' 등록 성공", eq.GymID, eq.Name)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	s.Repo.DeleteEquipment(id)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

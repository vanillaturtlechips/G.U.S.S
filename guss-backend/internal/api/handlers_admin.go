package api

import (
	"encoding/json"
	"guss-backend/internal/domain" // 추가
	"net/http"
	"strconv"
)

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_now": 12, "status": "Running", "total_users": 150, "total_revenue": 500000,
	})
}

func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	var req domain.Equipment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON 파싱 에러: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 실제 DB 저장 (mysql_repository의 AddEquipment 호출)
	if err := s.Repo.AddEquipment(&req); err != nil {
		http.Error(w, "DB 저장 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gym_id"), 10, 64)
	equipments, _ := s.Repo.GetEquipmentsByGymID(gymID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipments)
}

func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Path[len("/api/equipments/"):], 10, 64)
	s.Repo.DeleteEquipment(id)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) {}
func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request)        {}

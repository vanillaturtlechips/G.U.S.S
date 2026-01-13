package api

import (
	"encoding/json"
	"net/http"
)

// HandleDashboard: 관리자 메인 대시보드 통계 데이터 (Mock 데이터)
func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := map[string]interface{}{
		"status":        "Running",
		"active_now":    12,
		"total_revenue": 500000,
	}
	json.NewEncoder(w).Encode(stats)
}

// HandleGetReservations: 예약 현황 로그 조회 (Mock 데이터)
func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	logs := []map[string]interface{}{
		{"revs_number": 1, "user_name": "김철수", "user_phone": "010-1234-5678", "revs_status": "CONFIRMED"},
		{"revs_number": 2, "user_name": "이영희", "user_phone": "010-9876-5432", "revs_status": "CONFIRMED"},
	}
	json.NewEncoder(w).Encode(logs)
}

// HandleGetSales: 매출 로그 조회 (Mock 데이터)
func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	logs := []map[string]interface{}{
		{"type": "일일권", "amount": 10000, "date": "2026-01-13 14:00"},
		{"type": "PT 10회", "amount": 450000, "date": "2026-01-13 15:30"},
	}
	json.NewEncoder(w).Encode(logs)
}

// 주의: AuthMiddleware 함수는 middleware.go에 이미 있으므로 여기서 삭제했습니다.

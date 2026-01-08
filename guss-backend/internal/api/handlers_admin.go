package api

import (
	"encoding/json"
	"net/http"
	"guss-backend/internal/auth"
)

// HandleAddEquipment: 관리자가 새로운 기구를 등록하거나 상태를 변경할 때 호출됩니다.
// DynamoDB(LogRepo)에 기구 관련 로그를 남깁니다.
func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GussNumber int64  `json:"guss_number"`
		EquipID    string `json:"equip_id"`
		Status     string `json:"status"`
	}

	// 1. 요청 본문 파싱
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 JSON 형식입니다.", http.StatusBadRequest)
		return
	}

	// 2. 인증된 관리자 정보 추출 (미사용 변수 에러 방지)
	if userInfo, ok := r.Context().Value(UserContextKey).(*auth.Claims); ok {
		// 누가 기구를 등록하려 했는지 추적 로그 작성
		_ = s.LogRepo.SaveUserLog(userInfo.UserID, "ADMIN_EQUIP_REG_ACTION")
	}

	// 3. DynamoDB에 기구 로그 저장
	err := s.LogRepo.SaveEqLog(req.GussNumber, req.EquipID, req.Status)
	if err != nil {
		http.Error(w, "기구 로그 저장에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	// 4. 성공 응답
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "기구 정보가 성공적으로 로그에 저장되었습니다.",
		"equip_id": req.EquipID,
	})
}

// HandleGetAdminDashboard: 관리자 페이지의 메인 대시보드 데이터를 반환합니다.
// main.go의 registerRoutes에서 부르는 이름과 정확히 일치해야 합니다.
func (s *Server) HandleGetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	// 실제 운영 시에는 s.Repo.GetSalesStats() 등을 호출하여 MySQL 데이터를 가져옵니다.
	// 현재는 로컬 테스트를 위한 Mock 데이터를 반환합니다.
	dashboardData := map[string]interface{}{
		"daily_revenue": []map[string]interface{}{
			{"date": "2026-01-08", "amount": 1250000, "type": "MEMBERSHIP"},
			{"date": "2026-01-07", "amount": 980000, "type": "PT_SESSION"},
		},
		"active_users": 47,
		"system_status": "ONLINE",
	}

	w.Header().Set("Content-Type", "application/json")
	
	// HTTP 응답 본문 작성
	if err := json.NewEncoder(w).Encode(dashboardData); err != nil {
		http.Error(w, "응답 생성 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}
}
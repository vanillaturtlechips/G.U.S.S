package api // 1. 패키지 선언 누락 해결

import (
	"encoding/json"
	"guss-backend/internal/auth"
	"guss-backend/internal/infrastructure/aws" // 2. SendCheckInEvent를 쓰기 위한 임포트
	"net/http"
)

func (s *Server) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReservationID int64 `json:"reservation_id"`
		GymID         int64 `json:"gym_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorJSON(w, "잘못된 요청 양식입니다.", http.StatusBadRequest)
		return
	}

	// 토큰에서 인증된 Claims 추출
	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		s.errorJSON(w, "인증 정보를 찾을 수 없습니다.", http.StatusUnauthorized)
		return
	}

	// [수정 포인트] s.SQSURL을 첫 번째 인자로 추가하여 want (string, int64, string, string) 형식을 맞춥니다.
	err := aws.SendCheckInEvent(s.SQSURL, req.GymID, claims.UserID, "IN")
	if err != nil {
		s.errorJSON(w, "실시간 혼잡도 반영 실패", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "체크인 성공! 실시간 혼잡도가 업데이트됩니다.",
	})
}

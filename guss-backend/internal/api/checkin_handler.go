package api // 1. 패키지 선언 누락 해결

import (
	"fmt"
	"guss-backend/internal/infrastructure/aws" // 2. SendCheckInEvent를 쓰기 위한 임포트
	"log"
)

// HandleQRCheckIn: QR 스캔 시 호출되는 메인 로직
func HandleQRCheckIn(resID int64, gymID int64, userID string) error {
	// 1. 예약 유효성 검사 (실제 구현 시 SQL DB 조회 로직 추가)
	log.Printf("[CHECKIN] QR 스캔됨: 예약번호 %d, 유저 %s", resID, userID)

	// 2. [실시간 파이프라인] SQS로 체크인 이벤트 전송
	// 이전에 만든 infrastructure/aws 패키지의 함수를 호출합니다.
	err := aws.SendCheckInEvent(gymID, userID, "IN")
	if err != nil {
		log.Printf("[ERROR] SQS 전송 실패: %v", err)
		return fmt.Errorf("실시간 혼잡도 업데이트 실패: %v", err)
	}

	log.Printf("[SUCCESS] 체크인 완료 및 실시간 데이터 전송: 유저 %s", userID)
	return nil
}

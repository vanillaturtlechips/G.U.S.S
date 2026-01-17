package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// CheckInMessage: Worker(Lambda)에서 참조하는 이름과 동일하게 맞춤
type CheckInMessage struct {
	GymID  int64  `json:"gym_id"`
	UserID string `json:"user_id"`
	Action string `json:"action"` // "IN" 또는 "OUT"
}

// SendCheckInEvent: 체크인 발생 시 SQS로 메시지를 전송합니다.
func SendCheckInEvent(gymID int64, userID string, action string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("AWS 설정 로드 실패: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	// 구조체 이름을 CheckInMessage로 변경하여 생성
	event := CheckInMessage{
		GymID:  gymID,
		UserID: userID,
		Action: action,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("메시지 직렬화 실패: %v", err)
	}

	// 실제 본인의 SQS URL로 입력되어 있는지 확인하세요.
	queueURL := "https://sqs.ap-northeast-2.amazonaws.com/YOUR_ACCOUNT_ID/guss-checkin-queue"

	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})

	if err != nil {
		return fmt.Errorf("SQS 메시지 전송 실패: %v", err)
	}

	log.Printf("[SQS] 지점 %d 체크인 이벤트 전송 완료", gymID)
	return nil
}

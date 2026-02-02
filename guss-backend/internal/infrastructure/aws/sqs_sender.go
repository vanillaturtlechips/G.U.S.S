package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// CheckInMessage: SQS 메시지로 보낼 데이터 구조입니다.
type CheckInMessage struct {
	GymID  int64  `json:"gym_id"`
	UserID string `json:"user_id"`
	Action string `json:"action"` // "IN" 또는 "OUT"
}

// SendCheckInEvent: 환경별로 주입된 queueURL을 사용하여 FIFO 큐로 메시지를 전송합니다.
func SendCheckInEvent(queueURL string, gymID int64, userID string, action string) error {
	// AWS 기본 설정 로드
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("AWS 설정 로드 실패: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	// 메시지 본문 JSON 직렬화
	msg := CheckInMessage{
		GymID:  gymID,
		UserID: userID,
		Action: action,
	}
	msgBody, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("메시지 직렬화 실패: %v", err)
	}

	// FIFO 큐 전송 (MessageGroupId 필수)
	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(msgBody)),
		// 동일한 지점의 이벤트는 동일한 그룹 ID를 가져야 순차 처리가 보장됩니다.
		MessageGroupId: aws.String(fmt.Sprintf("gym-%d", gymID)),
		// 콘솔에서 '콘텐츠 기반 중복 제거'를 활성화했으므로 DeduplicationId는 생략 가능합니다.
	})

	return err
}

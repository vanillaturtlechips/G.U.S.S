package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CheckInMessage: SQS 메시지 구조체 (에러 방지를 위해 직접 정의)
type CheckInMessage struct {
	GymID  int64  `json:"gym_id"`
	UserID string `json:"user_id"`
	Action string `json:"action"` // "IN" 또는 "OUT"
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// 1. AWS SDK 설정 로드
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("AWS 설정 로드 실패: %v", err)
	}
	dbClient := dynamodb.NewFromConfig(cfg)

	// 2. SQS 레코드 반복 처리
	for _, message := range sqsEvent.Records {
		var checkIn CheckInMessage
		if err := json.Unmarshal([]byte(message.Body), &checkIn); err != nil {
			log.Printf("JSON 파싱 에러: %v", err)
			continue
		}

		// 3. DynamoDB 원자적 업데이트 (Atomic Increment)
		// Action이 "IN"이면 +1, 아니면 -1
		incrementValue := "1"
		if checkIn.Action == "OUT" {
			incrementValue = "-1"
		}

		// 테이블명은 본인의 DynamoDB 테이블명으로 수정하세요 (예: gym_status)
		tableName := "gym_status"

		_, err = dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"gym_id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", checkIn.GymID)},
			},
			// ADD 연산자를 사용하여 숫자를 안전하게 증가/감소시킴
			UpdateExpression: aws.String("ADD current_count :val"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":val": &types.AttributeValueMemberN{Value: incrementValue},
			},
		})

		if err != nil {
			log.Printf("DynamoDB 업데이트 실패 (지점 %d): %v", checkIn.GymID, err)
			return err
		}

		log.Printf("성공: 지점 %d 인원수 %s 처리 완료", checkIn.GymID, checkIn.Action)
	}

	return nil
}

func main() {
	// 람다 실행 시작
	lambda.Start(handler)
}

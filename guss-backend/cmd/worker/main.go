package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"google.golang.org/api/option"
)

// ReservationMessage: FCMToken 필드를 추가하여 기기를 식별합니다.
type ReservationMessage struct {
	ResID    string `json:"res_id" dynamodbav:"res_id"`
	UserID   string `json:"user_id" dynamodbav:"user_id"`
	GymID    int64  `json:"gym_id" dynamodbav:"gym_id"`
	Time     string `json:"time" dynamodbav:"visit_time"`
	EventAt  string `json:"event_at" dynamodbav:"event_at"`
	Amount   int    `json:"amount" dynamodbav:"amount"`
	Status   string `json:"status" dynamodbav:"status"`
	FCMToken string `json:"fcm_token" dynamodbav:"fcm_token"` // 알림을 받을 기기 토큰
}

var (
	dbClient  *dynamodb.Client
	tableName string
	fcmClient *messaging.Client
)

func init() {
	ctx := context.TODO()

	// 1. AWS 설정 초기화
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("AWS 설정 로드 실패: %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("TABLE_NAME")

	// 2. Firebase Admin SDK 초기화 (SSM에서 주입된 환경 변수 사용)
	// 파일 대신 메모리에 로드된 JSON 데이터를 직접 사용합니다.
	firebaseJSON := os.Getenv("FIREBASE_CONFIG")
	if firebaseJSON == "" {
		log.Fatal("FIREBASE_CONFIG 환경 변수가 설정되지 않았습니다.")
	}

	opt := option.WithCredentialsJSON([]byte(firebaseJSON))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Firebase 앱 초기화 실패: %v", err)
	}

	fcmClient, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalf("FCM 클라이언트 생성 실패: %v", err)
	}
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		fmt.Printf("[Lambda] 메시지 처리 시작: %s\n", record.MessageId)

		var msg ReservationMessage
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			log.Printf("메시지 파싱 실패: %v", err)
			continue
		}

		// 1. 혼잡도 계산
		congestionStatus := checkCongestion(msg.GymID)

		// 2. 실제 FCM 푸시 발송
		if msg.FCMToken != "" {
			sendRealPush(ctx, msg.FCMToken, congestionStatus)
		} else {
			log.Printf("유저 %s의 FCM 토큰이 없어 알림을 건너뜜", msg.UserID)
		}

		// 3. 데이터 보강 및 저장
		if msg.Status == "" {
			msg.Status = "RESERVED"
		}

		item, _ := attributevalue.MarshalMap(msg)
		_, err := dbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})

		if err != nil {
			log.Printf("DynamoDB 저장 에러: %v", err)
			return err
		}

		fmt.Printf("[SUCCESS] 예약 ID: %s 저장 및 FCM 발송 완료\n", msg.ResID)
	}
	return nil
}

func checkCongestion(gymID int64) string {
	now := time.Now().Hour()
	if now >= 18 && now <= 21 {
		return "HIGH"
	}
	return "NORMAL"
}

// sendRealPush: Firebase를 통해 실제 알림을 전송합니다.
func sendRealPush(ctx context.Context, token string, status string) {
	title := "GUSS 예약 알림"
	body := "✅ 예약이 완료되었습니다. 현장에서 QR을 스캔해주세요."

	if status == "HIGH" {
		title = "⚠️ 체육관 혼잡 알림"
		body = "현재 이용객이 많습니다. 입장이 조금 지연될 수 있어요!"
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	response, err := fcmClient.Send(ctx, message)
	if err != nil {
		log.Printf("[FCM ERROR] 발송 실패: %v", err)
		return
	}
	fmt.Printf("[FCM SUCCESS] 메시지 ID: %s\n", response)
}

func main() {
	lambda.Start(handler)
}

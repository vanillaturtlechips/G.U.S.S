package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

/* // [주의] 이 부분은 repository.go와 중복되므로 주석 처리합니다.
type LogRepository interface {
	SaveEqLog(gussNumber int64, equipID, status string) error
	SaveUserLog(userID, action string) error
}
*/

// dynamoLogRepo: DynamoDB를 사용하는 로그 저장소 실체
type dynamoLogRepo struct {
	client *dynamodb.Client
}

// NewDynamoLogRepository: DynamoDB 저장소 생성함수
func NewDynamoLogRepository(client *dynamodb.Client) LogRepository {
	return &dynamoLogRepo{
		client: client,
	}
}

// SaveEqLog: 기구 상태 변경 로그를 DynamoDB에 저장
func (d *dynamoLogRepo) SaveEqLog(gussNumber int64, equipID, status string) error {
	item := map[string]interface{}{
		"PK":          fmt.Sprintf("GYM#%d", gussNumber),
		"SK":          fmt.Sprintf("EQ#%s#%d", equipID, time.Now().Unix()),
		"GussNumber":  gussNumber,
		"EquipID":     equipID,
		"Status":      status,
		"Timestamp":   time.Now().Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("guss_logs"),
		Item:      av,
	})

	if err != nil {
		log.Printf("DynamoDB 기구 로그 저장 실패: %v", err)
	}
	return err
}

// SaveUserLog: 사용자 활동 로그를 DynamoDB에 저장
func (d *dynamoLogRepo) SaveUserLog(userID, action string) error {
	item := map[string]interface{}{
		"PK":        fmt.Sprintf("USER#%s", userID),
		"SK":        fmt.Sprintf("ACT#%d", time.Now().Unix()),
		"UserID":    userID,
		"Action":    action,
		"Timestamp": time.Now().Format(time.RFC3339),
		"TTL":       time.Now().Add(time.Hour * 24 * 30).Unix(), // 30일 후 자동 삭제
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("guss_logs"),
		Item:      av,
	})

	if err != nil {
		log.Printf("DynamoDB 유저 로그 저장 실패: %v", err)
	}
	return err
}
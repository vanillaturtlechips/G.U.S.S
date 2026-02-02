package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type CheckInMessage struct {
	GymID  int64  `json:"gym_id"`
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-2"))
	dbClient := dynamodb.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)

	for _, message := range sqsEvent.Records {
		var msg CheckInMessage
		json.Unmarshal([]byte(message.Body), &msg)

		val := "1"
		if msg.Action == "OUT" {
			val = "-1"
		}

		// 1. DynamoDB 숫자 업데이트 및 결과 받아오기 (ALL_NEW)
		res, err := dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String("gym_status"),
			Key: map[string]types.AttributeValue{
				"gym_id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", msg.GymID)},
			},
			UpdateExpression: aws.String("ADD current_count :v"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberN{Value: val},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			continue
		}

		// 2. SNS 알림 로직 (혼잡도 80% 체크)
		// DynamoDB에 max_size가 저장되어 있다고 가정
		current, _ := res.Attributes["current_count"].(*types.AttributeValueMemberN)
		max, _ := res.Attributes["max_size"].(*types.AttributeValueMemberN)

		if current != nil && max != nil {
			// 알림 조건: 입장(IN)일 때 80% 돌파 시
			if msg.Action == "IN" && current.Value >= "40" { // 예: 50명 중 40명(80%) 돌파 시
				snsClient.Publish(ctx, &sns.PublishInput{
					Message:  aws.String(fmt.Sprintf("[경고] 지점 %d 혼잡도 80%% 돌파! 현재 %s명", msg.GymID, current.Value)),
					TopicArn: aws.String("arn:aws:sns:ap-northeast-2:계정ID:guss-alert-topic"),
				})
			}
		}
	}
	return nil
}

func main() { lambda.Start(handler) }

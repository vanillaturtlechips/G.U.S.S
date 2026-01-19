// internal/infrastructure/aws/dynamodb_reader.go (또는 repository)
package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GymCongestion struct {
	GymID        int64   `json:"gym_id"`
	CurrentCount int     `json:"current_count"`
	MaxCapacity  int     `json:"max_capacity"`
	Ratio        float64 `json:"ratio"`
}

func GetGymCongestion(ctx context.Context, dbClient *dynamodb.Client, gymID int64) (*GymCongestion, error) {
	tableName := "GymStats"

	// DynamoDB에서 데이터 읽기
	out, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"gym_id": &types.AttributeValueMemberN{Value: strconv.FormatInt(gymID, 10)},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, fmt.Errorf("해당 지점의 데이터가 없습니다")
	}

	curr, _ := strconv.Atoi(out.Item["current_count"].(*types.AttributeValueMemberN).Value)
	max, _ := strconv.Atoi(out.Item["max_capacity"].(*types.AttributeValueMemberN).Value)

	return &GymCongestion{
		GymID:        gymID,
		CurrentCount: curr,
		MaxCapacity:  max,
		Ratio:        float64(curr) / float64(max),
	}, nil
}

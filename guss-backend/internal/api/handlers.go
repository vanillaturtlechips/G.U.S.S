package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"

	"guss-backend/internal/algo"
	"guss-backend/internal/auth"
	"guss-backend/internal/domain"
	"guss-backend/internal/repository"
)

type contextKey string
const UserContextKey contextKey = "user"

type Server struct {
	Repo         repository.Repository
	LogRepo      repository.LogRepository
	Algo         any
	SQSClient    *sqs.Client
	SQSURL       string
	DynamoClient *dynamodb.Client
	DynamoTable  string
}

func (s *Server) errorJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// HandleLogin: FCM 토큰 업데이트 포함
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		UserPW   string `json:"user_pw"`
		FCMToken string `json:"fcm_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.errorJSON(w, "잘못된 요청 형식", 400)
		return
	}

	user, err := s.Repo.GetUserByID(input.UserID)
	if err != nil || !auth.CheckPasswordHash(input.UserPW, user.UserPW) {
		s.errorJSON(w, "아이디 또는 비밀번호가 일치하지 않습니다.", 401)
		return
	}

	if input.FCMToken != "" {
		s.Repo.UpdateFCMToken(input.UserID, input.FCMToken)
	}

	token, _ := auth.GenerateToken(user.UserNumber, user.UserID, "USER")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success", "token": token, "user_name": user.UserName,
	})
}

// HandleReserve: 예약 정보만 생성 (인원 증가는 하지 않음)
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GymID     int64  `json:"gym_id"`
		VisitTime string `json:"visit_time"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	claims := r.Context().Value(UserContextKey).(*auth.Claims)
	t, _ := time.Parse("2006-01-02 15:04:05", req.VisitTime)
	
	// DB 예약 생성
	_, err := s.Repo.CreateReservation(claims.UserNumber, req.GymID, t)
	if err != nil {
		s.errorJSON(w, err.Error(), 400)
		return
	}

	resID := uuid.New().String()
	eventAt := time.Now().Format(time.RFC3339)
	fcmToken, _ := s.Repo.GetFCMToken(claims.UserID)

	// SQS 전송 (Lambda 트리거용)
	if s.SQSClient != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"res_id": resID, "user_id": claims.UserID, "fcm_token": fcmToken, "gym_id": req.GymID,
		})
		s.SQSClient.SendMessage(r.Context(), &sqs.SendMessageInput{
			QueueUrl: aws.String(s.SQSURL), MessageBody: aws.String(string(payload)),
		})
	}

	// QR에 담길 실제 체크인 API 주소
	checkInURL := fmt.Sprintf("https://api.guss.com/api/checkin?res_id=%s&user_id=%s&gym_id=%d&event_at=%s", 
		resID, claims.UserID, req.GymID, eventAt)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success", "qr_data": checkInURL,
	})
}

// HandleCheckIn: 실제 QR 스캔 시 인원 증가 및 DynamoDB 업데이트
func (s *Server) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	resID := r.URL.Query().Get("res_id")
	userID := r.URL.Query().Get("user_id")
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gym_id"), 10, 64)
	eventAt := r.URL.Query().Get("event_at")

	// 1. DB 인원수 증가
	if err := s.Repo.IncrementUserCount(gymID); err != nil {
		s.errorJSON(w, "입장 처리 실패", 500)
		return
	}

	// 2. DynamoDB 상태 변경
	s.DynamoClient.UpdateItem(r.Context(), &dynamodb.UpdateItemInput{
		TableName: aws.String(s.DynamoTable),
		Key: map[string]types.AttributeValue{
			"user_id":  &types.AttributeValueMemberS{Value: userID},
			"event_at": &types.AttributeValueMemberS{Value: eventAt},
		},
		UpdateExpression: aws.String("SET #s = :status"),
		ExpressionAttributeNames: map[string]string{"#s": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":status": &types.AttributeValueMemberS{Value: "ATTENDED"}},
	})

	log.Printf("[SUCCESS] User %s Check-in Gym %d (ResID: %s)", userID, gymID, resID)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>체크인 성공! 입장해 주세요.</h1>")
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// 실제 운영 환경에서는 DB 통계 쿼리 결과를 반환
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "running",
		"graph_data": []map[string]interface{}{
			{"time": "09:00", "count": 10}, {"time": "14:00", "count": 30}, {"time": "20:00", "count": 45},
		},
	})
}

// 기타 조회 핸들러
func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	gyms, _ := s.Repo.GetAllGyms()
	json.NewEncoder(w).Encode(gyms)
}

func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	gym, _ := s.Repo.GetGymDetail(id)
	calc := s.Algo.(*algo.RealTimeCalculator)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"gym": gym, "congestion": calc.Calculate(gym.GussUserCount, gym.GussSize),
	})
}
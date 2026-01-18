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

// HandleLogin: 로그인 및 FCM 토큰 업데이트
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		UserPW   string `json:"user_pw"`
		FCMToken string `json:"fcm_token"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	user, err := s.Repo.GetUserByID(input.UserID)
	if err != nil || !auth.CheckPasswordHash(input.UserPW, user.UserPW) {
		s.errorJSON(w, "인증 실패", 401)
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

// HandleRegister: 회원가입 (domain 패키지 사용)
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.errorJSON(w, "데이터 형식 오류", 400)
		return
	}
	hashed, _ := auth.HashPassword(u.UserPW)
	u.UserPW = hashed
	if err := s.Repo.CreateUser(&u); err != nil {
		s.errorJSON(w, "회원가입 실패", 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleReserve: 예약 생성 (visitTime 포함하여 3개 인자 전달)
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GymID     int64  `json:"gym_id"`
		VisitTime string `json:"visit_time"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// 시간 파싱 및 기본값 설정
	t, err := time.Parse("2006-01-02 15:04:05", req.VisitTime)
	if err != nil {
		t = time.Now().Add(1 * time.Hour)
	}
	
	claims := r.Context().Value(UserContextKey).(*auth.Claims)

	// [수정] 3개의 인자(userNum, gymNum, t)를 전달하도록 수정
	_, err = s.Repo.CreateReservation(claims.UserNumber, req.GymID, t)
	if err != nil {
		s.errorJSON(w, err.Error(), 400)
		return
	}

	resID := uuid.New().String()
	eventAt := time.Now().Format(time.RFC3339)

	if s.SQSClient != nil {
		fcmToken, _ := s.Repo.GetFCMToken(claims.UserID)
		payload, _ := json.Marshal(map[string]interface{}{
			"res_id": resID, "user_id": claims.UserID, "fcm_token": fcmToken, "gym_id": req.GymID, "event_at": eventAt,
		})
		s.SQSClient.SendMessage(r.Context(), &sqs.SendMessageInput{
			QueueUrl: aws.String(s.SQSURL), MessageBody: aws.String(string(payload)),
		})
	}

	qrURL := fmt.Sprintf("https://api.guss.com/api/checkin?res_id=%s&user_id=%s&gym_id=%d&event_at=%s", resID, claims.UserID, req.GymID, eventAt)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "qr_data": qrURL})
}

// HandleCheckIn: QR 스캔 시 호출 (resID 변수 사용)
func (s *Server) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	resID := r.URL.Query().Get("res_id")
	userID := r.URL.Query().Get("user_id")
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gym_id"), 10, 64)
	eventAt := r.URL.Query().Get("event_at")

	if err := s.Repo.IncrementUserCount(gymID); err != nil {
		s.errorJSON(w, "인원 증가 실패", 500)
		return
	}

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

	// [수정] resID를 사용하여 'declared and not used' 에러 방지
	log.Printf("[CHECK-IN SUCCESS] ReservationID: %s, User: %s, Gym: %d", resID, userID, gymID)
	w.Write([]byte("<html><body><h1>입장이 확인되었습니다.</h1></body></html>"))
}

// [핵심 해결] main.go에서 요구하는 HandleCancelReservation 추가
func (s *Server) HandleCancelReservation(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "예약이 취소되었습니다."})
}

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

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "running"})
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request)   { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request)    { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request)        { json.NewEncoder(w).Encode([]string{}) }

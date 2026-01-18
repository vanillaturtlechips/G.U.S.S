package api

import (
	"encoding/json"
	"fmt"
	"log" // 빌드 에러의 원인: 아래에서 반드시 사용해야 합니다.
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

// HandleReserve: 예약 생성 및 로그 기록
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GymID     int64  `json:"gym_id"`
		VisitTime string `json:"visit_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorJSON(w, "잘못된 요청 형식입니다.", 400)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", req.VisitTime)
	if err != nil { t = time.Now() }
	
	claims := r.Context().Value(UserContextKey).(*auth.Claims)

	status, err := s.Repo.CreateReservation(claims.UserNumber, req.GymID, t)
	if err != nil {
		if status == "DUPLICATE" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "DUPLICATE", "error": err.Error()})
			return
		}
		s.errorJSON(w, err.Error(), 400)
		return
	}

	resID := uuid.New().String()
	eventAt := time.Now().Format(time.RFC3339)

	// [log 사용] 예약 발생 로그 기록
	log.Printf("[RESERVE] User: %s, Gym: %d, Time: %s", claims.UserID, req.GymID, t.Format("15:04:05"))

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

// HandleCheckIn: 실제 QR 스캔 시에만 인원이 증가하며 로그를 남깁니다.
func (s *Server) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	resID := r.URL.Query().Get("res_id")
	userID := r.URL.Query().Get("user_id")
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gym_id"), 10, 64)
	eventAt := r.URL.Query().Get("event_at")

	if err := s.Repo.IncrementUserCount(gymID); err != nil {
		s.errorJSON(w, "인원 갱신 실패", 500)
		return
	}

	// [log 사용] 체크인 성공 로그 기록
	log.Printf("[CHECK-IN SUCCESS] ResID: %s, User: %s, Gym: %d", resID, userID, gymID)

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

	w.Write([]byte("<html><body><h1>Check-in Success! 반갑습니다.</h1></body></html>"))
}

func (s *Server) HandleCancelReservation(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// --- 아래는 빌드를 위한 공통 핸들러 (이전과 동일) ---
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct { UserID, UserPW, FCMToken string }
	json.NewDecoder(r.Body).Decode(&input)
	user, err := s.Repo.GetUserByID(input.UserID)
	if err != nil || !auth.CheckPasswordHash(input.UserPW, user.UserPW) {
		s.errorJSON(w, "인증 실패", 401); return
	}
	if input.FCMToken != "" { s.Repo.UpdateFCMToken(input.UserID, input.FCMToken) }
	token, _ := auth.GenerateToken(user.UserNumber, user.UserID, "USER")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "token": token, "user_name": user.UserName})
}
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	json.NewDecoder(r.Body).Decode(&u)
	h, _ := auth.HashPassword(u.UserPW); u.UserPW = h
	s.Repo.CreateUser(&u)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
func (s *Server) HandleGetGyms(w http.ResponseWriter, r *http.Request) {
	gyms, _ := s.Repo.GetAllGyms(); json.NewEncoder(w).Encode(gyms)
}
func (s *Server) HandleGetGymDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	gym, _ := s.Repo.GetGymDetail(id)
	calc := s.Algo.(*algo.RealTimeCalculator)
	json.NewEncoder(w).Encode(map[string]interface{}{"gym": gym, "congestion": calc.Calculate(gym.GussUserCount, gym.GussSize)})
}
func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "running"}) }
func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode([]string{}) }

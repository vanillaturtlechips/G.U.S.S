package api

import (
	"encoding/json"
	"fmt"
	"log" // 추가됨
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

// HandleLogin: 유저/관리자 로그인 시 FCM 토큰 업데이트 로직 추가
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		UserPW   string `json:"user_pw"`
		FCMToken string `json:"fcm_token"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	var userNumber int64
	var userName string
	var hashedPassword string
	var role string = "USER"
	var gymID int64 = 0

	user, err := s.Repo.GetUserByID(input.UserID)
	if err == nil {
		userNumber, userName, hashedPassword = user.UserNumber, user.UserName, user.UserPW
		if user.UserID == "admin" {
			role = "ADMIN"
		}
		// 로그인 성공 시 FCM 토큰 업데이트
		if input.FCMToken != "" {
			s.Repo.UpdateFCMToken(input.UserID, input.FCMToken)
		}
	} else {
		admin, err := s.Repo.GetAdminByID(input.UserID)
		if err != nil {
			s.errorJSON(w, "인증 실패", http.StatusUnauthorized)
			return
		}
		userNumber, userName, hashedPassword, role = admin.AdminNumber, "관리자", admin.AdminPW, "ADMIN"
		if admin.FKGussID.Valid {
			gymID = admin.FKGussID.Int64
		}
	}
	if !auth.CheckPasswordHash(input.UserPW, hashedPassword) {
		s.errorJSON(w, "비밀번호 불일치", http.StatusUnauthorized)
		return
	}
	token, _ := auth.GenerateToken(userNumber, input.UserID, role)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":     token,
		"user_role": role,
		"gym_id":    gymID,
		"user_name": userName,
		"status":    "success",
	})
}

// HandleReserve: SQS 메시지에 FCM 토큰 포함 및 성공 로그 추가
func (s *Server) HandleReserve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GymID        int64  `json:"gym_id"`
		FkGussNumber int64  `json:"fk_guss_number"`
		VisitTime    string `json:"visit_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorJSON(w, "잘못된 요청 형식입니다.", http.StatusBadRequest)
		return
	}

	targetID := req.GymID
	if targetID == 0 {
		targetID = req.FkGussNumber
	}
	t, _ := time.Parse("2006-01-02 15:04:05", req.VisitTime)
	claims := r.Context().Value(UserContextKey).(*auth.Claims)

	// DB에서 FCM 토큰 조회
	fcmToken, _ := s.Repo.GetFCMToken(claims.UserID)

	reservationID := uuid.New().String()
	eventAt := time.Now().Format(time.RFC3339)

	_, err := s.Repo.CreateReservation(claims.UserNumber, targetID, t)
	if err != nil {
		s.errorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if s.SQSClient != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"res_id":    reservationID,
			"user_id":   claims.UserID,
			"gym_id":    targetID,
			"fcm_token": fcmToken,
			"time":      req.VisitTime,
			"event_at":  eventAt,
		})

		_, err := s.SQSClient.SendMessage(r.Context(), &sqs.SendMessageInput{
			QueueUrl:               aws.String(s.SQSURL),
			MessageBody:            aws.String(string(payload)),
			MessageGroupId:         aws.String("GUSS-REV"),
			MessageDeduplicationId: aws.String(reservationID),
		})

		if err != nil {
			fmt.Printf("[SQS ERROR] 전송 실패: %v\n", err)
		} else {
			// userID -> claims.UserID로 수정 및 log 패키지 사용
			log.Printf("[SUCCESS] SQS 메시지 전송 완료! (ID: %s, User: %s)", reservationID, claims.UserID)
		}
	}

	checkInURL := fmt.Sprintf("https://api.guss.com/api/checkin?res_id=%s&user_id=%s&event_at=%s", reservationID, claims.UserID, eventAt)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "success",
		"reservation_id": reservationID,
		"qr_data":        checkInURL,
	})
}

func (s *Server) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	resID := r.URL.Query().Get("res_id")
	userID := r.URL.Query().Get("user_id")
	eventAt := r.URL.Query().Get("event_at")

	if resID == "" || userID == "" || eventAt == "" {
		s.errorJSON(w, "필수 체크인 정보가 누락되었습니다.", http.StatusBadRequest)
		return
	}

	_, err := s.DynamoClient.UpdateItem(r.Context(), &dynamodb.UpdateItemInput{
		TableName: aws.String(s.DynamoTable),
		Key: map[string]types.AttributeValue{
			"user_id":  &types.AttributeValueMemberS{Value: userID},
			"event_at": &types.AttributeValueMemberS{Value: eventAt},
		},
		UpdateExpression: aws.String("SET #s = :status"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "ATTENDED"},
		},
	})

	if err != nil {
		fmt.Printf("[DYNAMO ERROR] 체크인 실패: %v\n", err)
		s.errorJSON(w, "체크인 처리 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "체크인이 완료되었습니다. 즐거운 운동 되세요!",
	})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	json.NewDecoder(r.Body).Decode(&u)
	hashed, _ := auth.HashPassword(u.UserPW)
	u.UserPW = hashed
	s.Repo.CreateUser(&u)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
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
		"gym":        gym,
		"congestion": calc.Calculate(gym.GussUserCount, gym.GussSize),
	})
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request)   { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request)    { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]string{"status": "success"}) }
func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request)        { json.NewEncoder(w).Encode([]string{}) }
func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request)       { json.NewEncoder(w).Encode(map[string]string{"status": "running"}) }
func (s *Server) HandleCancelReservation(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

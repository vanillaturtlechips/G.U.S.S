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

	log.Printf("[RESERVE] User: %s, Gym: %d, Time: %s", claims.UserID, req.GymID, t.Format("15:04:05"))

	if s.SQSClient != nil {
		log.Printf("[SQS] Sending message to queue...")
		fcmToken, _ := s.Repo.GetFCMToken(claims.UserID)
		log.Printf("[SQS] FCM Token: %s", fcmToken)
		
		payload, _ := json.Marshal(map[string]interface{}{
			"res_id": resID, "user_id": claims.UserID, "fcm_token": fcmToken, "gym_id": req.GymID, "event_at": eventAt,
		})

		_, err := s.SQSClient.SendMessage(r.Context(), &sqs.SendMessageInput{
			QueueUrl:               aws.String(s.SQSURL),
			MessageBody:            aws.String(string(payload)),
			MessageGroupId:         aws.String("reservation"),
			MessageDeduplicationId: aws.String(resID),
		})

		if err != nil {
			log.Printf("[SQS ERROR] Failed to send: %v", err)
		} else {
			log.Printf("[SQS SUCCESS] Message sent!")
		}
	}

	qrURL := fmt.Sprintf("https://43.203.212.179/api/checkin?res_id=%s&user_id=%s&gym_id=%d&event_at=%s",
		resID, claims.UserID, req.GymID, eventAt)

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "qr_data": qrURL})
}

func (s *Server) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	resID := r.URL.Query().Get("res_id")
	userID := r.URL.Query().Get("user_id")
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gym_id"), 10, 64)
	eventAt := r.URL.Query().Get("event_at")

	if err := s.Repo.IncrementUserCount(gymID); err != nil {
		s.errorJSON(w, "인원 갱신 실패", 500)
		return
	}

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
	claims := r.Context().Value(UserContextKey).(*auth.Claims)

	var req struct {
		ReservationID int64 `json:"reservation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorJSON(w, "잘못된 요청 형식입니다.", 400)
		return
	}

	if err := s.Repo.CancelReservation(req.ReservationID, claims.UserNumber); err != nil {
		s.errorJSON(w, err.Error(), 400)
		return
	}

	log.Printf("[CANCEL] User: %s, ReservationID: %d", claims.UserID, req.ReservationID)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleGetActiveReservation(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(UserContextKey).(*auth.Claims)

	log.Printf("[GET ACTIVE] UserNumber: %d, UserID: %s", claims.UserNumber, claims.UserID)

	reservation, err := s.Repo.GetActiveReservationByUser(claims.UserNumber)

	log.Printf("[GET ACTIVE] Reservation: %+v, Error: %v", reservation, err)

	if err != nil || reservation == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"reservation": nil})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"reservation": reservation})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		UserPW   string `json:"user_pw"`
		FCMToken string `json:"fcm_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("[LOGIN] JSON decode error: %v", err)
		s.errorJSON(w, "잘못된 요청 형식", 400)
		return
	}

	log.Printf("[LOGIN] Attempting login for user: %s", input.UserID)

	user, err := s.Repo.GetUserByID(input.UserID)
	if err != nil {
		log.Printf("[LOGIN] GetUserByID error: %v", err)
		s.errorJSON(w, "인증 실패", 401)
		return
	}

	log.Printf("[LOGIN] User found: %s, checking password", user.UserID)

	if !auth.CheckPasswordHash(input.UserPW, user.UserPW) {
		log.Printf("[LOGIN] Password mismatch for user: %s", input.UserID)
		s.errorJSON(w, "인증 실패", 401)
		return
	}

	log.Printf("[LOGIN] Login successful for user: %s", input.UserID)

	if input.FCMToken != "" {
		s.Repo.UpdateFCMToken(input.UserID, input.FCMToken)
	}

	token, _ := auth.GenerateToken(user.UserNumber, user.UserID, "USER")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token": token,
		"user_name": user.UserName,
		"role": "USER",
	})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	json.NewDecoder(r.Body).Decode(&u)
	h, _ := auth.HashPassword(u.UserPW)
	u.UserPW = h
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
	json.NewEncoder(w).Encode(map[string]interface{}{"gym": gym, "congestion": calc.Calculate(gym.GussUserCount, gym.GussSize)})
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "running"})
}

func (s *Server) HandleGetReservations(w http.ResponseWriter, r *http.Request) {
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gymId"), 10, 64)

	log.Printf("[GET RESERVATIONS] GymID: %d", gymID)

	reservations, err := s.Repo.GetReservationsByGym(gymID)
	if err != nil {
		log.Printf("[GET RESERVATIONS] Error: %v", err)
		json.NewEncoder(w).Encode([]domain.Reservation{})
		return
	}

	log.Printf("[GET RESERVATIONS] Found %d reservations", len(reservations))

	if reservations == nil {
		reservations = []domain.Reservation{}
	}

	json.NewEncoder(w).Encode(reservations)
}

func (s *Server) HandleGetSales(w http.ResponseWriter, r *http.Request) {
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gymId"), 10, 64)

	log.Printf("[GET SALES] GymID: %d", gymID)

	sales, err := s.Repo.GetSalesByGym(gymID)
	if err != nil {
		log.Printf("[GET SALES] Error: %v", err)
		json.NewEncoder(w).Encode([]domain.Sale{})
		return
	}

	log.Printf("[GET SALES] Found %d sales", len(sales))

	if sales == nil {
		sales = []domain.Sale{}
	}

	json.NewEncoder(w).Encode(sales)
}

func (s *Server) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	gymID, _ := strconv.ParseInt(r.URL.Query().Get("gymId"), 10, 64)

	log.Printf("[GET EQUIPMENTS] GymID: %d", gymID)

	equipments, err := s.Repo.GetEquipmentsByGymID(gymID)
	if err != nil {
		log.Printf("[GET EQUIPMENTS] Error: %v", err)
		json.NewEncoder(w).Encode([]domain.Equipment{})
		return
	}

	log.Printf("[GET EQUIPMENTS] Found %d equipments", len(equipments))

	if equipments == nil {
		equipments = []domain.Equipment{}
	}

	json.NewEncoder(w).Encode(equipments)
}

func (s *Server) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	var eq domain.Equipment
	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		s.errorJSON(w, "잘못된 요청 형식입니다.", 400)
		return
	}

	log.Printf("[ADD EQUIPMENT] Gym: %d, Name: %s", eq.GymID, eq.Name)

	if err := s.Repo.AddEquipment(&eq); err != nil {
		s.errorJSON(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleUpdateEquipment(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	var eq domain.Equipment
	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		s.errorJSON(w, "잘못된 요청 형식입니다.", 400)
		return
	}

	eq.ID = id

	log.Printf("[UPDATE EQUIPMENT] ID: %d, Name: %s", eq.ID, eq.Name)

	if err := s.Repo.UpdateEquipment(&eq); err != nil {
		s.errorJSON(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) HandleDeleteEquipment(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	log.Printf("[DELETE EQUIPMENT] ID: %d", id)

	if err := s.Repo.DeleteEquipment(id); err != nil {
		s.errorJSON(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

package repository

import (
	"database/sql"
	"guss-backend/internal/domain"
	"log"
	"time"
)

// --- MockRepository ---
type MockRepository struct{}

func NewMockRepository() Repository {
	return &MockRepository{}
}

// 1. 유저 관련 Mock
func (m *MockRepository) CreateUser(u *domain.User) error {
	log.Printf("[MOCK] User Created: %v", u.UserName)
	return nil
}

func (m *MockRepository) GetUserByID(id string) (*domain.User, error) {
	return &domain.User{
		UserNumber: 1,
		UserID:     id,
		UserName:   "Mock일반유저",
		UserPW:     "$2a$10$Wp6S7Vf4X.0pGzXz9XyYduzI6.R8z8L5v5.m7Gz8z8z8z8z8z8z8z",
	}, nil
}

// 2. 관리자 관련 Mock
func (m *MockRepository) GetAdminByID(id string) (*domain.Admin, error) {
	return &domain.Admin{
		AdminNumber: 1,
		AdminID:     id,
		AdminPW:     "$2a$10$7cQkLrgVQGuNCvYyONufFOwO3EwmBl1H.1lJ1y906WRBaNTH2t1Fe",
		FKGussID:    sql.NullInt64{Int64: 1, Valid: true},
	}, nil
}

// 3. 체육관 관련 Mock
// [수정 완료] 인터페이스 요구사항에 맞춰 string 파라미터를 추가했습니다.
func (m *MockRepository) GetGyms(filter string) ([]domain.Gym, error) {
	log.Printf("[MOCK] GetGyms called with filter: %s", filter)
	return []domain.Gym{
		{GussNumber: 1, GussName: "Mock 강남점", GussStatus: "OPEN", GussSize: 50, GussUserCount: 10},
	}, nil
}

func (m *MockRepository) GetGymDetail(id int64) (*domain.Gym, error) {
	return &domain.Gym{
		GussNumber:    id,
		GussName:      "Mock 상세 지점",
		GussSize:      50,
		GussUserCount: 5,
	}, nil
}

// 4. 예약 관련 Mock
func (m *MockRepository) CreateReservation(userNum, gymNum int64, resTime time.Time) (string, error) {
	log.Printf("[MOCK] Reservation Created: User %d -> Gym %d at %v", userNum, gymNum, resTime)
	return "CONFIRMED", nil
}

func (m *MockRepository) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	return []domain.Reservation{}, nil
}

func (m *MockRepository) CancelReservation(resNum int64, userNum int64, status string) error {
	log.Printf("[MOCK] Reservation Cancelled: ResNum %d, UserNum %d, Status %s", resNum, userNum, status)
	return nil
}

// 5. 기구 관리 Mock
func (m *MockRepository) GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error) {
	return []domain.Equipment{
		{ID: 1, Name: "Mock 트레드밀", Category: "유산소", Quantity: 5, Status: "active"},
	}, nil
}

func (m *MockRepository) AddEquipment(eq *domain.Equipment) error {
	log.Printf("[MOCK] Equipment Added: %s", eq.Name)
	return nil
}

func (m *MockRepository) UpdateEquipment(eq *domain.Equipment) error {
	log.Printf("[MOCK] Equipment Updated: ID %d", eq.ID)
	return nil
}

func (m *MockRepository) DeleteEquipment(eqID int64) error {
	log.Printf("[MOCK] Equipment Deleted: ID %d", eqID)
	return nil
}

// 6. 매출 관련 Mock
func (m *MockRepository) GetSalesByGym(gymID int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"type": "Membership", "amount": 100000, "date": "2026-01-13"},
	}, nil
}

// --- MockLogRepository ---
type MockLogRepository struct{}

func NewMockLogRepository() LogRepository {
	return &MockLogRepository{}
}

func (m *MockLogRepository) SaveEqLog(gID int64, eID string, stat string) error {
	return nil
}

func (m *MockLogRepository) SaveUserLog(uID string, act string) error {
	return nil
}

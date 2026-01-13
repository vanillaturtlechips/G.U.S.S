package repository

import (
	"guss-backend/internal/domain"
	"log"
)

type MockRepository struct{}

func NewMockRepository() Repository {
	return &MockRepository{}
}

func (m *MockRepository) CreateUser(u *domain.User) error {
	log.Printf("[MOCK] User Created: %v", u.UserName)
	return nil
}

func (m *MockRepository) GetUserByID(id string) (*domain.User, error) {
	return &domain.User{UserNumber: 1, UserID: id, UserName: "MockUser"}, nil
}

func (m *MockRepository) GetGyms() ([]domain.Gym, error) {
	return []domain.Gym{
		{GussNumber: 1, GussName: "Mock 강남점", GussStatus: "OPEN"},
	}, nil
}

func (m *MockRepository) GetGymDetail(id int64) (*domain.Gym, error) {
	return &domain.Gym{GussNumber: id, GussName: "Mock 상세 지점"}, nil
}

func (m *MockRepository) CreateReservation(userNum, gymNum int64) (string, error) {
	return "CONFIRMED", nil
}

func (m *MockRepository) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	return []domain.Reservation{}, nil
}

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
	return nil
}

func (m *MockRepository) DeleteEquipment(eqID int64) error {
	return nil
}

func (m *MockRepository) GetSalesByGym(gymID int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"type": "Membership", "amount": 100000, "date": "2026-01-13"},
	}, nil
}

// LogRepository Mock
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

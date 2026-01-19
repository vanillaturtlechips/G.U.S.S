package repository

import (
	"guss-backend/internal/domain"
	"time"
)

type MockRepository struct{}

func NewMockRepository() Repository { return &MockRepository{} }

func (m *MockRepository) CreateUser(u *domain.User) error { return nil }
func (m *MockRepository) GetUserByID(id string) (*domain.User, error) { return &domain.User{}, nil }
func (m *MockRepository) UpdateFCMToken(uID, t string) error { return nil }
func (m *MockRepository) GetFCMToken(uID string) (string, error) { return "", nil }
func (m *MockRepository) GetAllGyms() ([]domain.Gym, error) { return []domain.Gym{}, nil }
func (m *MockRepository) GetGymDetail(id int64) (*domain.Gym, error) { return &domain.Gym{}, nil }
func (m *MockRepository) IncrementUserCount(id int64) error { return nil }
func (m *MockRepository) CreateReservation(u, g int64, t time.Time) (string, error) { return "SUCCESS", nil }
func (m *MockRepository) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return []domain.Reservation{}, nil }
func (m *MockRepository) GetActiveReservationByUser(userNum int64) (*domain.Reservation, error) { return nil, nil }
func (m *MockRepository) CancelReservation(resID int64, userNum int64) error { return nil }
func (m *MockRepository) GetAdminByID(id string) (*domain.Admin, error) { return &domain.Admin{}, nil }
func (m *MockRepository) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) { return []domain.Equipment{}, nil }
func (m *MockRepository) AddEquipment(e *domain.Equipment) error { return nil }
func (m *MockRepository) UpdateEquipment(e *domain.Equipment) error { return nil }
func (m *MockRepository) DeleteEquipment(id int64) error { return nil }
func (m *MockRepository) GetSalesByGym(id int64) ([]domain.Sale, error) { return []domain.Sale{}, nil }

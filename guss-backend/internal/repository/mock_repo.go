package repository

import (
	"errors"
	"guss-backend/internal/domain"
	"time"
)

type MockRepository struct {
	users map[string]*domain.User
}

func NewMockRepository() Repository {
	return &MockRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockRepository) CreateUser(u *domain.User) error { m.users[u.UserID] = u; return nil }
func (m *MockRepository) GetUserByID(id string) (*domain.User, error) {
	if u, ok := m.users[id]; ok { return u, nil }
	return nil, errors.New("not found")
}
func (m *MockRepository) GetAdminByID(id string) (*domain.Admin, error) { return &domain.Admin{}, nil }
func (m *MockRepository) GetAllGyms() ([]domain.Gym, error) { return []domain.Gym{}, nil }
func (m *MockRepository) GetGymDetail(id int64) (*domain.Gym, error) { return &domain.Gym{}, nil }
func (m *MockRepository) CreateReservation(uN, gN int64, vT time.Time) (string, error) { return "CONFIRMED", nil }
func (m *MockRepository) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) { return []domain.Equipment{}, nil }
func (m *MockRepository) AddEquipment(eq *domain.Equipment) error { return nil }
func (m *MockRepository) UpdateEquipment(eq *domain.Equipment) error { return nil }
func (m *MockRepository) DeleteEquipment(id int64) error { return nil }
func (r *MockRepository) CancelReservation(rN, uN int64, role string) error { return nil }
func (r *MockRepository) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return []domain.Reservation{}, nil }

// MockRepository가 Repository 인터페이스를 만족하도록 메서드 추가
func (m *MockRepository) UpdateFCMToken(userID string, token string) error {
    return nil
}

func (m *MockRepository) GetFCMToken(userID string) (string, error) {
    return "", nil
}

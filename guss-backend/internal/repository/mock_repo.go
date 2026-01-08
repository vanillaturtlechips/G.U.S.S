package repository

import (
	"errors"
	"fmt"
	"guss-backend/internal/domain"
	"time"
)

// MockMySQLRepo: 메모리에 데이터를 저장하는 가짜 DB
type MockMySQLRepo struct {
	users map[string]*domain.User
	gyms  map[int64]*domain.Gym
	revs  []domain.Reservation
}

func NewMockRepository() Repository {
	repo := &MockMySQLRepo{
		users: make(map[string]*domain.User),
		gyms:  make(map[int64]*domain.Gym),
	}
	// 테스트용 Mock 데이터 삽입
	repo.gyms[1] = &domain.Gym{GussNumber: 1, GussName: "GUSS 강남점", GussStatus: "open", GussUserCount: 15, GussSize: 50}
	repo.gyms[2] = &domain.Gym{GussNumber: 2, GussName: "GUSS 홍대점", GussStatus: "open", GussUserCount: 40, GussSize: 50}
	return repo
}

func (m *MockMySQLRepo) CreateUser(u *domain.User) error {
	m.users[u.UserID] = u
	return nil
}

func (m *MockMySQLRepo) GetUserByID(id string) (*domain.User, error) {
	if u, ok := m.users[id]; ok { return u, nil }
	return nil, errors.New("user not found")
}

func (m *MockMySQLRepo) GetAllGyms() ([]domain.Gym, error) {
	var list []domain.Gym
	for _, g := range m.gyms { list = append(list, *g) }
	return list, nil
}

func (m *MockMySQLRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	if g, ok := m.gyms[id]; ok { return g, nil }
	return nil, errors.New("gym not found")
}

func (m *MockMySQLRepo) CreateReservation(uID, gID int64) error {
	m.revs = append(m.revs, domain.Reservation{FKUserID: uID, FKGussID: gID, RevsTime: time.Now()})
	return nil
}

// MockDynamoLogRepo: 로그를 콘솔에 출력만 함
type MockDynamoLogRepo struct{}

func NewMockLogRepository() LogRepository { return &MockDynamoLogRepo{} }

func (m *MockDynamoLogRepo) SaveEqLog(gID int64, eID, stat string) error {
	fmt.Printf("[Mock Log] Gym %d, Equip %s: %s\n", gID, eID, stat)
	return nil
}

func (m *MockDynamoLogRepo) SaveUserLog(uID, act string) error {
	fmt.Printf("[Mock Log] User %s: %s\n", uID, act)
	return nil
}
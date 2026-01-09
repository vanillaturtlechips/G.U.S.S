package repository

import (
	"errors"
	"fmt"
	"guss-backend/internal/domain"
	"time"
)

// MockMySQLRepo: 메모리에 데이터를 저장하는 테스트용 가짜 DB
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
	// 테스트용 기본 데이터 삽입
	repo.gyms[1] = &domain.Gym{GussNumber: 1, GussName: "GUSS 강남점", GussStatus: "open", GussUserCount: 15, GussSize: 50}
	repo.gyms[2] = &domain.Gym{GussNumber: 2, GussName: "GUSS 홍대점", GussStatus: "open", GussUserCount: 40, GussSize: 50}

	// 테스트용 유저 데이터 (testuser03 포함)
	repo.users["testuser03"] = &domain.User{UserNumber: 3, UserID: "testuser03", UserName: "테스트유저03"}

	return repo
}

func (m *MockMySQLRepo) CreateUser(u *domain.User) error {
	m.users[u.UserID] = u
	return nil
}

func (m *MockMySQLRepo) GetUserByID(id string) (*domain.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockMySQLRepo) GetAllGyms() ([]domain.Gym, error) {
	var list []domain.Gym
	for _, g := range m.gyms {
		list = append(list, *g)
	}
	return list, nil
}

func (m *MockMySQLRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	if g, ok := m.gyms[id]; ok {
		return g, nil
	}
	return nil, errors.New("gym not found")
}

// [수정] 중복 예약 체크 로직 추가 및 (string, error) 반환
func (m *MockMySQLRepo) CreateReservation(uID, gID int64) (string, error) {
	// 1. 중복 예약 여부 전수 조사
	for _, r := range m.revs {
		if r.FKUserID == uID && r.FKGussID == gID {
			// 이미 예약이 있다면 DUPLICATE 반환
			return "DUPLICATE", nil
		}
	}

	// 2. 새로운 예약 추가
	m.revs = append(m.revs, domain.Reservation{
		FKUserID:   uID,
		FKGussID:   gID,
		RevsTime:   time.Now(),
		RevsStatus: "CONFIRMED",
	})

	return "CONFIRMED", nil
}

// MockDynamoLogRepo: 로그 기록 대행
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

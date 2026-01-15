package repository

import (
	"guss-backend/internal/domain"
	"time"
)

// --- 일반 데이터용 Mock ---
type mockRepo struct{}

func NewMockRepository() Repository {
	return &mockRepo{}
}

func (m *mockRepo) GetUserByID(id string) (*domain.User, error)                      { return &domain.User{}, nil }
func (m *mockRepo) GetAdminByID(id string) (*domain.Admin, error)                    { return &domain.Admin{}, nil }
func (m *mockRepo) CreateUser(u *domain.User) error                                  { return nil }
func (m *mockRepo) GetGyms() ([]domain.Gym, error)                                   { return []domain.Gym{}, nil }
func (m *mockRepo) GetGymDetail(id int64) (*domain.Gym, error)                       { return &domain.Gym{}, nil }
func (m *mockRepo) CreateReservation(userNum, gymNum int64) (string, error)          { return "SUCCESS", nil }
func (m *mockRepo) CreateReservationWithTime(u int64, g int64, s, e time.Time) error { return nil }
func (m *mockRepo) UpdateReservationStatus(r int64, u int64, s string) error         { return nil }
func (m *mockRepo) GetReservationsByGym(id int64) ([]domain.Reservation, error) {
	return []domain.Reservation{}, nil
}
func (m *mockRepo) GetHourlyReservationStats(id int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (m *mockRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) {
	return []domain.Equipment{}, nil
}
func (m *mockRepo) AddEquipment(e *domain.Equipment) error    { return nil }
func (m *mockRepo) UpdateEquipment(e *domain.Equipment) error { return nil }
func (m *mockRepo) DeleteEquipment(id int64) error            { return nil }
func (m *mockRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// --- 로그용 Mock (오류 해결 핵심 포인트) ---
type mockLogRepo struct{}

func NewMockLogRepository() LogRepository {
	return &mockLogRepo{}
}

// 만약 LogRepository 인터페이스에 메서드가 정의되어 있다면 여기에 빈 메서드를 추가해야 합니다.
// 현재 repository.go 기준으로는 빈 인터페이스이므로 이대로면 컴파일 에러가 사라집니다.

package repository

import "guss-backend/internal/domain"

type Repository interface {
	// User 관련
	CreateUser(u *domain.User) error
	GetUserByID(id string) (*domain.User, error)

	// Gym 관련 (GetGyms로 이름 변경하여 핸들러와 통일)
	GetGyms() ([]domain.Gym, error)
	GetGymDetail(id int64) (*domain.Gym, error)

	// Reservation 관련
	CreateReservation(userNum, gymNum int64) (string, error)
	GetReservationsByGym(gymID int64) ([]domain.Reservation, error)

	GetAdminByID(id string) (*domain.Admin, error)

	// Equipment 관련 (메서드 명칭 통일)
	GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error)
	AddEquipment(eq *domain.Equipment) error // domain 객체를 받도록 설정
	UpdateEquipment(eq *domain.Equipment) error
	DeleteEquipment(eqID int64) error

	// 매출 관련
	GetSalesByGym(gymID int64) ([]map[string]interface{}, error)
}

type LogRepository interface {
	SaveEqLog(gID int64, eID string, stat string) error
	SaveUserLog(uID string, act string) error
}

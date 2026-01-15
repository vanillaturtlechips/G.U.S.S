package repository

import (
	"guss-backend/internal/domain"
	"time"
)

type Repository interface {
	GetUserByID(id string) (*domain.User, error)
	GetAdminByID(id string) (*domain.Admin, error)
	CreateUser(u *domain.User) error

	GetGyms() ([]domain.Gym, error)
	GetGymDetail(id int64) (*domain.Gym, error)

	CreateReservation(userNum, gymNum int64) (string, error)
	CreateReservationWithTime(userNum int64, gymID int64, start, end time.Time) error
	UpdateReservationStatus(resID int64, userNum int64, status string) error
	GetReservationsByGym(gymID int64) ([]domain.Reservation, error)
	GetHourlyReservationStats(gymID int64) ([]map[string]interface{}, error)

	GetEquipmentsByGymID(id int64) ([]domain.Equipment, error)
	AddEquipment(eq *domain.Equipment) error
	UpdateEquipment(eq *domain.Equipment) error
	DeleteEquipment(id int64) error

	GetSalesByGym(id int64) ([]map[string]interface{}, error)
}

type LogRepository interface {
	// 로그 관련 로직 필요 시 추가
}

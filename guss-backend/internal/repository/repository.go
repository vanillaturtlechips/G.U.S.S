package repository

import (
	"guss-backend/internal/domain"
	"time"
)

type Repository interface {
	CreateUser(u *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetAdminByID(id string) (*domain.Admin, error)

	GetGyms(search string) ([]domain.Gym, error) // search 추가
	GetGymDetail(id int64) (*domain.Gym, error)

	CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error) // visitTime 추가
	GetReservationsByGym(gymID int64) ([]domain.Reservation, error)
	CancelReservation(revsNum, userNum int64, role string) error // 신규 추가

	GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error)
	AddEquipment(eq *domain.Equipment) error
	UpdateEquipment(eq *domain.Equipment) error
	DeleteEquipment(eqID int64) error
	GetSalesByGym(gymID int64) ([]map[string]interface{}, error)
}

type LogRepository interface {
	SaveEqLog(gID int64, eID string, stat string) error
	SaveUserLog(uID string, act string) error
}

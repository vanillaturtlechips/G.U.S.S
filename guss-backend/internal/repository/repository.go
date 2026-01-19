package repository

import (
	"guss-backend/internal/domain"
	"time"
)

type Repository interface {
	CreateUser(u *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	UpdateFCMToken(userID string, token string) error
	GetFCMToken(userID string) (string, error)

	GetAllGyms() ([]domain.Gym, error)
	GetGymDetail(id int64) (*domain.Gym, error)
	IncrementUserCount(gymID int64) error

	CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error)
	GetReservationsByGym(gymID int64) ([]domain.Reservation, error)
	GetActiveReservationByUser(userNum int64) (*domain.Reservation, error)
	CancelReservation(resID int64, userNum int64) error

	GetAdminByID(id string) (*domain.Admin, error)

	GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error)
	AddEquipment(eq *domain.Equipment) error
	UpdateEquipment(eq *domain.Equipment) error
	DeleteEquipment(eqID int64) error

	GetSalesByGym(gymID int64) ([]domain.Sale, error)
}

type LogRepository interface {
	SaveEqLog(gID int64, eID string, stat string) error
	SaveUserLog(uID string, act string) error
}

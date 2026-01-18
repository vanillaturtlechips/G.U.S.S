package repository

import (
	"guss-backend/internal/domain"
	"time"
)

type Repository interface {
	// 유저 및 관리자 관련
	CreateUser(u *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetAdminByID(id string) (*domain.Admin, error)

	// 체육관 및 장비 관련
	GetAllGyms() ([]domain.Gym, error)
	GetGymDetail(id int64) (*domain.Gym, error)
	GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error)
	AddEquipment(eq *domain.Equipment) error
	UpdateEquipment(eq *domain.Equipment) error
	DeleteEquipment(id int64) error

	// 예약 관련
	CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error)
	CancelReservation(revsNum, userNum int64, role string) error
	GetReservationsByGym(gymID int64) ([]domain.Reservation, error)

	UpdateFCMToken(userID string, token string) error // 추가
        GetFCMToken(userID string) (string, error)       // 추가
}

type LogRepository interface {
	SaveEqLog(gID int64, eID, stat string) error
	SaveUserLog(uID, act string) error
}

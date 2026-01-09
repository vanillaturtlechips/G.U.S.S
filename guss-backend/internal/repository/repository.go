// // internal/repository/repository.go
// package repository

// import (
// 	"database/sql"
// 	"guss-backend/internal/domain" // 이전에 정의한 모델 참조
// )

// type Repository interface {
// 	// User 관련
// 	CreateUser(u *domain.User) error
// 	GetUserByID(userID string) (*domain.User, error)

// 	// Gym 관련
// 	GetAllGyms() ([]domain.Gym, error)
// 	GetGymDetail(gussNumber int64) (*domain.Gym, error)

// 	// Reservation 관련
// 	CreateReservation(userID int64, gymID int64) error
// }

// type mysqlRepo struct {
// 	db *sql.DB
// }

// func NewMySQLRepository(db *sql.DB) Repository {
// 	return &mysqlRepo{db: db}
// }

// // CreateUser: 회원가입
// func (r *mysqlRepo) CreateUser(u *domain.User) error {
// 	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
// 	_, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
// 	return err
// }

// // GetUserByID: 로그인 검증용
// func (r *mysqlRepo) GetUserByID(userID string) (*domain.User, error) {
// 	u := &domain.User{}
// 	query := `SELECT user_number, user_name, user_id, user_pw FROM user_table WHERE user_id = ?`
// 	err := r.db.QueryRow(query, userID).Scan(&u.UserNumber, &u.UserName, &u.UserID, &u.UserPW)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return u, nil
// }

// // GetAllGyms: 대시보드 지도 마커용
// func (r *mysqlRepo) GetAllGyms() ([]domain.Gym, error) {
// 	query := `SELECT guss_number, guss_name, guss_status, guss_user_count, guss_size FROM guss_table`
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var gyms []domain.Gym
// 	for rows.Next() {
// 		var g domain.Gym
// 		rows.Scan(&g.GussNumber, &g.GussName, &g.GussStatus, &g.GussUserCount, &g.GussSize)
// 		gyms = append(gyms, g)
// 	}
// 	return gyms, nil
// }

// // CreateReservation: 예약 저장
// func (r *mysqlRepo) CreateReservation(userID int64, gymID int64) error {
// 	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) VALUES (?, ?, 'COMPLETED', NOW())`
// 	_, err := r.db.Exec(query, userID, gymID)
// 	return err
// }

// func (r *mysqlRepo) GetGymDetail(gussNumber int64) (*domain.Gym, error) {
// 	g := &domain.Gym{}
// 	query := `SELECT guss_number, guss_name, guss_address, guss_phone, guss_status, guss_user_count, guss_size FROM guss_table WHERE guss_number = ?`
// 	err := r.db.QueryRow(query, gussNumber).Scan(
// 		&g.GussNumber, &g.GussName, &g.GussAddress, &g.GussPhone,
// 		&g.GussStatus, &g.GussUserCount, &g.GussSize,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return g, nil
// }

package repository

import "guss-backend/internal/domain"

type Repository interface {
	CreateUser(u *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetAllGyms() ([]domain.Gym, error)
	GetGymDetail(id int64) (*domain.Gym, error)
	// 반환 타입을 (string, error)로 반드시 수정!
	CreateReservation(userNum, gymNum int64) (string, error)
}

// LogRepository: 로그용 규칙
type LogRepository interface {
	SaveEqLog(gID int64, eID, stat string) error
	SaveUserLog(uID, act string) error
}

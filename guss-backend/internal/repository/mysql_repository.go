package repository

import (
	"database/sql"
	"guss-backend/internal/domain"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

// 회원가입
func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	return err
}

// 로그인 (ID로 유저 정보 조회)
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT user_number, user_id, user_pw FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserID, &u.UserPW)
	return u, err
}

// 전체 체육관 목록 조회
func (r *mysqlRepo) GetAllGyms() ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_status, guss_user_count, guss_size FROM guss_table`
	rows, err := r.db.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var gyms []domain.Gym
	for rows.Next() {
		var g domain.Gym
		err := rows.Scan(&g.GussNumber, &g.GussName, &g.GussStatus, &g.GussUserCount, &g.GussSize)
		if err != nil { continue }
		gyms = append(gyms, g)
	}
	return gyms, nil
}

// 체육관 상세 정보 조회 (혼잡도 계산 핵심 데이터)
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	g := &domain.Gym{}
	query := `SELECT guss_number, guss_name, guss_user_count, guss_size, guss_ma_type, guss_ma_count FROM guss_table WHERE guss_number = ?`
	err := r.db.QueryRow(query, id).Scan(&g.GussNumber, &g.GussName, &g.GussUserCount, &g.GussSize, &g.GussMaType, &g.GussMaCount)
	return g, err
}

// 예약 생성 (revs_table 활용)
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64) error {
	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status) VALUES (?, ?, 'CONFIRMED')`
	_, err := r.db.Exec(query, userNum, gymNum)
	return err
}
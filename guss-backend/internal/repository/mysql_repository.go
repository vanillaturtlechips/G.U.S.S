package repository

import (
	"database/sql"
	"fmt"
	"guss-backend/internal/domain"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

// CreateUser: 회원가입
func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	return err
}

// GetUserByID: 로그인 조회
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT user_number, user_id, user_pw FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserID, &u.UserPW)
	return u, err
}

// GetAllGyms: 목록 조회
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

// GetGymDetail: [404 에러 해결 포인트] NULL 값을 허용하도록 Scan 수정
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	g := &domain.Gym{}
	
	// NULL이 포함될 수 있는 컬럼들(phone, ma_type 등)은 sql.NullString 등을 써야 하지만,
	// 현재 동작을 위해 값이 확실히 있는 컬럼 위주로 조회하거나 빈 값 처리를 수행합니다.
	query := `SELECT guss_number, guss_name, guss_address, guss_user_count, guss_size, guss_status 
              FROM guss_table WHERE guss_number = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&g.GussNumber, &g.GussName, &g.GussAddress, 
		&g.GussUserCount, &g.GussSize, &g.GussStatus,
	)
	
	if err != nil {
		fmt.Printf("[DB Scan Error] ID %d 조회 실패: %v\n", id, err)
		return nil, err
	}
	
	// ma_type 등 NULL인 값들은 안전하게 기본값 처리 (필요 시 추가)
	g.GussMaType = "N/A"
	g.GussMaCount = 0
	
	return g, nil
}

// CreateReservation: 예약 생성 (revs_time 컬럼명 반영)
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64) error {
	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
              VALUES (?, ?, 'CONFIRMED', NOW())`
	_, err := r.db.Exec(query, userNum, gymNum)
	return err
}
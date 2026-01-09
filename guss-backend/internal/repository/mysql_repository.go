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

// 1. 회원가입
func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	return err
}

// 2. 로그인 조회
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT user_number, user_id, user_pw, user_name FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserID, &u.UserPW, &u.UserName)
	return u, err
}

// 3. 체육관 목록 조회 (전화번호 포함)
func (r *mysqlRepo) GetAllGyms() ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_phone, guss_address, guss_status, guss_user_count, guss_size FROM guss_table`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gyms []domain.Gym
	for rows.Next() {
		var g domain.Gym
		err := rows.Scan(&g.GussNumber, &g.GussName, &g.GussPhone, &g.GussAddress, &g.GussStatus, &g.GussUserCount, &g.GussSize)
		if err != nil {
			continue
		}
		gyms = append(gyms, g)
	}
	return gyms, nil
}

// 4. 체육관 상세 조회 (기구 관리 데이터 포함)
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	g := &domain.Gym{}
	query := `SELECT guss_number, guss_name, guss_address, guss_phone, guss_user_count, guss_size, guss_status, 
                     guss_ma_type, guss_ma_count, guss_ma_state
              FROM guss_table WHERE guss_number = ?`

	err := r.db.QueryRow(query, id).Scan(
		&g.GussNumber, &g.GussName, &g.GussAddress, &g.GussPhone,
		&g.GussUserCount, &g.GussSize, &g.GussStatus,
		&g.GussMaType, &g.GussMaCount, &g.GussMaState,
	)

	if err != nil {
		return nil, err
	}
	return g, nil
}

// 5. 예약 생성 (중복 체크 + 정원 체크 + 상태 반환)
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64) (string, error) {
	// [추가] 중복 예약 확인 (아직 이용 중인 내역이 있는지)
	var exists int
	checkQuery := `SELECT COUNT(*) FROM revs_table 
                   WHERE fk_user_number = ? AND fk_guss_number = ? 
                   AND revs_status IN ('CONFIRMED', 'WAITING')`
	err := r.db.QueryRow(checkQuery, userNum, gymNum).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists > 0 {
		return "DUPLICATE", nil
	}

	// 1) 현재 체육관의 정원 상태 확인
	var currentCount, maxSize int
	err = r.db.QueryRow("SELECT guss_user_count, guss_size FROM guss_table WHERE guss_number = ?", gymNum).Scan(&currentCount, &maxSize)
	if err != nil {
		return "", err
	}

	// 2) 정원에 따라 상태 결정
	status := "CONFIRMED"
	if currentCount >= maxSize {
		status = "WAITING"
	}

	// 3) 예약 정보 삽입
	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
              VALUES (?, ?, ?, NOW())`
	_, err = r.db.Exec(query, userNum, gymNum, status)
	if err != nil {
		return "", err
	}

	// 4) CONFIRMED일 경우에만 실시간 이용 인원수 업데이트
	if status == "CONFIRMED" {
		_, updateErr := r.db.Exec("UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?", gymNum)
		if updateErr != nil {
			fmt.Printf("[Warning] 이용 인원 업데이트 실패: %v\n", updateErr)
		}
	}

	return status, nil
}

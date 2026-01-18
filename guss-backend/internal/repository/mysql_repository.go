package repository

import (
	"database/sql"
	"errors"
	"guss-backend/internal/domain"
	"time"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

// [핵심] 예약 생성: 인원수(UPDATE) 로직을 완전히 제거했습니다.
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error) {
	var count int
	// 1. 중복 체크: 이미 활성화된 예약이 있는지 확인
	checkQuery := `SELECT COUNT(*) FROM revs_table WHERE fk_user_number = ? AND revs_status = 'CONFIRMED'`
	err := r.db.QueryRow(checkQuery, userNum).Scan(&count)
	if err != nil { return "", err }
	
	if count > 0 { 
		return "DUPLICATE", errors.New("이미 활성화된 예약이 존재합니다.") 
	}

	// 2. 예약 데이터만 삽입 (guss_table 업데이트 코드는 여기 없습니다)
	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
              VALUES (?, ?, 'CONFIRMED', ?)`
	_, err = r.db.Exec(query, userNum, gymNum, visitTime)
	
	return "SUCCESS", err
}

// [핵심] 인원 증가: 이 함수는 나중에 QR 체크인 핸들러에서만 호출됩니다.
func (r *mysqlRepo) IncrementUserCount(gymID int64) error {
	query := `UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`
	_, err := r.db.Exec(query, gymID)
	return err
}

// --- 아래는 빌드 호환성을 위한 나머지 메서드들입니다 ---
func (r *mysqlRepo) UpdateFCMToken(userID, token string) error {
	_, err := r.db.Exec(`UPDATE user_table SET fcm_token = ? WHERE user_id = ?`, token, userID)
	return err
}
func (r *mysqlRepo) GetFCMToken(userID string) (string, error) {
	var token sql.NullString
	err := r.db.QueryRow(`SELECT fcm_token FROM user_table WHERE user_id = ?`, userID).Scan(&token)
	return token.String, err
}
func (r *mysqlRepo) CreateUser(u *domain.User) error {
	_, err := r.db.Exec(`INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	return err
}
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(`SELECT user_number, user_id, user_pw, user_name FROM user_table WHERE user_id = ?`, id).Scan(&u.UserNumber, &u.UserID, &u.UserPW, &u.UserName)
	return u, err
}
func (r *mysqlRepo) GetAdminByID(id string) (*domain.Admin, error) {
	a := &domain.Admin{}
	err := r.db.QueryRow(`SELECT admin_number, admin_id, admin_pw, fk_guss_number FROM admin_table WHERE admin_id = ?`, id).Scan(&a.AdminNumber, &a.AdminID, &a.AdminPW, &a.FKGussID)
	return a, err
}
func (r *mysqlRepo) GetAllGyms() ([]domain.Gym, error) {
	rows, err := r.db.Query(`SELECT guss_number, guss_name, guss_phone, guss_address, guss_status, guss_user_count, guss_size FROM guss_table`)
	if err != nil { return nil, err }
	defer rows.Close()
	var gyms []domain.Gym
	for rows.Next() {
		var g domain.Gym
		rows.Scan(&g.GussNumber, &g.GussName, &g.GussPhone, &g.GussAddress, &g.GussStatus, &g.GussUserCount, &g.GussSize)
		gyms = append(gyms, g)
	}
	return gyms, nil
}
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	g := &domain.Gym{}
	err := r.db.QueryRow(`SELECT guss_number, guss_name, guss_address, guss_phone, guss_user_count, guss_size, guss_status FROM guss_table WHERE guss_number = ?`, id).Scan(&g.GussNumber, &g.GussName, &g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize, &g.GussStatus)
	return g, err
}
func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) { return []domain.Equipment{}, nil }
func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error { return nil }
func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error { return nil }
func (r *mysqlRepo) DeleteEquipment(id int64) error { return nil }
func (r *mysqlRepo) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return []domain.Reservation{}, nil }
func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) { return []map[string]interface{}{}, nil }

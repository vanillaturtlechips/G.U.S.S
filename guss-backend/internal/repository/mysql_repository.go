package repository

import (
	"database/sql"
	"guss-backend/internal/domain"
	"time"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

// FCM 토큰 업데이트 로직 추가
func (r *mysqlRepo) UpdateFCMToken(userID string, token string) error {
	query := `UPDATE user_table SET fcm_token = ? WHERE user_id = ?`
	_, err := r.db.Exec(query, token, userID)
	return err
}

// FCM 토큰 조회 로직 추가
func (r *mysqlRepo) GetFCMToken(userID string) (string, error) {
	var token sql.NullString
	query := `SELECT fcm_token FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, userID).Scan(&token)
	return token.String, err
}

func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	return err
}

func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT user_number, user_id, user_pw, user_name FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserID, &u.UserPW, &u.UserName)
	return u, err
}

func (r *mysqlRepo) GetAdminByID(id string) (*domain.Admin, error) {
	a := &domain.Admin{}
	query := `SELECT admin_number, admin_id, admin_pw, fk_guss_number FROM admin_table WHERE admin_id = ?`
	err := r.db.QueryRow(query, id).Scan(&a.AdminNumber, &a.AdminID, &a.AdminPW, &a.FKGussID)
	return a, err
}

func (r *mysqlRepo) GetAllGyms() ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_phone, guss_address, guss_status, guss_user_count, guss_size FROM guss_table`
	rows, err := r.db.Query(query)
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
	query := `SELECT guss_number, guss_name, guss_address, guss_phone, guss_user_count, guss_size, guss_status FROM guss_table WHERE guss_number = ?`
	err := r.db.QueryRow(query, id).Scan(&g.GussNumber, &g.GussName, &g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize, &g.GussStatus)
	return g, err
}

func (r *mysqlRepo) CreateReservation(userNum, gymNum int64) (string, error) {
	var count int
	checkQuery := `SELECT COUNT(*) FROM revs_table WHERE fk_user_number = ? AND revs_status = 'CONFIRMED'`
	err := r.db.QueryRow(checkQuery, userNum).Scan(&count)
	if err != nil { return "", err }
	if count > 0 { return "DUPLICATE", errors.New("이미 활성화된 예약이 존재합니다.") }

	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
              VALUES (?, ?, 'CONFIRMED', NOW())`
	_, err = r.db.Exec(query, userNum, gymNum)
	return "SUCCESS", err
}

func (r *mysqlRepo) ProcessEntry(userNum, gymNum int64) error {
	tx, err := r.db.Begin()
	if err != nil { return err }
	defer tx.Rollback()

	// 1. 예약 상태를 USED로 변경
	res, err := tx.Exec(`UPDATE revs_table SET revs_status = 'USED' 
                         WHERE fk_user_number = ? AND fk_guss_number = ? AND revs_status = 'CONFIRMED' LIMIT 1`, 
                         userNum, gymNum)
	if err != nil { return err }
	affected, _ := res.RowsAffected()
	if affected == 0 { return errors.New("유효한 예약 내역이 없습니다.") }

	// 2. 체육관 실시간 인원 +1
	_, err = tx.Exec(`UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`, gymNum)
	if err != nil { return err }

	return tx.Commit()
}

func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) { return []domain.Equipment{}, nil }
func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error { return nil }
func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error { return nil }
func (r *mysqlRepo) DeleteEquipment(id int64) error { return nil }
func (r *mysqlRepo) CancelReservation(rN, uN int64, role string) error { return nil }
func (r *mysqlRepo) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return []domain.Reservation{}, nil }

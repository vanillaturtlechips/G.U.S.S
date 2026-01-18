package repository

import (
	"database/sql"
	"errors" // 추가
	"guss-backend/internal/domain"
	"time"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) UpdateFCMToken(userID string, token string) error {
	query := `UPDATE user_table SET fcm_token = ? WHERE user_id = ?`
	_, err := r.db.Exec(query, token, userID)
	return err
}

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

// 인터페이스에 맞춰 visitTime 추가
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error) {
	var count int
	checkQuery := `SELECT COUNT(*) FROM revs_table WHERE fk_user_number = ? AND revs_status = 'CONFIRMED'`
	err := r.db.QueryRow(checkQuery, userNum).Scan(&count)
	if err != nil { return "", err }
	if count > 0 { return "DUPLICATE", errors.New("이미 활성화된 예약이 존재합니다.") }

	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
              VALUES (?, ?, 'CONFIRMED', ?)`
	_, err = r.db.Exec(query, userNum, gymNum, visitTime)
	return "SUCCESS", err
}

func (r *mysqlRepo) IncrementUserCount(gymID int64) error {
	_, err := r.db.Exec(`UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`, gymID)
	return err
}

func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) { return []domain.Equipment{}, nil }
func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error { return nil }
func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error { return nil }
func (r *mysqlRepo) DeleteEquipment(id int64) error { return nil }
func (r *mysqlRepo) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return []domain.Reservation{}, nil }
func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) { return []map[string]interface{}{}, nil }

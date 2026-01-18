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

func (r *mysqlRepo) CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error) {
	var count int
	checkQuery := `SELECT COUNT(*) FROM revs_table WHERE fk_user_number = ? AND revs_status = 'CONFIRMED'`
	err := r.db.QueryRow(checkQuery, userNum).Scan(&count)
	if err != nil {
		return "", err
	}

	if count > 0 {
		return "DUPLICATE", errors.New("이미 활성화된 예약이 존재합니다.")
	}

	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
              VALUES (?, ?, 'CONFIRMED', ?)`
	_, err = r.db.Exec(query, userNum, gymNum, visitTime)

	return "SUCCESS", err
}

func (r *mysqlRepo) IncrementUserCount(gymID int64) error {
	query := `UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`
	_, err := r.db.Exec(query, gymID)
	return err
}

// 🔥 활성 예약 조회
func (r *mysqlRepo) GetActiveReservationByUser(userNum int64) (*domain.Reservation, error) {
	res := &domain.Reservation{}
	query := `SELECT revs_number, fk_user_number, fk_guss_number, revs_time, revs_status 
	          FROM revs_table 
	          WHERE fk_user_number = ? AND revs_status = 'CONFIRMED' 
	          ORDER BY revs_time DESC LIMIT 1`

	err := r.db.QueryRow(query, userNum).Scan(
		&res.RevsNumber,
		&res.FKUserID,
		&res.FKGussID,
		&res.VisitTime,
		&res.RevsStatus,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 🔥 예약 취소
func (r *mysqlRepo) CancelReservation(resID int64, userNum int64) error {
	query := `UPDATE revs_table 
	          SET revs_status = 'CANCELLED' 
	          WHERE revs_number = ? AND fk_user_number = ? AND revs_status = 'CONFIRMED'`

	result, err := r.db.Exec(query, resID, userNum)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("취소할 수 있는 예약이 없습니다.")
	}

	return nil
}

// 🔥 Admin - 예약 로그 조회
func (r *mysqlRepo) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	query := `SELECT r.revs_number, u.user_name, r.revs_time, r.revs_status, r.revs_time as visit_time
	          FROM revs_table r
	          JOIN user_table u ON r.fk_user_number = u.user_number
	          WHERE r.fk_guss_number = ?
	          ORDER BY r.revs_time DESC
	          LIMIT 100`

	rows, err := r.db.Query(query, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		err := rows.Scan(
			&res.RevsNumber,
			&res.UserName,
			&res.RevsTime,
			&res.RevsStatus,
			&res.VisitTime,
		)
		if err != nil {
			continue
		}
		reservations = append(reservations, res)
	}

	return reservations, nil
}

// 🔥 Admin - 매출 로그 조회
func (r *mysqlRepo) GetSalesByGym(gymID int64) ([]domain.Sale, error) {
	query := `SELECT sales_number, sales_date, sales_amount, sales_type
	          FROM sales_table
	          WHERE fk_guss_number = ?
	          ORDER BY sales_date DESC
	          LIMIT 100`

	rows, err := r.db.Query(query, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		var sale domain.Sale
		err := rows.Scan(
			&sale.SalesNumber,
			&sale.SalesDate,
			&sale.SalesAmount,
			&sale.SalesType,
		)
		if err != nil {
			continue
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

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
	_, err := r.db.Exec(`INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`,
		u.UserName, u.UserPhone, u.UserID, u.UserPW)
	return err
}

func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(`SELECT user_number, user_id, user_pw, user_name FROM user_table WHERE user_id = ?`, id).
		Scan(&u.UserNumber, &u.UserID, &u.UserPW, &u.UserName)
	return u, err
}

func (r *mysqlRepo) GetAdminByID(id string) (*domain.Admin, error) {
	a := &domain.Admin{}
	err := r.db.QueryRow(`SELECT admin_number, admin_id, admin_pw, fk_guss_number FROM admin_table WHERE admin_id = ?`, id).
		Scan(&a.AdminNumber, &a.AdminID, &a.AdminPW, &a.FKGussID)
	return a, err
}

func (r *mysqlRepo) GetAllGyms() ([]domain.Gym, error) {
	rows, err := r.db.Query(`SELECT guss_number, guss_name, guss_phone, guss_address, guss_status, guss_user_count, guss_size 
	                         FROM guss_table`)
	if err != nil {
		return nil, err
	}
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
	err := r.db.QueryRow(`SELECT guss_number, guss_name, guss_address, guss_phone, guss_user_count, guss_size, guss_status 
	                       FROM guss_table WHERE guss_number = ?`, id).
		Scan(&g.GussNumber, &g.GussName, &g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize, &g.GussStatus)
	return g, err
}

func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) {
	return []domain.Equipment{}, nil
}

func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error {
	return nil
}

func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error {
	return nil
}

func (r *mysqlRepo) DeleteEquipment(id int64) error {
	return nil
}

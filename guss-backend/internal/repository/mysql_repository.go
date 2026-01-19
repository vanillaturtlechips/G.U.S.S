package repository

import (
	"database/sql"
	"errors"
	"guss-backend/internal/domain"
	"log"
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

func (r *mysqlRepo) GetActiveReservationByUser(userNum int64) (*domain.Reservation, error) {
	res := &domain.Reservation{}
	
	query := `SELECT revs_number, fk_user_number, fk_guss_number, revs_time, revs_status 
	          FROM revs_table 
	          WHERE fk_user_number = ? AND revs_status = 'CONFIRMED' 
	          ORDER BY revs_time DESC LIMIT 1`

	var revsTime string
	
	err := r.db.QueryRow(query, userNum).Scan(
		&res.RevsNumber,
		&res.FKUserID,
		&res.FKGussID,
		&revsTime,
		&res.RevsStatus,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parsedTime, _ := time.Parse("2006-01-02 15:04:05", revsTime)
	res.VisitTime = parsedTime
	res.RevsTime = parsedTime

	return res, nil
}

func (r *mysqlRepo) CancelReservation(resID int64, userNum int64) error {
	var gymID int64
	checkQuery := `SELECT fk_guss_number FROM revs_table 
	               WHERE revs_number = ? AND fk_user_number = ? AND revs_status = 'CONFIRMED'`
	
	err := r.db.QueryRow(checkQuery, resID, userNum).Scan(&gymID)
	if err == sql.ErrNoRows {
		return errors.New("취소할 수 있는 예약이 없습니다.")
	}
	if err != nil {
		return err
	}

	updateQuery := `UPDATE revs_table 
	                SET revs_status = 'CANCELLED' 
	                WHERE revs_number = ? AND fk_user_number = ? AND revs_status = 'CONFIRMED'`

	_, err = r.db.Exec(updateQuery, resID, userNum)
	if err != nil {
		return err
	}

	decrementQuery := `UPDATE guss_table 
	                   SET guss_user_count = GREATEST(guss_user_count - 1, 0) 
	                   WHERE guss_number = ?`
	
	_, err = r.db.Exec(decrementQuery, gymID)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepo) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	query := `SELECT r.revs_number, u.user_name, r.revs_time, r.revs_status, r.revs_time as visit_time
	          FROM revs_table r
	          JOIN user_table u ON r.fk_user_number = u.user_number
	          WHERE r.fk_guss_number = ?
	          ORDER BY r.revs_time DESC
	          LIMIT 100`

	rows, err := r.db.Query(query, gymID)
	if err != nil {
		log.Printf("[GetReservationsByGym] Query Error: %v", err)
		return []domain.Reservation{}, nil
	}
	defer rows.Close()

	var reservations []domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		var revsTimeStr, visitTimeStr string
		
		err := rows.Scan(
			&res.RevsNumber,
			&res.UserName,
			&revsTimeStr,
			&res.RevsStatus,
			&visitTimeStr,
		)
		if err != nil {
			log.Printf("[GetReservationsByGym] Scan Error: %v", err)
			continue
		}

		res.RevsTime, _ = time.Parse("2006-01-02 15:04:05", revsTimeStr)
		res.VisitTime, _ = time.Parse("2006-01-02 15:04:05", visitTimeStr)
		
		reservations = append(reservations, res)
	}

	log.Printf("[GetReservationsByGym] Final count: %d", len(reservations))

	if reservations == nil {
		return []domain.Reservation{}, nil
	}

	return reservations, nil
}

func (r *mysqlRepo) GetSalesByGym(gymID int64) ([]domain.Sale, error) {
	query := `SELECT sales_number, sales_date, sales_amount, sales_type
	          FROM sales_table
	          WHERE fk_guss_number = ?
	          ORDER BY sales_date DESC
	          LIMIT 100`

	rows, err := r.db.Query(query, gymID)
	if err != nil {
		log.Printf("[GetSalesByGym] Query Error: %v", err)
		return []domain.Sale{}, nil
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		var sale domain.Sale
		var salesDateStr string
		
		err := rows.Scan(
			&sale.SalesNumber,
			&salesDateStr,
			&sale.SalesAmount,
			&sale.SalesType,
		)
		if err != nil {
			log.Printf("[GetSalesByGym] Scan Error: %v", err)
			continue
		}

		sale.SalesDate, _ = time.Parse("2006-01-02 15:04:05", salesDateStr)
		
		sales = append(sales, sale)
	}

	log.Printf("[GetSalesByGym] Final count: %d", len(sales))

	if sales == nil {
		return []domain.Sale{}, nil
	}

	return sales, nil
}

func (r *mysqlRepo) GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error) {
	query := `SELECT eq_number, fk_guss_number, eq_name, eq_category, eq_quantity, eq_status, eq_purchase_date
	          FROM equipments_table
	          WHERE fk_guss_number = ?
	          ORDER BY eq_number DESC`

	rows, err := r.db.Query(query, gymID)
	if err != nil {
		log.Printf("[GetEquipmentsByGymID] Query Error: %v", err)
		return []domain.Equipment{}, nil
	}
	defer rows.Close()

	var equipments []domain.Equipment
	for rows.Next() {
		var eq domain.Equipment
		var purchaseDateStr sql.NullString // 🔥 NullTime → NullString 변경

		err := rows.Scan(
			&eq.ID,
			&eq.GymID,
			&eq.Name,
			&eq.Category,
			&eq.Quantity,
			&eq.Status,
			&purchaseDateStr, // 🔥 문자열로 받기
		)
		if err != nil {
			log.Printf("[GetEquipmentsByGymID] Scan Error: %v", err)
			continue
		}

		// 🔥 문자열 그대로 할당
		if purchaseDateStr.Valid {
			eq.PurchaseDate = purchaseDateStr.String
		}

		equipments = append(equipments, eq)
	}

	log.Printf("[GetEquipmentsByGymID] Final count: %d", len(equipments))

	if equipments == nil {
		return []domain.Equipment{}, nil
	}

	return equipments, nil
}

func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error {
	query := `INSERT INTO equipments_table (fk_guss_number, eq_name, eq_category, eq_quantity, eq_status, eq_purchase_date)
	          VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, eq.GymID, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.PurchaseDate)
	return err
}

func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error {
	query := `UPDATE equipments_table 
	          SET eq_name = ?, eq_category = ?, eq_quantity = ?, eq_status = ?, eq_purchase_date = ?
	          WHERE eq_number = ?`

	result, err := r.db.Exec(query, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.PurchaseDate, eq.ID)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("수정할 기구를 찾을 수 없습니다.")
	}

	return nil
}

func (r *mysqlRepo) DeleteEquipment(eqID int64) error {
	query := `DELETE FROM equipments_table WHERE eq_number = ?`

	result, err := r.db.Exec(query, eqID)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("삭제할 기구를 찾을 수 없습니다.")
	}

	return nil
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

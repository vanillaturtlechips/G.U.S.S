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

func (r *mysqlRepo) GetGyms(search string) ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_status, 
               COALESCE(guss_address, ''), COALESCE(guss_phone, ''), 
               guss_user_count, guss_size,
               COALESCE(guss_open_time, ''), COALESCE(guss_close_time, '') 
               FROM guss_table WHERE guss_name LIKE ?`

	rows, err := r.db.Query(query, "%"+search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gyms := []domain.Gym{}
	for rows.Next() {
		var g domain.Gym
		err := rows.Scan(&g.GussNumber, &g.GussName, &g.GussStatus, &g.GussAddress, &g.GussPhone,
			&g.GussUserCount, &g.GussSize, &g.GussOpenTime, &g.GussCloseTime)
		if err == nil {
			gyms = append(gyms, g)
		}
	}
	return gyms, nil
}

func (r *mysqlRepo) CreateReservation(userNum, gymNum int64, visitTime time.Time) (string, error) {
	var count int
	r.db.QueryRow(`SELECT COUNT(*) FROM revs_table WHERE fk_user_number = ? AND revs_status = 'CONFIRMED'`, userNum).Scan(&count)
	if count > 0 {
		return "", errors.New("이미 활성화된 예약이 존재합니다")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, visit_time) VALUES (?, ?, 'CONFIRMED', ?)`, userNum, gymNum, visitTime)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`, gymNum)
	if err != nil {
		return "", err
	}

	return "SUCCESS", tx.Commit()
}

func (r *mysqlRepo) CancelReservation(revsNum, userNum int64, role string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var gymID int64
	var currentStatus string
	err = tx.QueryRow(`SELECT fk_guss_number, revs_status FROM revs_table WHERE revs_number = ? AND (fk_user_number = ? OR ? = 'ADMIN')`, revsNum, userNum, role).Scan(&gymID, &currentStatus)
	if err != nil {
		return errors.New("예약을 찾을 수 없거나 권한이 없습니다")
	}
	if currentStatus == "CANCELLED" {
		return errors.New("이미 취소된 예약입니다")
	}

	_, err = tx.Exec(`UPDATE revs_table SET revs_status = 'CANCELLED' WHERE revs_number = ?`, revsNum)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE guss_table SET guss_user_count = guss_user_count - 1 WHERE guss_number = ? AND guss_user_count > 0`, gymID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`SELECT user_number, user_name, user_phone, user_id, user_pw FROM user_table WHERE user_id = ?`, id).Scan(&u.UserNumber, &u.UserName, &u.UserPhone, &u.UserID, &u.UserPW)
	return &u, err
}

func (r *mysqlRepo) CreateUser(u *domain.User) error {
	res, err := r.db.Exec(`INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	if err == nil {
		u.UserNumber, _ = res.LastInsertId()
	}
	return err
}

func (r *mysqlRepo) GetAdminByID(id string) (*domain.Admin, error) {
	var a domain.Admin
	err := r.db.QueryRow(`SELECT admin_number, admin_id, admin_pw, fk_guss_number FROM admin_table WHERE admin_id = ?`, id).Scan(&a.AdminNumber, &a.AdminID, &a.AdminPW, &a.FKGussID)
	return &a, err
}

func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	var g domain.Gym
	err := r.db.QueryRow(`SELECT guss_number, guss_name, guss_status, COALESCE(guss_address,''), COALESCE(guss_phone,''), guss_user_count, guss_size, COALESCE(guss_open_time,''), COALESCE(guss_close_time,'') FROM guss_table WHERE guss_number = ?`, id).Scan(&g.GussNumber, &g.GussName, &g.GussStatus, &g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize, &g.GussOpenTime, &g.GussCloseTime)
	return &g, err
}

func (r *mysqlRepo) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	rows, err := r.db.Query(`SELECT r.revs_number, r.fk_user_number, r.fk_guss_number, r.revs_status, r.revs_time, u.user_name FROM revs_table r JOIN user_table u ON r.fk_user_number = u.user_number WHERE r.fk_guss_number = ?`, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []domain.Reservation{}
	for rows.Next() {
		var res domain.Reservation
		rows.Scan(&res.RevsNumber, &res.FKUserID, &res.FKGussID, &res.RevsStatus, &res.RevsTime, &res.UserName)
		list = append(list, res)
	}
	return list, nil
}

func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) {
	rows, err := r.db.Query(`SELECT equip_id, fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date FROM equipment_table WHERE fk_guss_number = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []domain.Equipment{}
	for rows.Next() {
		var e domain.Equipment
		rows.Scan(&e.ID, &e.GymID, &e.Name, &e.Category, &e.Quantity, &e.Status, &e.PurchaseDate)
		list = append(list, e)
	}
	return list, nil
}

func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error {
	_, err := r.db.Exec(`INSERT INTO equipment_table (fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date) VALUES (?, ?, ?, ?, ?, ?)`, eq.GymID, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.PurchaseDate)
	return err
}

func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error {
	_, err := r.db.Exec(`UPDATE equipment_table SET equip_name=?, equip_category=?, equip_quantity=?, equip_status=? WHERE equip_id=?`, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.ID)
	return err
}

func (r *mysqlRepo) DeleteEquipment(id int64) error {
	_, err := r.db.Exec(`DELETE FROM equipment_table WHERE equip_id = ?`, id)
	return err
}

func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

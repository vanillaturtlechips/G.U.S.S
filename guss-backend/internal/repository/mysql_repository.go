package repository

import (
	"database/sql"
	"fmt"
	"guss-backend/internal/domain"
	"time"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) GetGyms() ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_status, COALESCE(guss_address, ''), COALESCE(guss_phone, ''), guss_user_count, guss_size, COALESCE(guss_open_time, ''), COALESCE(guss_close_time, '') FROM guss_table`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	gyms := []domain.Gym{}
	for rows.Next() {
		var g domain.Gym
		rows.Scan(&g.GussNumber, &g.GussName, &g.GussStatus, &g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize, &g.GussOpenTime, &g.GussCloseTime)
		gyms = append(gyms, g)
	}
	return gyms, nil
}

func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	var g domain.Gym
	query := `SELECT guss_number, guss_name, guss_status, COALESCE(guss_address, ''), COALESCE(guss_phone, ''), guss_user_count, guss_size, COALESCE(guss_open_time, ''), COALESCE(guss_close_time, '') FROM guss_table WHERE guss_number = ?`
	err := r.db.QueryRow(query, id).Scan(&g.GussNumber, &g.GussName, &g.GussStatus, &g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize, &g.GussOpenTime, &g.GussCloseTime)
	return &g, err
}

func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT user_number, user_name, user_phone, user_id, user_pw FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserName, &u.UserPhone, &u.UserID, &u.UserPW)
	return &u, err
}

func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	u.UserNumber, _ = res.LastInsertId()
	return err
}

func (r *mysqlRepo) CreateReservation(u, g int64) (string, error) {
	return "SUCCESS", r.CreateReservationWithTime(u, g, time.Now(), time.Now().Add(30*time.Minute))
}

func (r *mysqlRepo) CreateReservationWithTime(uNum, gymID int64, start, end time.Time) error {
	query := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, start_time, end_time, revs_time) VALUES (?, ?, 'CONFIRMED', ?, ?, NOW())`
	_, err := r.db.Exec(query, uNum, gymID, start, end)
	return err
}

func (r *mysqlRepo) UpdateReservationStatus(resID, uNum int64, status string) error {
	query := `UPDATE revs_table SET revs_status = ? WHERE revs_number = ? AND fk_user_number = ?`
	_, err := r.db.Exec(query, status, resID, uNum)
	return err
}

func (r *mysqlRepo) GetHourlyReservationStats(gymID int64) ([]map[string]interface{}, error) {
	query := `SELECT HOUR(start_time) as h, COUNT(*) as c FROM revs_table WHERE fk_guss_number = ? AND revs_status = 'CONFIRMED' GROUP BY h ORDER BY h`
	rows, _ := r.db.Query(query, gymID)
	defer rows.Close()
	res := []map[string]interface{}{}
	for rows.Next() {
		var h, c int
		rows.Scan(&h, &c)
		res = append(res, map[string]interface{}{"hour": fmt.Sprintf("%02d", h), "count": c})
	}
	return res, nil
}

func (r *mysqlRepo) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	query := `SELECT r.revs_number, r.fk_user_number, r.fk_guss_number, r.revs_status, r.revs_time, u.user_name FROM revs_table r JOIN user_table u ON r.fk_user_number = u.user_number WHERE r.fk_guss_number = ?`
	rows, _ := r.db.Query(query, gymID)
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
	query := `SELECT equip_id, fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date FROM equipment_table WHERE fk_guss_number = ?`
	rows, _ := r.db.Query(query, id)
	defer rows.Close()
	list := []domain.Equipment{}
	for rows.Next() {
		var e domain.Equipment
		rows.Scan(&e.ID, &e.GymID, &e.Name, &e.Category, &e.Quantity, &e.Status, &e.PurchaseDate)
		list = append(list, e)
	}
	return list, nil
}

func (r *mysqlRepo) AddEquipment(e *domain.Equipment) error {
	query := `INSERT INTO equipment_table (fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, e.GymID, e.Name, e.Category, e.Quantity, e.Status, e.PurchaseDate)
	return err
}

func (r *mysqlRepo) UpdateEquipment(e *domain.Equipment) error {
	query := `UPDATE equipment_table SET equip_name=?, equip_category=?, equip_quantity=?, equip_status=? WHERE equip_id=?`
	_, err := r.db.Exec(query, e.Name, e.Category, e.Quantity, e.Status, e.ID)
	return err
}

func (r *mysqlRepo) DeleteEquipment(id int64) error {
	_, err := r.db.Exec(`DELETE FROM equipment_table WHERE equip_id = ?`, id)
	return err
}

func (r *mysqlRepo) GetAdminByID(id string) (*domain.Admin, error) {
	var a domain.Admin
	query := `SELECT admin_number, admin_id, admin_pw, fk_guss_number FROM admin_table WHERE admin_id = ?`
	err := r.db.QueryRow(query, id).Scan(&a.AdminNumber, &a.AdminID, &a.AdminPW, &a.FKGussID)
	return &a, err
}

func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

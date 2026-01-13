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

// 1. 모든 체육관 조회 (guss_table)
func (r *mysqlRepo) GetGyms() ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_status, 
               COALESCE(guss_address, ''), COALESCE(guss_phone, ''), 
               guss_user_count, guss_size FROM guss_table`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gyms []domain.Gym
	for rows.Next() {
		var g domain.Gym
		err := rows.Scan(&g.GussNumber, &g.GussName, &g.GussStatus,
			&g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize)
		if err != nil {
			continue
		}
		gyms = append(gyms, g)
	}
	return gyms, nil
}

// 2. 체육관 상세 조회
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	var g domain.Gym
	query := `SELECT guss_number, guss_name, guss_status, 
                     COALESCE(guss_address, ''), COALESCE(guss_phone, ''), 
                     guss_user_count, guss_size 
              FROM guss_table WHERE guss_number = ?`
	err := r.db.QueryRow(query, id).Scan(&g.GussNumber, &g.GussName, &g.GussStatus,
		&g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize)
	return &g, err
}

// 3. 유저 ID로 조회 (로그인용)
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT user_number, user_name, user_phone, user_id, user_pw FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserName, &u.UserPhone, &u.UserID, &u.UserPW)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// 4. 회원가입
func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	if err != nil {
		return err
	}
	u.UserNumber, _ = result.LastInsertId()
	return nil
}

// 5. 예약 생성
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
                      VALUES (?, ?, 'CONFIRMED', NOW())`, userNum, gymNum)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`, gymNum)
	if err != nil {
		return "", err
	}

	return "SUCCESS", tx.Commit()
}

// 6. 예약 목록 조회 (에러 났던 부분 수정)
func (r *mysqlRepo) GetReservationsByGym(gymID int64) ([]domain.Reservation, error) {
	query := `SELECT r.revs_number, r.fk_user_number, r.fk_guss_number, r.revs_status, r.revs_time, u.user_name
	          FROM revs_table r
	          JOIN user_table u ON r.fk_user_number = u.user_number
	          WHERE r.fk_guss_number = ?`

	rows, err := r.db.Query(query, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		// Scan 대상을 domain 필드명과 정확히 매칭 (FKUserID, FKGussID 등)
		err := rows.Scan(&res.RevsNumber, &res.FKUserID, &res.FKGussID, &res.RevsStatus, &res.RevsTime, &res.UserName)
		if err != nil {
			continue
		}
		list = append(list, res)
	}
	return list, nil
}

// 7. 기구 관련 (필드명 매칭: ID, GymID, Name 등)
func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) {
	query := `SELECT equip_id, fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date 
              FROM equipment_table WHERE fk_guss_number = ?`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Equipment
	for rows.Next() {
		var e domain.Equipment
		rows.Scan(&e.ID, &e.GymID, &e.Name, &e.Category, &e.Quantity, &e.Status, &e.PurchaseDate)
		list = append(list, e)
	}
	return list, nil
}

func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error {
	query := `INSERT INTO equipment_table (fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date) 
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, eq.GymID, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.PurchaseDate)
	return err
}

func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error {
	query := `UPDATE equipment_table SET equip_name=?, equip_category=?, equip_quantity=?, equip_status=? WHERE equip_id=?`
	_, err := r.db.Exec(query, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.ID)
	return err
}

func (r *mysqlRepo) DeleteEquipment(id int64) error {
	_, err := r.db.Exec(`DELETE FROM equipment_table WHERE equip_id = ?`, id)
	return err
}

func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) { return nil, nil }

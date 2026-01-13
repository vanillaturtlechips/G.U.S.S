package repository

import (
	"database/sql"
	"guss-backend/internal/domain"
	"log"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

// 1. 모든 체육관 조회 (guss_table 참조 및 컴파일 에러 필드 제거)
func (r *mysqlRepo) GetGyms() ([]domain.Gym, error) {
	// GussMaxSize를 쿼리 및 스캔에서 제외하여 컴파일 에러 해결
	query := `SELECT guss_number, guss_name, guss_status, 
		       COALESCE(guss_address, ''), COALESCE(guss_phone, ''), 
		       guss_user_count, guss_size FROM guss_table`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("[DB ERROR] GetGyms: %v", err)
		return nil, err
	}
	defer rows.Close()

	var gyms []domain.Gym
	for rows.Next() {
		var g domain.Gym
		// domain.Gym에 없는 GussMaxSize 필드 스캔 제거
		if err := rows.Scan(&g.GussNumber, &g.GussName, &g.GussStatus,
			&g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize); err != nil {
			log.Printf("[SCAN ERROR] GetGyms: %v", err)
			return nil, err
		}
		gyms = append(gyms, g)
	}
	return gyms, nil
}

// 2. 상세 조회 (guss_table 참조)
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error) {
	var g domain.Gym
	query := `SELECT guss_number, guss_name, guss_status, 
	                 COALESCE(guss_address, ''), COALESCE(guss_phone, ''), 
	                 guss_user_count, guss_size 
	          FROM guss_table WHERE guss_number = ?`
	err := r.db.QueryRow(query, id).Scan(&g.GussNumber, &g.GussName, &g.GussStatus,
		&g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize)
	if err != nil {
		log.Printf("[DB ERROR] GetGymDetail(%d): %v", id, err)
	}
	return &g, err
}

// 3. 기구 추가 (equipment_table 참조)
func (r *mysqlRepo) AddEquipment(eq *domain.Equipment) error {
	query := `INSERT INTO equipment_table (fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, eq.GymID, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.PurchaseDate)
	if err != nil {
		log.Printf("[DB ERROR] AddEquipment: %v", err)
	}
	return err
}

// 4. 기구 삭제 (equipment_table 참조)
func (r *mysqlRepo) DeleteEquipment(id int64) error {
	query := `DELETE FROM equipment_table WHERE equip_id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("[DB ERROR] DeleteEquipment(%d): %v", id, err)
	}
	return err
}

// 5. 기구 목록 조회 (equipment_table 참조)
func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) {
	query := `SELECT equip_id, fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date 
	          FROM equipment_table WHERE fk_guss_number = ?`
	rows, err := r.db.Query(query, id)
	if err != nil {
		log.Printf("[DB ERROR] GetEquipments: %v", err)
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

// 6. 예약 생성 (revs_table 인서트 및 guss_table 인원 업데이트)
// Error 1452 해결을 위해 부모 테이블 명칭(guss_table)을 정확히 사용
func (r *mysqlRepo) CreateReservation(userID, gymID int64) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 예약 내역 저장 (revs_table)
	queryInsert := `INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status) VALUES (?, ?, 'CONFIRMED')`
	_, err = tx.Exec(queryInsert, userID, gymID)
	if err != nil {
		log.Printf("[DB ERROR] CreateReservation (Insert): %v", err)
		return "", err
	}

	// 현재 인원(guss_user_count) 증가 (guss_table)
	queryUpdate := `UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`
	_, err = tx.Exec(queryUpdate, gymID)
	if err != nil {
		log.Printf("[DB ERROR] CreateReservation (Update): %v", err)
		return "", err
	}

	err = tx.Commit()
	return "RESERVATION_SUCCESS", err
}

// 나머지 인터페이스 만족용 빈 함수들
func (r *mysqlRepo) CreateUser(u *domain.User) error                             { return nil }
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error)                 { return nil, nil }
func (r *mysqlRepo) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return nil, nil }
func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error)    { return nil, nil }
func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error                  { return nil }

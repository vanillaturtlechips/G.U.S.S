package repository

import (
	"database/sql"
	"errors"
	"guss-backend/internal/domain"
	"log"
)

type mysqlRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

// 1. 모든 체육관 조회 (운영 시간 컬럼 포함 9개 필드 매칭)
func (r *mysqlRepo) GetGyms() ([]domain.Gym, error) {
	// [수정] DB 스키마 업데이트에 따라 guss_open_time, guss_close_time 추가
	query := `SELECT guss_number, guss_name, guss_status, 
               COALESCE(guss_address, ''), COALESCE(guss_phone, ''), 
               guss_user_count, guss_size,
               COALESCE(guss_open_time, ''), COALESCE(guss_close_time, '') 
               FROM guss_table`
	
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("[DB ERROR] GetGyms Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 프론트엔드에서 .map() 에러가 나지 않도록 nil이 아닌 빈 슬라이스로 초기화
	gyms := []domain.Gym{} 

	for rows.Next() {
		var g domain.Gym
		// [중요] Scan 인자 개수는 쿼리 컬럼 개수(9개)와 반드시 일치해야 함
		err := rows.Scan(
			&g.GussNumber, &g.GussName, &g.GussStatus,
			&g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize,
			&g.GussOpenTime, &g.GussCloseTime,
		)
		if err != nil {
			log.Printf("[DB ERROR] GetGyms Scan: %v", err)
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
                     guss_user_count, guss_size,
                     COALESCE(guss_open_time, ''), COALESCE(guss_close_time, '')
              FROM guss_table WHERE guss_number = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&g.GussNumber, &g.GussName, &g.GussStatus,
		&g.GussAddress, &g.GussPhone, &g.GussUserCount, &g.GussSize,
		&g.GussOpenTime, &g.GussCloseTime,
	)
	if err != nil {
		log.Printf("[DB ERROR] GetGymDetail(%d): %v", id, err)
		return nil, err
	}
	return &g, err
}

// 3. 유저 ID로 조회 (로그인 인증용)
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT user_number, user_name, user_phone, user_id, user_pw FROM user_table WHERE user_id = ?`
	err := r.db.QueryRow(query, id).Scan(&u.UserNumber, &u.UserName, &u.UserPhone, &u.UserID, &u.UserPW)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("사용자를 찾을 수 없습니다")
		}
		return nil, err
	}
	return &u, nil
}

// 4. 회원가입 처리
func (r *mysqlRepo) CreateUser(u *domain.User) error {
	query := `INSERT INTO user_table (user_name, user_phone, user_id, user_pw) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, u.UserName, u.UserPhone, u.UserID, u.UserPW)
	if err != nil {
		log.Printf("[DB ERROR] CreateUser: %v", err)
		return err
	}
	u.UserNumber, _ = result.LastInsertId()
	return nil
}

// 5. 예약 생성 (중복 예약 및 노쇼 방지 로직 포함)
func (r *mysqlRepo) CreateReservation(userNum, gymNum int64) (string, error) {
	// [체크] 이미 활성화된 예약이 있는지 확인
	var count int
	checkQuery := `SELECT COUNT(*) FROM revs_table WHERE fk_user_number = ? AND revs_status = 'CONFIRMED'`
	err := r.db.QueryRow(checkQuery, userNum).Scan(&count)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("이미 활성화된 예약이 존재합니다. 노쇼 방지를 위해 추가 예약은 불가합니다.")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 예약 정보 삽입
	_, err = tx.Exec(`INSERT INTO revs_table (fk_user_number, fk_guss_number, revs_status, revs_time) 
                      VALUES (?, ?, 'CONFIRMED', NOW())`, userNum, gymNum)
	if err != nil {
		return "", err
	}

	// 체육관 현재 이용 인원수 증가
	_, err = tx.Exec(`UPDATE guss_table SET guss_user_count = guss_user_count + 1 WHERE guss_number = ?`, gymNum)
	if err != nil {
		return "", err
	}

	return "SUCCESS", tx.Commit()
}

// 6. 예약 목록 조회 (관리자용)
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

	list := []domain.Reservation{} // 빈 배열 초기화
	for rows.Next() {
		var res domain.Reservation
		err := rows.Scan(&res.RevsNumber, &res.FKUserID, &res.FKGussID, &res.RevsStatus, &res.RevsTime, &res.UserName)
		if err != nil {
			continue
		}
		list = append(list, res)
	}
	return list, nil
}

// 7. 기구 관리 로직 (프론트엔드 map 에러 방지 적용)
func (r *mysqlRepo) GetEquipmentsByGymID(id int64) ([]domain.Equipment, error) {
	query := `SELECT equip_id, fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date 
              FROM equipment_table WHERE fk_guss_number = ?`
	
	rows, err := r.db.Query(query, id)
	if err != nil {
		log.Printf("[DB ERROR] GetEquipments: %v", err)
		return nil, err
	}
	defer rows.Close()

	// [중요] 데이터가 없어도 nil이 아닌 [] 슬라이스를 반환하여 프론트엔드 map 에러 방지
	list := []domain.Equipment{} 
	for rows.Next() {
		var e domain.Equipment
		err := rows.Scan(&e.ID, &e.GymID, &e.Name, &e.Category, &e.Quantity, &e.Status, &e.PurchaseDate)
		if err != nil {
			log.Printf("[DB ERROR] Scan Equipment: %v", err)
			continue
		}
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

func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil 
}
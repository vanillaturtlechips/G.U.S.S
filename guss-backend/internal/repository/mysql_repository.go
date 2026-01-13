package repository

import (
	"database/sql"
	"guss-backend/internal/domain"
)

type mysqlRepo struct {
	db *sql.DB
}

// NewMySQLRepository: main.go에서 호출할 수 있도록 추가
func NewMySQLRepository(db *sql.DB) Repository {
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) GetGyms() ([]domain.Gym, error) {
	query := `SELECT guss_number, guss_name, guss_address, guss_phone, guss_status, guss_user_count, guss_size FROM gyms`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gyms []domain.Gym
	for rows.Next() {
		var g domain.Gym
		rows.Scan(&g.GussNumber, &g.GussName, &g.GussAddress, &g.GussPhone, &g.GussStatus, &g.GussUserCount, &g.GussSize)
		gyms = append(gyms, g)
	}
	return gyms, nil
}

func (r *mysqlRepo) GetEquipmentsByGymID(gymID int64) ([]domain.Equipment, error) {
	query := `SELECT id, gym_id, name, category, quantity, status, purchase_date FROM equipments WHERE gym_id = ?`
	rows, err := r.db.Query(query, gymID)
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
	query := `INSERT INTO equipments (gym_id, name, category, quantity, status, purchase_date) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, eq.GymID, eq.Name, eq.Category, eq.Quantity, eq.Status, eq.PurchaseDate)
	return err
}

func (r *mysqlRepo) DeleteEquipment(eqID int64) error {
	_, err := r.db.Exec("DELETE FROM equipments WHERE id = ?", eqID)
	return err
}

// 나머지 인터페이스 스텁 (컴파일 에러 방지용)
func (r *mysqlRepo) CreateUser(u *domain.User) error                             { return nil }
func (r *mysqlRepo) GetUserByID(id string) (*domain.User, error)                 { return nil, nil }
func (r *mysqlRepo) GetGymDetail(id int64) (*domain.Gym, error)                  { return nil, nil }
func (r *mysqlRepo) CreateReservation(u, g int64) (string, error)                { return "", nil }
func (r *mysqlRepo) GetReservationsByGym(id int64) ([]domain.Reservation, error) { return nil, nil }
func (r *mysqlRepo) GetSalesByGym(id int64) ([]map[string]interface{}, error)    { return nil, nil }
func (r *mysqlRepo) UpdateEquipment(eq *domain.Equipment) error                  { return nil }

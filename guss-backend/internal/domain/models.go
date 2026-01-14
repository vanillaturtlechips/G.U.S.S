package domain

import (
	"database/sql"
	"time"
)

// 1. 사용자 정보 (user_table)
type User struct {
	UserNumber int64  `json:"user_number" db:"user_number"`
	UserName   string `json:"user_name"   db:"user_name"`
	UserPhone  string `json:"user_phone"  db:"user_phone"`
	UserID     string `json:"user_id"     db:"user_id"`
	UserPW     string `json:"user_pw"     db:"user_pw"`
}

// 2. 체육관 정보 (guss_table)
type Gym struct {
	GussNumber    int64  `json:"guss_number"    db:"guss_number"`
	GussName      string `json:"guss_name"      db:"guss_name"`
	GussAddress   string `json:"guss_address"   db:"guss_address"`
	GussPhone     string `json:"guss_phone"     db:"guss_phone"`
	GussStatus    string `json:"guss_status"    db:"guss_status"`
	GussUserCount int    `json:"guss_user_count" db:"guss_user_count"` // 현재 이용 인원
	GussSize      int    `json:"guss_size"       db:"guss_size"`       // 현재 최대 이용 인원
	GussOpenTime  string `json:"guss_open_time"  db:"guss_open_time"`
	GussCloseTime string `json:"guss_close_time" db:"guss_close_time"`
}

// 3. 기구 정보 (equipment_table)
type Equipment struct {
	ID           int64  `json:"id"             db:"equip_id"`
	GymID        int64  `json:"gym_id"         db:"fk_guss_number"` // fk_guss_number -> gym_id
	Name         string `json:"name"           db:"equip_name"`     // equip_name -> name
	Category     string `json:"category"       db:"equip_category"` // equip_category -> category
	Quantity     int    `json:"quantity"       db:"equip_quantity"` // equip_quantity -> quantity
	Status       string `json:"status"         db:"equip_status"`   // equip_status -> status
	PurchaseDate string `json:"purchaseDate"   db:"purchase_date"`  // purchase_date -> purchaseDate
}

// 4. 예약 정보 (revs_table)
type Reservation struct {
	RevsNumber int64     `json:"revs_number"    db:"revs_number"`
	FKUserID   int64     `json:"fk_user_number" db:"fk_user_number"`
	FKGussID   int64     `json:"fk_guss_number" db:"fk_guss_number"`
	RevsTime   time.Time `json:"revs_time"      db:"revs_time"`
	RevsStatus string    `json:"revs_status"    db:"revs_status"`
	UserName   string    `json:"user_name,omitempty"`
}

// 5. 관리자 정보 (admin_table)
type Admin struct {
	AdminNumber int64         `json:"admin_number"   db:"admin_number"`
	AdminID     string        `json:"admin_id"       db:"admin_id"`
	AdminPW     string        `json:"-"              db:"admin_pw"`
	FKGussID    sql.NullInt64 `json:"fk_guss_number"`
}

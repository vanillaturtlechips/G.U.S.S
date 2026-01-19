package domain

import (
	"database/sql"
	"time"
)

// User: 회원 정보
type User struct {
	UserNumber int64  `json:"user_number" db:"user_number"`
	UserName   string `json:"user_name"   db:"user_name"`
	UserPhone  string `json:"user_phone"  db:"user_phone"`
	UserID     string `json:"user_id"     db:"user_id"`
	UserPW     string `json:"user_pw"     db:"user_pw"` // 비밀번호는 JSON 제외
}

// Gym: 체육관 정보 (영업시간 추가)
type Gym struct {
	GussNumber    int64  `json:"guss_number"    db:"guss_number"`
	GussName      string `json:"guss_name"      db:"guss_name"`
	GussAddress   string `json:"guss_address"   db:"guss_address"`
	GussPhone     string `json:"guss_phone"     db:"guss_phone"`
	GussStatus    string `json:"guss_status"    db:"guss_status"`
	GussUserCount int    `json:"guss_user_count" db:"guss_user_count"`
	GussSize      int    `json:"guss_size"       db:"guss_size"`
	GussOpenTime  string `json:"guss_open_time"  db:"guss_open_time"`  // 추가
	GussCloseTime string `json:"guss_close_time" db:"guss_close_time"` // 추가
}

// Equipment: 기구 관리 (어드민 페이지용 추가)
type Equipment struct {
	ID            int64  `json:"id" db:"eq_number"`
	GymID         int64  `json:"gym_id" db:"fk_guss_number"`
	Name          string `json:"name" db:"eq_name"`
	Category      string `json:"category" db:"eq_category"`
	Quantity      int    `json:"quantity" db:"eq_quantity"`
	Status        string `json:"status" db:"eq_status"`
	PurchaseDate  string `json:"purchaseDate" db:"eq_purchase_date"`
}

// Reservation: 예약 정보 (방문 시간 추가)
type Reservation struct {
	RevsNumber int64     `json:"revs_number"    db:"revs_number"`
	FKUserID   int64     `json:"fk_user_number" db:"fk_user_number"`
	FKGussID   int64     `json:"fk_guss_number" db:"fk_guss_number"`
	VisitTime  time.Time `json:"visit_time"     db:"visit_time"` // 30분 단위 시간
	RevsTime   time.Time `json:"revs_time"      db:"revs_time"`
	RevsStatus string    `json:"revs_status"    db:"revs_status"`
	UserName   string    `json:"user_name,omitempty"` // 로그 출력용
}

// Admin: 관리자 정보
type Admin struct {
	AdminNumber int64         `json:"admin_number"   db:"admin_number"`
	AdminID     string        `json:"admin_id"       db:"admin_id"`
	AdminPW     string        `json:"admin_pw"       db:"admin_pw"`
	FKGussID    sql.NullInt64 `json:"fk_guss_number" db:"fk_guss_number"` // 특정 지점 관리 혹은 전체
}

type Sale struct {
	SalesNumber int64     `json:"sales_number" db:"sales_number"`
	SalesDate   time.Time `json:"sales_date"   db:"sales_date"`
	SalesAmount int       `json:"sales_amount" db:"sales_amount"`
	SalesType   string    `json:"sales_type"   db:"sales_type"`
}

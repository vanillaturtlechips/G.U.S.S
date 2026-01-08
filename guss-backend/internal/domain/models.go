package domain

import "time"

// User: 사용자 정보
type User struct {
	UserNumber int64  `json:"user_number"`
	UserName   string `json:"user_name"`
	UserPhone  string `json:"user_phone"`
	UserID     string `json:"user_id"`
	UserPW     string `json:"-"` // 비밀번호는 JSON 응답에서 제외
}

// Gym: 체육관 정보
type Gym struct {
	GussNumber    int64  `json:"guss_number"`
	GussName      string `json:"guss_name"`
	GussAddress   string `json:"guss_address"`
	GussPhone     string `json:"guss_phone"`
	GussStatus    string `json:"guss_status"` // open/close
	GussUserCount int    `json:"guss_user_count"`
	GussSize      int    `json:"guss_size"`

	GussMaType    string `json:"guss_ma_type"`   // 관리 기구 타입
	GussMaCount   int    `json:"guss_ma_count"`  // 관리 기구 수
	GussMaState   string `json:"guss_ma_state"`  // 기구 상태
}

// Reservation: 예약 정보
type Reservation struct {
	RevsNumber int64     `json:"revs_number"`
	FKUserID   int64     `json:"fk_user_number"`
	FKGussID   int64     `json:"fk_guss_number"`
	RevsTime   time.Time `json:"revs_time"`
	RevsStatus string    `json:"revs_status"`
}

// Admin: 관리자 정보
type Admin struct {
	AdminNumber int64  `json:"admin_number"`
	AdminID     string `json:"admin_id"`
	FKGussID    int64  `json:"fk_guss_id"`
}
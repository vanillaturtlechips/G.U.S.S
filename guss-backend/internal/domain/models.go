package domain

import "time"

type User struct {
	UserNumber int64  `json:"user_number" db:"user_number"`
	UserName   string `json:"user_name"   db:"user_name"`
	UserPhone  string `json:"user_phone"  db:"user_phone"`
	UserID     string `json:"user_id"     db:"user_id"`
	UserPW     string `json:"-"           db:"password"`
}

type Gym struct {
	GussNumber    int64  `json:"guss_number"`
	GussName      string `json:"guss_name"`
	GussAddress   string `json:"guss_address"`
	GussPhone     string `json:"guss_phone"`
	GussStatus    string `json:"guss_status"`
	GussUserCount int    `json:"guss_user_count"`
	GussSize      int    `json:"guss_size"`
}

type Equipment struct {
	ID           int64  `json:"id"`           // DB: id
	GymID        int64  `json:"gym_id"`       // DB: gym_id
	Name         string `json:"name"`         // DB: name
	Category     string `json:"category"`     // DB: category
	Quantity     int    `json:"quantity"`     // DB: quantity
	Status       string `json:"status"`       // DB: status
	PurchaseDate string `json:"purchaseDate"` // DB: purchase_date
}

type Reservation struct {
	RevsNumber int64     `json:"revs_number"`
	FKUserID   int64     `json:"fk_user_number"`
	FKGussID   int64     `json:"fk_guss_number"`
	RevsTime   time.Time `json:"revs_time"`
	RevsStatus string    `json:"revs_status"`
	MemberName string    `json:"member,omitempty"`
	Phone      string    `json:"phone,omitempty"`
}

type Admin struct {
	AdminNumber int64  `json:"admin_number"`
	AdminID     string `json:"admin_id"`
	AdminPW     string `json:"-"`
	FKGussID    int64  `json:"fk_guss_id"`
}

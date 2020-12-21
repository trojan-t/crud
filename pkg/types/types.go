package types

import (
	"errors"
	"time"
)

// Predefined errors
var (
	ErrNotFound        = errors.New("item not found")
	ErrInternal        = errors.New("internal error")
	ErrTokenNotFound   = errors.New("token not found")
	ErrNoSuchUser      = errors.New("no such user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrPhoneUsed       = errors.New("phone alredy registered")
	ErrTokenExpired    = errors.New("token expired")
)


// Manager is struct
type Manager struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Salary     int64     `json:"salary"`
	Plan       int64     `json:"plan"`
	BossID     int64     `json:"boss_id"`
	Department string    `json:"department"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	IsAdmin    bool      `json:"is_admin"`
	Created    time.Time `json:"created"`
}

// Product is struct
type Product struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Price   int       `json:"price"`
	Qty     int       `json:"qty"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

// Sale is struct
type Sale struct {
	ID         int64           `json:"id"`
	ManagerID  int64           `json:"manager_id"`
	CustomerID int64           `json:"customer_id"`
	Created    time.Time       `json:"created"`
	Positions  []*SalePosition `json:"positions"`
}

// SalePosition is struct
type SalePosition struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	SaleID    int64     `json:"sale_id"`
	Price     int       `json:"price"`
	Qty       int       `json:"qty"`
	Created   time.Time `json:"created"`
}

// Customer is struct
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}
